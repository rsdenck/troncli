package commands

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/rsdenck/nux/internal/output"
	"github.com/spf13/cobra"
)

var pkgCmd = &cobra.Command{
	Use:   "pkg",
	Short: "Universal package management",
	Long:  `Manage packages across distributions (apt, dnf, yum, pacman, apk, zypper).`,
}

// detectPackageManager detects the package manager for the current distribution
func detectPackageManager() string {
	// Check for package managers in order of preference
	managers := []struct {
		name    string
		command string
	}{
		{"apt", "apt"},
		{"dnf", "dnf"},
		{"yum", "yum"},
		{"pacman", "pacman"},
		{"apk", "apk"},
		{"zypper", "zypper"},
	}

	for _, m := range managers {
		if _, err := exec.LookPath(m.command); err == nil {
			return m.name
		}
	}

	return "unknown"
}

var pkgInstallCmd = &cobra.Command{
	Use:   "install [packages]",
	Short: "Install packages",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pkgManager := detectPackageManager()
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		var command string
		var cmdArgs []string

		switch pkgManager {
		case "apt":
			command = "apt"
			cmdArgs = append([]string{"install", "-y"}, args...)
		case "dnf":
			command = "dnf"
			cmdArgs = append([]string{"install", "-y"}, args...)
		case "yum":
			command = "yum"
			cmdArgs = append([]string{"install", "-y"}, args...)
		case "pacman":
			command = "pacman"
			cmdArgs = append([]string{"-S", "--noconfirm"}, args...)
		case "apk":
			command = "apk"
			cmdArgs = append([]string{"add"}, args...)
		case "zypper":
			command = "zypper"
			cmdArgs = append([]string{"install", "-y"}, args...)
		default:
			output.NewError("unsupported package manager", "PKG_UNSUPPORTED").Print()
			return
		}

		if dryRun {
			output.NewInfo(map[string]interface{}{
				"packages":        args,
				"package_manager": pkgManager,
				"dry_run":         true,
				"command":         fmt.Sprintf("%s %s", command, strings.Join(cmdArgs, " ")),
			}).Print()
			return
		}

		execCmd := exec.Command(command, cmdArgs...)
		out, err := execCmd.CombinedOutput()

		if err != nil {
			output.NewError(fmt.Sprintf("install failed: %s - %s", err.Error(), strings.TrimSpace(string(out))), "PKG_INSTALL_ERROR").Print()
			return
		}

		output.NewSuccess(map[string]interface{}{
			"packages":        args,
			"package_manager": pkgManager,
			"status":          "installed",
			"output":          strings.TrimSpace(string(out)),
		}).Print()
	},
}

var pkgUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update package lists",
	Run: func(cmd *cobra.Command, args []string) {
		pkgManager := detectPackageManager()
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		var command string
		var cmdArgs []string

		switch pkgManager {
		case "apt":
			command = "apt"
			cmdArgs = []string{"update"}
		case "dnf":
			command = "dnf"
			cmdArgs = []string{"check-update"}
		case "yum":
			command = "yum"
			cmdArgs = []string{"check-update"}
		case "pacman":
			command = "pacman"
			cmdArgs = []string{"-Sy"}
		case "apk":
			command = "apk"
			cmdArgs = []string{"update"}
		case "zypper":
			command = "zypper"
			cmdArgs = []string{"refresh"}
		default:
			output.NewError("unsupported package manager", "PKG_UNSUPPORTED").Print()
			return
		}

		if dryRun {
			output.NewInfo(map[string]interface{}{
				"package_manager": pkgManager,
				"dry_run":         true,
				"command":         fmt.Sprintf("%s %s", command, strings.Join(cmdArgs, " ")),
			}).Print()
			return
		}

		execCmd := exec.Command(command, cmdArgs...)
		out, err := execCmd.CombinedOutput()

		if err != nil && pkgManager != "dnf" && pkgManager != "yum" {
			output.NewError(fmt.Sprintf("update failed: %s - %s", err.Error(), strings.TrimSpace(string(out))), "PKG_UPDATE_ERROR").Print()
			return
		}

		output.NewSuccess(map[string]interface{}{
			"package_manager": pkgManager,
			"status":          "updated",
			"output":          strings.TrimSpace(string(out)),
		}).Print()
	},
}

var pkgListCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed packages",
	Run: func(cmd *cobra.Command, args []string) {
		pkgManager := detectPackageManager()

		var command string
		var cmdArgs []string

		switch pkgManager {
		case "apt":
			command = "dpkg"
			cmdArgs = []string{"-l"}
		case "dnf", "yum":
			command = "rpm"
			cmdArgs = []string{"-qa"}
		case "pacman":
			command = "pacman"
			cmdArgs = []string{"-Q"}
		case "apk":
			command = "apk"
			cmdArgs = []string{"info", "-v"}
		case "zypper":
			command = "zypper"
			cmdArgs = []string{"se", "--installed-only"}
		default:
			output.NewError("unsupported package manager", "PKG_UNSUPPORTED").Print()
			return
		}

		execCmd := exec.Command(command, cmdArgs...)
		out, err := execCmd.CombinedOutput()

		if err != nil {
			output.NewError(fmt.Sprintf("list failed: %s - %s", err.Error(), strings.TrimSpace(string(out))), "PKG_LIST_ERROR").Print()
			return
		}

		output.NewSuccess(map[string]interface{}{
			"package_manager": pkgManager,
			"output":          strings.TrimSpace(string(out)),
		}).Print()
	},
}

var pkgSearchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search for packages",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		query := args[0]
		pkgManager := detectPackageManager()

		var command string
		var cmdArgs []string

		switch pkgManager {
		case "apt":
			command = "apt"
			cmdArgs = []string{"search", query}
		case "dnf":
			command = "dnf"
			cmdArgs = []string{"search", query}
		case "yum":
			command = "yum"
			cmdArgs = []string{"search", query}
		case "pacman":
			command = "pacman"
			cmdArgs = []string{"-Ss", query}
		case "apk":
			command = "apk"
			cmdArgs = []string{"search", query}
		case "zypper":
			command = "zypper"
			cmdArgs = []string{"se", query}
		default:
			output.NewError("unsupported package manager", "PKG_UNSUPPORTED").Print()
			return
		}

		execCmd := exec.Command(command, cmdArgs...)
		out, err := execCmd.CombinedOutput()

		if err != nil {
			output.NewError(fmt.Sprintf("search failed: %s - %s", err.Error(), strings.TrimSpace(string(out))), "PKG_SEARCH_ERROR").Print()
			return
		}

		output.NewSuccess(map[string]interface{}{
			"query":           query,
			"package_manager": pkgManager,
			"output":          strings.TrimSpace(string(out)),
		}).Print()
	},
}

func init() {
	pkgCmd.AddCommand(pkgInstallCmd)
	pkgCmd.AddCommand(pkgUpdateCmd)
	pkgCmd.AddCommand(pkgListCmd)
	pkgCmd.AddCommand(pkgSearchCmd)
	rootCmd.AddCommand(pkgCmd)

	// Set GOOS for detection
	_ = runtime.GOOS
}
