# TRONCLI Command Reference (v1.0)

This document lists all available commands in TRONCLI v1.0.

## Global Flags
All commands support the following global flags:
- `--json`: Output in JSON format
- `--yaml`: Output in YAML format
- `--quiet`: Suppress output
- `--dry-run`: Simulate execution without making changes
- `--timeout`: Timeout in seconds (default 30)
- `--verbose`: Enable verbose logging
- `--no-color`: Disable color output

## System Management
### `troncli system`
Manage system information and profile.
- `troncli system info`: Display system information (OS, Kernel, etc.)

### `troncli doctor`
Check system health and prerequisites.
- `troncli doctor`: Run all health checks

### `troncli service`
Manage system services (systemd, openrc, etc.).
- `troncli service list`: List all services
- `troncli service start <name>`: Start a service
- `troncli service stop <name>`: Stop a service
- `troncli service restart <name>`: Restart a service
- `troncli service status <name>`: Get service status
- `troncli service enable <name>`: Enable a service at boot
- `troncli service disable <name>`: Disable a service at boot
- `troncli service logs <name>`: View service logs

### `troncli process`
Manage system processes.
- `troncli process tree`: Show process tree
- `troncli process kill <pid>`: Kill a process
- `troncli process renice <pid> <priority>`: Change process priority
- `troncli process open-files <pid>`: List open files by a process
- `troncli process ports <pid>`: List ports used by a process
- `troncli process listening`: List all listening ports
- `troncli process zombies`: Find and kill zombie processes

## Package Management
### `troncli pkg`
Universal package manager wrapper (apt, dnf, yum, pacman, apk, zypper).
- `troncli pkg install <package>`: Install a package
- `troncli pkg remove <package>`: Remove a package
- `troncli pkg update`: Update package lists
- `troncli pkg upgrade`: Upgrade all packages
- `troncli pkg search <query>`: Search for packages

## User & Group Management
### `troncli user`
Manage system users.
- `troncli user list`: List all users
- `troncli user add <username>`: Add a new user
- `troncli user delete <username>`: Delete a user
- `troncli user modify <username>`: Modify a user

### `troncli group`
Manage system groups.
- `troncli group list`: List all groups
- `troncli group add <groupname>`: Add a new group
- `troncli group delete <groupname>`: Delete a group

## Network & Security
### `troncli network`
Manage network configuration.
- `troncli network info`: Show network interfaces
- `troncli network set-state <interface> <up|down>`: Set interface state
- `troncli network sockets`: List open sockets
- `troncli network trace <target>`: Trace route to target
- `troncli network dig <domain>`: DNS lookup
- `troncli network scan <target>`: Port scan
- `troncli network capture <interface>`: Capture packets (tcpdump)

### `troncli firewall`
Manage firewall rules (nftables/iptables/ufw/firewalld).
- `troncli firewall list`: List rules
- `troncli firewall allow <port/proto>`: Allow traffic
- `troncli firewall deny <port/proto>`: Deny traffic
- `troncli firewall enable`: Enable firewall
- `troncli firewall disable`: Disable firewall

### `troncli audit`
Security auditing tools.
- `troncli audit logins`: List recent logins
- `troncli audit sudoers`: List sudoers
- `troncli audit file-changes <path>`: Monitor file changes
- `troncli audit commands`: Audit executed commands

## Storage
### `troncli disk`
Manage disks and filesystems.
- `troncli disk list`: List block devices
- `troncli disk usage <path>`: Show filesystem usage
- `troncli disk mounts`: List mount points
- `troncli disk mount <source> <target>`: Mount a device
- `troncli disk unmount <target>`: Unmount a device
- `troncli disk format <device> <fstype>`: Format a device

## Automation & Tools
### `troncli bash`
Execute bash commands and scripts.
- `troncli bash run <command>`: Run a single command
- `troncli bash script <file>`: Run a bash script

### `troncli remote`
Manage remote SSH connections.
- `troncli remote list`: List saved connections
- `troncli remote connect <name>`: Connect to a remote host
- `troncli remote exec <name> <command>`: Execute command remotely
- `troncli remote copy <name> <src> <dest>`: Copy files via SCP

### `troncli container`
Manage containers (Docker/Podman).
- `troncli container list`: List containers
- `troncli container start <id>`: Start a container
- `troncli container stop <id>`: Stop a container
- `troncli container logs <id>`: View container logs

### `troncli plugin`
Manage TRONCLI plugins.
- `troncli plugin list`: List installed plugins
- `troncli plugin install <url>`: Install a plugin
- `troncli plugin remove <name>`: Remove a plugin

### `troncli completion`
Generate shell autocompletion scripts.
- `troncli completion bash`: Generate bash completion
- `troncli completion zsh`: Generate zsh completion
- `troncli completion fish`: Generate fish completion
- `troncli completion powershell`: Generate powershell completion
