package commands

import (
	"bufio"
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
	Use:   "agent [intent]",
	Short: "TRON ROOT AGENT - Agente autônomo com IA",
	Long: `TRON ROOT AGENT - Agente autônomo que usa llama.cpp + Qwen2.5-Coder-7B.

Modo Interativo (RECOMENDADO):
  troncli agent
  
  Abre um prompt interativo onde você pode digitar comandos em linguagem natural.
  O agente analisa, gera comandos troncli e executa com confirmação.

Modo Direto:
  troncli agent "verificar saúde do sistema"
  troncli agent "instalar nginx"
  troncli agent "listar serviços ativos"

Primeira Execução:
  Na primeira vez, o agente baixa automaticamente:
  - llama.cpp (~50MB)
  - Modelo Qwen2.5-Coder-7B (~4GB)
  
  Isso leva ~10 minutos dependendo da conexão.`,
	Args: cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Auto-setup on first run
		if err := autoSetupIfNeeded(); err != nil {
			return fmt.Errorf("falha no setup automático: %w", err)
		}

		// Interactive mode if no args
		if len(args) == 0 {
			return runInteractiveAgent()
		}

		// Direct mode with intent
		intent := strings.Join(args, " ")
		return executeRootAgent(intent)
	},
}

// autoSetupIfNeeded checks and installs llama.cpp + model automatically
func autoSetupIfNeeded() error {
	home, _ := os.UserHomeDir()
	llamaPath := filepath.Join(home, ".troncli", "bin", "llama-cli")
	modelPath := filepath.Join(home, ".troncli", "models", "qwen2.5-coder-7b-instruct-q4_0.gguf")

	// Check if already setup
	llamaExists := false
	modelExists := false

	if _, err := os.Stat(llamaPath); err == nil {
		llamaExists = true
	}
	if _, err := os.Stat(modelPath); err == nil {
		modelExists = true
	}

	// If both exist, nothing to do
	if llamaExists && modelExists {
		return nil
	}

	// First time setup
	fmt.Printf("\n%s╔════════════════════════════════════════════════════════════╗%s\n", console.ColorCyan, console.ColorReset)
	fmt.Printf("%s║%s                                                            %s║%s\n", console.ColorCyan, console.ColorReset, console.ColorCyan, console.ColorReset)
	fmt.Printf("%s║%s  %s🚀 PRIMEIRA EXECUÇÃO - SETUP AUTOMÁTICO%s                  %s║%s\n", 
		console.ColorCyan, console.ColorReset, console.ColorBold, console.ColorReset, console.ColorCyan, console.ColorReset)
	fmt.Printf("%s║%s                                                            %s║%s\n", console.ColorCyan, console.ColorReset, console.ColorCyan, console.ColorReset)
	fmt.Printf("%s║%s  Instalando llama.cpp + Modelo Qwen2.5-Coder-7B           %s║%s\n", console.ColorCyan, console.ColorReset, console.ColorCyan, console.ColorReset)
	fmt.Printf("%s║%s  Tempo estimado: ~10 minutos                               %s║%s\n", console.ColorCyan, console.ColorReset, console.ColorCyan, console.ColorReset)
	fmt.Printf("%s║%s                                                            %s║%s\n", console.ColorCyan, console.ColorReset, console.ColorCyan, console.ColorReset)
	fmt.Printf("%s╚════════════════════════════════════════════════════════════╝%s\n\n", console.ColorCyan, console.ColorReset)

	// Run setup
	return setupRootAgent()
}

// runInteractiveAgent starts interactive mode
func runInteractiveAgent() error {
	fmt.Printf("\n%s╔════════════════════════════════════════════════════════════╗%s\n", console.ColorCyan, console.ColorReset)
	fmt.Printf("%s║%s                                                            %s║%s\n", console.ColorCyan, console.ColorReset, console.ColorCyan, console.ColorReset)
	fmt.Printf("%s║%s  %sTRON ROOT AGENT › MODO INTERATIVO%s                       %s║%s\n", 
		console.ColorCyan, console.ColorReset, console.ColorBold+console.ColorWhite, console.ColorReset, console.ColorCyan, console.ColorReset)
	fmt.Printf("%s║%s                                                            %s║%s\n", console.ColorCyan, console.ColorReset, console.ColorCyan, console.ColorReset)
	fmt.Printf("%s║%s  Digite comandos em linguagem natural                      %s║%s\n", console.ColorCyan, console.ColorReset, console.ColorCyan, console.ColorReset)
	fmt.Printf("%s║%s  Digite 'sair' ou 'exit' para encerrar                     %s║%s\n", console.ColorCyan, console.ColorReset, console.ColorCyan, console.ColorReset)
	fmt.Printf("%s║%s                                                            %s║%s\n", console.ColorCyan, console.ColorReset, console.ColorCyan, console.ColorReset)
	fmt.Printf("%s╚════════════════════════════════════════════════════════════╝%s\n\n", console.ColorCyan, console.ColorReset)

	reader := bufio.NewReader(os.Stdin)

	for {
		// Prompt
		fmt.Printf("%s❯%s ", console.ColorCyan+console.ColorBold, console.ColorReset)

		// Read input
		input, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("erro ao ler input: %w", err)
		}

		input = strings.TrimSpace(input)

		// Check exit
		if input == "sair" || input == "exit" || input == "quit" {
			fmt.Printf("\n%s👋 Até logo!%s\n\n", console.ColorCyan, console.ColorReset)
			return nil
		}

		// Skip empty input
		if input == "" {
			continue
		}

		// Execute intent
		fmt.Println()
		if err := executeRootAgent(input); err != nil {
			fmt.Printf("\n%s❌ Erro: %v%s\n\n", console.ColorRed, err, console.ColorReset)
		}
		fmt.Println()
	}
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

	// Always use non-streaming for cleaner output
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
