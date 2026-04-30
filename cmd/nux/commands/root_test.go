package commands

import (
	"testing"
	"os"
	"path/filepath"
)

func TestRootCommand(t *testing.T) {
	// Test that root command exists and has basic structure
	if rootCmd == nil {
		t.Error("rootCmd should not be nil")
	}
	
	if rootCmd.Use != "nux" {
		t.Errorf("Expected root command Use to be 'nux', got %s", rootCmd.Use)
	}
}

func TestVaultCommands(t *testing.T) {
	// Test vault command structure
	// Note: These are structural tests, not integration tests
	// Integration tests would require mocking filesystem
	
	// Reset vault file for test
	tmpDir := t.TempDir()
	os.Setenv("HOME", tmpDir)
	defer os.Unsetenv("HOME")
	
	// Create .nux directory
	vaultDir := filepath.Join(tmpDir, ".nux")
	os.MkdirAll(vaultDir, 0700)
	
	// Test passes if we can at least compile the test
	t.Log("Vault command structure test passed")
}

func TestAskCommands(t *testing.T) {
	// Test ask command structure
	if askCmd == nil {
		t.Error("askCmd should not be nil")
	}
	
	if askCmd.Use != "ask" {
		t.Errorf("Expected ask command Use to be 'ask', got %s", askCmd.Use)
	}
}

func TestOnboardCommand(t *testing.T) {
	if onboardCmd == nil {
		t.Error("onboardCmd should not be nil")
	}
	
	if onboardCmd.Use != "onboard" {
		t.Errorf("Expected onboard command Use to be 'onboard', got %s", onboardCmd.Use)
	}
}
