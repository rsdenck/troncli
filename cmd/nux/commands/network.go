package commands

import (

	"github.com/rsdenck/nux/internal/output"
	"github.com/spf13/cobra"
)

var networkCmd = &cobra.Command{
	Use:   "network",
	Short: "Network management",
	Long:  `Manage network configuration, interfaces, routes, and diagnostics.`,
}

var networkListCmd = &cobra.Command{
	Use:   "list",
	Short: "List network interfaces",
	Run: func(cmd *cobra.Command, args []string) {
		output.NewList([]map[string]interface{}{}, 0).WithMessage("Network interfaces").Print()
	},
}

func init() {
	networkCmd.AddCommand(networkListCmd)
	rootCmd.AddCommand(networkCmd)
}
