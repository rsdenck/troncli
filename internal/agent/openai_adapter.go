package agent

import "context"

type OpenAIAdapter struct {
	APIKey string
	Model  string
}

func NewOpenAIAdapter(apiKey, model string) *OpenAIAdapter {
	return &OpenAIAdapter{APIKey: apiKey, Model: model}
}

func (a *OpenAIAdapter) Name() string {
	return "openai"
}

func (a *OpenAIAdapter) SendPrompt(ctx context.Context, prompt string) (string, error) {
	// TODO: Implement OpenAI API call
	return "Thinking...", nil
}

func (a *OpenAIAdapter) ExecuteIntent(ctx context.Context, intent string) (string, error) {
	return "troncli help", nil
}
