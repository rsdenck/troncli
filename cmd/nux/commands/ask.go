package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/rsdenck/nux/internal/output"
	"github.com/rsdenck/nux/internal/vault"
	"github.com/spf13/cobra"
)

var askCmd = &cobra.Command{
	Use:   "ask",
	Short: "Ask questions to AI providers (Ollama, OpenAI, Claude)",
	Long:  `Query AI models for code assistance, system analysis, and automation.`,
}

var askQueryCmd = &cobra.Command{
	Use:   "query <question>",
	Short: "Ask a question to AI",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		question := strings.Join(args, " ")

		provider, _ := cmd.Flags().GetString("provider")
		model, _ := cmd.Flags().GetString("model")

		if provider == "" {
			provider = "ollama"
		}
		if model == "" {
			model = "qwen3-coder"
		}

		switch provider {
		case "ollama":
			askOllamaFunc(question, model, provider)
		case "openai":
			askOpenAI(question, model, provider)
		case "claude":
			askClaude(question, model, provider)
		default:
			output.NewError(fmt.Sprintf("unknown provider: %s", provider), "ASK_INVALID_PROVIDER").Print()
		}
	},
}

func askOllamaFunc(question, model, provider string) {
	v, err := vault.Load()
	if err != nil {
		output.NewError("failed to load vault, using defaults", "VAULT_ERROR").Print()
	}

	host := "http://localhost:11434"
	if v != nil && v.Config != nil {
		if h, ok := v.Config["ollama_host"].(string); ok {
			host = h
		}
	}

	url := fmt.Sprintf("%s/api/generate", host)

	payload := map[string]interface{}{
		"model":  model,
		"prompt": question,
		"stream": false,
	}

	data, _ := json.Marshal(payload)
	resp, err := http.Post(url, "application/json", strings.NewReader(string(data)))
	if err != nil {
		output.NewError(fmt.Sprintf("failed to connect to Ollama: %s", err.Error()), "ASK_OLLAMA_ERROR").Print()
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		output.NewError("failed to parse Ollama response", "ASK_PARSE_ERROR").Print()
		return
	}

	if response, ok := result["response"].(string); ok {
		output.NewSuccess(map[string]interface{}{
			"provider": provider,
			"model":    model,
			"response": response,
		}).WithMessage("Ollama Response").Print()
	} else {
		output.NewError("no response from Ollama", "ASK_NO_RESPONSE").Print()
	}
}

func askOpenAI(question, model, provider string) {
	v, err := vault.Load()
	if err != nil {
		output.NewError("failed to load vault", "VAULT_ERROR").Print()
		return
	}

	apiKey, ok := v.GetAPIKey("openai")
	if !ok {
		output.NewError("OpenAI API key not found. Use: nux vault set-key openai <key>", "ASK_NO_API_KEY").Print()
		return
	}

	url := "https://api.openai.com/v1/chat/completions"

	payload := map[string]interface{}{
		"model": model,
		"messages": []map[string]string{
			{"role": "user", "content": question},
		},
	}

	data, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", url, strings.NewReader(string(data)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		output.NewError(fmt.Sprintf("failed to connect to OpenAI: %s", err.Error()), "ASK_OPENAI_ERROR").Print()
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		output.NewError("failed to parse OpenAI response", "ASK_PARSE_ERROR").Print()
		return
	}

	if choices, ok := result["choices"].([]interface{}); ok && len(choices) > 0 {
		if choice, ok := choices[0].(map[string]interface{}); ok {
			if message, ok := choice["message"].(map[string]interface{}); ok {
				if content, ok := message["content"].(string); ok {
					output.NewSuccess(map[string]interface{}{
						"provider": "openai",
						"model":    model,
						"response": content,
					}).WithMessage("OpenAI Response").Print()
					return
				}
			}
		}
	}

	output.NewError("no response from OpenAI", "ASK_NO_RESPONSE").Print()
}

func askClaude(question, model, provider string) {
	output.NewInfo("Claude provider not yet implemented").Print()
}

var askConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure AI providers",
	Run: func(cmd *cobra.Command, args []string) {
		provider, _ := cmd.Flags().GetString("provider")
		apiKey, _ := cmd.Flags().GetString("api-key")
		host, _ := cmd.Flags().GetString("host")

		v, err := vault.Load()
		if err != nil {
			v = vault.NewVault()
		}

		if provider == "openai" && apiKey != "" {
			v.SetAPIKey("openai", apiKey)
		}

		if provider == "ollama" && host != "" {
			if v.Config == nil {
				v.Config = make(map[string]interface{})
			}
			v.Config["ollama_host"] = host
		}

		if err := vault.Save(v); err != nil {
			output.NewError(fmt.Sprintf("failed to save config: %s", err.Error()), "ASK_CONFIG_ERROR").Print()
			return
		}

		output.NewSuccess(map[string]interface{}{
			"provider": provider,
			"status":   "configured",
		}).Print()
	},
}

func init() {
	askQueryCmd.Flags().String("provider", "", "AI provider (ollama, openai, claude)")
	askQueryCmd.Flags().String("model", "", "Model name (default: qwen3-coder for ollama)")
	askCmd.AddCommand(askQueryCmd)

	askConfigCmd.Flags().String("provider", "", "Provider to configure")
	askConfigCmd.Flags().String("api-key", "", "API key for the provider")
	askConfigCmd.Flags().String("host", "", "Host for Ollama")
	askCmd.AddCommand(askConfigCmd)

	rootCmd.AddCommand(askCmd)
}
