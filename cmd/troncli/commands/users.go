package commands

import (
	"fmt"
	"os"

	"github.com/mascli/troncli/internal/core/adapter"
	"github.com/mascli/troncli/internal/core/ports"
	"github.com/mascli/troncli/internal/core/services"
	"github.com/mascli/troncli/internal/modules/users"
	"github.com/mascli/troncli/internal/ui/console"
	"github.com/spf13/cobra"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Gerenciamento de Usuários e Grupos",
	Long:  `Gerencie usuários e grupos do sistema (add, del, modify, list).`,
}

var userListCmd = &cobra.Command{
	Use:   "list",
	Short: "Listar usuários",
	Run: func(cmd *cobra.Command, args []string) {
		m, err := getUserManager()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		usersList, err := m.ListUsers()
		if err != nil {
			fmt.Printf("Error listing users: %v\n", err)
			return
		}

		table := console.NewBoxTable(os.Stdout)
		table.SetTitle("TRONCLI - LISTAGEM DE USUÁRIOS")
		table.SetHeaders([]string{"USERNAME", "UID", "GID", "SHELL", "HOME"})

		for _, u := range usersList {
			table.AddRow([]string{u.Username, u.UID, u.GID, u.Shell, u.HomeDir})
		}

		table.SetFooter(fmt.Sprintf("Total de usuários: %d", len(usersList)))
		table.Render()
	},
}

var userAddCmd = &cobra.Command{
	Use:   "add [username]",
	Short: "Adicionar novo usuário",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		m, err := getUserManager()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		username := args[0]
		opts := getUserOptions(cmd)

		if err := m.AddUser(username, opts); err != nil {
			fmt.Printf("Error adding user: %v\n", err)
			return
		}
		fmt.Printf("User %s added successfully.\n", username)
	},
}

var userDelCmd = &cobra.Command{
	Use:   "del [username]",
	Short: "Remover usuário",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		m, err := getUserManager()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		username := args[0]
		removeHome, _ := cmd.Flags().GetBool("remove-home")

		if err := m.DeleteUser(username, removeHome); err != nil {
			fmt.Printf("Error deleting user: %v\n", err)
			return
		}
		fmt.Printf("User %s deleted successfully.\n", username)
	},
}

var userModCmd = &cobra.Command{
	Use:   "modify [username]",
	Short: "Modificar usuário existente",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		m, err := getUserManager()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		username := args[0]
		opts := getUserOptions(cmd)

		if err := m.ModifyUser(username, opts); err != nil {
			fmt.Printf("Error modifying user: %v\n", err)
			return
		}
		fmt.Printf("User %s modified successfully.\n", username)
	},
}

// Group Commands
var groupCmd = &cobra.Command{
	Use:   "group",
	Short: "Gerenciamento de Grupos",
}

var groupListCmd = &cobra.Command{
	Use:   "list",
	Short: "Listar grupos",
	Run: func(cmd *cobra.Command, args []string) {
		m, err := getUserManager()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		groups, err := m.ListGroups()
		if err != nil {
			fmt.Printf("Error listing groups: %v\n", err)
			return
		}

		table := console.NewBoxTable(os.Stdout)
		table.SetTitle("TRONCLI - LISTAGEM DE GRUPOS")
		table.SetHeaders([]string{"GROUP", "GID", "MEMBERS"})

		for _, g := range groups {
			members := fmt.Sprintf("%v", g.Members)
			// Truncate members if too long?
			if len(members) > 50 {
				members = members[:47] + "..."
			}
			table.AddRow([]string{g.Groupname, g.GID, members})
		}

		table.SetFooter(fmt.Sprintf("Total de grupos: %d", len(groups)))
		table.Render()
	},
}

var groupAddCmd = &cobra.Command{
	Use:   "add [groupname]",
	Short: "Adicionar grupo",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		m, err := getUserManager()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		groupname := args[0]
		gid, _ := cmd.Flags().GetString("gid")

		if err := m.AddGroup(groupname, gid); err != nil {
			fmt.Printf("Error adding group: %v\n", err)
			return
		}
		fmt.Printf("Group %s added successfully.\n", groupname)
	},
}

var groupDelCmd = &cobra.Command{
	Use:   "del [groupname]",
	Short: "Remover grupo",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		m, err := getUserManager()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		groupname := args[0]
		if err := m.DeleteGroup(groupname); err != nil {
			fmt.Printf("Error deleting group: %v\n", err)
			return
		}
		fmt.Printf("Group %s deleted successfully.\n", groupname)
	},
}

func init() {
	rootCmd.AddCommand(userCmd)
	rootCmd.AddCommand(groupCmd)
	userCmd.AddCommand(userListCmd)
	userCmd.AddCommand(userAddCmd)
	userCmd.AddCommand(userDelCmd)
	userCmd.AddCommand(userModCmd)

	groupCmd.AddCommand(groupListCmd)
	groupCmd.AddCommand(groupAddCmd)
	groupCmd.AddCommand(groupDelCmd)

	// User Flags
	for _, cmd := range []*cobra.Command{userAddCmd, userModCmd} {
		cmd.Flags().String("uid", "", "User ID")
		cmd.Flags().String("gid", "", "Group ID")
		cmd.Flags().StringSlice("groups", nil, "Supplemental groups")
		cmd.Flags().String("shell", "", "Login shell")
		cmd.Flags().String("home", "", "Home directory")
		cmd.Flags().String("comment", "", "GECOS comment")
	}

	userDelCmd.Flags().Bool("remove-home", false, "Remove home directory")
	groupAddCmd.Flags().String("gid", "", "Group ID")
}

func getUserManager() (ports.UserManager, error) {
	executor := adapter.NewExecutor()
	profileEngine := services.NewProfileEngine(executor)
	profile, err := profileEngine.DetectProfile()
	if err != nil {
		return nil, fmt.Errorf("failed to detect system profile: %w", err)
	}

	return users.NewUniversalUserManager(executor, profile), nil
}

func getUserOptions(cmd *cobra.Command) ports.UserOptions {
	uid, _ := cmd.Flags().GetString("uid")
	gid, _ := cmd.Flags().GetString("gid")
	groups, _ := cmd.Flags().GetStringSlice("groups")
	shell, _ := cmd.Flags().GetString("shell")
	home, _ := cmd.Flags().GetString("home")
	comment, _ := cmd.Flags().GetString("comment")

	return ports.UserOptions{
		UID:     uid,
		GID:     gid,
		Groups:  groups,
		Shell:   shell,
		HomeDir: home,
		Comment: comment,
	}
}
