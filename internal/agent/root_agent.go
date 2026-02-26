package agent

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/mascli/troncli/internal/console"
	"github.com/mascli/troncli/internal/policy"
)

// AgentResponse represents the structured JSON response from the LLM
type AgentResponse struct {
	Analysis             string   `json:"analysis"`
	Commands             []string `json:"commands"`
	Risk                 string   `json:"risk"`
	Impact               string   `json:"impact"`
	ConfirmationRequired bool     `json:"confirmation_required"`
	Reasoning            string   `json:"reasoning"`
}

// RootAgent is the hardcore autonomous agent for TRONCLI
type RootAgent struct {
	ModelPath    string
	LlamaPath    string
	PolicyEngine *policy.PolicyEngine
	Streaming    bool
	MaxTokens    int
	Temperature  float64
	Threads      int
}

// NewRootAgent creates a new root agent instance
func NewRootAgent(modelPath, llamaPath string) *RootAgent {
	return &RootAgent{
		ModelPath:    modelPath,
		LlamaPath:    llamaPath,
		PolicyEngine: policy.NewPolicyEngine(),
		Streaming:    true,
		MaxTokens:    2048,
		Temperature:  0.2, // Low temperature for precise commands
		Threads:      8,
	}
}

// BuildSystemPrompt creates the hardcore system prompt for the root agent
func (a *RootAgent) BuildSystemPrompt(userIntent string) string {
	sysInfo := a.PolicyEngine.GetSystemInfo()

	prompt := fmt.Sprintf(`You are the TRON ROOT AGENT - the official autonomous AI agent for TRONCLI.

SYSTEM CONTEXT:
- Hostname: %s
- User: %s
- Kernel: %s
- Distribution: Linux
- Shell: bash

CRITICAL RULES:
1. You can ONLY generate commands that start with: troncli
2. You MUST respond in VALID JSON format (no markdown, no code blocks)
3. Analyze risk level: low, medium, high, critical
4. Provide clear impact analysis
5. Request confirmation for high-risk operations

USER INTENT: %s

RESPONSE FORMAT (STRICT JSON):
{
  "analysis": "Brief analysis of what needs to be done",
  "commands": ["troncli command1", "troncli command2"],
  "risk": "low|medium|high|critical",
  "impact": "Description of system impact",
  "confirmation_required": true|false,
  "reasoning": "Why these commands were chosen"
}

Respond with ONLY the JSON object, nothing else:`, 
		sysInfo["hostname"], sysInfo["user"], sysInfo["kernel"], userIntent)

	return prompt
}

// Execute runs the root agent with the given intent
func (a *RootAgent) Execute(ctx context.Context, userIntent string) error {
	// Validate llama.cpp binary
	if _, err := os.Stat(a.LlamaPath); os.IsNotExist(err) {
		return fmt.Errorf("❌ llama.cpp binary not found at: %s\n\nInstall with:\n  git clone https://github.com/ggerganov/llama.cpp\n  cd llama.cpp && make", a.LlamaPath)
	}

	// Validate model
	if _, err := os.Stat(a.ModelPath); os.IsNotExist(err) {
		return fmt.Errorf("❌ Model not found at: %s\n\nDownload with:\n  wget -O %s https://huggingface.co/Qwen/Qwen2.5-Coder-7B-Instruct-GGUF/resolve/main/qwen2.5-coder-7b-instruct-q4_0.gguf", a.ModelPath, a.ModelPath)
	}

	// Display agent header
	a.displayHeader(userIntent)

	// Build prompt
	prompt := a.BuildSystemPrompt(userIntent)

	// Get LLM response
	fmt.Printf("\n%s🧠 Analyzing intent...%s\n", console.ColorCyan, console.ColorReset)
	response, err := a.queryLLM(ctx, prompt)
	if err != nil {
		return fmt.Errorf("❌ LLM query failed: %w", err)
	}

	// Parse JSON response
	var agentResp AgentResponse
	if err := json.Unmarshal([]byte(response), &agentResp); err != nil {
		// Try to extract JSON from response if it's wrapped
		response = a.extractJSON(response)
		if err := json.Unmarshal([]byte(response), &agentResp); err != nil {
			return fmt.Errorf("❌ Failed to parse agent response as JSON: %w\n\nRaw response:\n%s", err, response)
		}
	}

	// Display analysis
	a.displayAnalysis(&agentResp)

	// Check if confirmation is required
	if agentResp.ConfirmationRequired {
		if !a.requestConfirmation(&agentResp) {
			fmt.Printf("\n%s⚠️  Operation cancelled by user%s\n", console.ColorYellow, console.ColorReset)
			return nil
		}
	}

	// Execute commands
	return a.executeCommands(ctx, &agentResp)
}

