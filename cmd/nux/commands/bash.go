package commands

import (
	"fmt"

	"github.com/rsdenck/nux/internal/core"
	"github.com/rsdenck/nux/internal/output"
	"github.com/spf13/cobra"
)

var bashExecutor core.Executor = &core.RealExecutor{}

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
		command := core.SanitizeInput(args[0])

		dryRun, _ := cmd.Flags().GetBool("dry-run")

		if dryRun {
			output.NewInfo(map[string]interface{}{
				"command": command,
				"dry_run": true,
				"shell":   "bash",
			}).Print()
			return
		}

		out, err := bashExecutor.CombinedOutput("bash", "-c", command)

		if err != nil {
			output.NewError(fmt.Sprintf("command failed: %s", err.Error()), "BASH_EXEC_ERROR").Print()
			return
		}

		output.NewSuccess(map[string]interface{}{
			"command": command,
			"output":  out,
			"status":  "executed",
		}).Print()
	},
}

var bashScriptCmd = &cobra.Command{
	Use:   "script [file]",
	Short: "Execute a bash script",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		scriptFile := core.SanitizeInput(args[0])

		dryRun, _ := cmd.Flags().GetBool("dry-run")

		if dryRun {
			output.NewInfo(map[string]interface{}{
				"script":  scriptFile,
				"dry_run": true,
				"command": fmt.Sprintf("bash %s", scriptFile),
			}).Print()
			return
		}

		out, err := bashExecutor.CombinedOutput("bash", scriptFile)

		if err != nil {
			output.NewError(fmt.Sprintf("script failed: %s", err.Error()), "BASH_SCRIPT_ERROR").Print()
			return
		}

		output.NewSuccess(map[string]interface{}{
			"script": scriptFile,
			"output": out,
			"status": "executed",
		}).Print()
	},
}

func init() {
	bashExecCmd.Flags().Bool("dry-run", false, "Simulate command")
	bashScriptCmd.Flags().Bool("dry-run", false, "Simulate command")
	
	bashCmd.AddCommand(bashExecCmd)
	bashCmd.AddCommand(bashScriptCmd)
	rootCmd.AddCommand(bashCmd)
}
