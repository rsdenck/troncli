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
	Run: func(cmd *cobra.Command, args []string) {
		// ANSI color codes - exact match to out.md #c9a01a (RGB: 201, 160, 26)
		colorYellow := "\033[38;2;201;160;26m"
		reset := "\033[0m"
		bold := "\033[1m"

		// ASCII Art Logo with color
		fmt.Println()
		fmt.Println(colorYellow + " ███╗   ██╗██╗   ██╗██╗  ██╗██╗" + reset)
		fmt.Println(colorYellow + " ████╗  ██║██║   ██║╚██╗██╔╝" + reset)
		fmt.Println(colorYellow + " ██╔██╗ ██║██║   ██║╚███╔╝" + reset)
		fmt.Println(colorYellow + " ██║╚██╗██║██║   ██║ ██╔██╗" + reset)
		fmt.Println(colorYellow + " ██║ ╚███║╚██████╔╝██╔╝ ██╗" + reset)
		fmt.Println(colorYellow + " ╚═╝  ╚═══╝ ╚═════╝ ╚═╝  ╚═╝" + reset)
		fmt.Println()

		// Title
		fmt.Println(bold + "NUX — Linux Operations Platform" + reset)
		fmt.Println("┌────────────────────────────────────────────────┐")
		fmt.Println("│ Manage Linux systems, automation, security, skills and AI agents │")
		fmt.Println("│ from one professional command-line interface.          │")
		fmt.Println("└────────────────────────────────────────────────┘")
		fmt.Println()

		// Core Modules
		fmt.Println("CORE MODULES")
		fmt.Println("┌────────────────────────────────────────────────┐")
		fmt.Println("│  doctor        Run full health diagnostics            │")
		fmt.Println("│  system        Show OS, kernel, CPU, RAM, uptime       │")
		fmt.Println("│  disk          Manage disks, partitions and filesystems │")
		fmt.Println("│  network       Network interfaces, routes and diagnostics │")
		fmt.Println("│  service       Manage system services                  │")
		fmt.Println("│  process       Process monitoring and control          │")
		fmt.Println("│  users         User and group management              │")
		fmt.Println("│  firewall      Firewall management                   │")
		fmt.Println("└────────────────────────────────────────────────┘")
		fmt.Println()

		// Automation & Extensions
		fmt.Println("AUTOMATION & EXTENSIONS")
		fmt.Println("┌────────────────────────────────────────────────┐")
		fmt.Println("│  skill         Install and manage external CLI skills  │")
		fmt.Println("│  plugin        Legacy plugin compatibility           │")
		fmt.Println("│  bash          Execute controlled shell commands      │")
		fmt.Println("│  completion    Generate shell completion scripts    │")
		fmt.Println("└────────────────────────────────────────────────┘")
		fmt.Println()

		// Security
		fmt.Println("SECURITY")
		fmt.Println("┌────────────────────────────────────────────────┐")
		fmt.Println("│  audit         Security auditing and investigations  │")
		fmt.Println("└────────────────────────────────────────────────┘")
		fmt.Println()

		// Containers & Remote
		fmt.Println("CONTAINERS & REMOTE")
		fmt.Println("┌────────────────────────────────────────────────┐")
		fmt.Println("│  container     Docker / Podman management           │")
		fmt.Println("│  remote        SSH remote operations               │")
		fmt.Println("└────────────────────────────────────────────────┘")
		fmt.Println()

		// AI
		fmt.Println("AI")
		fmt.Println("┌────────────────────────────────────────────────┐")
		fmt.Println("│  agent         Ollama Linux assistant              │")
		fmt.Println("└────────────────────────────────────────────────┘")
		fmt.Println()

		// Setup
		fmt.Println("SETUP")
		fmt.Println("┌────────────────────────────────────────────────┐")
		fmt.Println("│  onboard       First-time guided setup             │")
		fmt.Println("└────────────────────────────────────────────────┘")
		fmt.Println()

		// Global Flags
		fmt.Println("GLOBAL FLAGS")
		fmt.Println("┌────────────────────────────────────────────────┐")
		fmt.Println("│  --json        JSON output                         │")
		fmt.Println("│  --yaml        YAML output                         │")
		fmt.Println("│  --verbose     Verbose logs                        │")
		fmt.Println("│  --quiet       Quiet mode                          │")
		fmt.Println("│  --dry-run     Simulate actions                     │")
		fmt.Println("│  --timeout     Command timeout                     │")
		fmt.Println("│  --no-color    Disable colors                      │")
		fmt.Println("│  -v            Show version                        │")
		fmt.Println("└────────────────────────────────────────────────┘")
		fmt.Println()

		// Quick Start
		fmt.Println("QUICK START")
		fmt.Println("┌────────────────────────────────────────────────┐")
		fmt.Println("│  nux onboard                              │")
		fmt.Println("│  nux doctor                               │")
		fmt.Println("│  nux service list                         │")
		fmt.Println("│  nux skill list                          │")
		fmt.Println("│  nux agent status                        │")
		fmt.Println("└────────────────────────────────────────────────┘")
		fmt.Println()

		// Help
		fmt.Println("HELP")
		fmt.Println("┌────────────────────────────────────────────────┐")
		fmt.Println("│  nux <command> --help                      │")
		fmt.Println("└────────────────────────────────────────────────┘")
		fmt.Println()

		// Version
		fmt.Println("VERSION")
		fmt.Println("┌────────────────────────────────────────────────┐")
		fmt.Printf("│  NUX v%s%32s │\n", cmd.Version, " ")
		fmt.Println("└────────────────────────────────────────────────┘")
		fmt.Println()

		fmt.Println("Ready.")
		fmt.Println()
	},
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
