package commands

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/rsdenck/nux/internal/output"
	"github.com/spf13/cobra"
)

var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "User management",
	Long:  `Manage system users and groups.`,
}

var usersListCmd = &cobra.Command{
	Use:   "list",
	Short: "List users",
	Run: func(cmd *cobra.Command, args []string) {
		// Use getent passwd to list users
		getentCmd := exec.Command("getent", "passwd")
		out, err := getentCmd.CombinedOutput()
		
		if err != nil {
			output.NewError(fmt.Sprintf("failed to list users: %s", strings.TrimSpace(string(out))), "USERS_LIST_ERROR").Print()
			return
		}
		
		// Parse output
		lines := strings.Split(strings.TrimSpace(string(out)), "\n")
		items := []map[string]interface{}{}
		
		for _, line := range lines {
			parts := strings.Split(line, ":")
			if len(parts) < 7 {
				continue
			}
			
			item := map[string]interface{}{
				"username": parts[0],
				"uid":      parts[2],
				"gid":      parts[3],
				"home":     parts[5],
				"shell":    parts[6],
			}
			items = append(items, item)
		}
		
		output.NewList(items, len(items)).WithMessage("User list").Print()
	},
}

var usersAddCmd = &cobra.Command{
	Use:   "add <username>",
	Short: "Add a new user",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		username := args[0]
		
		createHome, _ := cmd.Flags().GetBool("create-home")
		shell, _ := cmd.Flags().GetString("shell")
		group, _ := cmd.Flags().GetString("group")
		
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		
		cmdArgs := []string{}
		if createHome {
			cmdArgs = append(cmdArgs, "-m")
		}
		if shell != "" {
			cmdArgs = append(cmdArgs, "-s", shell)
		}
		if group != "" {
			cmdArgs = append(cmdArgs, "-g", group)
		}
		cmdArgs = append(cmdArgs, username)
		
		if dryRun {
			output.NewInfo(map[string]interface{}{
				"username": username,
				"dry_run": true,
				"command":  fmt.Sprintf("useradd %s", strings.Join(cmdArgs, " ")),
			}).Print()
			return
		}
		
		useraddCmd := exec.Command("useradd", cmdArgs...)
		out, err := useraddCmd.CombinedOutput()
		
		if err != nil {
			output.NewError(fmt.Sprintf("failed to add user: %s - %s", err.Error(), strings.TrimSpace(string(out))), "USERS_ADD_ERROR").Print()
			return
		}
		
		output.NewSuccess(map[string]interface{}{
			"username": username,
			"status":   "created",
		}).Print()
	},
}

var usersDeleteCmd = &cobra.Command{
	Use:   "delete <username>",
	Short: "Delete a user",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		username := args[0]
		
		removeHome, _ := cmd.Flags().GetBool("remove-home")
		
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		
		cmdArgs := []string{}
		if removeHome {
			cmdArgs = append(cmdArgs, "-r")
		}
		cmdArgs = append(cmdArgs, username)
		
		if dryRun {
			output.NewInfo(map[string]interface{}{
				"username": username,
				"dry_run": true,
				"command":  fmt.Sprintf("userdel %s", strings.Join(cmdArgs, " ")),
			}).Print()
			return
		}
		
		userdelCmd := exec.Command("userdel", cmdArgs...)
		out, err := userdelCmd.CombinedOutput()
		
		if err != nil {
			output.NewError(fmt.Sprintf("failed to delete user: %s - %s", err.Error(), strings.TrimSpace(string(out))), "USERS_DELETE_ERROR").Print()
			return
		}
		
		output.NewSuccess(map[string]interface{}{
			"username": username,
			"status":   "deleted",
		}).Print()
	},
}

var groupsListCmd = &cobra.Command{
	Use:   "groups",
	Short: "List groups",
	Run: func(cmd *cobra.Command, args []string) {
		getentCmd := exec.Command("getent", "group")
		out, err := getentCmd.CombinedOutput()
		
		if err != nil {
			output.NewError(fmt.Sprintf("failed to list groups: %s", strings.TrimSpace(string(out))), "USERS_GROUPS_ERROR").Print()
			return
		}
		
		lines := strings.Split(strings.TrimSpace(string(out)), "\n")
		items := []map[string]interface{}{}
		
		for _, line := range lines {
			parts := strings.Split(line, ":")
			if len(parts) < 4 {
				continue
			}
			
			item := map[string]interface{}{
				"groupname": parts[0],
				"gid":       parts[2],
				"members":   parts[3],
			}
			items = append(items, item)
		}
		
		output.NewList(items, len(items)).WithMessage("Group list").Print()
	},
}

func init() {
	usersCmd.AddCommand(usersListCmd)
	usersAddCmd.Flags().Bool("create-home", true, "Create home directory")
	usersAddCmd.Flags().String("shell", "", "Login shell")
	usersAddCmd.Flags().String("group", "", "Primary group")
	usersCmd.AddCommand(usersAddCmd)
	usersDeleteCmd.Flags().Bool("remove-home", false, "Remove home directory")
	usersCmd.AddCommand(usersDeleteCmd)
	usersCmd.AddCommand(groupsListCmd)
	rootCmd.AddCommand(usersCmd)
}
