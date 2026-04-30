package commands

import (
	"fmt"

	"github.com/rsdenck/nux/internal/output"
	"github.com/rsdenck/nux/internal/vault"
	"github.com/spf13/cobra"
)

var onboardCmd = &cobra.Command{
	Use:   "onboard",
	Short: "Premium onboarding experience",
	Run: func(cmd *cobra.Command, args []string) {
		displayPremiumBanner()
		runOnboardFlow()
	},
}

func displayPremiumBanner() {
	banner := []string{
		"",
		"  ┌─────────────────────────────────────────────────────────────┐",
		"  │  NUX - Next-gen Unified eXecutor                         │",
		"  │  Linux CLI Moderno para Sysadmins                       │",
		"  └─────────────────────────────────────────────────────────────┘",
		"",
	}
	for _, line := range banner {
		fmt.Println(line)
	}
}

func runOnboardFlow() {
	fmt.Println("Selecione um perfil de configuracao:")
	fmt.Println()
	fmt.Println("  1) Minimal    - Apenas o essencial")
	fmt.Println("  2) Sysadmin    - Ferramentas de administracao")
	fmt.Println("  3) DevOps      - CI/CD, containers, k8s")
	fmt.Println("  4) Security    - Hardening, auditoria")
	fmt.Println("  5) Full        - Tudo incluido")
	fmt.Println("  6) Custom      - Personalizado")
	fmt.Println()

	v, err := vault.Load()
	if err != nil {
		v = vault.NewVault()
	}

	if v.Config == nil {
		v.Config = make(map[string]interface{})
	}
	v.Config["onboarded"] = true
	v.Config["onboard_date"] = "2026-04-30"

	if err := vault.Save(v); err != nil {
		output.NewError(fmt.Sprintf("failed to save: %s", err.Error()), "ONBOARD_SAVE_ERROR").Print()
		return
	}

	output.NewSuccess(map[string]interface{}{
		"status": "onboarded",
		"date":   "2026-04-30",
	}).Print()
}

func init() {
	rootCmd.AddCommand(onboardCmd)
}
