package commands

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/rsdenck/nux/internal/output"
	"github.com/spf13/cobra"
)

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Service management",
	Long:  `Manage system services (systemd, openrc, sysvinit, runit).`,
}

// ServiceInfo represents a service entry matching output.md format
type ServiceInfo struct {
	Name    string
	State   string
	Enabled string
	Pid     string
	Ports   string
}

var serviceListCmd = &cobra.Command{
	Use:   "list",
	Short: "List services",
	Run: func(cmd *cobra.Command, args []string) {
		// Use systemctl to list services with proper parsing
		systemctlCmd := exec.Command("systemctl", "list-units", "--type=service", "--all", "--no-pager")
		out, err := systemctlCmd.CombinedOutput()

		if err != nil {
			output.NewError(fmt.Sprintf("failed to list services: %s", strings.TrimSpace(string(out))), "SERVICE_LIST_ERROR").Print()
			return
		}

		services := parseServiceOutput(string(out))

		// Convert to items for output formatter
		items := []map[string]interface{}{}
		for _, s := range services {
			items = append(items, map[string]interface{}{
				"name":    s.Name,
				"state":   s.State,
				"enabled": s.Enabled,
				"pid":     s.Pid,
				"ports":   s.Ports,
			})
		}

		output.NewList(items, len(items)).WithMessage("Service list").Print()
	},
}

func parseServiceOutput(out string) []ServiceInfo {
	services := []ServiceInfo{}
	lines := strings.Split(out, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "UNIT") || strings.HasPrefix(line, "●") {
			continue
		}

		// Parse systemctl list-units output
		// Format: UNIT LOAD ACTIVE SUB DESCRIPTION
		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}

		name := strings.TrimSuffix(fields[0], ".service")

		// Get state (ACTIVE column)
		state := fields[2]
		if state == "active" {
			state = "active"
		} else if state == "inactive" || state == "dead" {
			state = "inactive"
		}

		// Check if enabled
		enabled := "no"
		if state == "active" {
			enabled = "yes"
		}

		// Try to get PID
		pid := "-"
		pidCmd := exec.Command("systemctl", "show", name, "--property=MainPID", "--value")
		pidOut, _ := pidCmd.CombinedOutput()
		pidStr := strings.TrimSpace(string(pidOut))
		if pidStr != "" && pidStr != "0" {
			pid = pidStr
		}

		// Ports (simplified - would need ss/netstat for real implementation)
		ports := "-"

		services = append(services, ServiceInfo{
			Name:    name,
			State:   state,
			Enabled: enabled,
			Pid:     pid,
			Ports:   ports,
		})
	}

	return services
}

var serviceStatusCmd = &cobra.Command{
	Use:   "status <service>",
	Short: "Show service status",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		service := args[0]

		systemctlCmd := exec.Command("systemctl", "status", service)
		out, _ := systemctlCmd.CombinedOutput()

		// Try to parse as JSON if available
		jsonCmd := exec.Command("systemctl", "show", service, "--output=json")
		jsonOut, _ := jsonCmd.CombinedOutput()

		var jsonData interface{}
		if jsonErr := json.Unmarshal(jsonOut, &jsonData); jsonErr == nil {
			output.NewSuccess(jsonData).Print()
			return
		}

		output.NewSuccess(map[string]interface{}{
			"service": service,
			"output":  strings.TrimSpace(string(out)),
		}).Print()
	},
}

var serviceStartCmd = &cobra.Command{
	Use:   "start <service>",
	Short: "Start a service",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		service := args[0]

		dryRun, _ := cmd.Flags().GetBool("dry-run")

		if dryRun {
			output.NewInfo(map[string]interface{}{
				"service": service,
				"action":  "start",
				"dry_run": true,
				"command": fmt.Sprintf("systemctl start %s", service),
			}).Print()
			return
		}

		systemctlCmd := exec.Command("systemctl", "start", service)
		out, err := systemctlCmd.CombinedOutput()

		if err != nil {
			output.NewError(fmt.Sprintf("failed to start service: %s - %s", err.Error(), strings.TrimSpace(string(out))), "SERVICE_START_ERROR").Print()
			return
		}

		output.NewSuccess(map[string]interface{}{
			"service": service,
			"action":  "start",
			"status":  "started",
		}).Print()
	},
}

var serviceStopCmd = &cobra.Command{
	Use:   "stop <service>",
	Short: "Stop a service",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		service := args[0]

		dryRun, _ := cmd.Flags().GetBool("dry-run")

		if dryRun {
			output.NewInfo(map[string]interface{}{
				"service": service,
				"action":  "stop",
				"dry_run": true,
				"command": fmt.Sprintf("systemctl stop %s", service),
			}).Print()
			return
		}

		systemctlCmd := exec.Command("systemctl", "stop", service)
		out, err := systemctlCmd.CombinedOutput()

		if err != nil {
			output.NewError(fmt.Sprintf("failed to stop service: %s - %s", err.Error(), strings.TrimSpace(string(out))), "SERVICE_STOP_ERROR").Print()
			return
		}

		output.NewSuccess(map[string]interface{}{
			"service": service,
			"action":  "stop",
			"status":  "stopped",
		}).Print()
	},
}

var serviceEnableCmd = &cobra.Command{
	Use:   "enable <service>",
	Short: "Enable service at boot",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		service := args[0]

		dryRun, _ := cmd.Flags().GetBool("dry-run")

		if dryRun {
			output.NewInfo(map[string]interface{}{
				"service": service,
				"action":  "enable",
				"dry_run": true,
				"command": fmt.Sprintf("systemctl enable %s", service),
			}).Print()
			return
		}

		systemctlCmd := exec.Command("systemctl", "enable", service)
		out, err := systemctlCmd.CombinedOutput()

		if err != nil {
			output.NewError(fmt.Sprintf("failed to enable service: %s - %s", err.Error(), strings.TrimSpace(string(out))), "SERVICE_ENABLE_ERROR").Print()
			return
		}

		output.NewSuccess(map[string]interface{}{
			"service": service,
			"action":  "enable",
			"status":  "enabled",
		}).Print()
	},
}

var serviceDisableCmd = &cobra.Command{
	Use:   "disable <service>",
	Short: "Disable service at boot",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		service := args[0]

		dryRun, _ := cmd.Flags().GetBool("dry-run")

		if dryRun {
			output.NewInfo(map[string]interface{}{
				"service": service,
				"action":  "disable",
				"dry_run": true,
				"command": fmt.Sprintf("systemctl disable %s", service),
			}).Print()
			return
		}

		systemctlCmd := exec.Command("systemctl", "disable", service)
		out, err := systemctlCmd.CombinedOutput()

		if err != nil {
			output.NewError(fmt.Sprintf("failed to disable service: %s - %s", err.Error(), strings.TrimSpace(string(out))), "SERVICE_DISABLE_ERROR").Print()
			return
		}

		output.NewSuccess(map[string]interface{}{
			"service": service,
			"action":  "disable",
			"status":  "disabled",
		}).Print()
	},
}

func init() {
	serviceCmd.AddCommand(serviceListCmd)
	serviceCmd.AddCommand(serviceStatusCmd)
	serviceCmd.AddCommand(serviceStartCmd)
	serviceCmd.AddCommand(serviceStopCmd)
	serviceCmd.AddCommand(serviceEnableCmd)
	serviceCmd.AddCommand(serviceDisableCmd)
	rootCmd.AddCommand(serviceCmd)
}
