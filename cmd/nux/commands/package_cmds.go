package commands

import (
	"fmt"
	"os/exec"

	"github.com/rsdenck/nux/internal/core"
	"github.com/rsdenck/nux/internal/output"
	"github.com/spf13/cobra"
)

var pkgExecutor core.Executor = &core.RealExecutor{}

func detectPackageManager() string {
	if _, err := exec.LookPath("apt"); err == nil {
		return "apt"
	}
	if _, err := exec.LookPath("dnf"); err == nil {
		return "dnf"
	}
	if _, err := exec.LookPath("yum"); err == nil {
		return "yum"
	}
	if _, err := exec.LookPath("pacman"); err == nil {
		return "pacman"
	}
	if _, err := exec.LookPath("zypper"); err == nil {
		return "zypper"
	}
	if _, err := exec.LookPath("apk"); err == nil {
		return "apk"
	}
	return ""
}

var installCmd = &cobra.Command{
	Use:     "install <packages>",
	Short:   "Instala pacotes (universal)",
	Aliases: []string{"i"},
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pm := detectPackageManager()
		if pm == "" {
			output.NewError("Nenhum gerenciador de pacotes suportado encontrado", "PKG_NO_PM").Print()
			return
		}

		dryRun, _ := cmd.Flags().GetBool("dry-run")

		switch pm {
		case "apt":
			if dryRun {
				output.NewInfo(map[string]interface{}{"action": "install", "packages": args, "dry_run": true}).Print()
				return
			}
			pkgExecutor.CombinedOutput("sudo", append([]string{"apt", "install", "-y"}, args...)...)
		case "pacman":
			if dryRun {
				output.NewInfo(map[string]interface{}{"action": "install", "packages": args, "dry_run": true}).Print()
				return
			}
			pkgExecutor.CombinedOutput("sudo", append([]string{"pacman", "-S", "--noconfirm"}, args...)...)
		default:
			output.NewError(fmt.Sprintf("Gerenciador não suportado: %s", pm), "PKG_UNSUPPORTED").Print()
		}

		output.NewSuccess(map[string]interface{}{"packages": args, "status": "installed", "manager": pm}).Print()
	},
}

var removeCmd = &cobra.Command{
	Use:     "remove <packages>",
	Short:   "Remove pacotes (universal)",
	Aliases: []string{"rm"},
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pm := detectPackageManager()
		if pm == "" {
			output.NewError("Nenhum gerenciador de pacotes suportado encontrado", "PKG_NO_PM").Print()
			return
		}

		dryRun, _ := cmd.Flags().GetBool("dry-run")

		switch pm {
		case "apt":
			if dryRun {
				output.NewInfo(map[string]interface{}{"action": "remove", "packages": args, "dry_run": true}).Print()
				return
			}
			pkgExecutor.CombinedOutput("sudo", append([]string{"apt", "remove", "-y"}, args...)...)
		case "pacman":
			if dryRun {
				output.NewInfo(map[string]interface{}{"action": "remove", "packages": args, "dry_run": true}).Print()
				return
			}
			pkgExecutor.CombinedOutput("sudo", append([]string{"pacman", "-R", "--noconfirm"}, args...)...)
		}

		output.NewSuccess(map[string]interface{}{"packages": args, "status": "removed", "manager": pm}).Print()
	},
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Atualiza lista de pacotes",
	Run: func(cmd *cobra.Command, args []string) {
		pm := detectPackageManager()
		if pm == "" {
			output.NewError("Nenhum gerenciador de pacotes suportado encontrado", "PKG_NO_PM").Print()
			return
		}

		dryRun, _ := cmd.Flags().GetBool("dry-run")

		switch pm {
		case "apt":
			if dryRun {
				output.NewInfo(map[string]interface{}{"action": "update", "manager": pm, "dry_run": true}).Print()
				return
			}
			exec.Command("sudo", "apt", "update").Run()
		case "pacman":
			if dryRun {
				output.NewInfo(map[string]interface{}{"action": "update", "manager": pm, "dry_run": true}).Print()
				return
			}
			exec.Command("sudo", "pacman", "-Sy").Run()
		}

		output.NewSuccess(map[string]interface{}{"action": "update", "manager": pm, "status": "updated"}).Print()
	},
}

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Atualiza pacotes instalados",
	Run: func(cmd *cobra.Command, args []string) {
		pm := detectPackageManager()
		if pm == "" {
			output.NewError("Nenhum gerenciador de pacotes suportado encontrado", "PKG_NO_PM").Print()
			return
		}

		dryRun, _ := cmd.Flags().GetBool("dry-run")

		switch pm {
		case "apt":
			if dryRun {
				output.NewInfo(map[string]interface{}{"action": "upgrade", "manager": pm, "dry_run": true}).Print()
				return
			}
			exec.Command("sudo", "apt", "upgrade", "-y").Run()
		case "pacman":
			if dryRun {
				output.NewInfo(map[string]interface{}{"action": "upgrade", "manager": pm, "dry_run": true}).Print()
				return
			}
			exec.Command("sudo", "pacman", "-Syu").Run()
		}

		output.NewSuccess(map[string]interface{}{"action": "upgrade", "manager": pm, "status": "upgraded"}).Print()
	},
}

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Limpa cache de pacotes",
	Run: func(cmd *cobra.Command, args []string) {
		pm := detectPackageManager()
		if pm == "" {
			output.NewError("Nenhum gerenciador de pacotes suportado encontrado", "PKG_NO_PM").Print()
			return
		}

		dryRun, _ := cmd.Flags().GetBool("dry-run")

		switch pm {
		case "apt":
			if dryRun {
				output.NewInfo(map[string]interface{}{"action": "clean", "manager": pm, "dry_run": true}).Print()
				return
			}
			exec.Command("sudo", "apt", "clean").Run()
		case "pacman":
			if dryRun {
				output.NewInfo(map[string]interface{}{"action": "clean", "manager": pm, "dry_run": true}).Print()
				return
			}
			exec.Command("sudo", "pacman", "-Sc").Run()
		}

		output.NewSuccess(map[string]interface{}{"action": "clean", "manager": pm, "status": "cleaned"}).Print()
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(upgradeCmd)
	rootCmd.AddCommand(cleanCmd)
}
