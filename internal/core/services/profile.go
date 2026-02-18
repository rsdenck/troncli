package services

// Package services implements core business logic services.

import (
	"context"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/mascli/troncli/internal/core/adapter"
	"github.com/mascli/troncli/internal/core/domain"
)

// ProfileEngine detects the system profile
type ProfileEngine struct {
	executor adapter.Executor
}

// NewProfileEngine creates a new profile engine
func NewProfileEngine(executor adapter.Executor) *ProfileEngine {
	return &ProfileEngine{
		executor: executor,
	}
}

// DetectProfile scans the system and returns the profile
func (e *ProfileEngine) DetectProfile() (*domain.SystemProfile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	profile := &domain.SystemProfile{}

	// 1. Detect Distro
	e.detectDistro(ctx, profile)

	// 2. Detect Init System
	e.detectInitSystem(ctx, profile)

	// 3. Detect Package Manager
	e.detectPackageManager(profile)

	// 4. Detect Firewall
	e.detectFirewall(ctx, profile)

	// 5. Detect Network Stack
	e.detectNetworkStack(profile)

	// 6. Detect Environment
	e.detectEnvironment(ctx, profile)

	return profile, nil
}

func (e *ProfileEngine) detectDistro(ctx context.Context, p *domain.SystemProfile) {
	// Try /etc/os-release
	out, err := e.executor.Exec(ctx, "cat", "/etc/os-release")
	if err == nil {
		lines := strings.Split(out.Stdout, "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "ID=") {
				p.Distro = strings.Trim(strings.TrimPrefix(line, "ID="), "\"")
			}
			if strings.HasPrefix(line, "VERSION_ID=") {
				p.Version = strings.Trim(strings.TrimPrefix(line, "VERSION_ID="), "\"")
			}
		}
	}
}

func (e *ProfileEngine) detectInitSystem(ctx context.Context, p *domain.SystemProfile) {
	// Check PID 1
	out, err := e.executor.Exec(ctx, "readlink", "/proc/1/exe")
	if err == nil {
		if strings.Contains(out.Stdout, "systemd") {
			p.InitSystem = "systemd"
			return
		}
		if strings.Contains(out.Stdout, "init") {
			// Could be sysv or openrc, check further
			if _, err := os.Stat("/run/openrc"); err == nil {
				p.InitSystem = "openrc"
				return
			}
			p.InitSystem = "sysvinit" // Fallback assumption
			return
		}
	}
	// Check for runit
	if _, err := os.Stat("/run/runit"); err == nil {
		p.InitSystem = "runit"
		return
	}
}

func (e *ProfileEngine) detectPackageManager(p *domain.SystemProfile) {
	managers := []string{"apt", "dnf", "yum", "pacman", "zypper", "apk"}
	for _, mgr := range managers {
		_, err := exec.LookPath(mgr)
		if err == nil {
			p.PackageManager = mgr
			return
		}
	}
}

func (e *ProfileEngine) detectFirewall(ctx context.Context, p *domain.SystemProfile) {
	// Priority: ufw > firewalld > nftables > iptables
	if _, err := exec.LookPath("ufw"); err == nil {
		// Check if active
		out, _ := e.executor.Exec(ctx, "ufw", "status")
		if strings.Contains(out.Stdout, "active") {
			p.Firewall = "ufw"
			return
		}
	}
	if _, err := exec.LookPath("firewall-cmd"); err == nil {
		out, _ := e.executor.Exec(ctx, "firewall-cmd", "--state")
		if strings.Contains(out.Stdout, "running") {
			p.Firewall = "firewalld"
			return
		}
	}
	if _, err := exec.LookPath("nft"); err == nil {
		p.Firewall = "nftables"
		return
	}
	if _, err := exec.LookPath("iptables"); err == nil {
		p.Firewall = "iptables"
		return
	}
}

func (e *ProfileEngine) detectNetworkStack(p *domain.SystemProfile) {
	if _, err := os.Stat("/etc/netplan"); err == nil {
		p.NetworkStack = "netplan"
		return
	}
	if _, err := exec.LookPath("nmcli"); err == nil {
		p.NetworkStack = "NetworkManager"
		return
	}
	if _, err := os.Stat("/etc/sysconfig/network-scripts"); err == nil {
		p.NetworkStack = "ifcfg" // RHEL/CentOS
		return
	}
	if _, err := os.Stat("/etc/network/interfaces"); err == nil {
		p.NetworkStack = "interfaces" // Debian/Old Ubuntu
		return
	}
	if _, err := exec.LookPath("networkctl"); err == nil {
		p.NetworkStack = "systemd-networkd"
		return
	}
}

func (e *ProfileEngine) detectEnvironment(ctx context.Context, p *domain.SystemProfile) {
	// WSL
	out, _ := e.executor.Exec(ctx, "cat", "/proc/version")
	if strings.Contains(strings.ToLower(out.Stdout), "microsoft") {
		p.Environment = "WSL"
		return
	}
	// Docker
	if _, err := os.Stat("/.dockerenv"); err == nil {
		p.Environment = "Docker"
		return
	}
	// Kubernetes
	// Check cgroups for kubepods
	cgroup, _ := e.executor.Exec(ctx, "cat", "/proc/1/cgroup")
	if strings.Contains(cgroup.Stdout, "kubepods") {
		p.Environment = "Kubernetes"
		return
	}
	// VM (DMI) - requires root usually, or try systemd-detect-virt
	virt, err := e.executor.Exec(ctx, "systemd-detect-virt")
	if err == nil && virt.Stdout != "none" {
		p.Environment = "VM (" + virt.Stdout + ")"
		return
	}
	p.Environment = "BareMetal"
}
