package commands

import (
	"fmt"

	"github.com/rsdenck/nux/internal/output"
	"github.com/rsdenck/nux/internal/vault"
	"github.com/spf13/cobra"
)

var vaultCmd = &cobra.Command{
	Use:   "vault",
	Short: "Manage NUX vault for secrets and API keys",
	Long:  `Secure storage for API keys, tokens, and configuration. Vault file is stored at ~/.nux/vault.json with 0600 permissions.`,
}

var vaultShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show vault contents (API keys masked)",
	Run: func(cmd *cobra.Command, args []string) {
		v, err := vault.Load()
		if err != nil {
			output.NewError(fmt.Sprintf("failed to load vault: %s", err.Error()), "VAULT_ERROR").Print()
			return
		}

		result := map[string]interface{}{
			"version":   v.Version,
			"installed": len(v.Installed),
			"enabled":   len(v.Enabled),
			"api_keys":  maskKeys(v.APIKeys),
			"tokens":    len(v.Tokens),
		}

		output.NewSuccess(result).WithMessage("NUX Vault").Print()
	},
}

func maskKeys(keys map[string]string) map[string]string {
	masked := make(map[string]string)
	for k, v := range keys {
		if len(v) > 4 {
			masked[k] = v[:4] + "****"
		} else {
			masked[k] = "****"
		}
	}
	return masked
}

var vaultSetKeyCmd = &cobra.Command{
	Use:   "set-key <service> <key>",
	Short: "Set API key for a service",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		service := args[0]
		key := args[1]

		v, err := vault.Load()
		if err != nil {
			v = vault.NewVault()
		}

		v.SetAPIKey(service, key)

		if err := vault.Save(v); err != nil {
			output.NewError(fmt.Sprintf("failed to save vault: %s", err.Error()), "VAULT_SAVE_ERROR").Print()
			return
		}

		output.NewSuccess(map[string]interface{}{
			"service": service,
			"status":  "key saved",
		}).Print()
	},
}

var vaultGetKeyCmd = &cobra.Command{
	Use:   "get-key <service>",
	Short: "Get API key for a service",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		service := args[0]

		v, err := vault.Load()
		if err != nil {
			output.NewError(fmt.Sprintf("failed to load vault: %s", err.Error()), "VAULT_ERROR").Print()
			return
		}

		key, ok := v.GetAPIKey(service)
		if !ok {
			output.NewError(fmt.Sprintf("no key found for service: %s", service), "VAULT_KEY_NOT_FOUND").Print()
			return
		}

		output.NewSuccess(map[string]interface{}{
			"service": service,
			"key":     key,
		}).Print()
	},
}

func init() {
	vaultCmd.AddCommand(vaultShowCmd)
	vaultCmd.AddCommand(vaultSetKeyCmd)
	vaultCmd.AddCommand(vaultGetKeyCmd)
	rootCmd.AddCommand(vaultCmd)
}
