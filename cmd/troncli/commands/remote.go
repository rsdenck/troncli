package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/mascli/troncli/internal/core/adapter"
	"github.com/mascli/troncli/internal/core/ports"
	"github.com/mascli/troncli/internal/core/services"
	"github.com/mascli/troncli/internal/modules/remote"
	"github.com/mascli/troncli/internal/ui/console"
	"github.com/spf13/cobra"
)

var remoteCmd = &cobra.Command{
	Use:   "remote",
	Short: "Gerenciar conexões remotas SSH",
	Long:  `Conectar, executar comandos e transferir arquivos via SSH.`,
}

var remoteConnectCmd = &cobra.Command{
	Use:   "connect [profile]",
	Short: "Conectar a um host remoto (interativo)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := getRemoteManager()
		if err != nil {
			fmt.Printf("Erro: %v\n", err)
			os.Exit(1)
		}
		if err := manager.Connect(args[0]); err != nil {
			fmt.Printf("Erro na conexão: %v\n", err)
			os.Exit(1)
		}
	},
}

var remoteExecCmd = &cobra.Command{
	Use:   "exec [profile] [command]",
	Short: "Executar comando em host remoto",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := getRemoteManager()
		if err != nil {
			fmt.Printf("Erro: %v\n", err)
			os.Exit(1)
		}
		output, err := manager.Execute(args[0], args[1])
		if err != nil {
			fmt.Printf("Erro na execução: %v\n", err)
			os.Exit(1)
		}

		table := console.NewBoxTable(os.Stdout)
		table.SetTitle(fmt.Sprintf("TRONCLI - REMOTE EXEC: %s (%s)", args[1], args[0]))
		table.SetHeaders([]string{"OUTPUT"})

		lines := strings.Split(output, "\n")
		for _, line := range lines {
			if strings.TrimSpace(line) != "" {
				if len(line) > 100 {
					line = line[:97] + "..."
				}
				table.AddRow([]string{line})
			}
		}
		table.Render()
	},
}

var remoteCopyCmd = &cobra.Command{
	Use:   "copy [src] [profile] [dest]",
	Short: "Copiar arquivo para host remoto (SCP)",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := getRemoteManager()
		if err != nil {
			fmt.Printf("Erro: %v\n", err)
			os.Exit(1)
		}
		if err := manager.CopyFile(args[1], args[0], args[2]); err != nil {
			fmt.Printf("Erro na cópia: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Arquivo copiado com sucesso.")
	},
}

var remoteListCmd = &cobra.Command{
	Use:   "list",
	Short: "Listar perfis SSH configurados",
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := getRemoteManager()
		if err != nil {
			fmt.Printf("Erro: %v\n", err)
			os.Exit(1)
		}
		profiles, err := manager.ListProfiles()
		if err != nil {
			fmt.Printf("Erro ao listar perfis: %v\n", err)
			os.Exit(1)
		}

		table := console.NewBoxTable(os.Stdout)
		table.SetTitle("TRONCLI - PERFIS SSH")
		table.SetHeaders([]string{"PROFILE"})

		for _, p := range profiles {
			table.AddRow([]string{p})
		}
		table.SetFooter(fmt.Sprintf("Total profiles: %d", len(profiles)))
		table.Render()
	},
}

func init() {
	rootCmd.AddCommand(remoteCmd)
	remoteCmd.AddCommand(remoteConnectCmd)
	remoteCmd.AddCommand(remoteExecCmd)
	remoteCmd.AddCommand(remoteCopyCmd)
	remoteCmd.AddCommand(remoteListCmd)
}

func getRemoteManager() (ports.SSHClient, error) {
	executor := adapter.NewExecutor()
	profileEngine := services.NewProfileEngine(executor)
	profile, err := profileEngine.DetectProfile()
	if err != nil {
		return nil, fmt.Errorf("failed to detect system profile: %w", err)
	}
	return remote.NewUniversalRemoteManager(executor, profile), nil
}
