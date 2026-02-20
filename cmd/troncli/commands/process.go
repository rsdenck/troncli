package commands

import (
	"fmt"
	"os"
	"strconv"

	"github.com/mascli/troncli/internal/core/adapter"
	"github.com/mascli/troncli/internal/core/services"
	"github.com/mascli/troncli/internal/modules/process"
	"github.com/mascli/troncli/internal/ui/console"
	"github.com/spf13/cobra"
)

var processCmd = &cobra.Command{
	Use:   "process",
	Short: "Gerenciamento de Processos do Sistema",
	Long:  `Visualiza, finaliza e gerencia prioridade de processos em execução.`,
}

func getProcessManager() (*process.UniversalProcessManager, error) {
	executor := adapter.NewExecutor()
	profileEngine := services.NewProfileEngine(executor)
	profile, err := profileEngine.DetectProfile()
	if err != nil {
		return nil, fmt.Errorf("failed to detect system profile: %w", err)
	}

	return process.NewUniversalProcessManager(executor, profile), nil
}

var treeCmd = &cobra.Command{
	Use:   "tree",
	Short: "Exibe a árvore de processos",
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := getProcessManager()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		nodes, err := manager.GetProcessTree()
		if err != nil {
			fmt.Printf("Error getting process tree: %v\n", err)
			os.Exit(1)
		}

		table := console.NewBoxTable(os.Stdout)
		table.SetTitle("TRONCLI - ÁRVORE DE PROCESSOS")
		table.SetHeaders([]string{"PID", "PPID", "USER", "STATE", "COMMAND"})

		for _, node := range nodes {
			// Truncate name if too long
			name := node.Name
			if len(name) > 50 {
				name = name[:47] + "..."
			}
			table.AddRow([]string{
				fmt.Sprintf("%d", node.PID),
				fmt.Sprintf("%d", node.PPID),
				node.User,
				node.State,
				name,
			})
		}
		table.SetFooter(fmt.Sprintf("Total processes: %d", len(nodes)))
		table.Render()
	},
}

var openFilesCmd = &cobra.Command{
	Use:   "open-files [pid]",
	Short: "Lista arquivos abertos por um processo",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pid, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Invalid PID")
			os.Exit(1)
		}

		manager, err := getProcessManager()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		files, err := manager.GetOpenFiles(pid)
		if err != nil {
			fmt.Printf("Error getting open files: %v\n", err)
			os.Exit(1)
		}

		table := console.NewBoxTable(os.Stdout)
		table.SetTitle(fmt.Sprintf("TRONCLI - ARQUIVOS ABERTOS (PID %d)", pid))
		table.SetHeaders([]string{"FILE"})

		for _, f := range files {
			if len(f) > 80 {
				f = f[:77] + "..."
			}
			table.AddRow([]string{f})
		}
		table.SetFooter(fmt.Sprintf("Total files: %d", len(files)))
		table.Render()
	},
}

var processPortsCmd = &cobra.Command{
	Use:   "ports [pid]",
	Short: "Lista portas ouvidas por um processo",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pid, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Invalid PID")
			os.Exit(1)
		}

		manager, err := getProcessManager()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		ports, err := manager.GetProcessPorts(pid)
		if err != nil {
			fmt.Printf("Error getting process ports: %v\n", err)
			os.Exit(1)
		}

		table := console.NewBoxTable(os.Stdout)
		table.SetTitle(fmt.Sprintf("TRONCLI - PORTAS DO PROCESSO (PID %d)", pid))
		table.SetHeaders([]string{"PORT/PROTOCOL"})

		for _, p := range ports {
			table.AddRow([]string{p})
		}
		table.SetFooter(fmt.Sprintf("Total ports: %d", len(ports)))
		table.Render()
	},
}

var listeningCmd = &cobra.Command{
	Use:   "listening",
	Short: "Lista todas as portas em escuta no sistema",
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := getProcessManager()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		ports, err := manager.GetAllListeningPorts()
		if err != nil {
			fmt.Printf("Error getting listening ports: %v\n", err)
			os.Exit(1)
		}

		table := console.NewBoxTable(os.Stdout)
		table.SetTitle("TRONCLI - PORTAS EM ESCUTA (LISTENING)")
		table.SetHeaders([]string{"PORT/PROTOCOL"})

		for _, p := range ports {
			table.AddRow([]string{p})
		}
		table.SetFooter(fmt.Sprintf("Total listening ports: %d", len(ports)))
		table.Render()
	},
}

var killCmd = &cobra.Command{
	Use:   "kill [pid] [signal]",
	Short: "Envia sinal para um processo (default SIGTERM)",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pid, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Invalid PID")
			os.Exit(1)
		}

		signal := "SIGTERM"
		if len(args) > 1 {
			signal = args[1]
		}

		manager, err := getProcessManager()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if err := manager.KillProcess(pid, signal); err != nil {
			fmt.Printf("Error killing process: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Signal %s sent to PID %d\n", signal, pid)
	},
}

var reniceCmd = &cobra.Command{
	Use:   "renice [pid] [priority]",
	Short: "Altera a prioridade de um processo",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		pid, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Invalid PID")
			os.Exit(1)
		}

		prio, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Invalid priority")
			os.Exit(1)
		}

		manager, err := getProcessManager()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if err := manager.ReniceProcess(pid, prio); err != nil {
			fmt.Printf("Error renicing process: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Process %d reniced to %d\n", pid, prio)
	},
}

var zombiesCmd = &cobra.Command{
	Use:   "zombies",
	Short: "Elimina processos zumbis",
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := getProcessManager()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		count, err := manager.KillZombies()
		if err != nil {
			fmt.Printf("Error killing zombies: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Killed %d zombie processes\n", count)
	},
}

func init() {
	rootCmd.AddCommand(processCmd)
	processCmd.AddCommand(treeCmd)
	processCmd.AddCommand(openFilesCmd)
	processCmd.AddCommand(processPortsCmd)
	processCmd.AddCommand(listeningCmd)
	processCmd.AddCommand(killCmd)
	processCmd.AddCommand(reniceCmd)
	processCmd.AddCommand(zombiesCmd)
}
