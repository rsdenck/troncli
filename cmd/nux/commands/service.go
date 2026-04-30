package commands

import (

	"github.com/rsdenck/nux/internal/output"
	"github.com/spf13/cobra"
)

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Service management",
	Long:  `Manage system services (systemd, openrc, sysvinit, runit).`,
}

var serviceListCmd = &cobra.Command{
	Use:   "list",
	Short: "List services",
	Run: func(cmd *cobra.Command, args []string) {
		output.NewList([]map[string]interface{}{}, 0).WithMessage("Service list").Print()
	},
}

func init() {
	serviceCmd.AddCommand(serviceListCmd)
	rootCmd.AddCommand(serviceCmd)
}
