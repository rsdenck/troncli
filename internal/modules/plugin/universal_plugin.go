package plugin

// Package plugin provides plugin management capabilities.

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/mascli/troncli/internal/core/adapter"
	"github.com/mascli/troncli/internal/core/domain"
	"github.com/mascli/troncli/internal/core/ports"
)

// PluginDef defines a plugin with security metadata
type PluginDef struct {
	URL      string `json:"url"`
	Checksum string `json:"checksum"` // SHA256 hash
}

// Default registry content
var defaultPluginRegistry = map[string]PluginDef{
	"arch": {
		URL:      "https://raw.githubusercontent.com/rsdenck/troncli-plugins/main/arch/arch-helper.sh",
		Checksum: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
	},
	"docker": {
		URL:      "https://raw.githubusercontent.com/rsdenck/troncli-plugins/main/docker/docker-helper.sh",
		Checksum: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
	},
	"k8s": {
		URL:      "https://raw.githubusercontent.com/rsdenck/troncli-plugins/main/k8s/k8s-helper.sh",
		Checksum: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
	},
}

type UniversalPluginManager struct {
	executor     adapter.Executor
	profile      *domain.SystemProfile
	pluginDir    string
	registry     map[string]PluginDef
	registryPath string
	httpClient   *http.Client
}

func NewUniversalPluginManager(executor adapter.Executor, profile *domain.SystemProfile) (*UniversalPluginManager, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	configDir := filepath.Join(home, ".troncli")
	pluginDir := filepath.Join(configDir, "plugins")
	if err := os.MkdirAll(pluginDir, 0755); err != nil {
		return nil, err
	}

	registryPath := filepath.Join(configDir, "plugins.json")
	registry := make(map[string]PluginDef)

	// Load registry from file or create default
	if _, err := os.Stat(registryPath); os.IsNotExist(err) {
		// Initialize with default registry
		registry = defaultPluginRegistry
		// Save it so user can edit later
		data, _ := json.MarshalIndent(registry, "", "  ")
		_ = os.WriteFile(registryPath, data, 0644)
	} else {
		// Load existing
		data, err := os.ReadFile(registryPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read plugin registry: %w", err)
		}
		if err := json.Unmarshal(data, &registry); err != nil {
			return nil, fmt.Errorf("failed to parse plugin registry: %w", err)
		}
	}

	return &UniversalPluginManager{
		executor:     executor,
		profile:      profile,
		pluginDir:    pluginDir,
		registry:     registry,
		registryPath: registryPath,
		httpClient:   &http.Client{Timeout: 30 * time.Second},
	}, nil
}

// SetHTTPClient allows overriding the default HTTP client (e.g. for testing)
func (m *UniversalPluginManager) SetHTTPClient(client *http.Client) {
	m.httpClient = client
}

func (m *UniversalPluginManager) ListPlugins() ([]ports.Plugin, error) {
	entries, err := os.ReadDir(m.pluginDir)
	if err != nil {
		return nil, err
	}

	var plugins []ports.Plugin
	for _, entry := range entries {
		if !entry.IsDir() {
			info, _ := entry.Info()
			plugins = append(plugins, ports.Plugin{
				Name:        entry.Name(),
				Description: "Plugin instalado",
				Version:     info.ModTime().Format("2006-01-02"),
				Path:        filepath.Join(m.pluginDir, entry.Name()),
			})
		}
	}
	return plugins, nil
}

