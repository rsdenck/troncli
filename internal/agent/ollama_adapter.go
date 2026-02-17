package agent

import (
	"context"
	"fmt"
)

type OllamaAdapter struct {
	Model    string
	Registry *CapabilityRegistry
}

func NewOllamaAdapter(model, registryPath string) *OllamaAdapter {
	reg := NewCapabilityRegistry(registryPath)
	if err := reg.Load(); err != nil {
		fmt.Printf("Warning: Failed to load capabilities: %v\n", err)
	}
	return &OllamaAdapter{
		Model:    model,
		Registry: reg,
	}
}

func (a *OllamaAdapter) Name() string {
	return "ollama"
}

func (a *OllamaAdapter) SendPrompt(ctx context.Context, prompt string) (string, error) {
	// TODO: Implement Ollama API call
	return fmt.Sprintf("Thinking with model %s...", a.Model), nil
}

func (a *OllamaAdapter) ExecuteIntent(ctx context.Context, intent string) (string, error) {
	// In a real implementation, the LLM would classify the intent first.
	// For now, we assume the input string contains the intent key if it matches one.
	// This is a simplification for the prototype.

	// Check if intent is allowed
	if !a.Registry.IsIntentAllowed(intent) {
		return "", fmt.Errorf("intent '%s' is not allowed by policy", intent)
	}

	// Mock implementation of returning a command based on intent
	switch intent {
	case "install_package":
		return "apt-get install -y <package>", nil
	case "audit_security":
		return "lynis audit system", nil
	default:
		return "", fmt.Errorf("unknown intent: %s", intent)
	}
}
