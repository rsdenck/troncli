package commands

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/mascli/troncli/internal/agent"
	"github.com/mascli/troncli/internal/console"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type AgentConfig struct {
	Provider string `yaml:"provider"`
	Model    string `yaml:"model"`
	APIKey   string `yaml:"api_key"`
}

var agentConfigPath string

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		slog.Warn("Could not get user home directory", "error", err)
	} else {
		configDir := filepath.Join(home, ".troncli")
		if err := os.MkdirAll(configDir, 0755); err != nil {
			slog.Warn("Could not create config directory", "error", err)
		}
		agentConfigPath = filepath.Join(configDir, "agent_config.yaml")
	}

	agentCmd.AddCommand(agentEnableCmd)
	agentCmd.AddCommand(agentSetModelCmd)
	agentCmd.AddCommand(agentAskCmd)
	agentCmd.AddCommand(agentStatusCmd)
	rootCmd.AddCommand(agentCmd)
}

var agentStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Exibir o status e configuração do agente",
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := loadAgentConfig()
		if err != nil {
			return fmt.Errorf("erro ao carregar configuração: %w", err)
		}

		table := console.NewBoxTable(os.Stdout)
		table.SetTitle("TRONCLI: AGENT STATUS")

		table.AddRow([]string{"Provider", config.Provider})
		table.AddRow([]string{"Model", config.Model})

		apiKeyStatus := "Not Set"
		if len(config.APIKey) > 4 {
			apiKeyStatus = "********" + config.APIKey[len(config.APIKey)-4:]
		} else if config.APIKey != "" {
			apiKeyStatus = "********"
		}
		table.AddRow([]string{"API Key", apiKeyStatus})

		table.RenderKeyValue()
		return nil
	},
}

func loadAgentConfig() (*AgentConfig, error) {
	data, err := os.ReadFile(agentConfigPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &AgentConfig{
				Provider: "local",
				Model:    "default",
			}, nil
		}
		return nil, err
	}
	var config AgentConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func saveAgentConfig(config *AgentConfig) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	return os.WriteFile(agentConfigPath, data, 0644)
}

func getAgentAdapter(config *AgentConfig) (agent.AgentAdapter, error) {
	home, _ := os.UserHomeDir()
	capabilitiesPath := filepath.Join(home, ".troncli", "capabilities.yaml")

	// Create default capabilities file if not exists
	if _, err := os.Stat(capabilitiesPath); os.IsNotExist(err) {
		// In a real app we might copy a default file or just create an empty one
	}

	switch config.Provider {
	case "ollama":
		return agent.NewOllamaAdapter(config.Model, capabilitiesPath)
	case "llamacpp":
		// Default paths for llama.cpp integration
		modelPath := filepath.Join(home, ".troncli", "models", "qwen2.5-coder-7b-instruct-q4_0.gguf")
		llamaPath := filepath.Join(home, ".troncli", "bin", "llama-cli")
		return agent.NewLlamaCppAdapter(modelPath, llamaPath, capabilitiesPath)
	case "claude":
		return agent.NewClaudeAdapter(config.APIKey, capabilitiesPath)
	case "openai":
		return agent.NewOpenAIAdapter(config.APIKey, config.Model, capabilitiesPath)
	case "local":
		return agent.NewLocalAgent(capabilitiesPath)
	default:
		return nil, fmt.Errorf("provedor desconhecido: %s", config.Provider)
	}
}

