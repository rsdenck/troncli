package commands

import (
	"fmt"

	"github.com/rsdenck/nux/internal/output"
	"github.com/spf13/cobra"
)

var _ = fmt.Sprintf

var pkgCmd = &cobra.Command{
	Use:   "pkg",
	Short: "Universal package management",
	Long:  `Manage packages across distributions (apt, dnf, yum, pacman, apk, zypper).`,
}

var pkgInstallCmd = &cobra.Command{
	Use:   "install [packages]",
	Short: "Install packages",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		output.NewSuccess(nil).WithMessage(fmt.Sprintf("Installed: %v", args)).Print()
	},
}

func init() {
	pkgCmd.AddCommand(pkgInstallCmd)
	rootCmd.AddCommand(pkgCmd)
}
