package commands

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/mascli/troncli/internal/core/adapter"
	"github.com/mascli/troncli/internal/core/ports"
	"github.com/mascli/troncli/internal/core/services"
	"github.com/mascli/troncli/internal/modules/network"
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

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "NAME\tSTATE\tMTU\tIPs\tMAC")
		for _, i := range ifaces {
			fmt.Fprintf(w, "%s\t%s\t%d\t%v\t%s\n", i.Name, i.State, i.MTU, i.IPAddresses, i.HardwareAddr)
		}
		w.Flush()

		// Hostname and DNS
		host, _ := m.GetHostname()
		dns, _ := m.GetDNSConfig()
		fmt.Printf("\nHostname: %s\n", host)
		fmt.Printf("DNS: %v\n", dns)
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

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "PROTO\tSTATE\tLOCAL\tREMOTE")
		for _, s := range stats {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", s.Protocol, s.State, s.Local, s.Remote)
		}
		w.Flush()
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
		fmt.Println(res)
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
		fmt.Println(res)
	},
}

var netNmapCmd = &cobra.Command{
	Use:   "scan [target]",
	Short: "Escanear portas (nmap)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		m, err := getNetworkManager()
		if err != nil {
			fmt.Println(err)
			return
		}
		opts, _ := cmd.Flags().GetString("opts")
		res, err := m.RunNmap(args[0], opts)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Println(res)
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
		fmt.Println(res)
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
