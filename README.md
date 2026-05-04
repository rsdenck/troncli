# NUX CLI

![NUX Logo](https://img.shields.io/badge/NUX-Professional-orange?style=for-the-badge)
![Version](https://img.shields.io/badge/version-0.3.0--beta-orange?style=for-the-badge)
![License](https://img.shields.io/badge/license-MIT-black?style=for-the-badge)
![Go Version](https://img.shields.io/badge/go-1.21+-orange?style=for-the-badge)
![Platform](https://img.shields.io/badge/platform-Linux-black?style=for-the-badge)

<p align="center">
  <img src="https://img.shields.io/badge/NUX-Big_Boss_for_Linux-orange?style=for-the-badge&labelColor=black">
</p>

<h1 align="center">NUX CLI - Big Boss for Linux & SysAdmins</h1>

<p align="center">
  <b>Professional CLI for Linux system administration with native AI integration</b><br/>
  <i>Ollama, NVIDIA Build, OpenAI, Claude - All in one place</i>
</p>

---

## About the Project

**NUX** is a next-generation CLI (Command Line Interface) for Linux system administration. Built in Go, it offers a modular architecture with native support for multiple AI providers.

### Key Features

- **Integrated AI Providers**: Ollama, NVIDIA Build, OpenAI, Claude
- **Skill Management**: 120+ skills for external tools
- **Automatic Installation**: Automatic download and deployment of tools
- **Professional Output**: Standardized GCX JSON format
- **Secure Vault**: Token and configuration management
- **Multi-distribution**: Compatible with apt, yum, dnf, pacman

---

## Installation

```bash
# Download binary (coming soon)
wget https://github.com/rsdenck/nux/releases/latest/download/nux-linux-amd64 -O nux
chmod +x nux
sudo mv nux /usr/local/bin/
nux --version
```

---

## Basic Usage

### AI Commands

```bash
# Query Ollama (default)
nux ask "How to list files in Linux?"

# Use NVIDIA Build
nux ask query --provider nvidia --model minimaxai/minimax-m2.7 "NVIDIA query"

# Configure providers
nux ask config --provider ollama --host http://192.168.130.25:11434
nux ask config --provider nvidia --api-key "your-token"
```

### Skill Management

```bash
# List available skills
nux skill list

# Install skill
nux skill install ansible

# Skill info
nux skill info terraform
```

---

## Architecture

```
nux/
├── cmd/nux/commands/    # Cobra commands
├── internal/            # Core (core, modules, vault)
├── skills/              # Skill definitions (.md)
└── tests/               # Go and k6 tests
```

---

## Supported Providers

| Provider | Status | Default Model |
|----------|--------|---------------|
| **Ollama** | Operational | `gemma:2b` |
| **NVIDIA Build** | Operational | `minimaxai/minimax-m2.7` |
| **OpenAI** | Configurable | `gpt-3.5-turbo` |
| **Claude** | Configurable | `claude-3-5-sonnet` |

---

## Available Skills

NUX has **120+ skills** organized by categories:

- `infrastructure/` - Terraform, Pulumi, Packer
- `automation/` - Ansible, Puppet, Chef
- `ci-cd/` - CircleCI, Jenkins, GitLab
- `git/` - GitHub CLI, GitLab CLI
- `testing/` - k6, JMeter
- `kubernetes/` - kubectl, Helm, K9s
- `containers/` - Docker, Podman, Buildah
- `security/` - Nmap, Metasploit, John

Installation scripts: [rsdenck/skillnux](https://github.com/rsdenck/skillnux)

---

## Testing

```bash
# Go tests
go test ./tests/...

# Load testing (k6)
k6 run tests/k6-load-test.js
```

---

## Contributing

1. Fork the project
2. Create your feature branch (`git checkout -b feature/new-skill`)
3. Commit your changes (`git commit -m 'feat: new skill'`)
4. Push to the branch (`git push origin feature/new-skill`)
5. Open a Pull Request

---

## License

This project is under the **MIT License** - see the [LICENSE](LICENSE) file for details.

---

<p align="center">
  <img src="https://img.shields.io/badge/Made_with-Go-black?style=for-the-badge&labelColor=orange">
  <img src="https://img.shields.io/badge/For-Linux-orange?style=for-the-badge&labelColor=black">
</p>
