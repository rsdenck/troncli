package plugin

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/mascli/troncli/internal/core/adapter"
	"github.com/mascli/troncli/internal/core/domain"
	"github.com/mascli/troncli/internal/core/ports"
)

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
			// Basic info for now
			plugins = append(plugins, ports.Plugin{
				Name:        entry.Name(),
				Description: "Plugin externo",
				Version:     info.ModTime().Format("2006-01-02"),
				Path:        filepath.Join(m.pluginDir, entry.Name()),
			})
		}
	}
	return plugins, nil
}

func (m *UniversalPluginManager) InstallPlugin(urlOrPath string) error {
	// If URL, download. If path, copy.
	name := filepath.Base(urlOrPath)
	destPath := filepath.Join(m.pluginDir, name)

	if strings.HasPrefix(urlOrPath, "http") {
		resp, err := http.Get(urlOrPath)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

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
		src, err := os.Open(urlOrPath)
		if err != nil {
			return err
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

	return os.Chmod(destPath, 0755)
}

func (m *UniversalPluginManager) RemovePlugin(name string) error {
	return os.Remove(filepath.Join(m.pluginDir, name))
}
