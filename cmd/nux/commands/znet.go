package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/rsdenck/nux/internal/output"
	"github.com/spf13/cobra"
)

var znetCmd = &cobra.Command{
	Use:   "znet",
	Short: "ZeroNet P2P network management",
	Long:  `Manage ZeroNet decentralized P2P networks. Connect, disconnect, list peers, and manage ZeroNet sites.`,
}

var znetInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install ZeroNet",
	Long:  `Download and install the latest ZeroNet from the official GitHub repository.`,
	Run: func(cmd *cobra.Command, args []string) {
		if isZeroNetInstalled() {
			output.NewInfo("ZeroNet is already installed").Print()
			return
		}
		output.NewInfo("Installing ZeroNet...").Print()

		home, _ := os.UserHomeDir()
		znetDir := filepath.Join(home, ".nux", "znet")
		os.MkdirAll(znetDir, 0755)

		arch := runtime.GOARCH
		dist := "linux64"
		if arch == "arm64" || arch == "aarch64" {
			dist = "linux arm64"
		} else if arch == "386" {
			dist = "linux32"
		}

		url := fmt.Sprintf("https://github.com/HelloZeroNet/ZeroNet-%s/archive/refs/heads/master.tar.gz", dist)
		output.NewInfo(fmt.Sprintf("Downloading ZeroNet for %s...", dist)).Print()

		cmdTar := exec.Command("bash", "-c", fmt.Sprintf(
			"curl -sL '%s' | tar -xzf - -C '%s' --strip-components=1", url, znetDir))
		if out, err := cmdTar.CombinedOutput(); err != nil {
			output.NewError(fmt.Sprintf("Failed to install ZeroNet: %s\n%s", err.Error(), string(out)), "ZNET_INSTALL_ERR").Print()
			return
		}

		output.NewInfo("ZeroNet installed successfully").Print()
		output.NewInfo(fmt.Sprintf("Location: %s", znetDir)).Print()
		output.NewInfo("Run 'nux znet start' to launch ZeroNet").Print()
	},
}

var znetStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start ZeroNet daemon",
	Long:  `Start the ZeroNet background process on port 43110.`,
	Run: func(cmd *cobra.Command, args []string) {
		if isZeroNetRunning() {
			output.NewInfo("ZeroNet is already running").Print()
			return
		}
		if !isZeroNetInstalled() {
			output.NewError("ZeroNet is not installed. Run 'nux znet install' first", "ZNET_NOT_FOUND").Print()
			return
		}
		home, _ := os.UserHomeDir()
		znetDir := filepath.Join(home, ".nux", "znet")
		znetPy := filepath.Join(znetDir, "zeronet.py")

		cmdStart := exec.Command("python3", znetPy)
		cmdStart.Dir = znetDir
		if err := cmdStart.Start(); err != nil {
			output.NewError(fmt.Sprintf("Failed to start ZeroNet: %s", err.Error()), "ZNET_START_ERR").Print()
			return
		}
		pidFile := filepath.Join(znetDir, ".pid")
		os.WriteFile(pidFile, []byte(fmt.Sprintf("%d", cmdStart.Process.Pid)), 0644)
		output.NewInfo(fmt.Sprintf("ZeroNet started (PID: %d)", cmdStart.Process.Pid)).Print()
		output.NewInfo("Web UI: http://127.0.0.1:43110").Print()
	},
}

var znetStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop ZeroNet daemon",
	Run: func(cmd *cobra.Command, args []string) {
		if !isZeroNetRunning() {
			output.NewInfo("ZeroNet is not running").Print()
			return
		}
		home, _ := os.UserHomeDir()
		pidFile := filepath.Join(home, ".nux", "znet", ".pid")
		pidBytes, err := os.ReadFile(pidFile)
		if err == nil {
			pid := strings.TrimSpace(string(pidBytes))
			exec.Command("kill", pid).Run()
			os.Remove(pidFile)
			output.NewInfo(fmt.Sprintf("ZeroNet (PID %s) stopped", pid)).Print()
			return
		}
		exec.Command("pkill", "-f", "zeronet.py").Run()
		output.NewInfo("ZeroNet stopped").Print()
	},
}

var znetStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show ZeroNet daemon status",
	Run: func(cmd *cobra.Command, args []string) {
		running := isZeroNetRunning()
		installed := isZeroNetInstalled()

		headers := []string{"KEY", "VALUE"}
		var rows [][]string
		if installed {
			rows = append(rows, []string{"installed", "yes"})
			home, _ := os.UserHomeDir()
			znetDir := filepath.Join(home, ".nux", "znet")
			rows = append(rows, []string{"path", znetDir})
		} else {
			rows = append(rows, []string{"installed", "no"})
		}
		if running {
			rows = append(rows, []string{"running", "yes"})
			rows = append(rows, []string{"web_ui", "http://127.0.0.1:43110"})
		} else {
			rows = append(rows, []string{"running", "no"})
		}
		output.PrintCompactTable(headers, rows)
	},
}

