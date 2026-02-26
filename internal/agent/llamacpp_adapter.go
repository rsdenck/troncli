package agent

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/mascli/troncli/internal/policy"
)

type LlamaCppAdapter struct {
	ModelPath    string
	LlamaPath    string
	Registry     *CapabilityRegistry
	PolicyEngine *policy.PolicyEngine
}

func NewLlamaCppAdapter(modelPath, llamaPath, registryPath string) (*LlamaCppAdapter, error) {
	reg := NewCapabilityRegistry(registryPath)
	if err := reg.Load(); err != nil {
		return nil, fmt.Errorf("failed to load capabilities: %w", err)
	}

	policyEngine := policy.NewPolicyEngine()

	return &LlamaCppAdapter{
		ModelPath:    modelPath,
		LlamaPath:    llamaPath,
		Registry:     reg,
		PolicyEngine: policyEngine,
	}, nil
}

func (a *LlamaCppAdapter) Name() string {
	return "llamacpp"
}

func (a *LlamaCppAdapter) SendPrompt(ctx context.Context, prompt string) (string, error) {
	// Check if llama.cpp binary exists
	if _, err := os.Stat(a.LlamaPath); os.IsNotExist(err) {
		return "", fmt.Errorf("llama.cpp binary not found at: %s", a.LlamaPath)
	}

	// Check if model exists
	if _, err := os.Stat(a.ModelPath); os.IsNotExist(err) {
		return "", fmt.Errorf("model file not found at: %s", a.ModelPath)
	}

	// Build llama.cpp command
	args := []string{
		"-m", a.ModelPath,
		"-p", prompt,
		"-n", "512", // max tokens
		"--ctx-size", "2048",
		"-t", "4", // threads
		"--temp", "0.7",
		"--repeat-penalty", "1.1",
		"-c", "2048", // context size
	}

	cmd := exec.CommandContext(ctx, a.LlamaPath, args...)

	// Set environment variables for optimal performance
	cmd.Env = append(os.Environ(),
		"OLLAMA_NUM_PARALLEL=1",
		"OLLAMA_MAX_LOADED_MODELS=1",
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("llama.cpp execution failed: %w\nOutput: %s", err, string(output))
	}

	// Clean up the output (remove llama.cpp formatting)
	response := a.cleanOutput(string(output))

	return response, nil
}

func (a *LlamaCppAdapter) ExecuteIntent(ctx context.Context, intent string) (string, error) {
	// Check if intent is allowed
	if !a.Registry.IsIntentAllowed(intent) {
		return "", fmt.Errorf("intent '%s' is not allowed by policy", intent)
	}

	// Build system prompt for qwen3-coder
	systemPrompt := `You are TRONCLI AI Agent, a Linux system administration expert.
You MUST respond with ONLY the exact shell command to execute the user's intent.
NO markdown, NO explanations, NO apologies, NO extra text.
Example:
Intent: "install nginx"
Response: apt install -y nginx

Now execute: `

	prompt := systemPrompt + intent

	// Get command from LLM
	command, err := a.SendPrompt(ctx, prompt)
	if err != nil {
		return "", fmt.Errorf("failed to get command from LLM: %w", err)
	}

	// Clean and validate command
	command = strings.TrimSpace(command)
	if command == "" {
		return "", fmt.Errorf("LLM returned empty command")
	}

	// Show system info for context
	sysInfo := a.PolicyEngine.GetSystemInfo()
	fmt.Printf("🖥️  System: %s@%s (Kernel: %s)\n",
		sysInfo["user"], sysInfo["hostname"], sysInfo["kernel"])

	// Show what will be executed
	fmt.Printf("🤖 AI Agent intends to execute: %s\n", command)
	fmt.Printf("⚡ Intent: %s\n", intent)

	// Ask for confirmation for high-risk operations
	if a.Registry.Capabilities.RiskLevel == "high" {
		fmt.Printf("⚠️  High-risk operation detected. Continue? (y/N): ")
		var response string
		fmt.Scanln(&response)
		if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
			return "", fmt.Errorf("operation cancelled by user")
		}
	}

	// Execute with policy engine
	fmt.Printf("🔧 Executing with policy engine...\n")

	// Add timeout for execution
	execCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	err = a.PolicyEngine.ExecuteWithPolicy(execCtx, command)
	if err != nil {
		return "", fmt.Errorf("policy engine blocked execution: %w", err)
	}

	fmt.Printf("✅ Command executed successfully\n")
	return fmt.Sprintf("Executed: %s", command), nil
}

func (a *LlamaCppAdapter) cleanOutput(output string) string {
	// Remove llama.cpp debug output and formatting
	lines := strings.Split(output, "\n")
	var cleanLines []string

	inResponse := false
	for _, line := range lines {
		// Skip llama.cpp debug lines
		if strings.Contains(line, "llm_load_tensors") ||
			strings.Contains(line, "llm_load_tensors") ||
			strings.Contains(line, "ggml_init") ||
			strings.Contains(line, "AVX") ||
			strings.Contains(line, "AVX2") ||
			strings.Contains(line, "AVX512") ||
			strings.Contains(line, "llm_build_attn") ||
			strings.Contains(line, "sampling") {
			continue
		}

		// Start collecting after we see the actual response
		if strings.TrimSpace(line) != "" && !inResponse {
			inResponse = true
		}

		if inResponse && strings.TrimSpace(line) != "" {
			cleanLines = append(cleanLines, strings.TrimSpace(line))
		}
	}

	return strings.Join(cleanLines, " ")
}

// DownloadModel downloads qwen3-coder model if not exists
func (a *LlamaCppAdapter) DownloadModel() error {
	if _, err := os.Stat(a.ModelPath); err == nil {
		return nil // Model already exists
	}

	fmt.Printf("📥 Downloading qwen3-coder model...\n")

	// Create directory if not exists
	modelDir := filepath.Dir(a.ModelPath)
	if err := os.MkdirAll(modelDir, 0755); err != nil {
		return fmt.Errorf("failed to create model directory: %w", err)
	}

	// Download using wget or curl
	downloadCmd := fmt.Sprintf("wget -O %s https://huggingface.co/Qwen/Qwen2.5-Coder-7B-Instruct-GGUF/resolve/main/qwen2.5-coder-7b-instruct-q4_0.gguf", a.ModelPath)

	cmd := exec.Command("sh", "-c", downloadCmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to download model: %w", err)
	}

	fmt.Printf("✅ Model downloaded to: %s\n", a.ModelPath)
	return nil
}
