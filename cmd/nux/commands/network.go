package commands

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/rsdenck/nux/internal/core"
	"github.com/rsdenck/nux/internal/output"
	"github.com/spf13/cobra"
)

var networkExecutor core.Executor = &core.RealExecutor{}

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

		// Calculate dynamic widths
		nameWidth := 4 // "NAME"
		typeWidth := 4 // "TYPE"
		stateWidth := 5 // "STATE"
		addrWidth := 7  // "ADDRESS"
		speedWidth := 5 // "SPEED"

		for _, iface := range interfaces {
			if len(iface.Name) > nameWidth {
				nameWidth = len(iface.Name)
			}
			if len(iface.Type) > typeWidth {
				typeWidth = len(iface.Type)
			}
			if len(iface.State) > stateWidth {
				stateWidth = len(iface.State)
			}
			if len(iface.Address) > addrWidth {
				addrWidth = len(iface.Address)
			}
			if len(iface.Speed) > speedWidth {
				speedWidth = len(iface.Speed)
			}
		}

		// Add padding
		nameWidth += 2
		typeWidth += 2
		stateWidth += 2
		addrWidth += 2
		speedWidth += 2

		// Print header
		fmt.Println("Network interfaces")
		printBorder(nameWidth, typeWidth, stateWidth, addrWidth, speedWidth, "┌", "┬", "┐")
		fmt.Printf("│ %-*s │ %-*s │ %-*s │ %-*s │ %-*s │\n",
			nameWidth-2, "NAME", typeWidth-2, "TYPE", stateWidth-2, "STATE", addrWidth-2, "ADDRESS", speedWidth-2, "SPEED")
		printBorder(nameWidth, typeWidth, stateWidth, addrWidth, speedWidth, "├", "┼", "┤")

		for _, iface := range interfaces {
			fmt.Printf("│ %-*s │ %-*s │ %-*s │ %-*s │ %-*s │\n",
				nameWidth-2, iface.Name, typeWidth-2, iface.Type, stateWidth-2, iface.State, addrWidth-2, iface.Address, speedWidth-2, iface.Speed)
		}

		printBorder(nameWidth, typeWidth, stateWidth, addrWidth, speedWidth, "└", "┴", "┘")
	},
}

func printBorder(w1, w2, w3, w4, w5 int, left, middle, right string) {
	fmt.Print(left)
	fmt.Print(strings.Repeat("─", w1))
	fmt.Print(middle)
	fmt.Print(strings.Repeat("─", w2))
	fmt.Print(middle)
	fmt.Print(strings.Repeat("─", w3))
	fmt.Print(middle)
	fmt.Print(strings.Repeat("─", w4))
	fmt.Print(middle)
	fmt.Print(strings.Repeat("─", w5))
	fmt.Println(right)
}

func padRight(s string, length int) string {
	if len(s) > length {
		return s
	}
	return s + strings.Repeat(" ", length-len(s))
}

func getNetworkInterfaces() ([]NetworkInterface, error) {
	linkOut, err := networkExecutor.CombinedOutput("ip", "-j", "link", "show")
	if err != nil {
		return nil, fmt.Errorf("ip link show failed: %s", err.Error())
	}

	addrOut, err := networkExecutor.CombinedOutput("ip", "-j", "addr", "show")
	if err != nil {
		return nil, fmt.Errorf("ip addr show failed: %s", err.Error())
	}

	linkOut = core.SanitizeInput(linkOut)
	addrOut = core.SanitizeInput(addrOut)

	var linkData []map[string]interface{}
	if err := json.Unmarshal([]byte(linkOut), &linkData); err != nil {
		return nil, err
	}

	var addrData []map[string]interface{}
	if err := json.Unmarshal([]byte(addrOut), &addrData); err != nil {
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
		} else if name == "docker0" || strings.HasPrefix(name, "br-") {
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
				speedOut, _ := networkExecutor.CombinedOutput("ethtool", name)
				if strings.Contains(speedOut, "Speed:") {
					lines := strings.Split(speedOut, "\n")
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
			iface := core.SanitizeInput(args[0])
			cmdArgs = append(cmdArgs, iface)
		}

		out, err := networkExecutor.CombinedOutput("ip", cmdArgs...)
		if err != nil {
			output.NewError(fmt.Sprintf("failed to show interface: %s", strings.TrimSpace(out)), "NETWORK_SHOW_ERROR").Print()
			return
		}

		var result interface{}
		if err := json.Unmarshal([]byte(out), &result); err == nil {
			output.NewSuccess(result).Print()
			return
		}

		output.NewSuccess(map[string]interface{}{
			"output": strings.TrimSpace(out),
		}).Print()
	},
}

var networkSetCmd = &cobra.Command{
	Use:   "set <interface> [flags]",
	Short: "Configure network interface",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		iface := core.SanitizeInput(args[0])

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
			ipAddr = core.SanitizeInput(ipAddr)
			if netmask != "" {
				commands = append(commands, fmt.Sprintf("ip addr add %s/%s dev %s", ipAddr, netmask, iface))
			} else {
				commands = append(commands, fmt.Sprintf("ip addr add %s dev %s", ipAddr, iface))
			}
		}
		if gateway != "" {
			gateway = core.SanitizeInput(gateway)
			commands = append(commands, fmt.Sprintf("ip route add default via %s", gateway))
		}

		if dryRun {
			output.NewInfo(map[string]interface{}{
				"interface": iface,
				"dry_run":   true,
				"commands":  commands,
			}).Print()
			return
		}

		for _, command := range commands {
			parts := strings.Fields(command)
			_, err := networkExecutor.Run(parts[0], parts[1:]...)
			if err != nil {
				output.NewError(fmt.Sprintf("command failed: %s", command), "NETWORK_SET_ERROR").Print()
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
