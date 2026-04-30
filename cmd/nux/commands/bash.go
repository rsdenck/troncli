package commands

import (
	"fmt"
	"os/exec"
	"strings"

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
		command := args[0]

		dryRun, _ := cmd.Flags().GetBool("dry-run")

		if dryRun {
			output.NewInfo(map[string]interface{}{
				"command": command,
				"dry_run": true,
				"shell":   "bash",
			}).Print()
			return
		}

		// Execute command via bash -c
		bashCmd := exec.Command("bash", "-c", command)
		out, err := bashCmd.CombinedOutput()

		if err != nil {
			output.NewError(fmt.Sprintf("command failed: %s - %s", err.Error(), strings.TrimSpace(string(out))), "BASH_EXEC_ERROR").Print()
			return
		}

		output.NewSuccess(map[string]interface{}{
			"command": command,
			"output":  strings.TrimSpace(string(out)),
			"status":  "executed",
		}).Print()
	},
}

var bashScriptCmd = &cobra.Command{
	Use:   "script [file]",
	Short: "Execute a bash script",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		scriptFile := args[0]

		dryRun, _ := cmd.Flags().GetBool("dry-run")

		if dryRun {
			output.NewInfo(map[string]interface{}{
				"script":  scriptFile,
				"dry_run": true,
				"command": fmt.Sprintf("bash %s", scriptFile),
			}).Print()
			return
		}

		// Execute script file
		bashCmd := exec.Command("bash", scriptFile)
		out, err := bashCmd.CombinedOutput()

		if err != nil {
			output.NewError(fmt.Sprintf("script failed: %s - %s", err.Error(), strings.TrimSpace(string(out))), "BASH_SCRIPT_ERROR").Print()
			return
		}

		output.NewSuccess(map[string]interface{}{
			"script": scriptFile,
			"output": strings.TrimSpace(string(out)),
			"status": "executed",
		}).Print()
	},
}

func init() {
	bashCmd.AddCommand(bashExecCmd)
	bashCmd.AddCommand(bashScriptCmd)
	rootCmd.AddCommand(bashCmd)
}
