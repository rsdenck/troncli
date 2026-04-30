package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/rsdenck/nux/internal/core"
	"github.com/rsdenck/nux/internal/output"
	"github.com/spf13/cobra"
)

var remoteExecutor core.Executor = &core.RealExecutor{}

var remoteCmd = &cobra.Command{
	Use:   "remote",
	Short: "Remote SSH connections",
	Long:  `Manage remote SSH connections.`,
}

var remoteListCmd = &cobra.Command{
	Use:   "list",
	Short: "List remote connections",
	Run: func(cmd *cobra.Command, args []string) {
		// Read SSH config file
		sshConfig := os.Getenv("HOME") + "/.ssh/config"
		file, err := os.Open(sshConfig)
		if err != nil {
			output.NewError(fmt.Sprintf("failed to read SSH config: %s", err.Error()), "REMOTE_LIST_ERROR").Print()
			return
		}
		defer file.Close()

		items := []map[string]interface{}{}
		scanner := bufio.NewScanner(file)
		currentHost := ""

		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if strings.HasPrefix(strings.ToLower(line), "host ") {
				currentHost = strings.TrimPrefix(strings.ToLower(line), "host ")
			} else if strings.HasPrefix(line, "HostName") && currentHost != "" {
				hostname := strings.TrimSpace(strings.TrimPrefix(line, "HostName"))
				items = append(items, map[string]interface{}{
					"host":     currentHost,
					"hostname": hostname,
				})
				currentHost = ""
			}
		}

		output.NewList(items, len(items)).WithMessage("Remote connections (SSH config)").Print()
	},
}

var remoteExecCmd = &cobra.Command{
	Use:   "exec <host> [command]",
	Short: "Execute command on remote host",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		host := core.SanitizeInput(args[0])
		command := core.SanitizeInput(strings.Join(args[1:], " "))

		dryRun, _ := cmd.Flags().GetBool("dry-run")

		if dryRun {
			output.NewInfo(map[string]interface{}{
				"host":    host,
				"command": command,
				"dry_run": true,
				"ssh_cmd": fmt.Sprintf("ssh %s '%s'", host, command),
			}).Print()
			return
		}

		out, err := remoteExecutor.CombinedOutput("ssh", host, command)

		if err != nil {
			output.NewError(fmt.Sprintf("remote execution failed: %s", err.Error()), "REMOTE_EXEC_ERROR").Print()
			return
		}

		output.NewSuccess(map[string]interface{}{
			"host":    host,
			"command": command,
			"output":  out,
		}).Print()
	},
}

func init() {
	remoteExecCmd.Flags().Bool("dry-run", false, "Simulate command")
	
	remoteCmd.AddCommand(remoteListCmd)
	remoteCmd.AddCommand(remoteExecCmd)
	rootCmd.AddCommand(remoteCmd)
}
