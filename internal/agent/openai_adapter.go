package agent

import (
	"context"
	"fmt"
)

type OpenAIAdapter struct {
	APIKey   string
	Model    string
	Registry *CapabilityRegistry
}

func NewOpenAIAdapter(apiKey, model, registryPath string) *OpenAIAdapter {
	reg := NewCapabilityRegistry(registryPath)
	if err := reg.Load(); err != nil {
		fmt.Printf("Warning: Failed to load capabilities: %v\n", err)
	}
	return &OpenAIAdapter{
		APIKey:   apiKey,
		Model:    model,
		Registry: reg,
	}
}

func (a *OpenAIAdapter) Name() string {
	return "openai"
}

func (a *OpenAIAdapter) SendPrompt(ctx context.Context, prompt string) (string, error) {
	// TODO: Implement OpenAI API call
	return fmt.Sprintf("Thinking with OpenAI model %s...", a.Model), nil
}

func (a *OpenAIAdapter) ExecuteIntent(ctx context.Context, intent string) (string, error) {
	if !a.Registry.IsIntentAllowed(intent) {
		return "", fmt.Errorf("intent '%s' is not allowed by policy", intent)
	}
	return fmt.Sprintf("Executing allowed intent via OpenAI: %s", intent), nil
}
