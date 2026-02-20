package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/mascli/troncli/internal/core/adapter"
	"github.com/mascli/troncli/internal/core/ports"
	"github.com/mascli/troncli/internal/core/services"
	"github.com/mascli/troncli/internal/modules/network"
	"github.com/mascli/troncli/internal/ui/console"
	"github.com/spf13/cobra"
)

var networkCmd = &cobra.Command{
	Use:   "network",
	Short: "Gerenciamento de Rede",
	Long:  `Gerencie interfaces, rotas, DNS e configurações de rede.`,
}

var netInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Informações detalhadas de rede",
	Run: func(cmd *cobra.Command, args []string) {
		m, err := getNetworkManager()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		ifaces, err := m.GetInterfaces()
		if err != nil {
			fmt.Printf("Error getting interfaces: %v\n", err)
			return
		}

		table := console.NewBoxTable(os.Stdout)
		table.SetTitle("TRONCLI - INFORMAÇÕES DE REDE")
		table.SetHeaders([]string{"NAME", "STATE", "MTU", "IPs", "MAC"})

		for _, i := range ifaces {
			ips := fmt.Sprintf("%v", i.IPAddresses)
			if len(ips) > 50 {
				ips = ips[:47] + "..."
			}
			table.AddRow([]string{i.Name, i.State, fmt.Sprintf("%d", i.MTU), ips, i.HardwareAddr})
		}
		table.Render()

		// Hostname and DNS
		host, _ := m.GetHostname()
		dns, _ := m.GetDNSConfig()

		fmt.Println() // Spacer

		infoTable := console.NewBoxTable(os.Stdout)
		infoTable.SetTitle("CONFIGURAÇÃO GERAL")
		infoTable.SetHeaders([]string{"KEY", "VALUE"})
		infoTable.AddRow([]string{"Hostname", host})
		infoTable.AddRow([]string{"DNS", fmt.Sprintf("%v", dns)})
		infoTable.Render()
	},
}

var netSetCmd = &cobra.Command{
	Use:   "set-state [interface] [up/down]",
	Short: "Alterar estado da interface",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		m, err := getNetworkManager()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		iface := args[0]
		state := args[1]
		up := state == "up"

		if err := m.SetInterfaceState(iface, up); err != nil {
			fmt.Printf("Error setting interface state: %v\n", err)
			return
		}
		fmt.Printf("Interface %s set to %s successfully.\n", iface, state)
	},
}

var netSocketCmd = &cobra.Command{
	Use:   "sockets",
	Short: "Listar sockets abertos (ss)",
	Run: func(cmd *cobra.Command, args []string) {
		m, err := getNetworkManager()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		stats, err := m.GetSocketStats()
		if err != nil {
			fmt.Printf("Error getting socket stats: %v\n", err)
			return
		}

		table := console.NewBoxTable(os.Stdout)
		table.SetTitle("TRONCLI - SOCKETS ABERTOS")
		table.SetHeaders([]string{"PROTO", "STATE", "LOCAL", "REMOTE"})

		for _, s := range stats {
			table.AddRow([]string{s.Protocol, s.State, s.Local, s.Remote})
		}

		table.SetFooter(fmt.Sprintf("Total sockets: %d", len(stats)))
		table.Render()
	},
}

var netTraceCmd = &cobra.Command{
	Use:   "trace [target]",
	Short: "Executar traceroute",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		m, err := getNetworkManager()
		if err != nil {
			fmt.Println(err)
			return
		}
		res, err := m.RunTraceRoute(args[0])
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		lines := strings.Split(res, "\n")
		table := console.NewBoxTable(os.Stdout)
		table.SetTitle(fmt.Sprintf("TRONCLI - TRACEROUTE: %s", args[0]))
		table.SetHeaders([]string{"RAW OUTPUT"})

		count := 0
		for _, line := range lines {
			if strings.TrimSpace(line) != "" {
				table.AddRow([]string{line})
				count++
			}
		}
		table.SetFooter(fmt.Sprintf("Lines: %d", count))
		table.Render()
	},
}

