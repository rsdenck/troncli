package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

func main() {
	if runtime.GOOS != "linux" {
		fmt.Println("❌ This script only works on Linux")
		return
	}

	fmt.Println("🔍 Detecting Linux distribution...")

	// Read /etc/os-release
	cmd := exec.Command("cat", "/etc/os-release")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("❌ Failed to detect distribution: %s\n", err)
		return
	}

	osRelease := strings.ToLower(string(out))
	var distro string

	switch {
	case strings.Contains(osRelease, "ubuntu") || strings.Contains(osRelease, "debian"):
		distro = "debian"
	case strings.Contains(osRelease, "fedora"):
		distro = "fedora"
	case strings.Contains(osRelease, "arch"):
		distro = "arch"
	case strings.Contains(osRelease, "centos") || strings.Contains(osRelease, "rhel") || strings.Contains(osRelease, "rocky") || strings.Contains(osRelease, "almalinux"):
		distro = "rhel"
	default:
		distro = "unknown"
	}

	fmt.Printf("✅ Detected distribution: %s\n", distro)

	switch distro {
	case "debian":
		installDebianUbuntu()
	case "fedora":
		installFedora()
	case "arch":
		installArch()
	default:
		fmt.Println("⚠ Unknown distribution. Attempting generic installation...")
		installGeneric()
	}
}

func installDebianUbuntu() {
	fmt.Println("📦 Installing ProtonVPN CLI for Debian/Ubuntu...")

	// Add Proton repository
	cmds := [][]string{
		{"wget", "-q", "https://protonvpn.com/download/protonvpn-stable-release/protonvpn-stable-release_1.0.0-1_all.deb"},
		{"sudo", "dpkg", "-i", "protonvpn-stable-release_1.0.0-1_all.deb"},
		{"sudo", "apt-get", "update"},
		{"sudo", "apt-get", "install", "-y", "protonvpn-cli"},
	}

	for _, cmd := range cmds {
		fmt.Printf("Running: %s\n", strings.Join(cmd, " "))
		if err := exec.Command(cmd[0], cmd[1:]...).Run(); err != nil {
			fmt.Printf("❌ Command failed: %s\n", err)
			return
		}
	}

	fmt.Println("✅ ProtonVPN CLI installed successfully!")
	fmt.Println("Run: protonvpn --help")
}

func installFedora() {
	fmt.Println("📦 Installing ProtonVPN CLI for Fedora...")

	cmds := [][]string{
		{"sudo", "dnf", "install", "-y", "https://protonvpn.com/download/protonvpn-stable-release/protonvpn-stable-release-1.0.0-1.noarch.rpm"},
		{"sudo", "dnf", "install", "-y", "protonvpn-cli"},
	}

	for _, cmd := range cmds {
		fmt.Printf("Running: %s\n", strings.Join(cmd, " "))
		if err := exec.Command(cmd[0], cmd[1:]...).Run(); err != nil {
			fmt.Printf("❌ Command failed: %s\n", err)
			return
		}
	}

	fmt.Println("✅ ProtonVPN CLI installed successfully!")
	fmt.Println("Run: protonvpn --help")
}

func installArch() {
	fmt.Println("📦 Installing ProtonVPN CLI for Arch Linux...")

	cmd := exec.Command("yay", "-S", "--noconfirm", "protonvpn-cli")
	if err := cmd.Run(); err != nil {
		fmt.Println("⚠ yay not found. Trying manual AUR installation...")
		// Fallback to git clone and makepkg
		cmds := [][]string{
			{"git", "clone", "https://aur.archlinux.org/protonvpn-cli.git"},
			{"sh", "-c", "cd protonvpn-cli && makepkg -si --noconfirm"},
		}
		for _, c := range cmds {
			if err := exec.Command(c[0], c[1:]...).Run(); err != nil {
				fmt.Printf("❌ Command failed: %s\n", err)
				return
			}
		}
	}

	fmt.Println("✅ ProtonVPN CLI installed successfully!")
	fmt.Println("Run: protonvpn --help")
}

func installGeneric() {
	fmt.Println("⚠ Attempting generic installation...")

	// Try to download and install the binary directly
	url := "https://protonvpn.com/download/protonvpn-cli_linux"
	fmt.Printf("Downloading from: %s\n", url)

	if err := exec.Command("wget", "-q", url, "-O", "/tmp/protonvpn-cli").Run(); err != nil {
		fmt.Printf("❌ Download failed: %s\n", err)
		return
	}

	if err := exec.Command("chmod", "+x", "/tmp/protonvpn-cli").Run(); err != nil {
		fmt.Printf("❌ Chmod failed: %s\n", err)
		return
	}

	if err := exec.Command("sudo", "mv", "/tmp/protonvpn-cli", "/usr/local/bin/protonvpn").Run(); err != nil {
		fmt.Printf("❌ Move failed: %s\n", err)
		return
	}

	fmt.Println("✅ ProtonVPN CLI installed successfully!")
	fmt.Println("Run: protonvpn --help")
}
