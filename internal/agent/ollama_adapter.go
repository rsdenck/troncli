package agent

import "context"

type OllamaAdapter struct {
	Model string
}

func NewOllamaAdapter(model string) *OllamaAdapter {
	return &OllamaAdapter{Model: model}
}

func (a *OllamaAdapter) Name() string {
	return "ollama"
}

func (a *OllamaAdapter) SendPrompt(ctx context.Context, prompt string) (string, error) {
	// TODO: Implement Ollama API call
	return "Thinking...", nil
}

func (a *OllamaAdapter) ExecuteIntent(ctx context.Context, intent string) (string, error) {
	return "troncli help", nil
}