var znetListCmd = &cobra.Command{
	Use:   "list",
	Short: "List ZeroNet sites and connected peers",
	Run: func(cmd *cobra.Command, args []string) {
		if !isZeroNetRunning() {
			output.NewInfo("ZeroNet is not running. Start it with 'nux znet start'").Print()
			return
		}
		sites, err := fetchZeroNetSites()
		if err != nil {
			output.NewError(fmt.Sprintf("Failed to query ZeroNet: %s", err.Error()), "ZNET_QUERY_ERR").Print()
			return
		}
		if len(sites) == 0 {
			output.NewInfo("No sites found in ZeroNet").Print()
			return
		}
		headers := []string{"ADDRESS", "PEERS", "SIZE", "CONTENT UPDATED"}
		var rows [][]string
		for _, s := range sites {
			rows = append(rows, []string{
				s.Address,
				fmt.Sprintf("%d", s.Peers),
				fmt.Sprintf("%.1f MB", float64(s.Size)/1048576),
				s.ContentUpdated,
			})
		}
		output.PrintCompactTable(headers, rows)
	},
}

var znetConnectCmd = &cobra.Command{
	Use:   "connect [site-address]",
	Short: "Connect to a ZeroNet site",
	Long:  `Open a specific ZeroNet site. If no address is given, opens the ZeroNet Web UI.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if !isZeroNetRunning() {
			output.NewInfo("Starting ZeroNet first...")
			znetStartCmd.Run(cmd, []string{})
		}
		if len(args) == 0 {
			output.NewInfo("ZeroNet Web UI: http://127.0.0.1:43110").Print()
			exec.Command("xdg-open", "http://127.0.0.1:43110").Start()
			return
		}
		addr := args[0]
		url := fmt.Sprintf("http://127.0.0.1:43110/%s", addr)
		output.NewInfo(fmt.Sprintf("Opening ZeroNet site: %s", url)).Print()
		exec.Command("xdg-open", url).Start()
	},
}

var znetDisconnectCmd = &cobra.Command{
	Use:   "disconnect",
	Short: "Disconnect from ZeroNet (stop daemon)",
	Run: func(cmd *cobra.Command, args []string) {
		znetStopCmd.Run(cmd, args)
	},
}

var znetPeersCmd = &cobra.Command{
	Use:   "peers",
	Short: "Show connected peers count",
	Run: func(cmd *cobra.Command, args []string) {
		if !isZeroNetRunning() {
			output.NewInfo("ZeroNet is not running").Print()
			return
		}
		sites, err := fetchZeroNetSites()
		if err != nil {
			output.NewError(fmt.Sprintf("Failed to query peers: %s", err.Error()), "ZNET_PEERS_ERR").Print()
			return
		}
		totalPeers := 0
		for _, s := range sites {
			totalPeers += s.Peers
		}
		output.NewInfo(fmt.Sprintf("Total peers: %d across %d sites", totalPeers, len(sites))).Print()
	},
}

type znetSite struct {
	Address        string `json:"address"`
	Peers          int    `json:"peers"`
	Size           int    `json:"size"`
	ContentUpdated string `json:"content_updated"`
}

func isZeroNetInstalled() bool {
	home, _ := os.UserHomeDir()
	znetPy := filepath.Join(home, ".nux", "znet", "zeronet.py")
	_, err := os.Stat(znetPy)
	return err == nil
}

func isZeroNetRunning() bool {
	out, err := exec.Command("pgrep", "-f", "zeronet.py").CombinedOutput()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(out)) != ""
}

func fetchZeroNetSites() ([]znetSite, error) {
	resp, err := http.Get("http://127.0.0.1:43110/stats")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var sites []znetSite
	if err := json.Unmarshal(body, &sites); err != nil {
		return nil, err
	}
	return sites, nil
}

func init() {
	znetCmd.AddCommand(znetInstallCmd)
	znetCmd.AddCommand(znetStartCmd)
	znetCmd.AddCommand(znetStopCmd)
	znetCmd.AddCommand(znetStatusCmd)
	znetCmd.AddCommand(znetListCmd)
	znetCmd.AddCommand(znetConnectCmd)
	znetCmd.AddCommand(znetDisconnectCmd)
	znetCmd.AddCommand(znetPeersCmd)
	rootCmd.AddCommand(znetCmd)
}
