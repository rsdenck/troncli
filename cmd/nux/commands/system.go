package commands

import (

	"github.com/rsdenck/nux/internal/output"
	"github.com/spf13/cobra"
)

var systemCmd = &cobra.Command{
	Use:   "system",
	Short: "System information",
	Long:  `Show system information and statistics.`,
}

var systemInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show system information",
	Run: func(cmd *cobra.Command, args []string) {
		output.NewSuccess(map[string]interface{}{
			"hostname": "localhost",
			"status":   "ok",
		}).Print()
	},
}

func init() {
	systemCmd.AddCommand(systemInfoCmd)
	rootCmd.AddCommand(systemCmd)
}
