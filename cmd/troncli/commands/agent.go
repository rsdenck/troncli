package commands

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/mascli/troncli/internal/agent"
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
	rootCmd.AddCommand(agentCmd)
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
	Long: `Comandos para configurar e utilizar agentes de IA (Ollama, Claude, OpenAI, Local).
Você pode passar uma intenção diretamente: troncli agent "instalar nginx"`,
	Args: cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			cmd.Help()
			return nil
		}

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

var agentEnableCmd = &cobra.Command{
	Use:   "enable [provider]",
	Short: "Habilitar um provedor de agente (ollama, claude, openai, local)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		provider := args[0]
		validProviders := map[string]bool{
			"ollama": true, "claude": true, "openai": true, "local": true,
		}
		if !validProviders[provider] {
			return fmt.Errorf("provedor inválido '%s'. Opções: ollama, claude, openai, local", provider)
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
