package commands

import (

	"github.com/rsdenck/nux/internal/output"
	"github.com/spf13/cobra"
)

var containerCmd = &cobra.Command{
	Use:   "container",
	Short: "Container management",
	Long:  `Manage containers (Docker/Podman).`,
}

var containerListCmd = &cobra.Command{
	Use:   "list",
	Short: "List containers",
	Run: func(cmd *cobra.Command, args []string) {
		output.NewList([]map[string]interface{}{}, 0).WithMessage("Container list").Print()
	},
}

func init() {
	containerCmd.AddCommand(containerListCmd)
	rootCmd.AddCommand(containerCmd)
}
