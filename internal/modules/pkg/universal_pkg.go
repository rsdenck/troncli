package pkg

// Package pkg provides package management capabilities.

import (
	"context"
	"fmt"
	"strings"

	"github.com/mascli/troncli/internal/core/adapter"
	"github.com/mascli/troncli/internal/core/domain"
	"github.com/mascli/troncli/internal/core/ports"
	"github.com/mascli/troncli/internal/policy"
)

const (
	pacmanCmd = "pacman"
	dnfCmd    = "dnf"
)

// UniversalPackageManager implements ports.PackageManager
type UniversalPackageManager struct {
	executor     adapter.Executor
	profile      *domain.SystemProfile
	policyEngine *policy.PolicyEngine
}

// NewUniversalPackageManager creates a new package manager
func NewUniversalPackageManager(executor adapter.Executor, profile *domain.SystemProfile, policyEngine *policy.PolicyEngine) *UniversalPackageManager {
	return &UniversalPackageManager{
		executor:     executor,
		profile:      profile,
		policyEngine: policyEngine,
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

	// Build command string for policy checking
	var fullCmd string
	switch cmd {
	case "apt":
		args = []string{"install", "-y", packageName}
		fullCmd = fmt.Sprintf("%s install -y %s", cmd, packageName)
	case dnfCmd, "yum":
		args = []string{"install", "-y", packageName}
		fullCmd = fmt.Sprintf("%s install -y %s", cmd, packageName)
	case pacmanCmd:
		args = []string{"-S", "--noconfirm", packageName}
		fullCmd = fmt.Sprintf("%s -S --noconfirm %s", cmd, packageName)
	case "zypper":
		args = []string{"install", "-y", packageName}
		fullCmd = fmt.Sprintf("%s install -y %s", cmd, packageName)
	case "apk":
		args = []string{"add", packageName}
		fullCmd = fmt.Sprintf("%s add %s", cmd, packageName)
	case "portage":
		cmd = "emerge"
		args = []string{packageName}
		fullCmd = fmt.Sprintf("emerge %s", packageName)
	case "xbps":
		cmd = "xbps-install"
		args = []string{"-y", packageName}
		fullCmd = fmt.Sprintf("xbps-install -y %s", packageName)
	default:
		return fmt.Errorf("unsupported package manager: %s", cmd)
	}

	// Check command against policy engine
	if err := m.policyEngine.CheckCommand(fullCmd); err != nil {
		return fmt.Errorf("policy engine blocked execution: %w", err)
	}

	_, err := m.executor.Exec(ctx, cmd, args...)
	return err
}

// Remove removes a package
func (m *UniversalPackageManager) Remove(packageName string) error {
	ctx := context.Background()
	var args []string
	cmd := m.profile.PackageManager

	// Build command string for policy checking
	var fullCmd string
	switch cmd {
	case "apt":
		args = []string{"remove", "-y", packageName}
		fullCmd = fmt.Sprintf("%s remove -y %s", cmd, packageName)
	case dnfCmd, "yum":
		args = []string{"remove", "-y", packageName}
		fullCmd = fmt.Sprintf("%s remove -y %s", cmd, packageName)
	case pacmanCmd:
		args = []string{"-Rs", "--noconfirm", packageName}
		fullCmd = fmt.Sprintf("%s -Rs --noconfirm %s", cmd, packageName)
	case "zypper":
		args = []string{"remove", "-y", packageName}
		fullCmd = fmt.Sprintf("%s remove -y %s", cmd, packageName)
	case "apk":
		args = []string{"del", packageName}
		fullCmd = fmt.Sprintf("%s del %s", cmd, packageName)
	case "portage":
		cmd = "emerge"
		args = []string{"--unmerge", packageName}
		fullCmd = fmt.Sprintf("emerge --unmerge %s", packageName)
	case "xbps":
		cmd = "xbps-remove"
		args = []string{"-y", packageName}
		fullCmd = fmt.Sprintf("xbps-remove -y %s", packageName)
	default:
		return fmt.Errorf("unsupported package manager: %s", cmd)
	}

	// Check command against policy engine
	if err := m.policyEngine.CheckCommand(fullCmd); err != nil {
		return fmt.Errorf("policy engine blocked execution: %w", err)
	}

	_, err := m.executor.Exec(ctx, cmd, args...)
	return err
}

// Update updates the package list
func (m *UniversalPackageManager) Update() error {
	ctx := context.Background()
	var args []string
	cmd := m.profile.PackageManager

	// Build command string for policy checking
	var fullCmd string
	switch cmd {
	case "apt":
		args = []string{"update"}
		fullCmd = fmt.Sprintf("%s update", cmd)
	case dnfCmd, "yum":
		args = []string{"check-update"}
		fullCmd = fmt.Sprintf("%s check-update", cmd)
	case pacmanCmd:
		args = []string{"-Sy"}
		fullCmd = fmt.Sprintf("%s -Sy", cmd)
	case "zypper":
		args = []string{"refresh"}
		fullCmd = fmt.Sprintf("%s refresh", cmd)
	case "apk":
		args = []string{"update"}
		fullCmd = fmt.Sprintf("%s update", cmd)
	case "portage":
		cmd = "emerge"
		args = []string{"--sync"}
		fullCmd = "emerge --sync"
	case "xbps":
		cmd = "xbps-install"
		args = []string{"-S"}
		fullCmd = "xbps-install -S"
	default:
		return fmt.Errorf("unsupported package manager: %s", cmd)
	}

	// Check command against policy engine
	if err := m.policyEngine.CheckCommand(fullCmd); err != nil {
		return fmt.Errorf("policy engine blocked execution: %w", err)
	}

	_, err := m.executor.Exec(ctx, cmd, args...)
	return err
}

// Upgrade upgrades the system packages
func (m *UniversalPackageManager) Upgrade() error {
	ctx := context.Background()
	var args []string
	cmd := m.profile.PackageManager

	// Build command string for policy checking
	var fullCmd string
	switch cmd {
	case "apt":
		args = []string{"upgrade", "-y"}
		fullCmd = fmt.Sprintf("%s upgrade -y", cmd)
	case "dnf", "yum":
		args = []string{"upgrade", "-y"}
		fullCmd = fmt.Sprintf("%s upgrade -y", cmd)
	case "pacman":
		args = []string{"-Syu", "--noconfirm"}
		fullCmd = fmt.Sprintf("%s -Syu --noconfirm", cmd)
	case "zypper":
		args = []string{"update", "-y"}
		fullCmd = fmt.Sprintf("%s update -y", cmd)
	case "apk":
		args = []string{"upgrade"}
		fullCmd = fmt.Sprintf("%s upgrade", cmd)
	case "portage":
		cmd = "emerge"
		args = []string{"--update", "--deep", "--newuse", "@world"}
		fullCmd = "emerge --update --deep --newuse @world"
	case "xbps":
		cmd = "xbps-install"
		args = []string{"-u"}
		fullCmd = "xbps-install -u"
	default:
		return fmt.Errorf("unsupported package manager: %s", cmd)
	}

	// Check command against policy engine
	if err := m.policyEngine.CheckCommand(fullCmd); err != nil {
		return fmt.Errorf("policy engine blocked execution: %w", err)
	}

	_, err := m.executor.Exec(ctx, cmd, args...)
	return err
}

// Search searches for a package
func (m *UniversalPackageManager) Search(query string) ([]ports.PackageInfo, error) {
	ctx := context.Background()
	var args []string
	cmdName := m.profile.PackageManager

	// Build command string for policy checking
	var fullCmd string
	switch cmdName {
	case "apt":
		args = []string{"search", query}
		fullCmd = fmt.Sprintf("%s search %s", cmdName, query)
	case "dnf", "yum":
		args = []string{"search", query}
		fullCmd = fmt.Sprintf("%s search %s", cmdName, query)
	case "pacman":
		args = []string{"-Ss", query}
		fullCmd = fmt.Sprintf("%s -Ss %s", cmdName, query)
	case "apk":
		args = []string{"search", "-v", query}
		fullCmd = fmt.Sprintf("%s search -v %s", cmdName, query)
	default:
		return nil, fmt.Errorf("search not implemented for %s", cmdName)
	}

	// Check command against policy engine
	if err := m.policyEngine.CheckCommand(fullCmd); err != nil {
		return nil, fmt.Errorf("policy engine blocked execution: %w", err)
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
