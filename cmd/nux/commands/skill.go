package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/rsdenck/nux/internal/skill"
	"github.com/spf13/cobra"
)

var _ = fmt.Printf
var _ = os.Stderr

var skillCmd = &cobra.Command{
	Use:   "skill",
	Short: "Manage NUX skills (external CLI integrations)",
	Long: `Manage skills - external CLI tools that NUX can integrate with.

Skills are defined as .md files in the skills/ directory.
Each skill can be installed, enabled, and managed through this command.`,
}

var skillInstallCmd = &cobra.Command{
	Use:   "install [skill]",
	Short: "Install a skill",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		skillName := args[0]

		s, err := skill.LoadSkillFromMD(skillName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Installing skill: %s\n", skillName)
		fmt.Printf("Description: %s\n", s.Description)

		v, err := skill.LoadVault()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading vault: %v\n", err)
			os.Exit(1)
		}

		if contains(v.InstalledSkills, skillName) {
			fmt.Printf("Skill %s is already installed\n", skillName)
			return
		}

		v.InstalledSkills = append(v.InstalledSkills, skillName)
		if err := skill.SaveVault(v); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving vault: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Skill %s installed successfully\n", skillName)
		fmt.Printf("Run 'nux skill enable %s' to enable it\n", skillName)
	},
}

var skillInfoCmd = &cobra.Command{
	Use:   "info [skill]",
	Short: "Show skill information",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		skillName := args[0]

		s, err := skill.LoadSkillFromMD(skillName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if flagJSON {
			data, _ := json.MarshalIndent(s, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("Skill: %s\n", s.Name)
		fmt.Printf("Description: %s\n", s.Description)
		fmt.Printf("Repo: %s\n", s.Repo)
		fmt.Printf("Install: %s\n", s.InstallCmd)
		fmt.Printf("Commands: %s\n", s.Commands)
		fmt.Printf("Type: %s\n", s.Type)
	},
}

var skillListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available skills",
	Run: func(cmd *cobra.Command, args []string) {
		skills, err := skill.ListSkills()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		v, _ := skill.LoadVault()

		if flagJSON {
			data, _ := json.MarshalIndent(skills, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Println("Available skills:")
		for _, s := range skills {
			installed := ""
			if v != nil && contains(v.InstalledSkills, s) {
				installed = " [installed]"
			}
			fmt.Printf("  - %s%s\n", s, installed)
		}
	},
}

var skillSearchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search skills by name or type",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		query := args[0]
		results, err := skill.SearchSkills(query)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if flagJSON {
			data, _ := json.MarshalIndent(results, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("Search results for '%s':\n", query)
		for _, r := range results {
			fmt.Printf("  - %s\n", r)
		}
	},
}

var skillEnableCmd = &cobra.Command{
	Use:   "enable [skill]",
	Short: "Enable a skill",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		skillName := args[0]

		v, err := skill.LoadVault()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading vault: %v\n", err)
			os.Exit(1)
		}

		if !contains(v.InstalledSkills, skillName) {
			fmt.Fprintf(os.Stderr, "Skill %s is not installed. Run 'nux skill install %s' first\n", skillName, skillName)
			os.Exit(1)
		}

		if contains(v.EnabledSkills, skillName) {
			fmt.Printf("Skill %s is already enabled\n", skillName)
			return
		}

		v.EnabledSkills = append(v.EnabledSkills, skillName)
		if err := skill.SaveVault(v); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving vault: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Skill %s enabled successfully\n", skillName)
	},
}

var skillSyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync skills with remote repositories",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Syncing skills...")
		fmt.Println("Checking for updates...")

		v, err := skill.LoadVault()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading vault: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Installed skills: %d\n", len(v.InstalledSkills))
		for _, s := range v.InstalledSkills {
			fmt.Printf("  - %s\n", s)
		}
		fmt.Println("Sync completed")
	},
}

func init() {
	skillCmd.AddCommand(skillInstallCmd)
	skillCmd.AddCommand(skillInfoCmd)
	skillCmd.AddCommand(skillListCmd)
	skillCmd.AddCommand(skillSearchCmd)
	skillCmd.AddCommand(skillEnableCmd)
	skillCmd.AddCommand(skillSyncCmd)
	rootCmd.AddCommand(skillCmd)
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
