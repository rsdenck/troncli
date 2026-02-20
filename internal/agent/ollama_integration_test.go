package agent

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestOllamaIntegration connects to a real Ollama instance
// Run with: OLLAMA_TEST_HOST=http://192.168.130.25:11434 go test -v ./internal/agent/ -run TestOllamaIntegration
func TestOllamaIntegration(t *testing.T) {
	// 1. Setup Environment
	targetURL := os.Getenv("OLLAMA_TEST_HOST")
	if targetURL == "" {
		t.Skip("Skipping integration test: OLLAMA_TEST_HOST not set")
	}
	t.Setenv("OLLAMA_HOST", targetURL)

	// 2. Setup Capability Registry
	tmpDir := t.TempDir()
	yamlContent := `allowed_intents:
  - test_intent
`
	configPath := filepath.Join(tmpDir, "capabilities.yaml")
	if err := os.WriteFile(configPath, []byte(yamlContent), 0600); err != nil {
		t.Fatal(err)
	}

	// 3. Initialize Adapter
	// Using 'nomic-embed-text:latest' - it exists on the server but doesn't support generate.
	// This allows us to verify connectivity quickly without waiting for model loading.
	adapter, err := NewOllamaAdapter("nomic-embed-text:latest", configPath)
	if err != nil {
		t.Fatalf("Failed to create adapter: %v", err)
	}

	// 4. Test Connection
	timeout := 10 * time.Second // Should be very fast
	fmt.Printf("Setting context timeout to: %v\n", timeout)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	fmt.Printf("Attempting to connect to Ollama at %s...\n", targetURL)
	start := time.Now()

	response, err := adapter.SendPrompt(ctx, "Hello")
	duration := time.Since(start)

	fmt.Printf("Request took: %v\n", duration)

	if err != nil {
		// If we get the specific error about the model not supporting generate, it means we connected!
		errStr := err.Error()
		if strings.Contains(errStr, "does not support generate") || strings.Contains(errStr, "400") {
			fmt.Printf("SUCCESS: Connected to Ollama and received expected error: %v\n", err)
			return
		}
		// Any other error (like timeout) is a failure
		t.Fatalf("Ollama connection failed: %v", err)
	}

	fmt.Printf("Ollama Response: %s\n", response)
}
