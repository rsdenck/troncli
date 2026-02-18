package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type OllamaAdapter struct {
	Model    string
	Registry *CapabilityRegistry
	BaseURL  string
}

func NewOllamaAdapter(model, registryPath string) (*OllamaAdapter, error) {
	reg := NewCapabilityRegistry(registryPath)
	if err := reg.Load(); err != nil {
		return nil, fmt.Errorf("failed to load capabilities: %w", err)
	}

	baseURL := os.Getenv("OLLAMA_HOST")
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}

	return &OllamaAdapter{
		Model:    model,
		Registry: reg,
		BaseURL:  baseURL,
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
	// In a real implementation, the LLM would classify the intent first.
	// We send the intent as a prompt to get the command.
	
	// Check if intent is allowed
	if !a.Registry.IsIntentAllowed(intent) {
		return "", fmt.Errorf("intent '%s' is not allowed by policy", intent)
	}

	prompt := fmt.Sprintf("You are a Linux command generator. Output ONLY the shell command to execute the following intent, with no markdown, no explanations. Intent: %s", intent)
	
	command, err := a.SendPrompt(ctx, prompt)
	if err != nil {
		return "", err
	}

	return command, nil
}
