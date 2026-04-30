# NUX - Linux CLI Manager

NUX is a powerful CLI Master/Manager for Linux that can incorporate functionalities from other CLIs through **Skills**.

## Concept

NUX acts as a **CLI Manager** - it doesn't just manage your system, it manages other CLIs! Through the skill system, you can:
- Install other CLI tools as "skills"
- Configure API keys and credentials (stored securely in vault)
- Use NUX as a proxy to control other tools

## Features

- **Multi-distribution support**: Works on any Linux (apt, dnf, yum, pacman, apk, zypper)
- **Skill Engine**: 145+ pre-defined skills ready to install
- **Vault System**: Secure storage for configurations and API keys (`~/.skills/.nux.json`)
- **Ollama Integration**: AI agent with `qwen3-coder` model
- **Core Linux Management**: Network, disk, LVM, NFS, services, firewall, users

## Installation

### From Release (Recommended)
```bash
# Download from https://github.com/rsdenck/nux/releases
# Debian/Ubuntu
dpkg -i nux_*.deb

# RHEL/Fedora
rpm -ivh nux-*.rpm

# Arch
pacman -U nux-*.pkg.tar.zst
```

### From Source
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

# List available skills
nux skill list

# Install a skill (e.g., docker)
nux skill install docker

# Enable the skill
nux skill enable docker

# Use Ollama AI
nux agent ask "create LVM with 50GB"
```

## Skill System

Skills are defined as `.md` files in the `skills/` directory. Each skill contains:
- Repository URL
- Description
- Install command
- Type (shell, tool, cloud, etc.)

### Example Skills
- **Shell**: bash, zsh, fish, tmux
- **System**: docker, kubectl, systemd, lvm
- **Network**: curl, wget, nmap, tcpdump
- **Cloud**: aws, gcloud, azure, terraform
- **AI**: ollama, openai, claude, aider

## Vault Configuration

NUX uses a vault file at `~/.skills/.nux.json` (permissions: 0600):
```json
{
  "version": "1.0.0",
  "installed_skills": ["docker", "kubectl"],
  "enabled_skills": ["docker"],
  "api_keys": {"github": "ghp_xxxxx"},
  "ollama": {
    "host": "http://localhost:11434",
    "model": "qwen3-coder",
    "enabled": true
  }
}
```

## Multi-Distribution Support

NUX automatically detects your distribution and uses the appropriate tools:
- **Package managers**: apt, dnf, yum, pacman, apk, zypper
- **Init systems**: systemd, openrc, sysvinit, runit
- **Firewalls**: nftables, iptables, firewalld, ufw

## License

MIT

## Repository

- **GitHub**: https://github.com/rsdenck/nux
- **Issues**: https://github.com/rsdenck/nux/issues
- **Releases**: https://github.com/rsdenck/nux/releases
