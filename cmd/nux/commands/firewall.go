package commands

import (

	"github.com/rsdenck/nux/internal/output"
	"github.com/spf13/cobra"
)

var firewallCmd = &cobra.Command{
	Use:   "firewall",
	Short: "Firewall management",
	Long:  `Manage firewall rules (nftables, iptables, firewalld, ufw).`,
}

var firewallListCmd = &cobra.Command{
	Use:   "list",
	Short: "List firewall rules",
	Run: func(cmd *cobra.Command, args []string) {
		output.NewList([]map[string]interface{}{}, 0).WithMessage("Firewall rules").Print()
	},
}

func init() {
	firewallCmd.AddCommand(firewallListCmd)
	rootCmd.AddCommand(firewallCmd)
}
