package plugin

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mascli/troncli/internal/core/adapter"
	"github.com/mascli/troncli/internal/core/domain"
	"github.com/mascli/troncli/internal/core/ports"
)

// Simple internal registry for Phase 4
var pluginRegistry = map[string]string{
	"arch":   "https://raw.githubusercontent.com/rsdenck/troncli-plugins/main/arch/arch-helper.sh",
	"docker": "https://raw.githubusercontent.com/rsdenck/troncli-plugins/main/docker/docker-helper.sh",
	"k8s":    "https://raw.githubusercontent.com/rsdenck/troncli-plugins/main/k8s/k8s-helper.sh",
}

type UniversalPluginManager struct {
	executor  adapter.Executor
	profile   *domain.SystemProfile
	pluginDir string
}

func NewUniversalPluginManager(executor adapter.Executor, profile *domain.SystemProfile) (*UniversalPluginManager, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	pluginDir := filepath.Join(home, ".troncli", "plugins")
	if err := os.MkdirAll(pluginDir, 0755); err != nil {
		return nil, err
	}

	return &UniversalPluginManager{
		executor:  executor,
		profile:   profile,
		pluginDir: pluginDir,
	}, nil
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

	// Check registry first
	if regUrl, ok := pluginRegistry[nameOrUrl]; ok {
		url = regUrl
		name = nameOrUrl
	} else {
		// Assume it's a direct URL or path
		url = nameOrUrl
		name = filepath.Base(url)
	}

	destPath := filepath.Join(m.pluginDir, name)
	fmt.Printf("Installing plugin '%s' from %s...\n", name, url)

	if strings.HasPrefix(url, "http") {
		client := &http.Client{Timeout: 30 * time.Second}
		resp, err := client.Get(url)
		if err != nil {
			return fmt.Errorf("failed to download plugin: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("bad status: %s", resp.Status)
		}

		out, err := os.Create(destPath)
		if err != nil {
			return err
		}
		defer out.Close()

		if _, err := io.Copy(out, resp.Body); err != nil {
			return err
		}
	} else {
		// Local file
		src, err := os.Open(url)
		if err != nil {
			return fmt.Errorf("failed to open local plugin: %w", err)
		}
		defer src.Close()

		out, err := os.Create(destPath)
		if err != nil {
			return err
		}
		defer out.Close()

		if _, err := io.Copy(out, src); err != nil {
			return err
		}
	}

	if err := os.Chmod(destPath, 0755); err != nil {
		return fmt.Errorf("failed to make plugin executable: %w", err)
	}

	fmt.Printf("Plugin '%s' installed successfully at %s\n", name, destPath)
	return nil
}

func (m *UniversalPluginManager) RemovePlugin(name string) error {
	path := filepath.Join(m.pluginDir, name)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("plugin '%s' not found", name)
	}
	return os.Remove(path)
}

func (m *UniversalPluginManager) ExecutePlugin(name string, args []string) error {
	path := filepath.Join(m.pluginDir, name)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("plugin '%s' not found", name)
	}

	// Use the executor to run the plugin
	// We use context.Background() for now as the interface doesn't pass context yet
	// In a real scenario, we should update the interface to accept context
	_, err := m.executor.Exec(context.Background(), path, args...)
	return err
}
