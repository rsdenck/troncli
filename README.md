# TRONCLI: The Ultimate Linux Sysadmin TUI

![TRON Legacy Theme](https://img.shields.io/badge/Theme-TRON_Legacy-00C3FF?style=for-the-badge)
![Go Version](https://img.shields.io/badge/Go-1.21+-0055AA?style=for-the-badge&logo=go)
![Platform](https://img.shields.io/badge/Platform-Linux-FCC624?style=for-the-badge&logo=linux)
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)

**troncli** is a modern, tabular, and elegant Terminal User Interface (TUI) designed for advanced Linux system administration. Inspired by the visual aesthetics of *TRON: Legacy*, it combines powerful functionality with a futuristic, neon-infused interface.

> "The most complete tool ever created for Linux Sysadmins."

## 🚀 Features

### 🖥️ Real-Time Dashboard
- **System Monitoring**: CPU, Memory, Swap, Load Average, and Disk I/O visualization.
- **Top Processes**: Live view of resource-consuming processes.
- **Network Status**: Real-time throughput monitoring.

### 🛡️ Security & Audit
- **Audit Logs**: Centralized view of SSH logins, sudo usage, and critical file permission changes.
- **Hardening Checks**: Automated detection of unsafe configurations (e.g., world-writable files, SUID binaries).
- **SSH Management**: Secure connection management via `rsd-sshm` integration.

### 💾 Storage Management (LVM)
- **Visual LVM**: Tabular view of Physical Volumes (PV), Volume Groups (VG), and Logical Volumes (LV).
- **Operations**: (Coming Soon) Create, Extend, Reduce, and Remove volumes directly from the TUI.

### 🌐 Connectivity
- **SSH Profiles**: Manage and connect to multiple remote servers.
- **RSD-SSHM Integration**: Enforces secure, standardized SSH connections.

## 🛠️ Architecture

`troncli` follows **Clean Architecture** principles to ensure scalability, maintainability, and testability.

```
/cmd/troncli         # Main entry point
/internal
    /core            # Business logic & Domain entities
        /ports       # Interfaces (Hexagonal Architecture)
    /infra           # Infrastructure implementations (LVM, SSH, Audit)
    /ui              # Presentation layer (tview/tcell)
        /views       # Screen definitions
        /components  # Reusable UI widgets
        /themes      # TRON visual style definitions
/pkg                 # Shared libraries
/config              # Configuration handling
```

## 📦 Installation

### Prerequisites
- Go 1.21 or higher
- Linux environment (or WSL)
- `rsd-sshm` (for SSH functionality)
- Root/Sudo privileges (for LVM and Audit)

### Build from Source
```bash
git clone https://github.com/mascli/troncli.git
cd troncli
go mod tidy
go build -o troncli ./cmd/troncli
```

## 🎮 Usage

Run the application with root privileges for full functionality:
```bash
sudo ./troncli
```

### Navigation
- **Keyboard**: Use `Arrow Keys` to navigate lists and tables.
- **Shortcuts**:
  - `d`: Dashboard
  - `l`: LVM Manager
  - `h`: SSH Manager
  - `a`: Audit Logs
  - `q` / `Ctrl+C`: Quit

## 🎨 Theme

The interface features a **High-Contrast Neon** theme:
- **Primary**: Neon Cyan (`#00FFFF`) & Blue (`#0000FF`)
- **Background**: Deep Black (`#000000`)
- **Alerts**: Red (`#FF0000`) & Yellow (`#FFFF00`)

## 🤝 Contributing

We welcome contributions! Please follow the [Code of Conduct](CODE_OF_CONDUCT.md) and submit Pull Requests following the Clean Architecture pattern.

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.

---
*End of Line.*
