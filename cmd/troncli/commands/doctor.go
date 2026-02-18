package commands

import (
	"fmt"
	"os"

	"github.com/mascli/troncli/internal/core/adapter"
	"github.com/mascli/troncli/internal/core/ports"
	"github.com/mascli/troncli/internal/core/services"
	"github.com/mascli/troncli/internal/modules/doctor"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Saúde do Sistema",
	Long:  `Executa verificações de saúde do sistema (Load, Swap, Disco, TCP, etc).`,
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := getDoctorManager()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("Running system health checks...")
		checks, err := manager.RunChecks()
		if err != nil {
			fmt.Printf("Error running checks: %v\n", err)
			os.Exit(1)
		}

		for _, check := range checks {
			statusIcon := "✅"
			if check.Status == ports.StatusWarning {
				statusIcon = "⚠️"
			} else if check.Status == ports.StatusCritical {
				statusIcon = "❌"
			}

			fmt.Printf("%s %-20s [%s] %s\n", statusIcon, check.Name, check.Value, check.Message)
		}
	},
}

func getDoctorManager() (*doctor.UniversalDoctorManager, error) {
	executor := adapter.NewExecutor()
	profileEngine := services.NewProfileEngine(executor)
	profile, err := profileEngine.DetectProfile()
	if err != nil {
		return nil, fmt.Errorf("failed to detect system profile: %w", err)
	}

	return doctor.NewUniversalDoctorManager(executor, profile), nil
}

func init() {
	rootCmd.AddCommand(doctorCmd)
	doctorCmd.Flags().Bool("full", false, "Executa todas as verificações")
	doctorCmd.Flags().Bool("security", false, "Executa verificações de segurança")
	doctorCmd.Flags().Bool("network", false, "Executa verificações de rede")
	doctorCmd.Flags().Bool("disk", false, "Executa verificações de disco")
}
