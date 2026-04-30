package commands

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/rsdenck/nux/internal/core"
	"github.com/rsdenck/nux/internal/output"
	"github.com/spf13/cobra"
)

var containerExecutor core.Executor = &core.RealExecutor{}

var containerCmd = &cobra.Command{
	Use:   "container",
	Short: "Container management",
	Long:  `Manage containers (Docker/Podman).`,
}

func detectContainerRuntime() string {
	runtimes := []struct {
		name    string
		command string
	}{
		{"docker", "docker"},
		{"podman", "podman"},
	}

	for _, r := range runtimes {
		if _, err := exec.LookPath(r.command); err == nil {
			return r.name
		}
	}

	return "none"
}

var containerListCmd = &cobra.Command{
	Use:   "list",
	Short: "List containers",
	Run: func(cmd *cobra.Command, args []string) {
		runtime := detectContainerRuntime()

		if runtime == "none" {
			output.NewError("no container runtime found (docker/podman)", "CONTAINER_RUNTIME_MISSING").Print()
			return
		}

		// Try JSON output first
		out, err := containerExecutor.CombinedOutput(runtime, "ps", "-a", "--format", "{{json .}}")

		if err != nil {
			// Fallback to text output
			textOut, _ := containerExecutor.CombinedOutput(runtime, "ps", "-a")
			output.NewSuccess(map[string]interface{}{
				"runtime": runtime,
				"output":  textOut,
			}).Print()
			return
		}

		// Parse JSON output
		lines := strings.Split(out, "\n")
		items := []map[string]interface{}{}

		for _, line := range lines {
			if line == "" {
				continue
			}

			var container map[string]interface{}
			if err := json.Unmarshal([]byte(line), &container); err == nil {
				items = append(items, container)
			}
		}

		output.NewList(items, len(items)).WithMessage("Container list").Print()
	},
}

var containerRunCmd = &cobra.Command{
	Use:   "run <image>",
	Short: "Run a container",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		image := core.SanitizeInput(args[0])
		runtime := detectContainerRuntime()

		if runtime == "none" {
			output.NewError("no container runtime found (docker/podman)", "CONTAINER_RUNTIME_MISSING").Print()
			return
		}

		name, _ := cmd.Flags().GetString("name")
		ports, _ := cmd.Flags().GetString("ports")
		detach, _ := cmd.Flags().GetBool("detach")

		dryRun, _ := cmd.Flags().GetBool("dry-run")

		cmdArgs := []string{"run"}
		if detach {
			cmdArgs = append(cmdArgs, "-d")
		}
		if name != "" {
			cmdArgs = append(cmdArgs, "--name", name)
		}
		if ports != "" {
			cmdArgs = append(cmdArgs, "-p", ports)
		}
		cmdArgs = append(cmdArgs, image)

		if dryRun {
			output.NewInfo(map[string]interface{}{
				"runtime": runtime,
				"image":   image,
				"dry_run": true,
				"command": fmt.Sprintf("%s %s", runtime, strings.Join(cmdArgs, " ")),
			}).Print()
			return
		}

		_, err := containerExecutor.CombinedOutput(runtime, cmdArgs...)

		if err != nil {
			output.NewError(fmt.Sprintf("failed to run container: %s", err.Error()), "CONTAINER_RUN_ERROR").Print()
			return
		}

		output.NewSuccess(map[string]interface{}{
			"runtime": runtime,
			"image":   image,
			"status":  "running",
		}).Print()
	},
}

var containerStopCmd = &cobra.Command{
	Use:   "stop <container>",
	Short: "Stop a container",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		container := core.SanitizeInput(args[0])
		runtime := detectContainerRuntime()

		if runtime == "none" {
			output.NewError("no container runtime found (docker/podman)", "CONTAINER_RUNTIME_MISSING").Print()
			return
		}

		dryRun, _ := cmd.Flags().GetBool("dry-run")

		if dryRun {
			output.NewInfo(map[string]interface{}{
				"runtime":   runtime,
				"container": container,
				"dry_run":   true,
				"command":   fmt.Sprintf("%s stop %s", runtime, container),
			}).Print()
			return
		}

		_, err := containerExecutor.CombinedOutput(runtime, "stop", container)

		if err != nil {
			output.NewError(fmt.Sprintf("failed to stop container: %s", err.Error()), "CONTAINER_STOP_ERROR").Print()
			return
		}

		output.NewSuccess(map[string]interface{}{
			"runtime":   runtime,
			"container": container,
			"status":    "stopped",
		}).Print()
	},
}

var containerRemoveCmd = &cobra.Command{
	Use:   "remove <container>",
	Short: "Remove a container",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		container := core.SanitizeInput(args[0])
		runtime := detectContainerRuntime()

		if runtime == "none" {
			output.NewError("no container runtime found (docker/podman)", "CONTAINER_RUNTIME_MISSING").Print()
			return
		}

		force, _ := cmd.Flags().GetBool("force")
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		cmdArgs := []string{"rm"}
		if force {
			cmdArgs = append(cmdArgs, "-f")
		}
		cmdArgs = append(cmdArgs, container)

		if dryRun {
			output.NewInfo(map[string]interface{}{
				"runtime":   runtime,
				"container": container,
				"dry_run":   true,
				"command":   fmt.Sprintf("%s %s", runtime, strings.Join(cmdArgs, " ")),
			}).Print()
			return
		}

		_, err := containerExecutor.CombinedOutput(runtime, cmdArgs...)

		if err != nil {
			output.NewError(fmt.Sprintf("failed to remove container: %s", err.Error()), "CONTAINER_REMOVE_ERROR").Print()
			return
		}

		output.NewSuccess(map[string]interface{}{
			"runtime":   runtime,
			"container": container,
			"status":    "removed",
		}).Print()
	},
}

func init() {
	containerRunCmd.Flags().String("name", "", "Container name")
	containerRunCmd.Flags().String("ports", "", "Port mappings (e.g., 8080:80)")
	containerRunCmd.Flags().Bool("detach", true, "Run in background")
	containerRunCmd.Flags().Bool("dry-run", false, "Simulate command")
	
	containerStopCmd.Flags().Bool("dry-run", false, "Simulate command")
	
	containerRemoveCmd.Flags().Bool("force", false, "Force removal")
	containerRemoveCmd.Flags().Bool("dry-run", false, "Simulate command")
	
	containerCmd.AddCommand(containerListCmd)
	containerCmd.AddCommand(containerRunCmd)
	containerCmd.AddCommand(containerStopCmd)
	containerCmd.AddCommand(containerRemoveCmd)
	rootCmd.AddCommand(containerCmd)
}
