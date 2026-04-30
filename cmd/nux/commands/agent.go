package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/rsdenck/nux/internal/core"
	"github.com/rsdenck/nux/internal/output"
	"github.com/rsdenck/nux/internal/vault"
	"github.com/spf13/cobra"
)

var agentExecutor core.Executor = &core.RealExecutor{}

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "AI agent for system analysis and automation",
	Long:  `Interact with AI providers (Ollama, OpenAI, Claude) for intelligent system management.`,
}

var agentStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show agent status and health",
	Run: func(cmd *cobra.Command, args []string) {
		// Load vault for configuration
		v, err := vault.Load()
		if err != nil {
			v = vault.NewVault()
		}

		// Get configuration
		provider := "ollama"
		host := "http://localhost:11434"
		model := "qwen3-coder"
		enabled := false

		if v != nil {
			if p, ok := v.Config["agent_provider"].(string); ok {
				provider = p
			}
			if h, ok := v.Config["ollama_host"].(string); ok {
				host = h
			}
			if m, ok := v.Config["agent_model"].(string); ok {
				model = m
			}
			if e, ok := v.Config["agent_enabled"].(bool); ok {
				enabled = e
			}
		}

		// Check if enabled
		if !enabled {
			output.NewInfo(map[string]interface{}{
				"provider": provider,
				"status":  "disabled",
			}).WithMessage("NUX Agent").Print()
			return
		}

		// Test connectivity and latency
		status := "offline"
		latency := "N/A"
		
		if provider == "ollama" {
			start := time.Now()
			url := fmt.Sprintf("%s/api/tags", host)
			resp, err := http.Get(url)
			if err == nil {
				defer resp.Body.Close()
				if resp.StatusCode == 200 {
					status = "online"
					latency = fmt.Sprintf("%dms", time.Since(start).Milliseconds())
				}
			}
		}

		// Display premium status
		fmt.Println()
		fmt.Println("◇ NUX Agent")
		fmt.Println()
		fmt.Printf("Provider:   %s\n", provider)
		fmt.Printf("Host:       %s\n", host)
		fmt.Printf("Model:      %s\n", model)
		fmt.Printf("Enabled:    %s\n", map[bool]string{true: "yes", false: "no"}[enabled])
		fmt.Printf("Latency:    %s\n", latency)
		fmt.Printf("Status:     %s\n", status)
		fmt.Println()

		if status == "online" {
			fmt.Println("✓ Ready for requests")
		} else {
			fmt.Println("✗ Agent unavailable")
		}
		fmt.Println("----------------------------")
	},
}

var agentQueryCmd = &cobra.Command{
	Use:   "query <question>",
	Short: "Ask a question to AI",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		question := strings.Join(args, " ")

		// Load vault for API keys
		v, _ := vault.Load()
		
		// Determine provider
		provider := "ollama"
		if v != nil {
			if p, ok := v.Config["agent_provider"].(string); ok {
				provider = p
			}
		}

		switch provider {
		case "ollama":
			agentQueryOllama(question, v, provider)
		case "openai":
			agentQueryOpenAI(question, v)
		default:
			output.NewError(fmt.Sprintf("unknown provider: %s", provider), "AGENT_INVALID_PROVIDER").Print()
		}
	},
}

func agentQueryOllama(question string, v *vault.Vault, provider string) {
	host := "http://localhost:11434"
	model := "qwen3-coder"

	if v != nil {
		if h, ok := v.Config["ollama_host"].(string); ok {
			host = h
		}
		if m, ok := v.Config["agent_model"].(string); ok {
			model = m
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
		output.NewError(fmt.Sprintf("failed to connect to Ollama: %s", err.Error()), "AGENT_OLLAMA_ERROR").Print()
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		output.NewError("failed to parse Ollama response", "AGENT_PARSE_ERROR").Print()
		return
	}

	if response, ok := result["response"].(string); ok {
		output.NewSuccess(map[string]interface{}{
			"provider": provider,
			"model":    model,
			"response": response,
		}).WithMessage("Ollama Response").Print()
	} else {
		output.NewError("no response from Ollama", "AGENT_NO_RESPONSE").Print()
	}
}

func agentQueryOpenAI(question string, v *vault.Vault) {
	// Similar to ask.go OpenAI implementation
	output.NewInfo("OpenAI provider not yet fully implemented").Print()
}

var agentConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure AI agent",
	Run: func(cmd *cobra.Command, args []string) {
		provider, _ := cmd.Flags().GetString("provider")
		apiKey, _ := cmd.Flags().GetString("api-key")
		host, _ := cmd.Flags().GetString("host")
		model, _ := cmd.Flags().GetString("model")

		v, err := vault.Load()
		if err != nil {
			v = vault.NewVault()
		}

		if provider != "" {
			if v.Config == nil {
				v.Config = make(map[string]interface{})
			}
			v.Config["agent_provider"] = provider
		}

		if apiKey != "" {
			v.SetAPIKey("openai", apiKey)
		}

		if host != "" {
			if v.Config == nil {
				v.Config = make(map[string]interface{})
			}
			v.Config["ollama_host"] = host
		}

		if model != "" {
			if v.Config == nil {
				v.Config = make(map[string]interface{})
			}
			v.Config["agent_model"] = model
		}

		if err := vault.Save(v); err != nil {
			output.NewError(fmt.Sprintf("failed to save config: %s", err.Error()), "AGENT_CONFIG_ERROR").Print()
			return
		}

		output.NewSuccess(map[string]interface{}{
			"provider": provider,
			"status":   "configured",
		}).Print()
	},
}

func init() {
	agentStatusCmd.Flags().Bool("json", false, "Output in JSON format")
	agentCmd.AddCommand(agentStatusCmd)
	
	agentQueryCmd.Flags().String("provider", "", "AI provider (ollama, openai, claude)")
	agentQueryCmd.Flags().String("model", "", "Model name")
	agentCmd.AddCommand(agentQueryCmd)
	
	agentConfigCmd.Flags().String("provider", "", "Provider to configure")
	agentConfigCmd.Flags().String("api-key", "", "API key for OpenAI")
	agentConfigCmd.Flags().String("host", "", "Host for Ollama")
	agentConfigCmd.Flags().String("model", "", "Model to use")
	agentCmd.AddCommand(agentConfigCmd)
	
	rootCmd.AddCommand(agentCmd)
}
