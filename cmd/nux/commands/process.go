package commands

import (

	"github.com/rsdenck/nux/internal/output"
	"github.com/spf13/cobra"
)

var processCmd = &cobra.Command{
	Use:   "process",
	Short: "Process management",
	Long:  `Manage system processes.`,
}

var processListCmd = &cobra.Command{
	Use:   "list",
	Short: "List processes",
	Run: func(cmd *cobra.Command, args []string) {
		output.NewList([]map[string]interface{}{}, 0).WithMessage("Process list").Print()
	},
}

func init() {
	processCmd.AddCommand(processListCmd)
	rootCmd.AddCommand(processCmd)
}
