package commands

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/rsdenck/nux/internal/output"
	"github.com/spf13/cobra"
)

var networkCmd = &cobra.Command{
	Use:   "network",
	Short: "Network management",
	Long:  `Manage network configuration, interfaces, routes, and diagnostics.`,
}

var networkListCmd = &cobra.Command{
	Use:   "list",
	Short: "List network interfaces",
	Run: func(cmd *cobra.Command, args []string) {
		// Use ip -j link show for JSON output
		ipCmd := exec.Command("ip", "-j", "link", "show")
		out, err := ipCmd.CombinedOutput()
		
		if err != nil {
			output.NewError(fmt.Sprintf("failed to list interfaces: %s", strings.TrimSpace(string(out))), "NETWORK_LIST_ERROR").Print()
			return
		}
		
		var interfaces []map[string]interface{}
		if err := json.Unmarshal(out, &interfaces); err != nil {
			// Fallback to text parsing if JSON fails
			output.NewList([]map[string]interface{}{
				{"output": strings.TrimSpace(string(out))},
			}, 1).WithMessage("Network interfaces (text)").Print()
			return
		}
		
		// Get addresses for each interface
		addrCmd := exec.Command("ip", "-j", "addr", "show")
		addrOut, _ := addrCmd.CombinedOutput()
		
		var addresses []map[string]interface{}
		json.Unmarshal(addrOut, &addresses)
		
		// Combine interface and address info
		items := []map[string]interface{}{}
		for _, iface := range interfaces {
			item := map[string]interface{}{
				"ifindex":  iface["ifindex"],
				"ifname":   iface["ifname"],
				"flags":    iface["flags"],
				"mtu":      iface["mtu"],
				"operstate": iface["operstate"],
				"address":  iface["address"],
			}
			
			// Find matching address
			for _, addr := range addresses {
				if addr["ifname"] == iface["ifname"] {
					item["addr_info"] = addr["addr_info"]
					break
				}
			}
			
			items = append(items, item)
		}
		
		output.NewList(items, len(items)).WithMessage("Network interfaces").Print()
	},
}

var networkShowCmd = &cobra.Command{
	Use:   "show [interface]",
	Short: "Show interface details",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cmdArgs := []string{"-j", "addr", "show"}
		if len(args) > 0 {
			cmdArgs = append(cmdArgs, args[0])
		}
		
		ipCmd := exec.Command("ip", cmdArgs...)
		out, err := ipCmd.CombinedOutput()
		
		if err != nil {
			output.NewError(fmt.Sprintf("failed to show interface: %s", strings.TrimSpace(string(out))), "NETWORK_SHOW_ERROR").Print()
			return
		}
		
		var result interface{}
		if err := json.Unmarshal(out, &result); err != nil {
			output.NewSuccess(map[string]interface{}{
				"output": strings.TrimSpace(string(out)),
			}).Print()
			return
		}
		
		output.NewSuccess(result).Print()
	},
}

var networkSetCmd = &cobra.Command{
	Use:   "set <interface> [flags]",
	Short: "Configure network interface",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		iface := args[0]
		
		ipAddr, _ := cmd.Flags().GetString("ip")
		netmask, _ := cmd.Flags().GetString("netmask")
		gateway, _ := cmd.Flags().GetString("gateway")
		up, _ := cmd.Flags().GetBool("up")
		down, _ := cmd.Flags().GetBool("down")
		
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		
		commands := []string{}
		
		if up {
			commands = append(commands, fmt.Sprintf("ip link set %s up", iface))
		}
		if down {
			commands = append(commands, fmt.Sprintf("ip link set %s down", iface))
		}
		if ipAddr != "" {
			if netmask != "" {
				commands = append(commands, fmt.Sprintf("ip addr add %s/%s dev %s", ipAddr, netmask, iface))
			} else {
				commands = append(commands, fmt.Sprintf("ip addr add %s dev %s", ipAddr, iface))
			}
		}
		if gateway != "" {
			commands = append(commands, fmt.Sprintf("ip route add default via %s", gateway))
		}
		
		if dryRun {
			output.NewInfo(map[string]interface{}{
				"interface": iface,
				"dry_run":  true,
				"commands": commands,
			}).Print()
			return
		}
		
		for _, command := range commands {
			parts := strings.Fields(command)
			execCmd := exec.Command(parts[0], parts[1:]...)
			out, err := execCmd.CombinedOutput()
			
			if err != nil {
				output.NewError(fmt.Sprintf("command failed: %s - %s", command, strings.TrimSpace(string(out))), "NETWORK_SET_ERROR").Print()
				return
			}
		}
		
		output.NewSuccess(map[string]interface{}{
			"interface": iface,
			"status":    "configured",
			"commands":  commands,
		}).Print()
	},
}

func init() {
	networkCmd.AddCommand(networkListCmd)
	networkCmd.AddCommand(networkShowCmd)
	networkSetCmd.Flags().String("ip", "", "IP address (CIDR or with --netmask)")
	networkSetCmd.Flags().String("netmask", "", "Netmask (optional)")
	networkSetCmd.Flags().String("gateway", "", "Default gateway")
	networkSetCmd.Flags().Bool("up", false, "Bring interface up")
	networkSetCmd.Flags().Bool("down", false, "Bring interface down")
	networkCmd.AddCommand(networkSetCmd)
	rootCmd.AddCommand(networkCmd)
}
