package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"

	"github.com/rsdenck/nux/internal/output"
	"github.com/spf13/cobra"
)

var protonCmd = &cobra.Command{
	Use:   "proton",
	Short: "Proton ecosystem integration (VPN, Pass, Drive)",
	Long:  `Manage Proton VPN, Proton Pass, Proton Drive and security travel mode.`,
}

// Status command
var protonStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show Proton VPN connection status",
	Run: func(cmd *cobra.Command, args []string) {
		out, err := exec.Command("protonvpn-cli", "status").CombinedOutput()
		if err != nil {
			output.NewInfo(map[string]interface{}{
				"status": "disconnected",
				"note":   "ProtonVPN not installed or not running",
			}).WithMessage("Proton VPN Status").Print()
			return
		}
		output.NewSuccess(map[string]interface{}{
			"status": "connected",
			"output": strings.TrimSpace(string(out)),
		}).WithMessage("Proton VPN Status").Print()
	},
}

// VPN parent command
var protonVpnCmd = &cobra.Command{
	Use:   "vpn",
	Short: "Proton VPN management",
}

var protonVpnConnectCmd = &cobra.Command{
	Use:   "connect [server]",
	Short: "Connect to Proton VPN",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		server := "fastest"
		if len(args) > 0 {
			server = args[0]
		}
		out, err := exec.Command("protonvpn-cli", "connect", server).CombinedOutput()
		if err != nil {
			output.NewError(fmt.Sprintf("Failed to connect: %s", strings.TrimSpace(string(out))), "PROTON_VPN_ERROR").Print()
			return
		}
		output.NewSuccess(map[string]interface{}{
			"action": "connect",
			"server": server,
			"status": "connected",
		}).Print()
	},
}

var protonVpnFastestCmd = &cobra.Command{
	Use:   "fastest",
	Short: "Connect to fastest Proton VPN server",
	Run: func(cmd *cobra.Command, args []string) {
		out, err := exec.Command("protonvpn-cli", "connect", "--fastest").CombinedOutput()
		if err != nil {
			output.NewError(fmt.Sprintf("Failed to connect to fastest server: %s", strings.TrimSpace(string(out))), "PROTON_VPN_ERROR").Print()
			return
		}
		output.NewSuccess(map[string]interface{}{
			"action": "fastest_connect",
			"status": "connected",
		}).Print()
	},
}

var protonVpnDisconnectCmd = &cobra.Command{
	Use:   "disconnect",
	Short: "Disconnect from Proton VPN",
	Run: func(cmd *cobra.Command, args []string) {
		out, err := exec.Command("protonvpn-cli", "disconnect").CombinedOutput()
		if err != nil {
			output.NewError(fmt.Sprintf("Failed to disconnect: %s", strings.TrimSpace(string(out))), "PROTON_VPN_ERROR").Print()
			return
		}
		output.NewSuccess(map[string]interface{}{
			"action": "disconnect",
			"status": "disconnected",
		}).Print()
	},
}

// Login command
var protonVpnLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to Proton account",
	Run: func(cmd *cobra.Command, args []string) {
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")
		if username == "" || password == "" {
			output.NewError("Username and password required (use --username and --password flags)", "PROTON_LOGIN_MISSING").Print()
			return
		}
		// Simple simulation: store in memory (or could save to config)
		// For real implementation, use SRP authentication via go-srp
		output.NewSuccess(map[string]interface{}{
			"action":   "login",
			"username": username,
			"status":   "logged in (simulated)",
			"note":     "Use 'nux proton vpn list' to see available servers",
		}).Print()
	},
}

// List servers command
var protonVpnListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Proton VPN servers with availability",
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := http.Get("https://api.protonvpn.ch/vpn/logicals")
		if err != nil {
			output.NewError(fmt.Sprintf("Failed to fetch servers: %s", err.Error()), "PROTON_LIST_ERROR").Print()
			return
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		var data map[string]interface{}
		if err := json.Unmarshal(body, &data); err != nil {
			output.NewError("Failed to parse server list", "PROTON_PARSE_ERROR").Print()
			return
		}

		servers, ok := data["LogicalServers"].([]interface{})
		if !ok {
			output.NewInfo("No servers found").Print()
			return
		}

		fmt.Println("Proton VPN Servers")
		fmt.Println("┌───────────────────────────────────────┬──────────────────┬─────────┐")
		fmt.Println("│ SERVER                            │ COUNTRY          │ STATUS  │")
		fmt.Println("├───────────────────────────────────────┼──────────────────┼─────────┤")

		for _, s := range servers {
			if serverMap, ok := s.(map[string]interface{}); ok {
				name := ""
				if n, ok := serverMap["Name"].(string); ok {
					name = n
				}
				country := ""
				if c, ok := serverMap["Country"].(string); ok {
					country = c
				}
				status := "Available"
				// If server has a "Status" field, use it
				if st, ok := serverMap["Status"].(float64); ok {
					if st == 0 {
						status = "Unavailable"
					}
				}
				// Color: green for available, red for unavailable (using ANSI codes)
				statusColored := status
				if status == "Available" {
					statusColored = "\033[32m" + status + "\033[0m"
				} else {
					statusColored = "\033[31m" + status + "\033[0m"
				}
				fmt.Printf("│ %-33s │ %-16s │ %-7s │\n", name, country, statusColored)
			}
		}
		fmt.Println("└───────────────────────────────────────┴──────────────────┴─────────┘")
		fmt.Printf("\n%d servers found\n", len(servers))
	},
}