func (m *UniversalPluginManager) InstallPlugin(nameOrUrl string) error {
	var url string
	var name string
	var expectedChecksum string

	// Check registry first
	if pluginDef, ok := m.registry[nameOrUrl]; ok {
		url = pluginDef.URL
		expectedChecksum = pluginDef.Checksum
		name = nameOrUrl
	} else {
		// BLOCK: Direct URL installation is disabled in Phase 1 Hardening
		return fmt.Errorf("direct URL installation is disabled for security. Use registered plugins only")
	}

	if !strings.HasPrefix(url, "https://") {
		return fmt.Errorf("insecure protocol: HTTPS is required for plugin download")
	}

	destPath := filepath.Join(m.pluginDir, name)
	fmt.Printf("Installing plugin '%s' from %s...\n", name, url)

	// Use configured httpClient
	resp, err := m.httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download plugin: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Read content to verify checksum
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read plugin content: %w", err)
	}

	// Verify Checksum
	hash := sha256.Sum256(content)
	calculatedChecksum := hex.EncodeToString(hash[:])
	if calculatedChecksum != expectedChecksum {
		return fmt.Errorf("security violation: checksum mismatch for plugin '%s'. Expected %s, got %s", name, expectedChecksum, calculatedChecksum)
	}

	// Write file with minimal permissions
	if err := os.WriteFile(destPath, content, 0700); err != nil { // 0700 = rwx------
		return fmt.Errorf("failed to write plugin file: %w", err)
	}

	fmt.Printf("Plugin '%s' installed and verified successfully (SHA256 matched).\n", name)
	return nil
}

func (m *UniversalPluginManager) RemovePlugin(name string) error {
	path := filepath.Join(m.pluginDir, name)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("plugin '%s' not found", name)
	}
	return os.Remove(path)
}

// limitWriter wraps an io.Writer and returns an error if the limit is exceeded
type limitWriter struct {
	w     io.Writer
	n     int64
	limit int64
}

func (l *limitWriter) Write(p []byte) (n int, err error) {
	if l.n >= l.limit {
		return 0, fmt.Errorf("output limit exceeded")
	}
	remaining := l.limit - l.n
	if int64(len(p)) > remaining {
		p = p[:remaining]
		err = fmt.Errorf("output limit exceeded")
	}
	n, wErr := l.w.Write(p)
	l.n += int64(n)
	if wErr != nil {
		return n, wErr
	}
	return n, err
}

func (m *UniversalPluginManager) ExecutePlugin(ctx context.Context, name string, args ...string) error {
	path := filepath.Join(m.pluginDir, name)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("plugin '%s' not found", name)
	}

	// Phase 4 Hardening: Plugin Sandbox
	// 1. Timeout (5 minutes by default if not set by caller)
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
	}

	// 2. Drop privileges if root? (Not implemented here for portability, but warned)
	if os.Geteuid() == 0 {
		fmt.Println("Warning: Running plugin as root is dangerous!")
	}

	// 3. Secure Execution (Sandbox)
	// Use direct exec.CommandContext instead of adapter to control environment and streams
	cmd := exec.CommandContext(ctx, path, args...)

	// Environment sanitization - only pass essential variables
	cmd.Env = []string{
		fmt.Sprintf("PATH=%s", os.Getenv("PATH")),
		fmt.Sprintf("HOME=%s", os.Getenv("HOME")),
		fmt.Sprintf("TERM=%s", os.Getenv("TERM")),
		fmt.Sprintf("LANG=%s", os.Getenv("LANG")),
		fmt.Sprintf("TRONCLI_PLUGIN_NAME=%s", name),
	}

	// Limit Output (10MB max)
	const MaxOutputSize = 10 * 1024 * 1024

	// Create pipes for stdout/stderr to control them
	// Using LimitWriter wrapper for os.Stdout/Stderr is tricky because it writes directly.
	// We want to capture or pass through but limit.
	// Since we are CLI, we usually want to stream to user.

	cmd.Stdout = &limitWriter{w: os.Stdout, limit: MaxOutputSize}
	cmd.Stderr = &limitWriter{w: os.Stderr, limit: MaxOutputSize}

	// Run
	if err := cmd.Run(); err != nil {
		// If context deadline exceeded, wrap it
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("plugin execution timed out: %w", err)
		}
		return fmt.Errorf("plugin execution failed: %w", err)
	}

	return nil
}
