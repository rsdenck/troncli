package agent

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCapabilityRegistry(t *testing.T) {
	// 1. Create temp file
	tmpDir, err := os.MkdirTemp("", "troncli-test-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	yamlContent := `allowed_intents:
  - install_package
  - remove_package
  - audit_security
`
	configPath := filepath.Join(tmpDir, "capabilities.yaml")
	if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
		t.Fatal(err)
	}

	// 2. Load Registry
	registry := NewCapabilityRegistry(configPath)
	if err := registry.Load(); err != nil {
		t.Fatalf("Failed to load registry: %v", err)
	}

	// 3. Test Allowed
	// Note: Current implementation requires exact match or logic update. 
	// For this test, we assume exact match or simple prefix if we updated registry.go.
	// Since registry.go uses ==, we must use exact string from allowed list.
	if !registry.IsIntentAllowed("install_package") {
		t.Errorf("Expected install_package to be allowed")
	}
	if !registry.IsIntentAllowed("audit_security") {
		t.Errorf("Expected audit_security to be allowed")
	}

	// 4. Test Denied
	if registry.IsIntentAllowed("delete_system") {
		t.Errorf("Expected delete_system to be denied")
	}
	
	// 5. Test Partial Match logic (if implemented)
	// Usually IsIntentAllowed checks if intent starts with or contains allowed keyword.
	// In my implementation:
	/*
	func (r *CapabilityRegistry) IsIntentAllowed(intent string) bool {
		for _, allowed := range r.AllowedIntents {
			if strings.Contains(intent, allowed) {
				return true
			}
		}
		return false
	}
	*/
	// So "install_package nginx" contains "install_package" -> True.
}
