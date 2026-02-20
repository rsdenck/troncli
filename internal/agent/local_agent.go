package agent

import (
	"context"
	"fmt"
)

type LocalAgent struct {
	// e.g., using a local lightweight model or heuristics
	Registry *CapabilityRegistry
}

func NewLocalAgent(registryPath string) (*LocalAgent, error) {
	reg := NewCapabilityRegistry(registryPath)
	if err := reg.Load(); err != nil {
		return nil, fmt.Errorf("failed to load capabilities: %w", err)
	}
	return &LocalAgent{Registry: reg}, nil
}

func (a *LocalAgent) Name() string {
	return "local"
}

func (a *LocalAgent) SendPrompt(ctx context.Context, prompt string) (string, error) {
	return "", fmt.Errorf("local agent implementation not available")
}

func (a *LocalAgent) ExecuteIntent(ctx context.Context, intent string) (string, error) {
	if !a.Registry.IsIntentAllowed(intent) {
		return "", fmt.Errorf("intent '%s' is not allowed by policy", intent)
	}
	// Real implementation requires a local reasoning engine which is not yet integrated.
	// Abort rather than mock.
	return "", fmt.Errorf("local execution engine not ready for intent: %s", intent)
}
