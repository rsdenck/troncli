# 🗺️ TRONCLI Roadmap

> "The most complete tool ever created for Linux Sysadmins."

Our vision is to build a unified, enterprise-grade TUI that centralizes all critical Linux administration tasks.

---

## Phase 1: Foundation (Current)
- [x] **Project Structure**: Clean Architecture implementation.
- [x] **UI Framework**: `tview` based TUI with TRON: Legacy theme.
- [x] **Core Modules**:
  - [x] **Dashboard**: Real-time monitoring (CPU, Memory, Load, Top).
  - [x] **LVM Manager**: List PVs, VGs, LVs.
  - [x] **Audit**: Log viewer for SSH and Sudo.
  - [x] **SSH**: Profile list and basic `rsd-sshm` integration.

## Phase 2: Enhanced Functionality (Next)
- [ ] **LVM Operations**:
  - Implement Create/Extend/Reduce/Remove logic.
  - Confirmation modals for destructive actions.
- [ ] **SSH Advanced**:
  - Multi-session management (tmux integration?).
  - Parallel command execution on multiple hosts.
- [ ] **Network Manager**:
  - `netplan` / `NetworkManager` integration.
  - Interface configuration (IP, DNS, Routes).
  - Bandwidth monitoring per interface.

## Phase 3: Security & Compliance
- [ ] **Firewall Management**:
  - `nftables` / `iptables` rule editor.
  - `fail2ban` status and unban actions.
- [ ] **User Management**:
  - Create/Delete users and groups.
  - Permission auditing (SUID/SGID finder).
  - Password policy enforcement checks.

## Phase 4: Enterprise Features
- [ ] **Remote Mode**:
  - Run `troncli` locally but manage remote servers via SSH agent forwarding.
- [ ] **Plugins System**:
  - Load custom Go plugins for specific environments (AWS, K8s).
- [ ] **Report Generation**:
  - Export audit logs and system health reports to PDF/HTML.

## Phase 5: The Ultimate Tool
- [ ] **Container Management**:
  - Basic Docker/Podman container status and logs.
- [ ] **Systemd Services**:
  - Service status, start/stop/restart/enable/disable.
- [ ] **Package Management**:
  - Unified interface for `apt`, `dnf`, `pacman`.

---

*Note: This roadmap is subject to change based on community feedback and project priorities.*
