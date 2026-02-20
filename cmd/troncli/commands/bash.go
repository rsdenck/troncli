package commands

import (
	"fmt"
	"os"

	"github.com/mascli/troncli/internal/core/adapter"
	"github.com/mascli/troncli/internal/core/ports"
	"github.com/mascli/troncli/internal/core/services"
	"github.com/mascli/troncli/internal/modules/bash"
	"github.com/mascli/troncli/internal/ui/console"
	"github.com/spf13/cobra"
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

		table := console.NewBoxTable(os.Stdout)
		table.SetTitle(fmt.Sprintf("TRONCLI - BASH EXEC: %s", args[0]))
		table.SetHeaders([]string{"OUTPUT"})

		// Split output by lines to fit in table
		// This is a simple approach; might need better handling for very long outputs
		lines := splitLines(output)
		for _, line := range lines {
			if len(line) > 100 {
				line = line[:97] + "..."
			}
			table.AddRow([]string{line})
		}
		table.Render()
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

		table := console.NewBoxTable(os.Stdout)
		table.SetTitle(fmt.Sprintf("TRONCLI - BASH SCRIPT: %s", args[0]))
		table.SetHeaders([]string{"OUTPUT"})

		lines := splitLines(output)
		for _, line := range lines {
			if len(line) > 100 {
				line = line[:97] + "..."
			}
			table.AddRow([]string{line})
		}
		table.Render()
	},
}

func splitLines(s string) []string {
	var lines []string
	current := ""
	for _, c := range s {
		if c == '\n' {
			lines = append(lines, current)
			current = ""
		} else {
			current += string(c)
		}
	}
	if current != "" {
		lines = append(lines, current)
	}
	return lines
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
