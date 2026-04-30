# NUX - Linux CLI Manager

NUX is a production-grade CLI for comprehensive Linux system administration.
It manages system resources and provides a skill engine for extending functionality.

## Features

- Multi-distribution support (apt, dnf, yum, pacman, apk, zypper)
- Service Management (systemd, openrc, sysvinit, runit)
- Network Configuration and Diagnostics
- Disk and LVM Management
- User and Group Management
- Security Auditing and Compliance
- Container Management (Docker/Podman)
- Skill Engine (manage external CLIs)
- Agent Integration (Ollama AI with qwen3-coder)
- Shell Autocompletion

## Installation

### From Release

Download the latest release from:
https://github.com/rsdenck/nux/releases

#### Debian/Ubuntu
```bash
dpkg -i nux_*.deb
```

#### RHEL/Fedora/CentOS
```bash
rpm -ivh nux-*.rpm
```

#### Arch Linux
```bash
pacman -U nux-*.pkg.tar.zst
```

#### From Source
```bash
git clone https://github.com/rsdenck/nux.git
cd nux
go build -o nux ./cmd/nux
sudo mv nux /usr/local/bin/
```

## Quick Start

```bash
# First-time setup
nux onboard

# Check system health
nux doctor

# List available skills
nux skill list

# Install a skill
nux skill install docker

# Enable skill
nux skill enable docker

# Use Ollama AI agent
nux agent ask "create LVM with 50GB"
```

## Command Output

All commands follow GCX-style JSON output by default:

```json
{
  "status": "success",
  "data": {
    "items": [...]
  },
  "total": 10
}
```

Use --json flag for force JSON output, or --quiet to suppress output.

## Skill System

Skills are defined as .md files in the skills/ directory.
Each skill contains repository URL, description, install command, and type.

Available skill categories:
- Shell: bash, zsh, fish, tmux
- System: docker, kubectl, systemd, lvm
- Network: curl, wget, nmap, tcpdump
- Cloud: aws, gcloud, azure, terraform
- AI: ollama, openai, claude, aider

## Configuration Vault

NUX uses a secure vault at ~/.skills/.nux.json (permissions: 0600):

```json
{
  "version": "1.0.0",
  "installed_skills": [],
  "enabled_skills": [],
  "api_keys": {},
  "ollama": {
    "host": "http://localhost:11434",
    "model": "qwen3-coder",
    "enabled": false
  },
  "vault_mode": true
}
```

## Multi-Distribution Support

NUX automatically detects your distribution:

Package Managers:
- Debian/Ubuntu: apt
- RHEL/Fedora: dnf/yum
- Arch: pacman
- Alpine: apk
- openSUSE: zypper

Init Systems:
- systemd
- openrc
- sysvinit
- runit

Firewalls:
- nftables
- iptables
- firewalld
- ufw

## CI/CD

The project uses GitHub Actions for:
- Automated builds on multiple platforms
- GoReleaser for multi-distribution releases
- Dependabot for dependency updates

## License

MIT

## Repository

GitHub: https://github.com/rsdenck/nux
Issues: https://github.com/rsdenck/nux/issues
Releases: https://github.com/rsdenck/nux/releases
