package pkg

import (
	"context"
	"fmt"
	"strings"

	"github.com/mascli/troncli/internal/core/adapter"
	"github.com/mascli/troncli/internal/core/domain"
	"github.com/mascli/troncli/internal/core/ports"
)

// UniversalPackageManager implements ports.PackageManager
type UniversalPackageManager struct {
	executor adapter.Executor
	profile  *domain.SystemProfile
}

// NewUniversalPackageManager creates a new package manager
func NewUniversalPackageManager(executor adapter.Executor, profile *domain.SystemProfile) *UniversalPackageManager {
	return &UniversalPackageManager{
		executor: executor,
		profile:  profile,
	}
}

// DetectManager is now handled by ProfileEngine, but we keep this for interface compliance or internal check
func (m *UniversalPackageManager) DetectManager() (string, error) {
	if m.profile.PackageManager == "" {
		return "", fmt.Errorf("no supported package manager detected in profile")
	}
	return m.profile.PackageManager, nil
}

// Install installs a package
func (m *UniversalPackageManager) Install(packageName string) error {
	ctx := context.Background()
	var args []string
	cmd := m.profile.PackageManager

	switch cmd {
	case "apt":
		args = []string{"install", "-y", packageName}
	case "dnf", "yum":
		args = []string{"install", "-y", packageName}
	case "pacman":
		args = []string{"-S", "--noconfirm", packageName}
	case "zypper":
		args = []string{"install", "-y", packageName}
	case "apk":
		args = []string{"add", packageName}
	default:
		return fmt.Errorf("unsupported package manager: %s", cmd)
	}

	_, err := m.executor.Exec(ctx, cmd, args...)
	return err
}

// Remove removes a package
func (m *UniversalPackageManager) Remove(packageName string) error {
	ctx := context.Background()
	var args []string
	cmd := m.profile.PackageManager

	switch cmd {
	case "apt":
		args = []string{"remove", "-y", packageName}
	case "dnf", "yum":
		args = []string{"remove", "-y", packageName}
	case "pacman":
		args = []string{"-Rs", "--noconfirm", packageName}
	case "zypper":
		args = []string{"remove", "-y", packageName}
	case "apk":
		args = []string{"del", packageName}
	default:
		return fmt.Errorf("unsupported package manager: %s", cmd)
	}

	_, err := m.executor.Exec(ctx, cmd, args...)
	return err
}

// Update updates the package list
func (m *UniversalPackageManager) Update() error {
	ctx := context.Background()
	var args []string
	cmd := m.profile.PackageManager

	switch cmd {
	case "apt":
		args = []string{"update"}
	case "dnf", "yum":
		args = []string{"check-update"}
	case "pacman":
		args = []string{"-Sy"}
	case "zypper":
		args = []string{"refresh"}
	case "apk":
		args = []string{"update"}
	default:
		return fmt.Errorf("unsupported package manager: %s", cmd)
	}

	_, err := m.executor.Exec(ctx, cmd, args...)
	return err
}

// Search searches for a package
func (m *UniversalPackageManager) Search(query string) ([]ports.PackageInfo, error) {
	ctx := context.Background()
	var args []string
	cmdName := m.profile.PackageManager

	switch cmdName {
	case "apt":
		args = []string{"search", query}
	case "dnf", "yum":
		args = []string{"search", query}
	case "pacman":
		args = []string{"-Ss", query}
	case "apk":
		args = []string{"search", "-v", query}
	default:
		return nil, fmt.Errorf("search not implemented for %s", cmdName)
	}

	res, err := m.executor.Exec(ctx, cmdName, args...)
	if err != nil {
		return []ports.PackageInfo{}, nil // Return empty on error/not found
	}

	return m.parseSearch(res.Stdout), nil
}

