package commands

import (
	"fmt"
	"strings"

	"github.com/rsdenck/nux/internal/core"
	"github.com/rsdenck/nux/internal/output"
	"github.com/spf13/cobra"
)

var usersExecutor core.Executor = &core.RealExecutor{}

var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "User management",
	Long:  `Manage system users and groups.`,
}

var usersListCmd = &cobra.Command{
	Use:   "list",
	Short: "List users",
	Run: func(cmd *cobra.Command, args []string) {
		out, err := usersExecutor.CombinedOutput("getent", "passwd")
		if err != nil {
			output.NewError(fmt.Sprintf("failed to list users: %s", err.Error()), "USERS_LIST_ERROR").Print()
			return
		}

		lines := strings.Split(strings.TrimSpace(out), "\n")
		if len(lines) == 0 {
			output.NewInfo("No users found").Print()
			return
		}

		// Fixed columns order
		headers := []string{"USERNAME", "UID", "GID", "HOME", "SHELL"}
		colWidths := make([]int, len(headers))
		for i, h := range headers {
			colWidths[i] = len(h) + 2
		}

		type userEntry struct {
			username string
			uid      string
			gid      string
			home     string
			shell    string
		}
		var entries []userEntry

		for _, line := range lines {
			parts := strings.Split(line, ":")
			if len(parts) < 7 {
				continue
			}
			entry := userEntry{
				username: parts[0],
				uid:      parts[2],
				gid:      parts[3],
				home:     parts[5],
				shell:    parts[6],
			}
			entries = append(entries, entry)

			// Update widths
			if len(entry.username)+2 > colWidths[0] {
				colWidths[0] = len(entry.username) + 2
			}
			if len(entry.uid)+2 > colWidths[1] {
				colWidths[1] = len(entry.uid) + 2
			}
			if len(entry.gid)+2 > colWidths[2] {
				colWidths[2] = len(entry.gid) + 2
			}
			if len(entry.home)+2 > colWidths[3] {
				colWidths[3] = len(entry.home) + 2
			}
			if len(entry.shell)+2 > colWidths[4] {
				colWidths[4] = len(entry.shell) + 2
			}
		}

		fmt.Println("User list")
		printBorderUsers(colWidths, "┌", "┬", "┐")
		printRowUsers(headers, colWidths)
		printBorderUsers(colWidths, "├", "┼", "┤")

		for _, e := range entries {
			printRowUsers([]string{e.username, e.uid, e.gid, e.home, e.shell}, colWidths)
		}

		printBorderUsers(colWidths, "└", "┴", "┘")
		fmt.Printf("\n%d users found\n", len(entries))
	},
}

func printBorderUsers(widths []int, left, middle, right string) {
	fmt.Print(left)
	for i, w := range widths {
		fmt.Print(strings.Repeat("─", w))
		if i < len(widths)-1 {
			fmt.Print(middle)
		}
	}
	fmt.Println(right)
}

func printRowUsers(values []string, widths []int) {
	fmt.Print("│")
	for i, v := range values {
		fmt.Printf(" %-*s ", widths[i]-2, v)
		fmt.Print("│")
	}
	fmt.Println()
}

var usersAddCmd = &cobra.Command{
	Use:   "add <username>",
	Short: "Add a new user",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		username := core.SanitizeInput(args[0])

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
				"dry_run":  true,
				"command":  fmt.Sprintf("useradd %s", strings.Join(cmdArgs, " ")),
			}).Print()
			return
		}

		_, err := usersExecutor.CombinedOutput("useradd", cmdArgs...)
		if err != nil {
			output.NewError(fmt.Sprintf("failed to add user: %s", err.Error()), "USERS_ADD_ERROR").Print()
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
		username := core.SanitizeInput(args[0])

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
				"dry_run":  true,
				"command":  fmt.Sprintf("userdel %s", strings.Join(cmdArgs, " ")),
			}).Print()
			return
		}

		_, err := usersExecutor.CombinedOutput("userdel", cmdArgs...)
		if err != nil {
			output.NewError(fmt.Sprintf("failed to delete user: %s", err.Error()), "USERS_DELETE_ERROR").Print()
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
		out, err := usersExecutor.CombinedOutput("getent", "group")
		if err != nil {
			output.NewError(fmt.Sprintf("failed to list groups: %s", err.Error()), "USERS_GROUPS_ERROR").Print()
			return
		}

		lines := strings.Split(strings.TrimSpace(out), "\n")
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
	usersAddCmd.Flags().Bool("create-home", true, "Create home directory")
	usersAddCmd.Flags().String("shell", "", "Login shell")
	usersAddCmd.Flags().String("group", "", "Primary group")
	usersAddCmd.Flags().Bool("dry-run", false, "Simulate command")

	usersDeleteCmd.Flags().Bool("remove-home", false, "Remove home directory")
	usersDeleteCmd.Flags().Bool("dry-run", false, "Simulate command")

	usersCmd.AddCommand(usersListCmd)
	usersCmd.AddCommand(usersAddCmd)
	usersCmd.AddCommand(usersDeleteCmd)
	usersCmd.AddCommand(groupsListCmd)
	rootCmd.AddCommand(usersCmd)
}
