package commands

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/rsdenck/nux/internal/output"
	"github.com/spf13/cobra"
)

var firewallCmd = &cobra.Command{
	Use:   "firewall",
	Short: "Firewall management",
	Long:  `Manage firewall rules (nftables, iptables, firewalld, ufw).`,
}

// detectFirewall detects the firewall tool available
func detectFirewall() string {
	tools := []struct {
		name    string
		command string
	}{
		{"nft", "nft"},
		{"iptables", "iptables"},
		{"firewalld", "firewall-cmd"},
		{"ufw", "ufw"},
	}

	for _, t := range tools {
		if _, err := exec.LookPath(t.command); err == nil {
			return t.name
		}
	}

	return "unknown"
}

var firewallListCmd = &cobra.Command{
	Use:   "list",
	Short: "List firewall rules",
	Run: func(cmd *cobra.Command, args []string) {
		fw := detectFirewall()

		var command string
		var cmdArgs []string

		switch fw {
		case "nft":
			command = "nft"
			cmdArgs = []string{"list", "ruleset"}
		case "iptables":
			command = "iptables"
			cmdArgs = []string{"-L", "-n", "-v"}
		case "firewalld":
			command = "firewall-cmd"
			cmdArgs = []string{"--list-all"}
		case "ufw":
			command = "ufw"
			cmdArgs = []string{"status", "verbose"}
		default:
			output.NewError("no supported firewall found", "FIREWALL_UNSUPPORTED").Print()
			return
		}

		fwCmd := exec.Command(command, cmdArgs...)
		out, err := fwCmd.CombinedOutput()

		if err != nil {
			output.NewError(fmt.Sprintf("failed to list rules: %s - %s", err.Error(), strings.TrimSpace(string(out))), "FIREWALL_LIST_ERROR").Print()
			return
		}

		output.NewSuccess(map[string]interface{}{
			"firewall": fw,
			"rules":    strings.TrimSpace(string(out)),
		}).Print()
	},
}

var firewallAddCmd = &cobra.Command{
	Use:   "add [flags]",
	Short: "Add firewall rule",
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetString("port")
		protocol, _ := cmd.Flags().GetString("protocol")
		action, _ := cmd.Flags().GetString("action")

		if port == "" {
			output.NewError("port is required", "FIREWALL_PORT_REQUIRED").Print()
			return
		}

		fw := detectFirewall()
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		var command string
		var cmdArgs []string

		switch fw {
		case "nft":
			command = "nft"
			cmdArgs = []string{"add", "rule", "inet", "filter", "input", "tcp", "dport", port, action}
		case "iptables":
			command = "iptables"
			if action == "allow" {
				action = "ACCEPT"
			} else {
				action = "DROP"
			}
			cmdArgs = []string{"-A", "INPUT", "-p", protocol, "--dport", port, "-j", action}
		case "firewalld":
			command = "firewall-cmd"
			cmdArgs = []string{"--add-port=" + port + "/" + protocol}
		case "ufw":
			command = "ufw"
			cmdArgs = []string{action, port + "/" + protocol}
		default:
			output.NewError("no supported firewall found", "FIREWALL_UNSUPPORTED").Print()
			return
		}

		if dryRun {
			output.NewInfo(map[string]interface{}{
				"firewall": fw,
				"dry_run":  true,
				"command":  fmt.Sprintf("%s %s", command, strings.Join(cmdArgs, " ")),
			}).Print()
			return
		}

		fwCmd := exec.Command(command, cmdArgs...)
		out, err := fwCmd.CombinedOutput()

		if err != nil {
			output.NewError(fmt.Sprintf("failed to add rule: %s - %s", err.Error(), strings.TrimSpace(string(out))), "FIREWALL_ADD_ERROR").Print()
			return
		}

		output.NewSuccess(map[string]interface{}{
			"firewall": fw,
			"port":     port,
			"protocol": protocol,
			"action":   action,
			"status":   "added",
		}).Print()
	},
}

var firewallRemoveCmd = &cobra.Command{
	Use:   "remove [flags]",
	Short: "Remove firewall rule",
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetString("port")
		protocol, _ := cmd.Flags().GetString("protocol")

		if port == "" {
			output.NewError("port is required", "FIREWALL_PORT_REQUIRED").Print()
			return
		}

		fw := detectFirewall()
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		var command string
		var cmdArgs []string

		switch fw {
		case "iptables":
			command = "iptables"
			cmdArgs = []string{"-D", "INPUT", "-p", protocol, "--dport", port, "-j", "ACCEPT"}
		case "firewalld":
			command = "firewall-cmd"
			cmdArgs = []string{"--remove-port=" + port + "/" + protocol}
		case "ufw":
			command = "ufw"
			cmdArgs = []string{"delete", "allow", port + "/" + protocol}
		default:
			output.NewError("firewall does not support rule removal or is unsupported", "FIREWALL_REMOVE_UNSUPPORTED").Print()
			return
		}

		if dryRun {
			output.NewInfo(map[string]interface{}{
				"firewall": fw,
				"dry_run":  true,
				"command":  fmt.Sprintf("%s %s", command, strings.Join(cmdArgs, " ")),
			}).Print()
			return
		}

		fwCmd := exec.Command(command, cmdArgs...)
		out, err := fwCmd.CombinedOutput()

		if err != nil {
			output.NewError(fmt.Sprintf("failed to remove rule: %s - %s", err.Error(), strings.TrimSpace(string(out))), "FIREWALL_REMOVE_ERROR").Print()
			return
		}

		output.NewSuccess(map[string]interface{}{
			"firewall": fw,
			"port":     port,
			"protocol": protocol,
			"status":   "removed",
		}).Print()
	},
}

func init() {
	firewallCmd.AddCommand(firewallListCmd)
	firewallAddCmd.Flags().String("port", "", "Port number")
	firewallAddCmd.Flags().String("protocol", "tcp", "Protocol (tcp/udp)")
	firewallAddCmd.Flags().String("action", "allow", "Action (allow/deny)")
	firewallCmd.AddCommand(firewallAddCmd)
	firewallRemoveCmd.Flags().String("port", "", "Port number")
	firewallRemoveCmd.Flags().String("protocol", "tcp", "Protocol (tcp/udp)")
	firewallCmd.AddCommand(firewallRemoveCmd)
	rootCmd.AddCommand(firewallCmd)
}
