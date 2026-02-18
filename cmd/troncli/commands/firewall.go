package commands

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/mascli/troncli/internal/core/adapter"
	"github.com/mascli/troncli/internal/core/ports"
	"github.com/mascli/troncli/internal/core/services"
	"github.com/mascli/troncli/internal/modules/firewall"
	"github.com/spf13/cobra"
)

var firewallCmd = &cobra.Command{
	Use:   "firewall",
	Short: "Gerenciamento de Firewall",
	Long:  `Controlar regras de firewall (ufw, firewalld, iptables, nftables).`,
}

var fwAllowCmd = &cobra.Command{
	Use:   "allow [port] [protocol]",
	Short: "Permitir tráfego na porta",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		m, err := getFirewallManager()
		if err != nil {
			fmt.Println(err)
			return
		}
		if err := m.AllowPort(args[0], args[1]); err != nil {
			fmt.Printf("Error allowing port: %v\n", err)
			return
		}
		fmt.Printf("Port %s/%s allowed.\n", args[0], args[1])
	},
}

var fwBlockCmd = &cobra.Command{
	Use:   "deny [port] [protocol]",
	Short: "Bloquear tráfego na porta",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		m, err := getFirewallManager()
		if err != nil {
			fmt.Println(err)
			return
		}
		if err := m.BlockPort(args[0], args[1]); err != nil {
			fmt.Printf("Error blocking port: %v\n", err)
			return
		}
		fmt.Printf("Port %s/%s blocked.\n", args[0], args[1])
	},
}

var fwListCmd = &cobra.Command{
	Use:   "list",
	Short: "Listar regras de firewall",
	Run: func(cmd *cobra.Command, args []string) {
		m, err := getFirewallManager()
		if err != nil {
			fmt.Println(err)
			return
		}
		rules, err := m.ListRules()
		if err != nil {
			fmt.Printf("Error listing rules: %v\n", err)
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "ID\tACTION\tPROTO\tPORT\tCOMMENT")
		for _, r := range rules {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", r.ID, r.Action, r.Protocol, r.Port, r.Comment)
		}
		w.Flush()
	},
}

var fwEnableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Habilitar firewall",
	Run: func(cmd *cobra.Command, args []string) {
		m, err := getFirewallManager()
		if err != nil {
			fmt.Println(err)
			return
		}
		if err := m.Enable(); err != nil {
			fmt.Printf("Error enabling firewall: %v\n", err)
			return
		}
		fmt.Println("Firewall enabled.")
	},
}

var fwDisableCmd = &cobra.Command{
	Use:   "disable",
	Short: "Desabilitar firewall",
	Run: func(cmd *cobra.Command, args []string) {
		m, err := getFirewallManager()
		if err != nil {
			fmt.Println(err)
			return
		}
		if err := m.Disable(); err != nil {
			fmt.Printf("Error disabling firewall: %v\n", err)
			return
		}
		fmt.Println("Firewall disabled.")
	},
}

func init() {
	rootCmd.AddCommand(firewallCmd)
	firewallCmd.AddCommand(fwAllowCmd)
	firewallCmd.AddCommand(fwBlockCmd)
	firewallCmd.AddCommand(fwListCmd)
	firewallCmd.AddCommand(fwEnableCmd)
	firewallCmd.AddCommand(fwDisableCmd)
}

func getFirewallManager() (ports.FirewallManager, error) {
	executor := adapter.NewExecutor()
	profileEngine := services.NewProfileEngine(executor)
	profile, err := profileEngine.DetectProfile()
	if err != nil {
		return nil, fmt.Errorf("failed to detect system profile: %w", err)
	}
	return firewall.NewUniversalFirewallManager(executor, profile), nil
}
