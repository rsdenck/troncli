package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/mascli/troncli/internal/policy"
)

type OllamaAdapter struct {
	Model               string
	Registry            *CapabilityRegistry
	BaseURL             string
	PolicyEngine        *policy.PolicyEngine
	NotificationManager *NotificationManager
}

func NewOllamaAdapter(model, registryPath string) (*OllamaAdapter, error) {
	reg := NewCapabilityRegistry(registryPath)
	if err := reg.Load(); err != nil {
		return nil, fmt.Errorf("failed to load capabilities: %w", err)
	}

	policyEngine := policy.NewPolicyEngine()

	// Setup notification manager
	home, _ := os.UserHomeDir()
	logFile := home + "/.troncli/agent_notifications.log"
	notifManager := NewNotificationManager(logFile)

	baseURL := os.Getenv("OLLAMA_HOST")
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}

	// Default to qwen3-coder if not specified
	if model == "" {
		model = "qwen2.5-coder:latest"
	}

	return &OllamaAdapter{
		Model:               model,
		Registry:            reg,
		BaseURL:             baseURL,
		PolicyEngine:        policyEngine,
		NotificationManager: notifManager,
	}, nil
}

func (a *OllamaAdapter) Name() string {
	return "ollama"
}

type ollamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type ollamaResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

func (a *OllamaAdapter) SendPrompt(ctx context.Context, prompt string) (string, error) {
	reqBody := ollamaRequest{
		Model:  a.Model,
		Prompt: prompt,
		Stream: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/api/generate", a.BaseURL), bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("ollama connection failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ollama API error (status %d): %s", resp.StatusCode, string(body))
	}

	var ollamaResp ollamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return ollamaResp.Response, nil
}

func (a *OllamaAdapter) ExecuteIntent(ctx context.Context, intent string) (string, error) {
	// Check if intent is allowed
	if !a.Registry.IsIntentAllowed(intent) {
		a.NotificationManager.PolicyViolationNotification(intent, "Intent not allowed by policy")
		return "", fmt.Errorf("intent '%s' is not allowed by policy", intent)
	}

	// Show system info for context
	sysInfo := a.PolicyEngine.GetSystemInfo()
	a.NotificationManager.SystemStatusNotification(
		fmt.Sprintf("System: %s@%s (Kernel: %s)",
			sysInfo["user"], sysInfo["hostname"], sysInfo["kernel"]))

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
		a.NotificationManager.PostExecutionNotification(intent, "LLM Query", "", err)
		return "", fmt.Errorf("failed to get command from LLM: %w", err)
	}

	// Clean and validate command
	command = strings.TrimSpace(command)
	if command == "" {
		return "", fmt.Errorf("LLM returned empty command")
	}

	// Show what will be executed
	a.NotificationManager.PreExecutionWarning(intent, command)
	fmt.Printf("🤖 AI Agent intends to execute: %s\n", command)
	fmt.Printf("⚡ Intent: %s\n", intent)

	// Ask for confirmation for high-risk operations
	if a.Registry.Capabilities.RiskLevel == "high" {
		fmt.Printf("⚠️  High-risk operation detected. Continue? (y/N): ")
		var response string
		fmt.Scanln(&response)
		if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
			a.NotificationManager.Notify(NotificationWarning, "Operation Cancelled",
				fmt.Sprintf("User cancelled operation: %s", intent), "TRONCLI-Agent")
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
		a.NotificationManager.PolicyViolationNotification(command, err.Error())
		return "", fmt.Errorf("policy engine blocked execution: %w", err)
	}

	a.NotificationManager.PostExecutionNotification(intent, command, "Success", nil)
	fmt.Printf("✅ Command executed successfully\n")
	return fmt.Sprintf("Executed: %s", command), nil
}
