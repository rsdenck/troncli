package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDefaultVault(t *testing.T) {
	v := defaultVault()
	
	if v.Version != "1.0.0" {
		t.Errorf("Expected version 1.0.0, got %s", v.Version)
	}
	
	if v.Installed == nil {
		t.Error("Installed map should be initialized")
	}
	
	if v.Enabled == nil {
		t.Error("Enabled map should be initialized")
	}
	
	if v.APIKeys == nil {
		t.Error("APIKeys map should be initialized")
	}
}

func TestSetAndGetAPIKey(t *testing.T) {
	v := defaultVault()
	
	v.SetAPIKey("openai", "sk-test123")
	
	key, ok := v.GetAPIKey("openai")
	if !ok {
		t.Error("Expected to find openai key")
	}
	
	if key != "sk-test123" {
		t.Errorf("Expected sk-test123, got %s", key)
	}
	
	_, ok = v.GetAPIKey("nonexistent")
	if ok {
		t.Error("Should not find nonexistent key")
	}
}

func TestEnableAndDisableSkill(t *testing.T) {
	v := defaultVault()
	
	v.EnableSkill("docker")
	if !v.Enabled["docker"] {
		t.Error("Docker should be enabled")
	}
	
	if v.Installed["docker"].Status != "active" {
		t.Error("Docker status should be active")
	}
	
	v.DisableSkill("docker")
	if v.Enabled["docker"] {
		t.Error("Docker should be disabled")
	}
}

func TestSaveAndLoad(t *testing.T) {
	// Create temp dir
	tmpDir := t.TempDir()
	
	// Override home dir for test
	os.Setenv("HOME", tmpDir)
	defer os.Unsetenv("HOME")
	
	// Create and save vault
	v := defaultVault()
	v.SetAPIKey("test", "key123")
	v.EnableSkill("bash")
	
	if err := Save(v); err != nil {
		t.Fatalf("Failed to save vault: %v", err)
	}
	
	// Check file exists
	vaultPath := filepath.Join(tmpDir, ".nux", "vault.json")
	if _, err := os.Stat(vaultPath); os.IsNotExist(err) {
		t.Error("Vault file should exist")
	}
	
	// Load vault
	loaded, err := Load()
	if err != nil {
		t.Fatalf("Failed to load vault: %v", err)
	}
	
	if loaded.Version != v.Version {
		t.Errorf("Version mismatch: expected %s, got %s", v.Version, loaded.Version)
	}
}