func (m *UniversalPackageManager) parseSearch(output string) []ports.PackageInfo {
	// Simplified parsing dispatch
	switch m.profile.PackageManager {
	case "apt":
		return m.parseAptSearch(output)
	case "pacman":
		return m.parsePacmanSearch(output)
	case "dnf", "yum":
		return m.parseDnfSearch(output)
	case "zypper":
		return m.parseZypperSearch(output)
	case "apk":
		return m.parseApkSearch(output)
	default:
		// Default simple line parser
		return []ports.PackageInfo{}
	}
}

func (m *UniversalPackageManager) parseDnfSearch(output string) []ports.PackageInfo {
	var pkgs []ports.PackageInfo
	lines := strings.Split(output, "\n")
	
	for _, line := range lines {
		// DNF format: package.arch   version   repo
		// Skip headers or metadata
		if strings.Contains(line, "=") || strings.Contains(line, "Matched:") {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) >= 3 {
			nameParts := strings.Split(parts[0], ".")
			name := nameParts[0]
			version := parts[1]
			
			pkgs = append(pkgs, ports.PackageInfo{
				Name:    name,
				Version: version,
				Manager: m.profile.PackageManager,
			})
		}
	}
	return pkgs
}

func (m *UniversalPackageManager) parseZypperSearch(output string) []ports.PackageInfo {
	var pkgs []ports.PackageInfo
	lines := strings.Split(output, "\n")
	
	// S | Name | Summary | Type
	// --+------+---------+-------
	// i | pkg  | desc    | package
	
	for _, line := range lines {
		if strings.HasPrefix(line, "S |") || strings.HasPrefix(line, "--") {
			continue
		}
		
		parts := strings.Split(line, "|")
		if len(parts) >= 4 {
			status := strings.TrimSpace(parts[0])
			name := strings.TrimSpace(parts[1])
			desc := strings.TrimSpace(parts[2])
			
			installed := (status == "i" || status == "i+")
			
			pkgs = append(pkgs, ports.PackageInfo{
				Name:        name,
				Description: desc,
				Installed:   installed,
				Manager:     "zypper",
			})
		}
	}
	return pkgs
}

func (m *UniversalPackageManager) parseAptSearch(output string) []ports.PackageInfo {
	var pkgs []ports.PackageInfo
	lines := strings.Split(output, "\n")
	var currentPkg *ports.PackageInfo

	for _, line := range lines {
		if strings.Contains(line, "/") && strings.Contains(line, "now") { // apt search output format varies
			parts := strings.Fields(line)
			if len(parts) > 0 {
				name := strings.Split(parts[0], "/")[0]
				currentPkg = &ports.PackageInfo{Name: name, Manager: "apt"}
				if len(parts) > 1 {
					currentPkg.Version = parts[1]
				}
				if strings.Contains(line, "[installed]") {
					currentPkg.Installed = true
				}
				pkgs = append(pkgs, *currentPkg)
			}
		}
	}
	return pkgs
}

func (m *UniversalPackageManager) parsePacmanSearch(output string) []ports.PackageInfo {
	var pkgs []ports.PackageInfo
	lines := strings.Split(output, "\n")
	var currentPkg *ports.PackageInfo

	for _, line := range lines {
		if !strings.HasPrefix(line, "    ") && strings.TrimSpace(line) != "" {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				name := parts[0]
				if strings.Contains(name, "/") {
					name = strings.Split(name, "/")[1]
				}
				version := parts[1]
				installed := strings.Contains(line, "[installed]")

				currentPkg = &ports.PackageInfo{
					Name:      name,
					Version:   version,
					Installed: installed,
					Manager:   "pacman",
				}
				pkgs = append(pkgs, *currentPkg)
			}
		}
	}
	return pkgs
}

func (m *UniversalPackageManager) parseApkSearch(output string) []ports.PackageInfo {
	var pkgs []ports.PackageInfo
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		pkgs = append(pkgs, ports.PackageInfo{
			Name:    line,
			Manager: "apk",
		})
	}
	return pkgs
}
