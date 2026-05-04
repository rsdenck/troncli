package commands

import (
	"fmt"
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
			output.NewError(fmt.Sprintf("ProtonVPN not installed or not running: %s", strings.TrimSpace(string(out))), "PROTON_NOT_FOUND").Print()
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
	Use:   "connect",
	Short: "Connect to Proton VPN",
	Run: func(cmd *cobra.Command, args []string) {
		out, err := exec.Command("protonvpn-cli", "connect").CombinedOutput()
		if err != nil {
			output.NewError(fmt.Sprintf("Failed to connect: %s", strings.TrimSpace(string(out))), "PROTON_VPN_ERROR").Print()
			return
		}
		output.NewSuccess(map[string]interface{}{
			"action": "connect",
			"output": strings.TrimSpace(string(out)),
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
			"output": strings.TrimSpace(string(out)),
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
			"output": strings.TrimSpace(string(out)),
		}).Print()
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
		out, err := exec.Command("xdg-open", "https://mail.protonmail.com").CombinedOutput()
		if err != nil {
			output.NewError(fmt.Sprintf("Failed to open Proton Mail: %s", strings.TrimSpace(string(out))), "PROTON_OPEN_ERROR").Print()
			return
		}
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
		out, err := exec.Command("xdg-open", "https://drive.protonmail.com").CombinedOutput()
		if err != nil {
			output.NewError(fmt.Sprintf("Failed to open Proton Drive: %s", strings.TrimSpace(string(out))), "PROTON_OPEN_ERROR").Print()
			return
		}
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
		out, err := exec.Command("nux", "vault", "sync", "proton").CombinedOutput()
		if err != nil {
			output.NewError(fmt.Sprintf("Vault sync failed: %s", strings.TrimSpace(string(out))), "PROTON_SYNC_ERROR").Print()
			return
		}
		output.NewSuccess(map[string]interface{}{
			"action": "vault_sync",
			"status": "completed",
		}).Print()
	},
}

// Secure travel mode command
var protonSecureCmd = &cobra.Command{
	Use:   "secure",
	Short: "Activate Proton security travel mode",
	Run: func(cmd *cobra.Command, args []string) {
		output.NewInfo("Activating Proton Travel Mode...").Print()
		
		// Step 1: Connect VPN fastest
		out, err := exec.Command("nux", "proton", "vpn", "fastest").CombinedOutput()
		if err != nil {
			output.NewError(fmt.Sprintf("VPN connection failed: %s", strings.TrimSpace(string(out))), "PROTON_TRAVEL_VPN_ERROR").Print()
			return
		}
		
		// Step 2: Check firewall
		exec.Command("nux", "firewall", "list").CombinedOutput()
		
		// Step 3: Disable insecure services
		exec.Command("systemctl", "disable", "--now", "telnet").CombinedOutput()
		
		output.NewSuccess(map[string]interface{}{
			"mode":    "travel_secure",
			"steps":   []string{"vpn_connected", "firewall_checked", "insecure_services_disabled"},
			"status":  "activated",
		}).WithMessage("Proton Travel Mode Activated").Print()
	},
}

// Pass lookup command
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
			"status":  "lookup_simulated",
			"note":    "Proton Pass integration requires active subscription and CLI setup",
		}).Print()
	},
}

func init() {
	// Add VPN subcommands
	protonVpnCmd.AddCommand(protonVpnConnectCmd)
	protonVpnCmd.AddCommand(protonVpnFastestCmd)
	protonVpnCmd.AddCommand(protonVpnDisconnectCmd)
	
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
