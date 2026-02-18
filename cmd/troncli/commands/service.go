package commands

import (
	"fmt"
	"os"
	"strconv"

	"github.com/mascli/troncli/internal/core/adapter"
	"github.com/mascli/troncli/internal/core/ports"
	"github.com/mascli/troncli/internal/core/services"
	"github.com/mascli/troncli/internal/modules/service"
	"github.com/spf13/cobra"
)

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Gerenciar serviços do sistema",
	Long:  `Controlar serviços (systemd, sysvinit, openrc, runit) de forma unificada.`,
}

var serviceListCmd = &cobra.Command{
	Use:   "list",
	Short: "Listar serviços",
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := getServiceManager()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		units, err := manager.ListServices()
		if err != nil {
			fmt.Printf("Error listing services: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("%-30s %-10s %-10s %s\n", "UNIT", "LOAD", "ACTIVE", "DESCRIPTION")
		for _, u := range units {
			// Filter out too many entries if needed, but for now list all
			fmt.Printf("%-30s %-10s %-10s %s\n", u.Name, u.LoadState, u.ActiveState, u.Description)
		}
	},
}

var serviceStartCmd = &cobra.Command{
	Use:   "start [service]",
	Short: "Iniciar um serviço",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := getServiceManager()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		if err := manager.StartService(args[0]); err != nil {
			fmt.Printf("Error starting service: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Service started.")
	},
}

var serviceStopCmd = &cobra.Command{
	Use:   "stop [service]",
	Short: "Parar um serviço",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := getServiceManager()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		if err := manager.StopService(args[0]); err != nil {
			fmt.Printf("Error stopping service: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Service stopped.")
	},
}

var serviceRestartCmd = &cobra.Command{
	Use:   "restart [service]",
	Short: "Reiniciar um serviço",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := getServiceManager()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		if err := manager.RestartService(args[0]); err != nil {
			fmt.Printf("Error restarting service: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Service restarted.")
	},
}

var serviceStatusCmd = &cobra.Command{
	Use:   "status [service]",
	Short: "Ver status detalhado de um serviço",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := getServiceManager()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		status, err := manager.GetServiceStatus(args[0])
		// Status often returns error code for inactive services, but we might have output
		if status != "" {
			fmt.Println(status)
		} else if err != nil {
			fmt.Printf("Error getting status: %v\n", err)
		}
	},
}

var serviceLogsCmd = &cobra.Command{
	Use:   "logs [service]",
	Short: "Ver logs do serviço",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := getServiceManager()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		linesStr, _ := cmd.Flags().GetString("lines")
		lines := 20
		if l, err := strconv.Atoi(linesStr); err == nil {
			lines = l
		}

		logs, err := manager.GetServiceLogs(args[0], lines)
		if err != nil {
			fmt.Printf("Error getting logs: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(logs)
	},
}

var serviceEnableCmd = &cobra.Command{
	Use:   "enable [service]",
	Short: "Habilitar serviço na inicialização",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := getServiceManager()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		if err := manager.EnableService(args[0]); err != nil {
			fmt.Printf("Error enabling service: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Service enabled.")
	},
}

var serviceDisableCmd = &cobra.Command{
	Use:   "disable [service]",
	Short: "Desabilitar serviço na inicialização",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := getServiceManager()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		if err := manager.DisableService(args[0]); err != nil {
			fmt.Printf("Error disabling service: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Service disabled.")
	},
}

func init() {
	rootCmd.AddCommand(serviceCmd)
	serviceCmd.AddCommand(serviceListCmd)
	serviceCmd.AddCommand(serviceStartCmd)
	serviceCmd.AddCommand(serviceStopCmd)
	serviceCmd.AddCommand(serviceRestartCmd)
	serviceCmd.AddCommand(serviceStatusCmd)
	serviceCmd.AddCommand(serviceLogsCmd)
	serviceCmd.AddCommand(serviceEnableCmd)
	serviceCmd.AddCommand(serviceDisableCmd)

	serviceLogsCmd.Flags().StringP("lines", "n", "20", "Número de linhas")
}

func getServiceManager() (ports.ServiceManager, error) {
	executor := adapter.NewExecutor()
	profileEngine := services.NewProfileEngine(executor)
	profile, err := profileEngine.DetectProfile()
	if err != nil {
		return nil, fmt.Errorf("failed to detect system profile: %w", err)
	}
	return service.NewUniversalServiceManager(executor, profile), nil
}
