package linux

import (
	"fmt"
	"os/exec"
	"strings"
)

// Package linux provides Linux-specific abstractions for NUX
// This centralizes distro detection, command availability, and common operations

type Distro string

const (
	DistroDebian Distro = "debian"
	DistroUbuntu Distro = "ubuntu"
	DistroRHEL    Distro = "rhel"
	DistroCentOS  Distro = "centos"
	DistroFedora  Distro = "fedora"
	DistroRocky   Distro = "rocky"
	DistroAlma    Distro = "almalinux"
	DistroArch    Distro = "arch"
	DistroAlpine  Distro = "alpine"
	DistroSUSE    Distro = "suse"
)

func DetectDistro() Distro {
	// Check /etc/os-release first
	data, err := exec.Command("grep", "^ID=", "/etc/os-release").CombinedOutput()
	if err == nil {
		line := strings.TrimSpace(string(data))
		id := strings.TrimPrefix(line, "ID=")
		id = strings.Trim(id, "\"")

		switch id {
		case "debian":
			return DistroDebian
		case "ubuntu":
			return DistroUbuntu
		case "rhel", "redhat":
			return DistroRHEL
		case "centos":
			return DistroCentOS
		case "fedora":
			return DistroFedora
		case "rocky":
			return DistroRocky
		case "almalinux":
			return DistroAlma
		case "arch":
			return DistroArch
		case "alpine":
			return DistroAlpine
		case "suse", "opensuse":
			return DistroSUSE
		}
	}

	return ""
}

func (d Distro) String() string {
	return string(d)
}

func (d Distro) IsDebianBased() bool {
	return d == DistroDebian || d == DistroUbuntu
}

func (d Distro) IsRedHatBased() bool {
	return d == DistroRHEL || d == DistroCentOS || d == DistroFedora || d == DistroRocky || d == DistroAlma
}

func (d Distro) IsArchBased() bool {
	return d == DistroArch
}

func (d Distro) IsAlpine() bool {
	return d == DistroAlpine
}

func (d Distro) PackageManager() string {
	if d.IsDebianBased() {
		return "apt"
	}
	if d.IsRedHatBased() {
		return "dnf"
	}
	if d.IsArchBased() {
		return "pacman"
	}
	if d.IsAlpine() {
		return "apk"
	}
	if d == DistroSUSE {
		return "zypper"
	}
	return "unknown"
}

func CommandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func RunCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("command %s failed: %s - %s", name, err.Error(), strings.TrimSpace(string(out)))
	}
	return strings.TrimSpace(string(out)), nil
}

func RunCommandSilent(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	return cmd.Run()
}
