package commands

import (
	"fmt"
	"os"
	"time"

	"github.com/rsdenck/nux/internal/core/adapter"
	"github.com/rsdenck/nux/internal/core/services"
	"github.com/rsdenck/nux/internal/modules/audit"
	"github.com/rsdenck/nux/internal/output"
	"github.com/spf13/cobra"
)

var _ = fmt.Printf
var _ = os.Exit

var auditCmd = &cobra.Command{
	Use:   "audit",
	Short: "Security auditing",
	Long:  `Security auditing tools for logs and system integrity.`,
}

var auditLoginsCmd = &cobra.Command{
	Use:   "logins",
	Short: "List login history",
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := getAuditManager()
		if err != nil {
			output.NewError(err.Error(), "ERR_AUDIT").Print()
			os.Exit(1)
		}

		// Default lookback 24h
		events, err := manager.AnalyzeLogins(24 * time.Hour)
		if err != nil {
			output.NewError(err.Error(), "ERR_LOGINS").Print()
			os.Exit(1)
		}

		items := make([]map[string]interface{}, 0)
		for _, e := range events {
			items = append(items, map[string]interface{}{
				"user":    e.User,
				"ip":      e.IP,
				"message": e.Message,
			})
		}
		output.NewList(items, len(items)).WithMessage("Login history (24h)").Print()
	},
}

var auditPrivilegedCmd = &cobra.Command{
	Use:   "privileged",
	Short: "List privileged users",
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := getAuditManager()
		if err != nil {
			output.NewError(err.Error(), "ERR_AUDIT").Print()
			os.Exit(1)
		}

		users, err := manager.CheckPrivilegedGroups()
		if err != nil {
			output.NewError(err.Error(), "ERR_PRIVILEGED").Print()
			os.Exit(1)
		}

		items := make([]map[string]interface{}, 0)
		for _, u := range users {
			items = append(items, map[string]interface{}{
				"user": u,
			})
		}
		output.NewList(items, len(items)).WithMessage("Privileged users (sudo/wheel/root)").Print()
	},
}

var auditChangesCmd = &cobra.Command{
	Use:   "changes",
	Short: "Show file changes",
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := getAuditManager()
		if err != nil {
			output.NewError(err.Error(), "ERR_AUDIT").Print()
			os.Exit(1)
		}

		paths := []string{"/etc", "/var"}
		changes, err := manager.AnalyzeFileChanges(paths, 24*time.Hour)
		if err != nil {
			output.NewError(err.Error(), "ERR_CHANGES").Print()
			os.Exit(1)
		}

		items := make([]map[string]interface{}, 0)
		for _, c := range changes {
			items = append(items, map[string]interface{}{
				"user":    c.User,
				"message": c.Message,
				"time":    c.Timestamp,
			})
		}
		output.NewList(items, len(items)).WithMessage("File changes (24h)").Print()
	},
}

var auditCommandsCmd = &cobra.Command{
	Use:   "commands",
	Short: "Show command history",
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := getAuditManager()
		if err != nil {
			output.NewError(err.Error(), "ERR_AUDIT").Print()
			os.Exit(1)
		}

		commands, err := manager.AnalyzeSudo(24 * time.Hour)
		if err != nil {
			output.NewError(err.Error(), "ERR_COMMANDS").Print()
			os.Exit(1)
		}

		items := make([]map[string]interface{}, 0)
		for _, c := range commands {
			items = append(items, map[string]interface{}{
				"user":    c.User,
				"message": c.Message,
				"time":    c.Timestamp,
			})
		}
		output.NewList(items, len(items)).WithMessage("Command history (sudo/logs)").Print()
	},
}

func getAuditManager() (*audit.UniversalAuditManager, error) {
	executor := adapter.NewExecutor()
	profileEngine := services.NewProfileEngine(executor)
	profile, err := profileEngine.DetectProfile()
	if err != nil {
		return nil, fmt.Errorf("failed to detect system profile: %w", err)
	}
	return audit.NewUniversalAuditManager(executor, profile), nil
}

func init() {
	auditCmd.AddCommand(auditLoginsCmd)
	auditCmd.AddCommand(auditPrivilegedCmd)
	auditCmd.AddCommand(auditChangesCmd)
	auditCmd.AddCommand(auditCommandsCmd)
	rootCmd.AddCommand(auditCmd)
}
