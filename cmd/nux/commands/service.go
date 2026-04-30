package commands

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/rsdenck/nux/internal/core"
	"github.com/rsdenck/nux/internal/output"
	"github.com/spf13/cobra"
)

var serviceExecutor core.Executor = &core.RealExecutor{}

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Service management",
	Long:  `Manage system services (systemd, openrc, sysvinit, runit).`,
}

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
		out, err := serviceExecutor.CombinedOutput("systemctl", "list-units", "--type=service", "--all", "--no-pager")
		
		if err != nil {
			output.NewError(fmt.Sprintf("failed to list services: %s", strings.TrimSpace(out)), "SERVICE_LIST_ERROR").Print()
			return
		}

		services := parseServiceOutput(out)

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

		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}

		name := strings.TrimSuffix(fields[0], ".service")
		state := fields[2]
		
		if state == "active" {
			state = "active"
		} else if state == "inactive" || state == "dead" {
			state = "inactive"
		}

		enabled := "no"
		if state == "active" {
			enabled = "yes"
		}

		pid := "-"
		pidOut, _ := serviceExecutor.CombinedOutput("systemctl", "show", name, "--property=MainPID", "--value")
		pidStr := strings.TrimSpace(pidOut)
		if pidStr != "" && pidStr != "0" {
			pid = pidStr
		}

		services = append(services, ServiceInfo{
			Name:    name,
			State:   state,
			Enabled: enabled,
			Pid:     pid,
			Ports:   "-",
		})
	}

	return services
}

var serviceStatusCmd = &cobra.Command{
	Use:   "status <service>",
	Short: "Show service status",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		service := core.SanitizeInput(args[0])
		
		out, _ := serviceExecutor.CombinedOutput("systemctl", "status", service)
		jsonOut, _ := serviceExecutor.CombinedOutput("systemctl", "show", service, "--output=json")

		var jsonData interface{}
		if jsonErr := json.Unmarshal([]byte(jsonOut), &jsonData); jsonErr == nil {
			output.NewSuccess(jsonData).Print()
			return
		}

		output.NewSuccess(map[string]interface{}{
			"service": service,
			"output":  strings.TrimSpace(out),
		}).Print()
	},
}

var serviceStartCmd = &cobra.Command{
	Use:   "start <service>",
	Short: "Start a service",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		service := core.SanitizeInput(args[0])
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

		_, err := serviceExecutor.CombinedOutput("systemctl", "start", service)
		
		if err != nil {
			output.NewError(fmt.Sprintf("failed to start service: %s", err.Error()), "SERVICE_START_ERROR").Print()
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
		service := core.SanitizeInput(args[0])
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

		_, err := serviceExecutor.CombinedOutput("systemctl", "stop", service)
		
		if err != nil {
			output.NewError(fmt.Sprintf("failed to stop service: %s", err.Error()), "SERVICE_STOP_ERROR").Print()
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
		service := core.SanitizeInput(args[0])
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

		_, err := serviceExecutor.CombinedOutput("systemctl", "enable", service)
		
		if err != nil {
			output.NewError(fmt.Sprintf("failed to enable service: %s", err.Error()), "SERVICE_ENABLE_ERROR").Print()
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
		service := core.SanitizeInput(args[0])
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

		_, err := serviceExecutor.CombinedOutput("systemctl", "disable", service)
		
		if err != nil {
			output.NewError(fmt.Sprintf("failed to disable service: %s", err.Error()), "SERVICE_DISABLE_ERROR").Print()
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
	serviceListCmd.Flags().Bool("json", false, "Output in JSON format")
	serviceCmd.AddCommand(serviceListCmd)
	serviceCmd.AddCommand(serviceStatusCmd)
	serviceCmd.AddCommand(serviceStartCmd)
	serviceCmd.AddCommand(serviceStopCmd)
	serviceCmd.AddCommand(serviceEnableCmd)
	serviceCmd.AddCommand(serviceDisableCmd)
	
	serviceStartCmd.Flags().Bool("dry-run", false, "Simulate command")
	serviceStopCmd.Flags().Bool("dry-run", false, "Simulate command")
	serviceEnableCmd.Flags().Bool("dry-run", false, "Simulate command")
	serviceDisableCmd.Flags().Bool("dry-run", false, "Simulate command")
	
	rootCmd.AddCommand(serviceCmd)
}