var agentCmd = &cobra.Command{
	Use:   "agent [command] or [intent]",
	Short: "Gerenciar e interagir com agentes de IA",
	Long: `Comandos para configurar e utilizar agentes de IA (Ollama, Claude, OpenAI, Local, LlamaCpp).
Você pode passar uma intenção diretamente: troncli agent "instalar nginx"

TRON ROOT AGENT (llama.cpp):
  O agente root é um agente autônomo hardcore que usa llama.cpp diretamente.
  Ele analisa riscos, gera comandos troncli e executa com confirmação.
  
  Exemplo: troncli agent root "verificar saúde do sistema"`,
	Args: cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			cmd.Help()
			return nil
		}

		// Check if first arg is "root" for root agent
		if args[0] == "root" {
			if len(args) < 2 {
				return fmt.Errorf("root agent requires an intent. Example: troncli agent root \"install nginx\"")
			}
			intent := strings.Join(args[1:], " ")
			return executeRootAgent(intent)
		}

		// Regular agent execution
		intent := strings.Join(args, " ")
		config, err := loadAgentConfig()
		if err != nil {
			return fmt.Errorf("erro ao carregar configuração: %w", err)
		}

		adapter, err := getAgentAdapter(config)
		if err != nil {
			return fmt.Errorf("erro ao inicializar agente: %w", err)
		}

		slog.Info("Agente executando intenção", "provider", adapter.Name(), "intent", intent)
		result, err := adapter.ExecuteIntent(context.Background(), intent)
		if err != nil {
			return fmt.Errorf("erro ao executar intenção: %w", err)
		}
		fmt.Println(result)
		return nil
	},
}

// executeRootAgent executes the TRON ROOT AGENT
func executeRootAgent(intent string) error {
	home, _ := os.UserHomeDir()
	modelPath := filepath.Join(home, ".troncli", "models", "qwen2.5-coder-7b-instruct-q4_0.gguf")
	llamaPath := filepath.Join(home, ".troncli", "bin", "llama-cli")

	// Check for system-wide llama-cli
	if _, err := os.Stat(llamaPath); os.IsNotExist(err) {
		// Try system paths
		systemPaths := []string{
			"/usr/local/bin/llama-cli",
			"/usr/bin/llama-cli",
			"/opt/llama.cpp/llama-cli",
		}
		for _, path := range systemPaths {
			if _, err := os.Stat(path); err == nil {
				llamaPath = path
				break
			}
		}
	}

	rootAgent := agent.NewRootAgent(modelPath, llamaPath)

	// Use streaming for better UX
	streaming, _ := os.LookupEnv("TRONCLI_AGENT_STREAMING")
	if streaming == "true" {
		return rootAgent.StreamingExecute(context.Background(), intent)
	}

	return rootAgent.Execute(context.Background(), intent)
}

var agentEnableCmd = &cobra.Command{
	Use:   "enable [provider]",
	Short: "Habilitar um provedor de agente (ollama, llamacpp, claude, openai, local)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		provider := args[0]
		validProviders := map[string]bool{
			"ollama": true, "llamacpp": true, "claude": true, "openai": true, "local": true,
		}
		if !validProviders[provider] {
			return fmt.Errorf("provedor inválido '%s'. Opções: ollama, llamacpp, claude, openai, local", provider)
		}

		config, err := loadAgentConfig()
		if err != nil {
			return fmt.Errorf("erro ao carregar configuração: %w", err)
		}

		config.Provider = provider
		if err := saveAgentConfig(config); err != nil {
			return fmt.Errorf("erro ao salvar configuração: %w", err)
		}

		slog.Info("Agente habilitado com sucesso", "provider", provider)
		return nil
	},
}

var agentSetModelCmd = &cobra.Command{
	Use:   "set-model [model]",
	Short: "Definir o modelo para o agente atual",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		model := args[0]
		config, err := loadAgentConfig()
		if err != nil {
			return fmt.Errorf("erro ao carregar configuração: %w", err)
		}

		config.Model = model
		if err := saveAgentConfig(config); err != nil {
			return fmt.Errorf("erro ao salvar configuração: %w", err)
		}

		slog.Info("Modelo definido", "model", model)
		return nil
	},
}

var agentAskCmd = &cobra.Command{
	Use:   "ask [prompt]",
	Short: "Enviar um prompt para o agente",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		prompt := args[0]
		config, err := loadAgentConfig()
		if err != nil {
			return fmt.Errorf("erro ao carregar configuração: %w", err)
		}

		adapter, err := getAgentAdapter(config)
		if err != nil {
			return fmt.Errorf("erro ao inicializar agente: %w", err)
		}

		slog.Info("Enviando prompt para agente", "provider", adapter.Name())
		response, err := adapter.SendPrompt(context.Background(), prompt)
		if err != nil {
			return fmt.Errorf("erro ao comunicar com o agente: %w", err)
		}
		fmt.Println(response)
		return nil
	},
}
