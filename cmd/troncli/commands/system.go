package commands

import (
	"fmt"
	"os"

	"github.com/mascli/troncli/internal/core/adapter"
	"github.com/mascli/troncli/internal/core/services"
	"github.com/mascli/troncli/internal/ui/console"
	"github.com/spf13/cobra"
)

var systemCmd = &cobra.Command{
	Use:   "system",
	Short: "Informações e Perfil do Sistema",
	Long:  `Exibe informações detalhadas sobre o sistema, kernel, uptime e ambiente.`,
}

var systemInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Exibe informações gerais do sistema",
	Run: func(cmd *cobra.Command, args []string) {
		executor := adapter.NewExecutor()
		profile, err := services.NewProfileEngine(executor).DetectProfile()
		if err != nil {
			fmt.Printf("Error detecting system profile: %v\n", err)
			os.Exit(1)
		}

		table := console.NewBoxTable(os.Stdout)
		table.SetTitle("TRONCLI - INFORMAÇÕES DO SISTEMA")
		table.SetHeaders([]string{"PROPERTY", "VALUE"})

		table.AddRow([]string{"OS", fmt.Sprintf("%s %s", profile.Distro, profile.Version)})
		table.AddRow([]string{"Init System", profile.InitSystem})
		table.AddRow([]string{"Package Manager", profile.PackageManager})
		table.AddRow([]string{"Firewall", profile.Firewall})
		table.AddRow([]string{"Network Stack", profile.NetworkStack})
		table.AddRow([]string{"Environment", profile.Environment})

		table.Render()
	},
}

var systemProfileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Exibe o perfil completo do sistema",
	Run: func(cmd *cobra.Command, args []string) {
		executor := adapter.NewExecutor()
		profile, err := services.NewProfileEngine(executor).DetectProfile()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		table := console.NewBoxTable(os.Stdout)
		table.SetTitle("TRONCLI - PERFIL DETALHADO")
		table.SetHeaders([]string{"KEY", "VALUE"})

		table.AddRow([]string{"Distro", profile.Distro})
		table.AddRow([]string{"Version", profile.Version})
		table.AddRow([]string{"InitSystem", profile.InitSystem})
		table.AddRow([]string{"PackageManager", profile.PackageManager})
		table.AddRow([]string{"Firewall", profile.Firewall})
		table.AddRow([]string{"NetworkStack", profile.NetworkStack})
		table.AddRow([]string{"Environment", profile.Environment})

		table.Render()
	},
}

var systemKernelCmd = &cobra.Command{
	Use:   "kernel",
	Short: "Exibe versão do kernel",
	Run: func(cmd *cobra.Command, args []string) {
		// Simple uname -r implementation via executor if needed,
		// but let's stick to profile service if it has kernel info.
		// If not, we might need to extend profile or run uname.
		// For now, placeholder or quick exec.
		fmt.Println("Kernel version check not fully implemented in profile service yet.")
	},
}

func init() {
	rootCmd.AddCommand(systemCmd)
	systemCmd.AddCommand(systemInfoCmd)
	systemCmd.AddCommand(systemProfileCmd)
	systemCmd.AddCommand(systemKernelCmd)
	// Add other subcommands: uptime, hostname, env
}
