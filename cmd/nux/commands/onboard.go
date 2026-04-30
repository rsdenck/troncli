package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rsdenck/nux/internal/output"
	"github.com/rsdenck/nux/internal/skill"
	"github.com/spf13/cobra"
)

var onboardCmd = &cobra.Command{
	Use:   "onboard",
	Short: "First-time setup and skill selection",
	Long:  `Interactive setup that runs on first install. Lists available skills and allows sysadmin to select which to enable.`,
	Run: func(cmd *cobra.Command, args []string) {
		output.NewInfo("Welcome to NUX - Linux Master CLI").Print()
		output.NewInfo("Running first-time onboard process...").Print()

		vault, err := skill.LoadVault()
		if err != nil {
			vault = &skill.Vault{
				Version:         "1.0.0",
				InstalledSkills: []string{},
				EnabledSkills:   []string{},
				APIKeys:         make(map[string]string),
			}
		}

		skillsDir := "/opt/cli/nux/skills"
		files, err := filepath.Glob(filepath.Join(skillsDir, "*.md"))
		if err != nil {
			output.NewError(fmt.Sprintf("failed to read skills directory: %s", err.Error()), "ONBOARD_ERROR").Print()
			return
		}

		output.NewInfo(fmt.Sprintf("Found %d skills available", len(files))).Print()
		output.NewInfo("For each skill, choose 'yes' to enable or 'no' to skip").Print()

		enabled := 0
		skipped := 0

		for _, file := range files {
			skillName := strings.TrimSuffix(filepath.Base(file), ".md")

			if skillName == "geoip" || skillName == "onboard" {
				continue
			}

			isEnabled := false
			for _, s := range vault.EnabledSkills {
				if s == skillName {
					isEnabled = true
					break
				}
			}

			status := "inactive"
			if isEnabled {
				status = "active"
			}

			output.NewInfo(fmt.Sprintf("\nSkill: %s", skillName)).Print()
			output.NewInfo(fmt.Sprintf("Current status: %s", status)).Print()

			fmt.Printf("Enable skill %s? (yes/no): ", skillName)
			var answer string
			fmt.Scanln(&answer)

			if strings.ToLower(answer) == "yes" || strings.ToLower(answer) == "y" {
				vault.EnabledSkills = append(vault.EnabledSkills, skillName)
				enabled++
				output.NewSuccess(fmt.Sprintf("Skill %s enabled", skillName)).Print()
			} else {
				skipped++
				output.NewInfo(fmt.Sprintf("Skill %s skipped", skillName)).Print()
			}
		}

		geoipPath := filepath.Join(skillsDir, "geoip.md")
		if _, err := os.Stat(geoipPath); err == nil {
			output.NewInfo("\nSpecial skill detected: geoip").Print()
			output.NewInfo("This skill provides IP geolocation and security features").Print()
			fmt.Printf("Enable geoip skill? (yes/no): ")
			var answer string
			fmt.Scanln(&answer)

			if strings.ToLower(answer) == "yes" || strings.ToLower(answer) == "y" {
				vault.EnabledSkills = append(vault.EnabledSkills, "geoip")
				enabled++
				output.NewSuccess("Skill geoip enabled").Print()
			} else {
				skipped++
				output.NewInfo("Skill geoip skipped").Print()
			}
		}

		vault.InstalledSkills = append(vault.InstalledSkills, "onboard")

		if err := skill.SaveVault(vault); err != nil {
			output.NewError(fmt.Sprintf("failed to save vault: %s", err.Error()), "VAULT_ERROR").Print()
			return
		}

		output.NewSuccess(map[string]interface{}{
			"onboard_completed": true,
			"skills_enabled":    enabled,
			"skills_skipped":    skipped,
			"total_skills":      len(files),
		}).Print()

		output.NewInfo("\nOnboard completed! You can now use NUX.").Print()
		output.NewInfo("To see available commands: nux --help").Print()
	},
}

func init() {
	rootCmd.AddCommand(onboardCmd)
}
