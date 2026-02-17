# TRONCLI

> **SYSTEM IDENTITY**
> 
> HOST: GITHUB
> OS: LINUX
> KERNEL: 6.8.0-RC1
> UPTIME: 99.99%

## OVERVIEW

TRONCLI is a production-grade Linux System Administration TUI (Text User Interface) written 100% in Go. It provides real-time system monitoring, LVM management, network analysis, and security auditing with a visual identity inspired by high-tech grid systems.

**NO MOCKS. NO SIMULATION. REAL KERNEL CONTROL.**

![License](https://img.shields.io/badge/LICENSE-MIT-00d9ff?style=for-the-badge&labelColor=000000)
![Go Version](https://img.shields.io/badge/GO-1.22+-00d9ff?style=for-the-badge&labelColor=000000&logo=go)
![Platform](https://img.shields.io/badge/PLATFORM-LINUX-00d9ff?style=for-the-badge&labelColor=000000&logo=linux)
![Build](https://img.shields.io/badge/BUILD-PASSING-00d9ff?style=for-the-badge&labelColor=000000)

## CORE MODULES

### [01] SYSTEM DASHBOARD
Real-time metrics from `/proc` and `/sys`.
- CPU Usage (User/System/Idle)
- Memory (Ram/Swap)
- Load Average (1/5/15)
- Disk I/O & Network Throughput

### [02] LVM MANAGER
Direct interface for Logical Volume Manager.
- Physical Volumes (PV)
- Volume Groups (VG)
- Logical Volumes (LV)
- Extend/Reduce operations

### [03] NETWORK MATRIX
Advanced network stack analysis.
- Interface statistics
- Real-time RX/TX rates
- DNS Configuration
- Socket states

### [04] SECURITY AUDIT
System hardening and user management.
- User/Group enumeration
- SSH session management
- File permission matrix
- Audit logs

## INSTALLATION

### PRE-REQUISITES
- Linux Kernel 5.4+
- Root privileges (for LVM/Audit)

### BUILD FROM SOURCE

```bash
git clone https://github.com/rsdenck/troncli.git
cd troncli
go build -ldflags="-s -w" -o troncli cmd/troncli/main.go
./troncli
```

## ARCHITECTURE

The system follows Clean Architecture principles with strict separation of concerns.

```text
cmd/
  troncli/       # Entry point
internal/
  core/          # Domain logic & Ports
  modules/       # Implementations (Linux specific)
  ui/            # TUI Layer (tview/tcell)
```

## CONTRIBUTION

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'feat: Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## SECURITY

Please report vulnerabilities to `security@troncli.local`.
See [SECURITY.md](SECURITY.md) for details.

---
TRONCLI | SYSTEM END OF LINE
