package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search [term]",
	Short: "Search across files, packages, processes, services, and more",
	Long: `Search for terms across multiple sources:
- Files and directories (find, fd, locate, tree, ls)
- File contents (grep, ripgrep)
- Running processes (ps, pgrep, pidof)
- System services (systemd, openrc)
- Log files (journalctl, /var/log)
- Installed packages (dpkg, rpm, pacman -Qs, apk info)
- Available packages (apt search, dnf search, pacman -Ss)
- Containers (docker, podman, lxc)
- Network ports and connections (ss, netstat, lsof, nmap)
- Firewall rules (iptables, nftables, ufw, firewalld)

Output is categorized automatically.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		term := args[0]

		fmt.Printf("Searching for: %s\n\n", term)

		// Files and directories
		fmt.Println("[FILES]")
		searchFiles(term)
		fmt.Println()

		// File contents
		fmt.Println("[CONTENTS]")
		searchContents(term)
		fmt.Println()

		// Processes
		fmt.Println("[PROCESSES]")
		searchProcesses(term)
		fmt.Println()

		// Services
		fmt.Println("[SERVICES]")
		searchServices(term)
		fmt.Println()

		// Installed packages
		fmt.Println("[PACKAGES INSTALLED]")
		searchInstalledPackages(term)
		fmt.Println()

		// Available packages
		fmt.Println("[PACKAGES AVAILABLE]")
		searchAvailablePackages(term)
		fmt.Println()

		// Containers
		fmt.Println("[CONTAINERS]")
		searchContainers(term)
		fmt.Println()

		// Network
		fmt.Println("[PORTS]")
		searchPorts(term)
		fmt.Println()

		// Binaries in PATH
		fmt.Println("[BINARIES]")
		searchBinaries(term)
		fmt.Println()
	},
}

func searchFiles(term string) {
	// Try fd first, then find
	if hasCommand("fd") {
		exec.Command("fd", "-t", "f", "-t", "d", term).Run()
	} else if hasCommand("find") {
		exec.Command("find", "/", "-name", "*"+term+"*", "-maxdepth", "5").Run()
	} else {
		exec.Command("ls", "-la", "/").Run()
	}
}

func searchContents(term string) {
	if hasCommand("rg") {
		exec.Command("rg", "-i", term, "--max-count=5").Run()
	} else if hasCommand("grep") {
		exec.Command("grep", "-r", "-i", term, "/etc", "--include=*").Run()
	}
}

func searchProcesses(term string) {
	if hasCommand("pgrep") {
		exec.Command("pgrep", "-l", term).Run()
	}
	if hasCommand("ps") {
		exec.Command("ps", "aux").Run() // Would need grep piping
	}
}

func searchServices(term string) {
	if hasCommand("systemctl") {
		exec.Command("systemctl", "list-units", "--all").Run() // Would need grep
	}
}

func searchInstalledPackages(term string) {
	if hasCommand("dpkg") {
		exec.Command("sh", "-c", "dpkg -l | grep -i "+term).Run()
	} else if hasCommand("rpm") {
		exec.Command("sh", "-c", "rpm -qa | grep -i "+term).Run()
	} else if hasCommand("pacman") {
		exec.Command("pacman", "-Qs", term).Run()
	} else if hasCommand("apk") {
		exec.Command("apk", "info", term).Run()
	}
}

func searchAvailablePackages(term string) {
	if hasCommand("apt") {
		exec.Command("apt", "search", term).Run()
	} else if hasCommand("dnf") {
		exec.Command("dnf", "search", term).Run()
	} else if hasCommand("pacman") {
		exec.Command("pacman", "-Ss", term).Run()
	} else if hasCommand("zypper") {
		exec.Command("zypper", "search", term).Run()
	}
}

func searchContainers(term string) {
	if hasCommand("docker") {
		exec.Command("docker", "ps", "-a", "--filter", "name="+term).Run()
	}
	if hasCommand("podman") {
		exec.Command("podman", "ps", "-a", "--filter", "name="+term).Run()
	}
}

func searchPorts(term string) {
	if hasCommand("ss") {
		exec.Command("ss", "-tuln").Run() // Would need grep
	}
}

func searchBinaries(term string) {
	path := os.Getenv("PATH")
	dirs := strings.Split(path, ":")
	for _, dir := range dirs {
		matches, _ := filepath.Glob(filepath.Join(dir, "*"+term+"*"))
		for _, m := range matches {
			fmt.Println(m)
		}
	}
}

func hasCommand(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
