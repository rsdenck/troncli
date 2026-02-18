package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ClaudeAdapter struct {
	APIKey   string
	Registry *CapabilityRegistry
}

func NewClaudeAdapter(apiKey, registryPath string) (*ClaudeAdapter, error) {
	reg := NewCapabilityRegistry(registryPath)
	if err := reg.Load(); err != nil {
		return nil, fmt.Errorf("failed to load capabilities: %w", err)
	}
	return &ClaudeAdapter{
		APIKey:   apiKey,
		Registry: reg,
	}, nil
}

func (a *ClaudeAdapter) Name() string {
	return "claude"
}

type claudeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type claudeRequest struct {
	Model     string          `json:"model"`
	Messages  []claudeMessage `json:"messages"`
	MaxTokens int             `json:"max_tokens"`
}

type claudeContent struct {
	Text string `json:"text"`
	Type string `json:"type"`
}

type claudeResponse struct {
	Content []claudeContent `json:"content"`
}

func (a *ClaudeAdapter) SendPrompt(ctx context.Context, prompt string) (string, error) {
	if a.APIKey == "" {
		return "", fmt.Errorf("claude API key is missing")
	}

	reqBody := claudeRequest{
		Model: "claude-3-opus-20240229", // Defaulting to Opus for now
		Messages: []claudeMessage{
			{Role: "user", Content: prompt},
		},
		MaxTokens: 1024,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", a.APIKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Claude connection failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Claude API error (status %d): %s", resp.StatusCode, string(body))
	}

	var claudeResp claudeResponse
	if err := json.NewDecoder(resp.Body).Decode(&claudeResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(claudeResp.Content) == 0 {
		return "", fmt.Errorf("no content in Claude response")
	}

	return claudeResp.Content[0].Text, nil
}

func (a *ClaudeAdapter) ExecuteIntent(ctx context.Context, intent string) (string, error) {
	if !a.Registry.IsIntentAllowed(intent) {
		return "", fmt.Errorf("intent '%s' is not allowed by policy", intent)
	}

	prompt := fmt.Sprintf("You are a Linux command generator. Output ONLY the shell command to execute the following intent, with no markdown, no explanations. Intent: %s", intent)

	return a.SendPrompt(ctx, prompt)
}
