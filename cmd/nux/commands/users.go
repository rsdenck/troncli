package commands

import (

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
		output.NewList([]map[string]interface{}{}, 0).WithMessage("User list").Print()
	},
}

func init() {
	usersCmd.AddCommand(usersListCmd)
	rootCmd.AddCommand(usersCmd)
}
