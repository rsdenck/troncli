package commands

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/mascli/troncli/internal/core/adapter"
	"github.com/mascli/troncli/internal/core/ports"
	"github.com/mascli/troncli/internal/core/services"
	"github.com/mascli/troncli/internal/modules/plugin"
	"github.com/mascli/troncli/internal/ui/console"
	"github.com/spf13/cobra"
)

var pluginCmd = &cobra.Command{
	Use:   "plugin",
	Short: "Gerenciar plugins do TRONCLI",
	Long:  `Instalar, listar e remover plugins (scripts ou binários) do TRONCLI.`,
}

var pluginListCmd = &cobra.Command{
	Use:   "list",
	Short: "Listar plugins instalados",
	RunE: func(cmd *cobra.Command, args []string) error {
		manager, err := getPluginManager()
		if err != nil {
			return fmt.Errorf("erro: %w", err)
		}
		plugins, err := manager.ListPlugins()
		if err != nil {
			return fmt.Errorf("erro ao listar: %w", err)
		}

		table := console.NewBoxTable(os.Stdout)
		table.SetTitle("TRONCLI - PLUGINS INSTALADOS")
		table.SetHeaders([]string{"NAME", "VERSION", "DESCRIPTION"})

		for _, p := range plugins {
			desc := p.Description
			if len(desc) > 50 {
				desc = desc[:47] + "..."
			}
			table.AddRow([]string{p.Name, p.Version, desc})
		}
		table.SetFooter(fmt.Sprintf("Total plugins: %d", len(plugins)))
		table.Render()

		return nil
	},
}

var pluginInstallCmd = &cobra.Command{
	Use:   "install [url|path|name]",
	Short: "Instalar um plugin (URL, caminho local ou nome registrado)",
	Long: `Instale plugins fornecendo uma URL direta, um caminho de arquivo local ou um nome de plugin registrado (ex: arch, docker, k8s).
	
Exemplos:
  troncli plugin install arch
  troncli plugin install https://example.com/my-plugin.sh
  troncli plugin install ./local-script.py`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		manager, err := getPluginManager()
		if err != nil {
			return fmt.Errorf("erro: %w", err)
		}
		slog.Info("Instalando plugin", "plugin", args[0])
		if err := manager.InstallPlugin(args[0]); err != nil {
			return fmt.Errorf("erro ao instalar: %w", err)
		}
		slog.Info("Plugin instalado com sucesso", "plugin", args[0])
		return nil
	},
}

var pluginRemoveCmd = &cobra.Command{
	Use:   "remove [name]",
	Short: "Remover um plugin",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		manager, err := getPluginManager()
		if err != nil {
			return fmt.Errorf("erro: %w", err)
		}
		slog.Info("Removendo plugin", "plugin", args[0])
		if err := manager.RemovePlugin(args[0]); err != nil {
			return fmt.Errorf("erro ao remover: %w", err)
		}
		slog.Info("Plugin removido com sucesso", "plugin", args[0])
		return nil
	},
}

var pluginExecCmd = &cobra.Command{
	Use:   "exec [name] [args...]",
	Short: "Executar um plugin instalado",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		manager, err := getPluginManager()
		if err != nil {
			return fmt.Errorf("erro: %w", err)
		}
		name := args[0]
		pluginArgs := args[1:]

		slog.Info("Executando plugin", "plugin", name, "args", pluginArgs)
		if err := manager.ExecutePlugin(context.Background(), name, pluginArgs...); err != nil {
			return fmt.Errorf("erro ao executar plugin '%s': %w", name, err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(pluginCmd)
	pluginCmd.AddCommand(pluginListCmd)
	pluginCmd.AddCommand(pluginInstallCmd)
	pluginCmd.AddCommand(pluginRemoveCmd)
	pluginCmd.AddCommand(pluginExecCmd)
}

func getPluginManager() (ports.PluginManager, error) {
	executor := adapter.NewExecutor()
	profileEngine := services.NewProfileEngine(executor)
	profile, err := profileEngine.DetectProfile()
	if err != nil {
		return nil, fmt.Errorf("failed to detect system profile: %w", err)
	}
	return plugin.NewUniversalPluginManager(executor, profile)
}
