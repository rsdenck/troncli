package agent

import (
	"context"
	"fmt"
)

type ClaudeAdapter struct {
	APIKey   string
	Registry *CapabilityRegistry
}

func NewClaudeAdapter(apiKey, registryPath string) *ClaudeAdapter {
	reg := NewCapabilityRegistry(registryPath)
	if err := reg.Load(); err != nil {
		fmt.Printf("Warning: Failed to load capabilities: %v\n", err)
	}
	return &ClaudeAdapter{
		APIKey:   apiKey,
		Registry: reg,
	}
}

func (a *ClaudeAdapter) Name() string {
	return "claude"
}

func (a *ClaudeAdapter) SendPrompt(ctx context.Context, prompt string) (string, error) {
	// TODO: Implement Claude API call
	return "Thinking with Claude...", nil
}

func (a *ClaudeAdapter) ExecuteIntent(ctx context.Context, intent string) (string, error) {
	if !a.Registry.IsIntentAllowed(intent) {
		return "", fmt.Errorf("intent '%s' is not allowed by policy", intent)
	}
	return fmt.Sprintf("Executing allowed intent via Claude: %s", intent), nil
}