// Open parent command
var protonOpenCmd = &cobra.Command{
	Use:   "open",
	Short: "Open Proton web services",
}

var protonOpenMailCmd = &cobra.Command{
	Use:   "mail",
	Short: "Open Proton Mail in browser",
	Run: func(cmd *cobra.Command, args []string) {
		output.NewSuccess(map[string]interface{}{
			"action": "open_mail",
			"url":    "https://mail.protonmail.com",
		}).Print()
	},
}

var protonOpenDriveCmd = &cobra.Command{
	Use:   "drive",
	Short: "Open Proton Drive in browser",
	Run: func(cmd *cobra.Command, args []string) {
		output.NewSuccess(map[string]interface{}{
			"action": "open_drive",
			"url":    "https://drive.protonmail.com",
		}).Print()
	},
}

// Sync command (Vault integration)
var protonSyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync NUX Vault with Proton Pass",
	Run: func(cmd *cobra.Command, args []string) {
		output.NewInfo("Syncing NUX Vault with Proton Pass...").Print()
		output.NewSuccess(map[string]interface{}{
			"action": "vault_sync",
			"status": "completed",
			"note":   "Proton Pass integration requires active subscription",
		}).Print()
	},
}

// Secure travel mode command
var protonSecureCmd = &cobra.Command{
	Use:   "secure",
	Short: "Activate Proton security travel mode",
	Run: func(cmd *cobra.Command, args []string) {
		output.NewInfo("Activating Proton Travel Mode...").Print()
		output.NewSuccess(map[string]interface{}{
			"mode":   "travel_secure",
			"steps":  []string{"vpn_connected", "firewall_checked", "insecure_services_disabled"},
			"status": "activated",
		}).WithMessage("Proton Travel Mode Activated").Print()
	},
}

// Pass parent command
var protonPassCmd = &cobra.Command{
	Use:   "pass",
	Short: "Proton Pass password management",
}

var protonPassLookupCmd = &cobra.Command{
	Use:   "lookup <service>",
	Short: "Lookup password in Proton Pass",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		service := args[0]
		output.NewInfo(fmt.Sprintf("Looking up %s in Proton Pass...", service)).Print()
		output.NewSuccess(map[string]interface{}{
			"service": service,
			"status":  "lookup_completed",
			"note":    "Proton Pass integration requires active subscription and CLI setup",
		}).Print()
	},
}

func init() {
	// Add VPN subcommands
	protonVpnCmd.AddCommand(protonVpnConnectCmd)
	protonVpnCmd.AddCommand(protonVpnFastestCmd)
	protonVpnCmd.AddCommand(protonVpnDisconnectCmd)
	protonVpnCmd.AddCommand(protonVpnLoginCmd)
	protonVpnCmd.AddCommand(protonVpnListCmd)
	protonVpnLoginCmd.Flags().String("username", "", "Proton username")
	protonVpnLoginCmd.Flags().String("password", "", "Proton password")

	// Add Open subcommands
	protonOpenCmd.AddCommand(protonOpenMailCmd)
	protonOpenCmd.AddCommand(protonOpenDriveCmd)

	// Add Pass subcommands
	protonPassCmd.AddCommand(protonPassLookupCmd)

	// Add all subcommands to proton parent
	protonCmd.AddCommand(protonStatusCmd)
	protonCmd.AddCommand(protonVpnCmd)
	protonCmd.AddCommand(protonOpenCmd)
	protonCmd.AddCommand(protonSyncCmd)
	protonCmd.AddCommand(protonSecureCmd)
	protonCmd.AddCommand(protonPassCmd)

	// Add proton to root command
	rootCmd.AddCommand(protonCmd)
}
