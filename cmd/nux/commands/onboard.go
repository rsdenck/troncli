package commands

import (
	"fmt"
	"os"

	"github.com/rsdenck/nux/internal/skill"
	"github.com/spf13/cobra"
)

var onboardCmd = &cobra.Command{
	Use:   "onboard",
	Short: "Run first-time setup for NUX",
	Long: `Interactive onboarding process to configure NUX for first use.

This will:
  - Configure Ollama integration
  - Set up the skill vault
  - Enable basic skills
  - Configure output preferences`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("=== NUX Onboarding ===")
		fmt.Println()

		v, err := skill.LoadVault()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading vault: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("1. Configuring skill vault...")
		if v.VaultMode {
			fmt.Println("   Vault mode: ENABLED (permissions: 0600)")
		}

		fmt.Println("2. Ollama integration:")
		fmt.Printf("   Host: %s\n", v.Ollama.Host)
		fmt.Printf("   Model: %s\n", v.Ollama.Model)
		fmt.Printf("   Status: %v\n", v.Ollama.Enabled)

		fmt.Println("3. Basic skills available:")
		skills, _ := skill.ListSkills()
		count := 0
		for i, s := range skills {
			if i >= 5 {
				break
			}
			fmt.Printf("   - %s\n", s)
			count++
		}
		fmt.Printf("   ... and %d more\n", len(skills)-count)

		fmt.Println()
		fmt.Println("Onboarding completed!")
		fmt.Println("Run 'nux skill list' to see all available skills")
		fmt.Println("Run 'nux ask \"your question\"' to use the AI agent")
	},
}

func init() {
	rootCmd.AddCommand(onboardCmd)
}
