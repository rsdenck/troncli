package commands

import (
	"fmt"
	"os"

	"github.com/mascli/troncli/internal/core/logger"
	"github.com/mascli/troncli/internal/ui"
	"github.com/spf13/cobra"
)

var (
	// Global flags
	flagJSON    bool
	flagYAML    bool
	flagQuiet   bool
	flagDryRun  bool
	flagTimeout int
	flagVerbose bool
	flagNoColor bool
	flagLogFile string
)

var rootCmd = &cobra.Command{
	Use:   "troncli",
	Short: "TRONCLI — Universal Linux Administration Platform",
	Long: `TRONCLI — Universal Linux Administration Platform (v1.0)
-------------------------------------------------------
A production-grade CLI and TUI for comprehensive Linux system administration.
Designed for DevOps, SREs, and SysAdmins who need reliable, idempotent,
and verifiable system modifications.

Features:
  - Universal Package Management (apt, dnf, yum, pacman, apk, zypper)
  - Service Management (systemd, openrc, sysvinit)
  - Process & Resource Monitoring
  - Network Configuration & Diagnostics
  - User & Group Management
  - Disk & Filesystem Operations
  - Security Auditing & Compliance
  - Container Management (Docker/Podman)
  - Remote Execution (SSH)
  - Plugin System
  - Shell Autocompletion

Launch without arguments to start the interactive TUI mode.`,
	Version: "v0.1.1-beta",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		opts := logger.Options{
			Debug:    flagVerbose,
			UseColor: !flagNoColor,
			LogFile:  flagLogFile,
		}
		if err := logger.Init(opts); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Default behavior: Launch TUI if no subcommand is provided
		if len(args) == 0 {
			startTUI()
		} else {
			cmd.Help()
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Global flags available to all commands
	rootCmd.PersistentFlags().BoolVar(&flagJSON, "json", false, "Output in JSON format")
	rootCmd.PersistentFlags().BoolVar(&flagYAML, "yaml", false, "Output in YAML format")
	rootCmd.PersistentFlags().BoolVar(&flagQuiet, "quiet", false, "Suppress output")
	rootCmd.PersistentFlags().BoolVar(&flagDryRun, "dry-run", false, "Simulate execution without making changes")
	rootCmd.PersistentFlags().IntVar(&flagTimeout, "timeout", 30, "Timeout in seconds")
	rootCmd.PersistentFlags().BoolVar(&flagVerbose, "verbose", false, "Enable verbose logging")
	rootCmd.PersistentFlags().StringVar(&flagLogFile, "log-file", "", "Log file path (enables debug logging to file)")
	rootCmd.PersistentFlags().BoolVar(&flagNoColor, "no-color", false, "Disable color output")
}

func startTUI() {
	app, err := ui.NewApp()
	if err != nil {
		fmt.Printf("Failed to initialize application: %v\n", err)
		os.Exit(1)
	}
	if err := app.Run(); err != nil {
		fmt.Printf("Error running application: %v\n", err)
		os.Exit(1)
	}
}
