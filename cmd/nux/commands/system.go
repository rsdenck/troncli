package commands

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/rsdenck/nux/internal/core"
	"github.com/rsdenck/nux/internal/output"
	"github.com/spf13/cobra"
)

var systemExecutor core.Executor = &core.RealExecutor{}

var systemCmd = &cobra.Command{
	Use:   "system",
	Short: "System information",
	Long:  `Show system information and statistics.`,
}

var systemInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show system information",
	Run: func(cmd *cobra.Command, args []string) {
		info := map[string]interface{}{
			"hostname": getHostname(),
			"kernel":   getKernelVersion(),
			"arch":     runtime.GOARCH,
			"os":       getOSInfo(),
			"uptime":   getUptime(),
			"load":     getLoadAverage(),
		}

		output.NewSuccess(info).Print()
	},
}

var systemUptimeCmd = &cobra.Command{
	Use:   "uptime",
	Short: "Show system uptime",
	Run: func(cmd *cobra.Command, args []string) {
		out, err := systemExecutor.CombinedOutput("uptime")

		if err != nil {
			output.NewError(fmt.Sprintf("failed to get uptime: %s", err.Error()), "SYSTEM_UPTIME_ERROR").Print()
			return
		}

		output.NewSuccess(map[string]interface{}{
			"uptime": out,
		}).Print()
	},
}

var systemLoadCmd = &cobra.Command{
	Use:   "load",
	Short: "Show load average",
	Run: func(cmd *cobra.Command, args []string) {
		load := getLoadAverage()
		output.NewSuccess(map[string]interface{}{
			"load_average": load,
		}).Print()
	},
}

func getHostname() string {
	hostname, _ := os.Hostname()
	return hostname
}

func getKernelVersion() string {
	out, err := systemExecutor.CombinedOutput("uname", "-r")
	if err != nil {
		return "unknown"
	}
	return out
}

func getOSInfo() string {
	if file, err := os.Open("/etc/os-release"); err == nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "PRETTY_NAME=") {
				return strings.Trim(strings.TrimPrefix(line, "PRETTY_NAME="), "\"")
			}
		}
	}

	out, err := systemExecutor.CombinedOutput("lsb_release", "-ds")
	if err == nil {
		return strings.Trim(out, "\"")
	}

	return "unknown"
}

func getUptime() string {
	out, err := systemExecutor.CombinedOutput("uptime", "-p")
	if err != nil {
		return "unknown"
	}
	return out
}

func getLoadAverage() string {
	out, err := systemExecutor.CombinedOutput("cat", "/proc/loadavg")
	if err != nil {
		return "unknown"
	}
	return out
}

func init() {
	systemCmd.AddCommand(systemInfoCmd)
	systemCmd.AddCommand(systemUptimeCmd)
	systemCmd.AddCommand(systemLoadCmd)
	rootCmd.AddCommand(systemCmd)
}