// queryLLM sends the prompt to llama.cpp and returns the response
func (a *RootAgent) queryLLM(ctx context.Context, prompt string) (string, error) {
	args := []string{
		"-m", a.ModelPath,
		"-p", prompt,
		"-n", fmt.Sprintf("%d", a.MaxTokens),
		"--ctx-size", "4096",
		"-t", fmt.Sprintf("%d", a.Threads),
		"--temp", fmt.Sprintf("%.2f", a.Temperature),
		"--repeat-penalty", "1.1",
		"-c", "4096",
		"--no-display-prompt", // Don't echo the prompt
	}

	cmd := exec.CommandContext(ctx, a.LlamaPath, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Set environment for optimal performance
	cmd.Env = append(os.Environ(),
		"LLAMA_NATIVE=1",
	)

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("llama.cpp execution failed: %w\nStderr: %s", err, stderr.String())
	}

	// Clean output
	output := a.cleanLlamaOutput(stdout.String())
	return output, nil
}

// cleanLlamaOutput removes llama.cpp debug output
func (a *RootAgent) cleanLlamaOutput(output string) string {
	lines := strings.Split(output, "\n")
	var cleanLines []string

	for _, line := range lines {
		// Skip debug lines
		if strings.Contains(line, "llm_load") ||
			strings.Contains(line, "ggml_") ||
			strings.Contains(line, "AVX") ||
			strings.Contains(line, "sampling") ||
			strings.Contains(line, "llama_") ||
			strings.TrimSpace(line) == "" {
			continue
		}
		cleanLines = append(cleanLines, line)
	}

	return strings.Join(cleanLines, "\n")
}

// extractJSON tries to extract JSON from a response that might have extra text
func (a *RootAgent) extractJSON(response string) string {
	// Find first { and last }
	start := strings.Index(response, "{")
	end := strings.LastIndex(response, "}")

	if start != -1 && end != -1 && end > start {
		return response[start : end+1]
	}

	return response
}

// displayHeader shows the agent header
func (a *RootAgent) displayHeader(intent string) {
	table := console.NewBoxTable(os.Stdout)
	table.SetTitle("TRON ROOT AGENT › AUTONOMOUS MODE")
	table.AddRow([]string{"Intent", intent})
	table.AddRow([]string{"Model", "Qwen2.5-Coder-7B"})
	table.AddRow([]string{"Engine", "llama.cpp"})
	table.AddRow([]string{"Mode", "Hardcore Linux"})
	table.RenderKeyValue()
}

// displayAnalysis shows the agent's analysis
func (a *RootAgent) displayAnalysis(resp *AgentResponse) {
	fmt.Printf("\n%s┌── AGENT ANALYSIS ────────────────────────────────────────┐%s\n", console.ColorCyan, console.ColorReset)
	fmt.Printf("%s│%s                                                          %s│%s\n", console.ColorCyan, console.ColorReset, console.ColorCyan, console.ColorReset)
	fmt.Printf("%s│%s  %s%s%s\n", console.ColorCyan, console.ColorReset, resp.Analysis, strings.Repeat(" ", 56-len(resp.Analysis)), console.ColorCyan+"│"+console.ColorReset)
	fmt.Printf("%s│%s                                                          %s│%s\n", console.ColorCyan, console.ColorReset, console.ColorCyan, console.ColorReset)
	fmt.Printf("%s└──────────────────────────────────────────────────────────┘%s\n", console.ColorCyan, console.ColorReset)

	// Risk level with color
	riskColor := console.ColorGreen
	switch strings.ToLower(resp.Risk) {
	case "medium":
		riskColor = console.ColorYellow
	case "high":
		riskColor = console.ColorRed
	case "critical":
		riskColor = console.ColorRed + console.ColorBold
	}

	table := console.NewBoxTable(os.Stdout)
	table.SetTitle("RISK ASSESSMENT")
	table.AddRow([]string{"Risk Level", fmt.Sprintf("%s%s%s", riskColor, strings.ToUpper(resp.Risk), console.ColorReset)})
	table.AddRow([]string{"Impact", resp.Impact})
	table.AddRow([]string{"Confirmation", fmt.Sprintf("%v", resp.ConfirmationRequired)})
	table.RenderKeyValue()

	// Commands
	if len(resp.Commands) > 0 {
		fmt.Printf("\n%s┌── COMMANDS TO EXECUTE ───────────────────────────────────┐%s\n", console.ColorCyan, console.ColorReset)
		fmt.Printf("%s│%s                                                          %s│%s\n", console.ColorCyan, console.ColorReset, console.ColorCyan, console.ColorReset)
		for i, cmd := range resp.Commands {
			padding := 54 - len(cmd)
			if padding < 0 {
				padding = 0
			}
			fmt.Printf("%s│%s  %d. %s%s%s  %s│%s\n", 
				console.ColorCyan, console.ColorReset, i+1, 
				console.ColorGreen, cmd, console.ColorReset,
				strings.Repeat(" ", padding), console.ColorCyan+"│"+console.ColorReset)
		}
		fmt.Printf("%s│%s                                                          %s│%s\n", console.ColorCyan, console.ColorReset, console.ColorCyan, console.ColorReset)
		fmt.Printf("%s└──────────────────────────────────────────────────────────┘%s\n\n", console.ColorCyan, console.ColorReset)
	}

	// Reasoning
	if resp.Reasoning != "" {
		fmt.Printf("%s💡 Reasoning:%s %s\n\n", console.ColorCyan, console.ColorReset, resp.Reasoning)
	}
}

