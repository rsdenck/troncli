package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rsdenck/nux/internal/output"
	"github.com/spf13/cobra"
)

var pluginCmd = &cobra.Command{
	Use:   "plugin",
	Short: "Plugin management",
	Long:  `Manage NUX plugins and skills.`,
}

var pluginListCmd = &cobra.Command{
	Use:   "list",
	Short: "List plugins",
	Run: func(cmd *cobra.Command, args []string) {
		// List plugins from skills directory
		skillsDir := "skills"
		if _, err := os.Stat(skillsDir); os.IsNotExist(err) {
			// Try absolute path
			skillsDir = "/opt/cli/nux/skills"
		}

		files, err := os.ReadDir(skillsDir)
		if err != nil {
			output.NewError(fmt.Sprintf("failed to read skills directory: %s", err.Error()), "PLUGIN_LIST_ERROR").Print()
			return
		}

		items := []map[string]interface{}{}
		for _, file := range files {
			if file.IsDir() || !strings.HasSuffix(file.Name(), ".md") {
				continue
			}

			skillName := strings.TrimSuffix(file.Name(), ".md")
			item := map[string]interface{}{
				"name":      skillName,
				"type":      "skill",
				"file":      file.Name(),
				"installed": checkSkillInstalled(skillName),
			}
			items = append(items, item)
		}

		output.NewList(items, len(items)).WithMessage("Plugin/Skill list").Print()
	},
}

var pluginInfoCmd = &cobra.Command{
	Use:   "info <plugin>",
	Short: "Show plugin information",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pluginName := args[0]

		// Try to read skill file
		skillFile := filepath.Join("skills", pluginName+".md")
		if _, err := os.Stat(skillFile); err != nil {
			skillFile = filepath.Join("/opt/cli/nux/skills", pluginName+".md")
		}

		content, err := os.ReadFile(skillFile)
		if err != nil {
			output.NewError(fmt.Sprintf("plugin not found: %s", pluginName), "PLUGIN_NOT_FOUND").Print()
			return
		}

		output.NewSuccess(map[string]interface{}{
			"name":    pluginName,
			"content": strings.TrimSpace(string(content)),
		}).Print()
	},
}

func checkSkillInstalled(skillName string) bool {
	// Check if skill is in vault as installed
	vaultPath := os.Getenv("HOME") + "/.skills/.nux.json"
	if _, err := os.Stat(vaultPath); err != nil {
		return false
	}

	data, err := os.ReadFile(vaultPath)
	if err != nil {
		return false
	}

	var vault map[string]interface{}
	if err := json.Unmarshal(data, &vault); err != nil {
		return false
	}

	installed, ok := vault["installed_skills"].([]interface{})
	if !ok {
		return false
	}

	for _, s := range installed {
		if s == skillName {
			return true
		}
	}

	return false
}

func init() {
	pluginCmd.AddCommand(pluginListCmd)
	pluginCmd.AddCommand(pluginInfoCmd)
	rootCmd.AddCommand(pluginCmd)
}
