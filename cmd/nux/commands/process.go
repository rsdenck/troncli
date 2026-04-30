package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/rsdenck/nux/internal/core"
	"github.com/rsdenck/nux/internal/output"
	"github.com/spf13/cobra"
)

var processExecutor core.Executor = &core.RealExecutor{}

var processCmd = &cobra.Command{
	Use:   "process",
	Short: "Process management",
	Long:  `Manage system processes.`,
}

var processListCmd = &cobra.Command{
	Use:   "list",
	Short: "List processes",
	Run: func(cmd *cobra.Command, args []string) {
		out, err := processExecutor.CombinedOutput("ps", "aux")

		if err != nil {
			output.NewError(fmt.Sprintf("failed to list processes: %s", err.Error()), "PROCESS_LIST_ERROR").Print()
			return
		}

		// Parse ps output
		lines := strings.Split(out, "\n")
		if len(lines) < 1 {
			output.NewList([]map[string]interface{}{}, 0).WithMessage("Process list").Print()
			return
		}

		items := []map[string]interface{}{}
		for _, line := range lines[1:] {
			fields := strings.Fields(line)
			if len(fields) < 11 {
				continue
			}

			item := map[string]interface{}{
				"user":    fields[0],
				"pid":     fields[1],
				"cpu":     fields[2],
				"mem":     fields[3],
				"vsz":     fields[4],
				"rss":     fields[5],
				"tty":     fields[6],
				"stat":    fields[7],
				"start":   fields[8],
				"time":    fields[9],
				"command": strings.Join(fields[10:], " "),
			}
			items = append(items, item)
		}

		output.NewList(items, len(items)).WithMessage("Process list").Print()
	},
}

var processKillCmd = &cobra.Command{
	Use:   "kill <pid> [signal]",
	Short: "Kill a process",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pid := core.SanitizeInput(args[0])
		signal := "TERM"

		if len(args) > 1 {
			signal = core.SanitizeInput(args[1])
		}

		dryRun, _ := cmd.Flags().GetBool("dry-run")

		if dryRun {
			output.NewInfo(map[string]interface{}{
				"pid":     pid,
				"signal":  signal,
				"dry_run": true,
				"command": fmt.Sprintf("kill -%s %s", signal, pid),
			}).Print()
			return
		}

		_, err := processExecutor.CombinedOutput("kill", "-"+signal, pid)

		if err != nil {
			output.NewError(fmt.Sprintf("failed to kill process: %s", err.Error()), "PROCESS_KILL_ERROR").Print()
			return
		}

		output.NewSuccess(map[string]interface{}{
			"pid":    pid,
			"signal": signal,
			"status": "killed",
		}).Print()
	},
}

var processInfoCmd = &cobra.Command{
	Use:   "info <pid>",
	Short: "Show process information",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pid := core.SanitizeInput(args[0])

		// Validate PID is numeric
		if _, err := strconv.Atoi(pid); err != nil {
			output.NewError("invalid PID: must be numeric", "PROCESS_INVALID_PID").Print()
			return
		}

		// Get process info using ps
		out, err := processExecutor.CombinedOutput("ps", "-p", pid, "-o", "pid,user,%cpu,%mem,vsz,rss,tty,stat,start,time,comm")

		if err != nil {
			output.NewError(fmt.Sprintf("failed to get process info: %s", err.Error()), "PROCESS_INFO_ERROR").Print()
			return
		}

		// Also get full command line
		cmdlineOut, _ := processExecutor.CombinedOutput("cat", fmt.Sprintf("/proc/%s/cmdline", pid))

		output.NewSuccess(map[string]interface{}{
			"pid":     pid,
			"ps_info": out,
			"cmdline": cmdlineOut,
		}).Print()
	},
}

func init() {
	processKillCmd.Flags().Bool("dry-run", false, "Simulate command")
	
	processCmd.AddCommand(processListCmd)
	processCmd.AddCommand(processKillCmd)
	processCmd.AddCommand(processInfoCmd)
	rootCmd.AddCommand(processCmd)
}
