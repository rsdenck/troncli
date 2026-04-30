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

type NetworkInterface struct {
	Name    string
	Type    string
	State   string
	Address string
	Speed   string
}

var networkListCmd = &cobra.Command{
	Use:   "list",
	Short: "List network interfaces",
	Run: func(cmd *cobra.Command, args []string) {
		interfaces, err := getNetworkInterfaces()
		if err != nil {
			output.NewError(fmt.Sprintf("failed to list interfaces: %s", err.Error()), "NETWORK_LIST_ERROR").Print()
			return
		}

		items := []map[string]interface{}{}
		for _, iface := range interfaces {
			items = append(items, map[string]interface{}{
				"name":    iface.Name,
				"type":    iface.Type,
				"state":   iface.State,
				"address": iface.Address,
				"speed":   iface.Speed,
			})
		}

		output.NewList(items, len(items)).WithMessage("Network interfaces").Print()
	},
}

func getNetworkInterfaces() ([]NetworkInterface, error) {
	linkCmd := exec.Command("ip", "-j", "link", "show")
	linkOut, err := linkCmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("ip link show failed: %s", string(linkOut))
	}

	addrCmd := exec.Command("ip", "-j", "addr", "show")
	addrOut, err := addrCmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("ip addr show failed: %s", string(addrOut))
	}

	var linkData []map[string]interface{}
	if err := json.Unmarshal(linkOut, &linkData); err != nil {
		return nil, err
	}

	var addrData []map[string]interface{}
	if err := json.Unmarshal(addrOut, &addrData); err != nil {
		return nil, err
	}

	linkMap := make(map[string]map[string]interface{})
	for _, link := range linkData {
		if name, ok := link["ifname"].(string); ok {
			linkMap[name] = link
		}
	}

	result := []NetworkInterface{}

	for _, addr := range addrData {
		if name, ok := addr["ifname"].(string); ok {
			linkInfo := linkMap[name]

			ifType := "ethernet"
			if name == "lo" {
				ifType = "loopback"
			} else if strings.HasPrefix(name, "br-") {
				ifType = "bridge"
			} else if strings.HasPrefix(name, "veth") {
				ifType = "veth"
			} else if strings.HasPrefix(name, "wg") {
				ifType = "wireguard"
			}

			state := "unknown"
			if linkInfo != nil {
				if operstate, ok := linkInfo["operstate"].(string); ok {
					state = operstate
				}
			}

			var addresses []string
			if addrInfo, ok := addr["addr_info"].([]interface{}); ok {
				count := 0
				for _, ai := range addrInfo {
					if count >= 2 {
						break
					}
					if addrMap, ok := ai.(map[string]interface{}); ok {
						if family, ok := addrMap["family"].(string); ok && (family == "inet" || family == "inet6") {
							if local, ok := addrMap["local"].(string); ok {
								prefixLen := 0.0
								if pl, ok := addrMap["prefixlen"].(float64); ok {
									prefixLen = pl
								}
								addresses = append(addresses, fmt.Sprintf("%s/%v", local, prefixLen))
								count++
							}
						}
					}
				}
			}

			speed := "-"
			if state == "up" && name != "lo" && !strings.HasPrefix(name, "veth") && !strings.HasPrefix(name, "br-") {
				speedCmd := exec.Command("ethtool", name)
				speedOut, _ := speedCmd.CombinedOutput()
				if strings.Contains(string(speedOut), "Speed:") {
					lines := strings.Split(string(speedOut), "\n")
					for _, line := range lines {
						if strings.Contains(line, "Speed:") {
							parts := strings.Fields(line)
							if len(parts) >= 2 {
								speed = strings.Trim(parts[1], ":")
								break
							}
						}
					}
				}
			}

			result = append(result, NetworkInterface{
				Name:    name,
				Type:    ifType,
				State:   state,
				Address: strings.Join(addresses, ", "),
				Speed:   speed,
			})
		}
	}

	return result, nil
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
		if err := json.Unmarshal(out, &result); err == nil {
			output.NewSuccess(result).Print()
			return
		}

		output.NewSuccess(map[string]interface{}{
			"output": strings.TrimSpace(string(out)),
		}).Print()
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
