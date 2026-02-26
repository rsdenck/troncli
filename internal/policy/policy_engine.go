package policy

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// PolicyEngine provides minimal security policies for agent operations
type PolicyEngine struct {
	Enabled bool
}

// NewPolicyEngine creates a new policy engine
func NewPolicyEngine() *PolicyEngine {
	return &PolicyEngine{
		Enabled: true,
	}
}

// CheckCommand validates if a command is safe to execute
func (p *PolicyEngine) CheckCommand(cmd string) error {
	if !p.Enabled {
		return nil
	}

	// Split command to analyze
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return fmt.Errorf("empty command")
	}

	// CRITICAL: Never allow kernel deletion
	if p.isKernelDeletionCommand(cmd) {
		return fmt.Errorf("CRITICAL: Command attempts to delete kernel - BLOCKED")
	}

	// WARNING: Dangerous operations that require explicit confirmation
	if p.isDangerousCommand(cmd) {
		fmt.Printf("⚠️  WARNING: Dangerous command detected: %s\n", cmd)
		fmt.Printf("This operation could cause system instability.\n")
		fmt.Printf("Continue? (y/N): ")

		var response string
		fmt.Scanln(&response)
		if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
			return fmt.Errorf("Operation cancelled by user")
		}
	}

	return nil
}

// isKernelDeletionCommand checks if command attempts to delete kernel
func (p *PolicyEngine) isKernelDeletionCommand(cmd string) bool {
	dangerousPatterns := []string{
		"rm -rf /boot",
		"rm -rf /lib/modules",
		"dpkg --remove linux-image",
		"apt-get remove linux-image",
		"yum remove kernel",
		"dnf remove kernel",
		"pacman -R linux",
		"mkfs.ext4 /dev/sda1",
		"mkfs.ext4 /dev/vda1",
		"format /",
		"rm -rf /",
	}

	lowerCmd := strings.ToLower(cmd)
	for _, pattern := range dangerousPatterns {
		if strings.Contains(lowerCmd, pattern) {
			return true
		}
	}

	return false
}

// isDangerousCommand checks if command is potentially dangerous
func (p *PolicyEngine) isDangerousCommand(cmd string) bool {
	dangerousPatterns := []string{
		"rm -rf /",
		"mkfs.",
		"format",
		"fdisk",
		"parted",
		"shutdown",
		"reboot",
		"halt",
		"poweroff",
		"iptables -F",
		"systemctl stop",
		"kill -9",
		"killall",
		"pkill -9",
	}

	lowerCmd := strings.ToLower(cmd)
	for _, pattern := range dangerousPatterns {
		if strings.Contains(lowerCmd, pattern) {
			return true
		}
	}

	return false
}

// ExecuteWithPolicy executes a command with policy checks
func (p *PolicyEngine) ExecuteWithPolicy(ctx context.Context, cmd string) error {
	if err := p.CheckCommand(cmd); err != nil {
		return err
	}

	// Execute with proper privileges
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return fmt.Errorf("empty command")
	}

	// Check if command needs root privileges
	if p.needsRoot(cmd) {
		parts = append([]string{"sudo"}, parts...)
		cmd = strings.Join(parts, " ")
	}

	fmt.Printf("🔧 Executing: %s\n", cmd)

	execCmd := exec.Command("sh", "-c", cmd)
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr
	execCmd.Stdin = os.Stdin

	return execCmd.Run()
}

// needsRoot checks if command requires root privileges
func (p *PolicyEngine) needsRoot(cmd string) bool {
	rootCommands := []string{
		"apt-get", "apt", "yum", "dnf", "pacman", "zypper",
		"systemctl", "service", "iptables", "nftables",
		"fdisk", "parted", "mkfs", "mount", "umount",
		"useradd", "userdel", "usermod", "groupadd", "groupdel",
		"crontab", "insmod", "rmmod", "modprobe",
		"sysctl", "dmesg", "lsblk", "fdisk -l",
	}

	lowerCmd := strings.ToLower(cmd)
	for _, rootCmd := range rootCommands {
		if strings.Contains(lowerCmd, rootCmd) {
			return true
		}
	}

	return false
}

// GetSystemInfo returns system information for context
func (p *PolicyEngine) GetSystemInfo() map[string]string {
	info := make(map[string]string)

	// Get user info
	if os.Getuid() == 0 {
		info["user"] = "root"
	} else {
		info["user"] = "non-root"
	}

	// Get hostname
	if hostname, err := os.Hostname(); err == nil {
		info["hostname"] = hostname
	}

	// Get kernel version
	if runtime.GOOS == "linux" {
		info["kernel"] = "Linux"
	} else {
		info["kernel"] = runtime.GOOS
	}

	return info
}
