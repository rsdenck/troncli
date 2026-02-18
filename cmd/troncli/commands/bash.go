package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/mascli/troncli/internal/core/adapter"
	"github.com/mascli/troncli/internal/core/ports"
	"github.com/mascli/troncli/internal/core/services"
	"github.com/mascli/troncli/internal/modules/bash"
)

var bashCmd = &cobra.Command{
	Use:   "bash",
	Short: "Executar comandos e scripts Bash",
	Long:  `Executa comandos Bash diretamente ou scripts de arquivos, gerenciando permissões e execução.`,
}

var bashRunCmd = &cobra.Command{
	Use:   "run [command]",
	Short: "Executar um comando Bash",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := getBashManager()
		if err != nil {
			fmt.Printf("Erro ao inicializar gerenciador Bash: %v\n", err)
			os.Exit(1)
		}

		output, err := manager.RunCommand(args[0])
		if err != nil {
			fmt.Printf("Erro ao executar comando: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(output)
	},
}

var bashScriptCmd = &cobra.Command{
	Use:   "script [file]",
	Short: "Executar um script Bash de arquivo",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := getBashManager()
		if err != nil {
			fmt.Printf("Erro ao inicializar gerenciador Bash: %v\n", err)
			os.Exit(1)
		}

		output, err := manager.RunScript(args[0])
		if err != nil {
			fmt.Printf("Erro ao executar script: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(output)
	},
}

func init() {
	rootCmd.AddCommand(bashCmd)
	bashCmd.AddCommand(bashRunCmd)
	bashCmd.AddCommand(bashScriptCmd)
}

func getBashManager() (ports.BashManager, error) {
	executor := adapter.NewExecutor()
	profileEngine := services.NewProfileEngine(executor)
	profile, err := profileEngine.DetectProfile()
	if err != nil {
		return nil, fmt.Errorf("failed to detect system profile: %w", err)
	}
	return bash.NewUniversalBashManager(executor, profile), nil
}
