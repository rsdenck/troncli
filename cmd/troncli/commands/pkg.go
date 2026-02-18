package commands

import (
	"fmt"
	"os"

	"github.com/mascli/troncli/internal/core/adapter"
	"github.com/mascli/troncli/internal/core/services"
	"github.com/mascli/troncli/internal/modules/pkg"
	"github.com/spf13/cobra"
)

var pkgCmd = &cobra.Command{
	Use:   "pkg",
	Short: "Gerenciador de Pacotes Universal",
	Long:  `Instala, remove e gerencia pacotes de forma transparente em apt, dnf, yum, pacman, apk e zypper.`,
}

func getPkgManager() (*pkg.UniversalPackageManager, error) {
	// Need to detect profile first to init manager
	executor := adapter.NewExecutor()
	profileEngine := services.NewProfileEngine(executor)
	profile, err := profileEngine.DetectProfile()
	if err != nil {
		return nil, fmt.Errorf("failed to detect system profile: %w", err)
	}

	return pkg.NewUniversalPackageManager(executor, profile), nil
}

var pkgInstallCmd = &cobra.Command{
	Use:   "install [package]",
	Short: "Instala um pacote",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := getPkgManager()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("Installing %s...\n", args[0])
		if err := manager.Install(args[0]); err != nil {
			fmt.Printf("Error installing package: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Package installed successfully.")
	},
}

var pkgRemoveCmd = &cobra.Command{
	Use:   "remove [package]",
	Short: "Remove um pacote",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := getPkgManager()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("Removing %s...\n", args[0])
		if err := manager.Remove(args[0]); err != nil {
			fmt.Printf("Error removing package: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Package removed successfully.")
	},
}

var pkgUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Atualiza a lista de pacotes",
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := getPkgManager()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("Updating package lists...")
		if err := manager.Update(); err != nil {
			fmt.Printf("Error updating lists: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Package lists updated.")
	},
}

var pkgSearchCmd = &cobra.Command{
	Use:   "search [term]",
	Short: "Pesquisa por pacotes",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := getPkgManager()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		results, err := manager.Search(args[0])
		if err != nil {
			fmt.Printf("Error searching: %v\n", err)
			os.Exit(1)
		}

		for _, res := range results {
			fmt.Println(res)
		}
	},
}

var pkgUpgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Atualiza todos os pacotes do sistema",
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := getPkgManager()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("Upgrading system packages...")
		if err := manager.Upgrade(); err != nil {
			fmt.Printf("Error upgrading packages: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("System upgraded successfully.")
	},
}

func init() {
	rootCmd.AddCommand(pkgCmd)
	pkgCmd.AddCommand(pkgInstallCmd)
	pkgCmd.AddCommand(pkgRemoveCmd)
	pkgCmd.AddCommand(pkgUpdateCmd)
	pkgCmd.AddCommand(pkgUpgradeCmd)
	pkgCmd.AddCommand(pkgSearchCmd)
}
