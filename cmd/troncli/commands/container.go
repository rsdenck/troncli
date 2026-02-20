package commands

import (
	"fmt"
	"os"
	"strconv"

	"github.com/mascli/troncli/internal/core/adapter"
	"github.com/mascli/troncli/internal/core/ports"
	"github.com/mascli/troncli/internal/core/services"
	"github.com/mascli/troncli/internal/modules/container"
	"github.com/mascli/troncli/internal/ui/console"
	"github.com/spf13/cobra"
)

var containerCmd = &cobra.Command{
	Use:   "container",
	Short: "Gerenciar containers (Docker/Podman)",
	Long:  `Gerenciar ciclo de vida de containers, suportando Docker e Podman automaticamente.`,
}

var containerListCmd = &cobra.Command{
	Use:   "list",
	Short: "Listar containers",
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := getContainerManager()
		if err != nil {
			fmt.Printf("Erro: %v\n", err)
			os.Exit(1)
		}
		
		all, _ := cmd.Flags().GetBool("all")
		containers, err := manager.ListContainers(all)
		if err != nil {
			fmt.Printf("Erro ao listar containers: %v\n", err)
			os.Exit(1)
		}
		
		table := console.NewBoxTable(os.Stdout)
		table.SetTitle("TRONCLI - LISTAGEM DE CONTAINERS")
		table.SetHeaders([]string{"ID", "IMAGE", "STATE", "STATUS", "RUNTIME"})
		
		for _, c := range containers {
			name := ""
			if len(c.Names) > 0 {
				name = c.Names[0]
			}
			// Truncate ID
			id := c.ID
			if len(id) > 12 {
				id = id[:12]
			}
			table.AddRow([]string{id, name, c.State, c.Status, c.Runtime})
		}
		
		table.SetFooter(fmt.Sprintf("Total containers: %d", len(containers)))
		table.Render()
	},
}

var containerStartCmd = &cobra.Command{
	Use:   "start [id]",
	Short: "Iniciar um container",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := getContainerManager()
		if err != nil {
			fmt.Printf("Erro: %v\n", err)
			os.Exit(1)
		}
		if err := manager.StartContainer(args[0]); err != nil {
			fmt.Printf("Erro ao iniciar: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Container iniciado.")
	},
}

var containerStopCmd = &cobra.Command{
	Use:   "stop [id]",
	Short: "Parar um container",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := getContainerManager()
		if err != nil {
			fmt.Printf("Erro: %v\n", err)
			os.Exit(1)
		}
		if err := manager.StopContainer(args[0]); err != nil {
			fmt.Printf("Erro ao parar: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Container parado.")
	},
}

var containerLogsCmd = &cobra.Command{
	Use:   "logs [id]",
	Short: "Ver logs de um container",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := getContainerManager()
		if err != nil {
			fmt.Printf("Erro: %v\n", err)
			os.Exit(1)
		}

		tailStr, _ := cmd.Flags().GetString("tail")
		tail := 20
		if tailStr != "" {
			t, err := strconv.Atoi(tailStr)
			if err == nil {
				tail = t
			}
		}

		logs, err := manager.GetContainerLogs(args[0], tail)
		if err != nil {
			fmt.Printf("Erro ao obter logs: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(logs)
	},
}

func init() {
	rootCmd.AddCommand(containerCmd)
	containerCmd.AddCommand(containerListCmd)
	containerCmd.AddCommand(containerStartCmd)
	containerCmd.AddCommand(containerStopCmd)
	containerCmd.AddCommand(containerLogsCmd)

	containerListCmd.Flags().BoolP("all", "a", false, "Mostrar todos os containers (padrão mostra apenas em execução)")
	containerLogsCmd.Flags().String("tail", "20", "Número de linhas para mostrar")
}

func getContainerManager() (ports.ContainerManager, error) {
	executor := adapter.NewExecutor()
	profileEngine := services.NewProfileEngine(executor)
	profile, err := profileEngine.DetectProfile()
	if err != nil {
		return nil, fmt.Errorf("failed to detect system profile: %w", err)
	}
	return container.NewUniversalContainerManager(executor, profile), nil
}
