package agent

import (
	"context"
	"fmt"
)

type LocalAgent struct {
	// e.g., using a local lightweight model or heuristics
	Registry *CapabilityRegistry
}

func NewLocalAgent(registryPath string) *LocalAgent {
	reg := NewCapabilityRegistry(registryPath)
	if err := reg.Load(); err != nil {
		// Log warning or handle error
		fmt.Printf("Warning: Failed to load capabilities: %v\n", err)
	}
	return &LocalAgent{Registry: reg}
}

func (a *LocalAgent) Name() string {
	return "local"
}

func (a *LocalAgent) SendPrompt(ctx context.Context, prompt string) (string, error) {
	return "Processing locally...", nil
}

func (a *LocalAgent) ExecuteIntent(ctx context.Context, intent string) (string, error) {
	if !a.Registry.IsIntentAllowed(intent) {
		return "", fmt.Errorf("intent '%s' is not allowed by policy", intent)
	}
	return fmt.Sprintf("Executing allowed intent: %s", intent), nil
}
