package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

const vaultDir = ".nux"
const vaultFile = "vault.json"

type Vault struct {
	Version   string                 `json:"version"`
	Installed map[string]SkillStatus `json:"installed"`
	Enabled   map[string]bool       `json:"enabled"`
	APIKeys   map[string]string     `json:"api_keys"`
	Tokens    map[string]TokenInfo  `json:"tokens"`
	Config    map[string]interface{} `json:"config"`
	mu        sync.RWMutex
}

type SkillStatus struct {
	InstalledAt string `json:"installed_at"`
	Version     string `json:"version"`
	Status      string `json:"status"`
}

type TokenInfo struct {
	Value    string `json:"value"`
	Expires  string `json:"expires,omitempty"`
	Provider string `json:"provider"`
}

func getVaultPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(home, vaultDir, vaultFile), nil
}

func Load() (*Vault, error) {
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

func Save(v *Vault) error {
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

func NewVault() *Vault {
	return &Vault{
		Version:   "1.0.0",
		Installed: make(map[string]SkillStatus),
		Enabled:   make(map[string]bool),
		APIKeys:   make(map[string]string),
		Tokens:    make(map[string]TokenInfo),
		Config:    make(map[string]interface{}),
	}
}

func defaultVault() *Vault {
	return NewVault()
}

func (v *Vault) SetAPIKey(service, key string) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.APIKeys[service] = key
}

func (v *Vault) GetAPIKey(service string) (string, bool) {
	v.mu.RLock()
	defer v.mu.RUnlock()
	key, ok := v.APIKeys[service]
	return key, ok
}

func (v *Vault) EnableSkill(name string) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.Enabled[name] = true
	v.Installed[name] = SkillStatus{
		InstalledAt: "now",
		Status:      "active",
	}
}

func (v *Vault) DisableSkill(name string) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.Enabled[name] = false
	if s, ok := v.Installed[name]; ok {
		s.Status = "inactive"
		v.Installed[name] = s
	}
}
