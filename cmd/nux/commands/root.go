package commands

import (
	"fmt"
	"os"

	"github.com/rsdenck/nux/internal/core/logger"
	"github.com/rsdenck/nux/internal/output"
	"github.com/spf13/cobra"
)

var _ = fmt.Println
var _ = os.Stderr

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
	Use:   "nux",
	Short: "NUX — Linux CLI Manager",
	Long: `NUX — Linux CLI Manager
-------------------------
A production-grade CLI for comprehensive Linux system administration and CLI management.

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
  - Skill Engine (manage external CLIs)
  - Agent Integration (Ollama AI)
  - Shell Autocompletion

Run 'nux onboard' for first-time setup.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		opts := logger.Options{
			Debug:    flagVerbose,
			UseColor: !flagNoColor,
			LogFile:  flagLogFile,
		}
		if err := logger.Init(opts); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		}
		output.SetFormat(flagJSON, flagYAML)
	},
}

func Execute(version, commit, date string) {
	if version != "" {
		rootCmd.Version = version
	}
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
