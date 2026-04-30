package skill

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const vaultDir = ".skills"
const vaultFile = ".nux.json"

type Vault struct {
	Version       string            `json:"version"`
	InstalledSkills []string        `json:"installed_skills"`
	EnabledSkills  []string        `json:"enabled_skills"`
	APIKeys       map[string]string `json:"api_keys"`
	Ollama        OllamaConfig     `json:"ollama"`
	VaultMode     bool             `json:"vault_mode"`
}

type OllamaConfig struct {
	Host    string `json:"host"`
	Model   string `json:"model"`
	Enabled bool   `json:"enabled"`
}

func getVaultPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(home, vaultDir, vaultFile), nil
}

func LoadVault() (*Vault, error) {
	path, err := getVaultPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return defaultVault(), nil
		}
		return nil, fmt.Errorf("failed to read vault: %w", err)
	}

	var v Vault
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, fmt.Errorf("failed to parse vault: %w", err)
	}
	return &v, nil
}

func SaveVault(v *Vault) error {
	path, err := getVaultPath()
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create vault directory: %w", err)
	}

	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal vault: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write vault: %w", err)
	}
	return nil
}

func defaultVault() *Vault {
	return &Vault{
		Version:        "1.0.0",
		InstalledSkills: []string{},
		EnabledSkills:  []string{},
		APIKeys:        make(map[string]string),
		Ollama: OllamaConfig{
			Host:    "http://localhost:11434",
			Model:   "qwen3-coder",
			Enabled: false,
		},
		VaultMode: true,
	}
}
