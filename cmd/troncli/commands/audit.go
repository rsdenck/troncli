package commands

import (
	"fmt"
	"os"
	"time"

	"github.com/mascli/troncli/internal/core/adapter"
	"github.com/mascli/troncli/internal/core/ports"
	"github.com/mascli/troncli/internal/core/services"
	"github.com/mascli/troncli/internal/modules/audit"
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

		fmt.Printf("%-20s %-15s %-30s\n", "USER", "IP", "TIMESTAMP")
		fmt.Println("----------------------------------------------------------------")
		for _, e := range events {
			fmt.Printf("%-20s %-15s %-30s\n", e.User, e.IP, e.Message)
		}
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

		fmt.Println("Privileged Users (sudo/wheel/root):")
		for _, u := range users {
			fmt.Println("- " + u)
		}
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

		for _, e := range events {
			fmt.Printf("[%s] %s\n", e.Severity, e.Message)
		}
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

		fmt.Printf("%-20s %-50s\n", "USER", "COMMAND")
		fmt.Println("----------------------------------------------------------------")
		for _, e := range events {
			fmt.Printf("%-20s %-50s\n", e.User, e.Message)
		}
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
