package agent

// Package agent provides AI agent adapters and capabilities.

import "context"

// AgentAdapter defines the interface for AI agents
type AgentAdapter interface {
	// Name returns the agent name (e.g., "ollama", "openai")
	Name() string

	// SendPrompt sends a text prompt to the agent and returns the response
	SendPrompt(ctx context.Context, prompt string) (string, error)

	// ExecuteIntent translates a natural language request into a CLI command
	ExecuteIntent(ctx context.Context, intent string) (string, error)
}
