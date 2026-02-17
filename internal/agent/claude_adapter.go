package agent

import "context"

type ClaudeAdapter struct {
	APIKey string
}

func NewClaudeAdapter(apiKey string) *ClaudeAdapter {
	return &ClaudeAdapter{APIKey: apiKey}
}

func (a *ClaudeAdapter) Name() string {
	return "claude"
}

func (a *ClaudeAdapter) SendPrompt(ctx context.Context, prompt string) (string, error) {
	// TODO: Implement Claude API call
	return "Thinking...", nil
}

func (a *ClaudeAdapter) ExecuteIntent(ctx context.Context, intent string) (string, error) {
	return "troncli help", nil
}
