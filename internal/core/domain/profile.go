package domain

// Package domain defines the core domain models.

// SystemProfile represents the detected system environment
type SystemProfile struct {
	Distro         string // Ubuntu, Fedora, Arch, etc.
	Version        string // 22.04, 39, etc.
	InitSystem     string // systemd, openrc, runit, sysvinit
	PackageManager string // apt, dnf, yum, pacman, zypper, apk
	Firewall       string // nftables, iptables, firewalld, ufw
	NetworkStack   string // netplan, ifcfg, interfaces, NetworkManager, systemd-networkd
	Environment    string // WSL, Docker, Kubernetes, VM, BareMetal
}
