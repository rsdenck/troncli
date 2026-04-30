package commands

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/rsdenck/nux/internal/output"
	"github.com/spf13/cobra"
)

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
		uptimeCmd := exec.Command("uptime")
		out, err := uptimeCmd.CombinedOutput()
		
		if err != nil {
			output.NewError(fmt.Sprintf("failed to get uptime: %s", strings.TrimSpace(string(out))), "SYSTEM_UPTIME_ERROR").Print()
			return
		}
		
		output.NewSuccess(map[string]interface{}{
			"uptime": strings.TrimSpace(string(out)),
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
	unameCmd := exec.Command("uname", "-r")
	out, err := unameCmd.CombinedOutput()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(out))
}

func getOSInfo() string {
	// Try /etc/os-release first
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
	
	// Fallback to lsb_release
	lsbCmd := exec.Command("lsb_release", "-ds")
	out, err := lsbCmd.CombinedOutput()
	if err == nil {
		return strings.TrimSpace(strings.Trim(string(out), "\""))
	}
	
	return "unknown"
}

func getUptime() string {
	uptimeCmd := exec.Command("uptime", "-p")
	out, err := uptimeCmd.CombinedOutput()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(out))
}

func getLoadAverage() string {
	loadCmd := exec.Command("cat", "/proc/loadavg")
	out, err := loadCmd.CombinedOutput()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(out))
}

func init() {
	systemCmd.AddCommand(systemInfoCmd)
	systemCmd.AddCommand(systemUptimeCmd)
	systemCmd.AddCommand(systemLoadCmd)
	rootCmd.AddCommand(systemCmd)
}