var netDigCmd = &cobra.Command{
	Use:   "dig [target]",
	Short: "Consultar DNS (dig)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		m, err := getNetworkManager()
		if err != nil {
			fmt.Println(err)
			return
		}
		res, err := m.RunDig(args[0])
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		lines := strings.Split(res, "\n")
		table := console.NewBoxTable(os.Stdout)
		table.SetTitle(fmt.Sprintf("TRONCLI - DNS LOOKUP: %s", args[0]))
		table.SetHeaders([]string{"DNS RECORDS"})

		for _, line := range lines {
			if strings.TrimSpace(line) != "" {
				table.AddRow([]string{line})
			}
		}
		table.Render()
	},
}

var netNmapCmd = &cobra.Command{
	Use:   "scan [target]",
	Short: "Escanear portas (nmap)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		m, err := getNetworkManager()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		opts, _ := cmd.Flags().GetString("opts")
		results, err := m.RunNmap(args[0], opts)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		table := console.NewBoxTable(os.Stdout)
		table.SetTitle(fmt.Sprintf("TRONCLI - PORT SCAN: %s", args[0]))
		table.SetHeaders([]string{"PORT", "PROTO", "STATE", "SERVICE"})

		for _, r := range results {
			table.AddRow([]string{
				fmt.Sprintf("%d", r.Port),
				r.Protocol,
				r.State,
				r.Service,
			})
		}

		table.SetFooter(fmt.Sprintf("Open ports: %d", len(results)))
		table.Render()
	},
}

var netTcpdumpCmd = &cobra.Command{
	Use:   "capture [interface]",
	Short: "Capturar pacotes (tcpdump)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		m, err := getNetworkManager()
		if err != nil {
			fmt.Println(err)
			return
		}
		filter, _ := cmd.Flags().GetString("filter")
		duration, _ := cmd.Flags().GetInt("duration")

		fmt.Printf("Capturing on %s for %ds...\n", args[0], duration)
		res, err := m.RunTcpdump(args[0], filter, duration)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		lines := strings.Split(res, "\n")
		table := console.NewBoxTable(os.Stdout)
		table.SetTitle(fmt.Sprintf("TRONCLI - PACKET CAPTURE: %s", args[0]))
		table.SetHeaders([]string{"PACKETS"})

		count := 0
		for _, line := range lines {
			if strings.TrimSpace(line) != "" {
				if len(line) > 100 {
					line = line[:97] + "..."
				}
				table.AddRow([]string{line})
				count++
			}
		}
		table.SetFooter(fmt.Sprintf("Packets captured: %d", count))
		table.Render()
	},
}

func init() {
	rootCmd.AddCommand(networkCmd)
	networkCmd.AddCommand(netInfoCmd)
	networkCmd.AddCommand(netSetCmd)
	networkCmd.AddCommand(netSocketCmd)
	networkCmd.AddCommand(netTraceCmd)
	networkCmd.AddCommand(netDigCmd)
	networkCmd.AddCommand(netNmapCmd)
	networkCmd.AddCommand(netTcpdumpCmd)

	netNmapCmd.Flags().String("opts", "", "Opções adicionais para o nmap")
	netTcpdumpCmd.Flags().String("filter", "", "Filtro tcpdump (ex: 'port 80')")
	netTcpdumpCmd.Flags().Int("duration", 10, "Duração da captura em segundos")
}

func getNetworkManager() (ports.NetworkManager, error) {
	executor := adapter.NewExecutor()
	profileEngine := services.NewProfileEngine(executor)
	profile, err := profileEngine.DetectProfile()
	if err != nil {
		return nil, fmt.Errorf("failed to detect system profile: %w", err)
	}

	return network.NewUniversalNetworkManager(executor, profile), nil
}
