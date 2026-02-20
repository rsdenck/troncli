package commands

import (
	"fmt"
	"os"
	"time"

	"github.com/mascli/troncli/internal/core/adapter"
	"github.com/mascli/troncli/internal/core/ports"
	"github.com/mascli/troncli/internal/core/services"
	"github.com/mascli/troncli/internal/modules/audit"
	"github.com/mascli/troncli/internal/ui/console"
	"github.com/spf13/cobra"
)

var auditCmd = &cobra.Command{
	Use:   "audit",
	Short: "Auditoria de Segurança",
	Long:  `Ferramentas para auditoria de segurança, logs e integridade do sistema.`,
}

var auditLoginsCmd = &cobra.Command{
	Use:   "logins",
	Short: "Lista histórico de logins",
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := getAuditManager()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Default lookback 24h
		events, err := manager.AnalyzeLogins(24 * time.Hour)
		if err != nil {
			fmt.Printf("Error analyzing logins: %v\n", err)
			os.Exit(1)
		}

		table := console.NewBoxTable(os.Stdout)
		table.SetTitle("TRONCLI - HISTÓRICO DE LOGINS (24H)")
		table.SetHeaders([]string{"USER", "IP", "MESSAGE"})

		for _, e := range events {
			// Truncate message if needed
			msg := e.Message
			if len(msg) > 50 {
				msg = msg[:47] + "..."
			}
			table.AddRow([]string{e.User, e.IP, msg})
		}
		table.SetFooter(fmt.Sprintf("Total logins: %d", len(events)))
		table.Render()
	},
}

var auditSudoersCmd = &cobra.Command{
	Use:   "sudoers",
	Short: "Lista usuários com privilégios sudo",
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := getAuditManager()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		users, err := manager.CheckPrivilegedGroups()
		if err != nil {
			fmt.Printf("Error checking privileged groups: %v\n", err)
			os.Exit(1)
		}

		table := console.NewBoxTable(os.Stdout)
		table.SetTitle("TRONCLI - USUÁRIOS PRIVILEGIADOS (SUDO/WHEEL/ROOT)")
		table.SetHeaders([]string{"USERNAME"})

		for _, u := range users {
			table.AddRow([]string{u})
		}
		table.SetFooter(fmt.Sprintf("Total privileged users: %d", len(users)))
		table.Render()
	},
}

var auditFileChangesCmd = &cobra.Command{
	Use:   "file-changes [paths...]",
	Short: "Monitora alterações recentes em arquivos",
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := getAuditManager()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		paths := args
		if len(paths) == 0 {
			paths = []string{"/etc", "/bin", "/sbin", "/usr/bin", "/usr/sbin"}
		}

		fmt.Printf("Scanning for file changes in last 24h: %v\n", paths)
		events, err := manager.AnalyzeFileChanges(paths, 24*time.Hour)
		if err != nil {
			fmt.Printf("Error analyzing file changes: %v\n", err)
			os.Exit(1)
		}

		if len(events) == 0 {
			fmt.Println("No recent file changes detected.")
			return
		}

		table := console.NewBoxTable(os.Stdout)
		table.SetTitle("TRONCLI - ALTERAÇÕES EM ARQUIVOS (24H)")
		table.SetHeaders([]string{"SEVERITY", "MESSAGE"})

		for _, e := range events {
			msg := e.Message
			if len(msg) > 60 {
				msg = msg[:57] + "..."
			}
			table.AddRow([]string{e.Severity, msg})
		}
		table.SetFooter(fmt.Sprintf("Total events: %d", len(events)))
		table.Render()
	},
}

var auditCommandsCmd = &cobra.Command{
	Use:   "commands",
	Short: "Lista comandos executados (via sudo/logs)",
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := getAuditManager()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		events, err := manager.AnalyzeSudo(24 * time.Hour)
		if err != nil {
			fmt.Printf("Error analyzing commands: %v\n", err)
			os.Exit(1)
		}

		table := console.NewBoxTable(os.Stdout)
		table.SetTitle("TRONCLI - HISTÓRICO DE COMANDOS (SUDO/LOGS)")
		table.SetHeaders([]string{"USER", "COMMAND"})

		for _, e := range events {
			cmd := e.Message
			if len(cmd) > 60 {
				cmd = cmd[:57] + "..."
			}
			table.AddRow([]string{e.User, cmd})
		}
		table.SetFooter(fmt.Sprintf("Total commands: %d", len(events)))
		table.Render()
	},
}

func init() {
	rootCmd.AddCommand(auditCmd)
	auditCmd.AddCommand(auditLoginsCmd, auditSudoersCmd, auditFileChangesCmd, auditCommandsCmd)
}

func getAuditManager() (ports.AuditManager, error) {
	executor := adapter.NewExecutor()
	profileEngine := services.NewProfileEngine(executor)
	profile, err := profileEngine.DetectProfile()
	if err != nil {
		return nil, fmt.Errorf("failed to detect system profile: %w", err)
	}

	return audit.NewUniversalAuditManager(executor, profile), nil
}
