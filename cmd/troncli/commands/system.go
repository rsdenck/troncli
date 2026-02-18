package commands

import (
	"fmt"
	"os"

	"github.com/mascli/troncli/internal/core/adapter"
	"github.com/mascli/troncli/internal/core/services"
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
		fmt.Printf("OS: %s %s\n", profile.Distro, profile.Version)
		fmt.Printf("Init System: %s\n", profile.InitSystem)
		fmt.Printf("Package Manager: %s\n", profile.PackageManager)
		fmt.Printf("Firewall: %s\n", profile.Firewall)
		fmt.Printf("Network Stack: %s\n", profile.NetworkStack)
		fmt.Printf("Environment: %s\n", profile.Environment)
	},
}

var systemProfileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Exibe o perfil completo do sistema (JSON)",
	Run: func(cmd *cobra.Command, args []string) {
		// Similar to info but maybe just JSON dump
		executor := adapter.NewExecutor()
		profile, err := services.NewProfileEngine(executor).DetectProfile()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		// TODO: Implement proper JSON serialization if requested via flag,
		// but for now just print struct
		fmt.Printf("%+v\n", profile)
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
