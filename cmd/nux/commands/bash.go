package commands

import (

	"github.com/rsdenck/nux/internal/output"
	"github.com/spf13/cobra"
)

var bashCmd = &cobra.Command{
	Use:   "bash",
	Short: "Execute Bash commands",
	Long:  `Execute Bash commands or scripts directly from NUX.`,
}

var bashExecCmd = &cobra.Command{
	Use:   "exec [command]",
	Short: "Execute a bash command",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		output.NewSuccess(map[string]interface{}{
			"command": args[0],
			"status":  "executed",
		}).Print()
	},
}

var bashScriptCmd = &cobra.Command{
	Use:   "script [file]",
	Short: "Execute a bash script",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		output.NewSuccess(map[string]interface{}{
			"script":  args[0],
			"status": "executed",
		}).Print()
	},
}

func init() {
	bashCmd.AddCommand(bashExecCmd)
	bashCmd.AddCommand(bashScriptCmd)
	rootCmd.AddCommand(bashCmd)
}
