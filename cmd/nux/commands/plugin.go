package commands

import (

	"github.com/rsdenck/nux/internal/output"
	"github.com/spf13/cobra"
)

var pluginCmd = &cobra.Command{
	Use:   "plugin",
	Short: "Plugin management",
	Long:  `Manage NUX plugins.`,
}

var pluginListCmd = &cobra.Command{
	Use:   "list",
	Short: "List plugins",
	Run: func(cmd *cobra.Command, args []string) {
		output.NewList([]map[string]interface{}{}, 0).WithMessage("Plugin list").Print()
	},
}

func init() {
	pluginCmd.AddCommand(pluginListCmd)
	rootCmd.AddCommand(pluginCmd)
}
