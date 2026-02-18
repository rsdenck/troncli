package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type OpenAIAdapter struct {
	APIKey   string
	Model    string
	Registry *CapabilityRegistry
}

func NewOpenAIAdapter(apiKey, model, registryPath string) (*OpenAIAdapter, error) {
	reg := NewCapabilityRegistry(registryPath)
	if err := reg.Load(); err != nil {
		return nil, fmt.Errorf("failed to load capabilities: %w", err)
	}
	return &OpenAIAdapter{
		APIKey:   apiKey,
		Model:    model,
		Registry: reg,
	}, nil
}

func (a *OpenAIAdapter) Name() string {
	return "openai"
}

type openAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIRequest struct {
	Model    string          `json:"model"`
	Messages []openAIMessage `json:"messages"`
}

type openAIChoice struct {
	Message openAIMessage `json:"message"`
}

type openAIResponse struct {
	Choices []openAIChoice `json:"choices"`
}

func (a *OpenAIAdapter) SendPrompt(ctx context.Context, prompt string) (string, error) {
	if a.APIKey == "" {
		return "", fmt.Errorf("OpenAI API key is missing")
	}

	reqBody := openAIRequest{
		Model: a.Model,
		Messages: []openAIMessage{
			{Role: "user", Content: prompt},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+a.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("OpenAI connection failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("OpenAI API error (status %d): %s", resp.StatusCode, string(body))
	}

	var openAIResp openAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&openAIResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(openAIResp.Choices) == 0 {
		return "", fmt.Errorf("no response choices from OpenAI")
	}

	return openAIResp.Choices[0].Message.Content, nil
}

func (a *OpenAIAdapter) ExecuteIntent(ctx context.Context, intent string) (string, error) {
	if !a.Registry.IsIntentAllowed(intent) {
		return "", fmt.Errorf("intent '%s' is not allowed by policy", intent)
	}

	prompt := fmt.Sprintf("You are a Linux command generator. Output ONLY the shell command to execute the following intent, with no markdown, no explanations. Intent: %s", intent)

	return a.SendPrompt(ctx, prompt)
}
