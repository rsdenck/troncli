package commands

import (

	"github.com/rsdenck/nux/internal/output"
	"github.com/spf13/cobra"
)

var remoteCmd = &cobra.Command{
	Use:   "remote",
	Short: "Remote SSH connections",
	Long:  `Manage remote SSH connections.`,
}

var remoteListCmd = &cobra.Command{
	Use:   "list",
	Short: "List remote connections",
	Run: func(cmd *cobra.Command, args []string) {
		output.NewList([]map[string]interface{}{}, 0).WithMessage("Remote connections").Print()
	},
}

func init() {
	remoteCmd.AddCommand(remoteListCmd)
	rootCmd.AddCommand(remoteCmd)
}