// requestConfirmation asks the user for confirmation
func (a *RootAgent) requestConfirmation(resp *AgentResponse) bool {
	fmt.Printf("%s⚠️  CONFIRMATION REQUIRED%s\n", console.ColorYellow+console.ColorBold, console.ColorReset)
	fmt.Printf("This operation has been classified as %s%s%s risk.\n", 
		console.ColorRed, strings.ToUpper(resp.Risk), console.ColorReset)
	fmt.Printf("\nDo you want to proceed? (yes/no): ")

	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))

	return response == "yes" || response == "y"
}

// executeCommands executes the commands from the agent response
func (a *RootAgent) executeCommands(ctx context.Context, resp *AgentResponse) error {
	fmt.Printf("%s🚀 Executing commands...%s\n\n", console.ColorGreen+console.ColorBold, console.ColorReset)

	for i, command := range resp.Commands {
		fmt.Printf("%s[%d/%d]%s Executing: %s%s%s\n", 
			console.ColorCyan, i+1, len(resp.Commands), console.ColorReset,
			console.ColorGreen, command, console.ColorReset)

		// Execute with timeout
		execCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
		defer cancel()

		cmd := exec.CommandContext(execCtx, "bash", "-c", command)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			fmt.Printf("%s❌ Command failed: %v%s\n", console.ColorRed, err, console.ColorReset)
			return fmt.Errorf("command execution failed: %w", err)
		}

		fmt.Printf("%s✅ Command completed successfully%s\n\n", console.ColorGreen, console.ColorReset)
	}

	fmt.Printf("%s🎉 All commands executed successfully!%s\n", console.ColorGreen+console.ColorBold, console.ColorReset)
	return nil
}

// StreamingExecute runs the agent with streaming output (real-time response)
func (a *RootAgent) StreamingExecute(ctx context.Context, userIntent string) error {
	// Display header
	a.displayHeader(userIntent)

	// Build prompt
	prompt := a.BuildSystemPrompt(userIntent)

	fmt.Printf("\n%s🧠 Analyzing intent (streaming)...%s\n\n", console.ColorCyan, console.ColorReset)

	// Build command with streaming
	args := []string{
		"-m", a.ModelPath,
		"-p", prompt,
		"-n", fmt.Sprintf("%d", a.MaxTokens),
		"--ctx-size", "4096",
		"-t", fmt.Sprintf("%d", a.Threads),
		"--temp", fmt.Sprintf("%.2f", a.Temperature),
		"--repeat-penalty", "1.1",
		"-c", "4096",
		"--no-display-prompt",
	}

	cmd := exec.CommandContext(ctx, a.LlamaPath, args...)

	// Get stdout pipe for streaming
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %w", err)
	}

	// Start command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start llama.cpp: %w", err)
	}

	// Read streaming output
	var response strings.Builder
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		// Skip debug lines
		if strings.Contains(line, "llm_load") || strings.Contains(line, "ggml_") {
			continue
		}
		fmt.Print(line)
		response.WriteString(line + "\n")
	}

	// Wait for command to finish
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("llama.cpp execution failed: %w", err)
	}

	fmt.Println()

	// Parse and execute
	var agentResp AgentResponse
	responseStr := a.extractJSON(response.String())
	if err := json.Unmarshal([]byte(responseStr), &agentResp); err != nil {
		return fmt.Errorf("failed to parse agent response: %w", err)
	}

	// Display analysis and execute
	a.displayAnalysis(&agentResp)

	if agentResp.ConfirmationRequired {
		if !a.requestConfirmation(&agentResp) {
			return nil
		}
	}

	return a.executeCommands(ctx, &agentResp)
}
