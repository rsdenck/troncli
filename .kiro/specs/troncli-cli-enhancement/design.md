# TRONCLI CLI Enhancement - Design Document

**Feature:** Remove TUI and Enhance Multi-Distribution CLI Support  
**Spec Type:** Feature  
**Workflow:** Requirements-First  
**Created:** 2025-02-25

---

## Overview

This design document specifies the technical approach for removing all Terminal User Interface (TUI) components from TRONCLI and enhancing it as a pure Command-Line Interface (CLI) tool with improved multi-distribution Linux support and direct Linux subsystem integration.

### Design Goals

1. **Complete TUI Removal**: Eliminate all TUI code, dependencies, and references while preserving CLI functionality
2. **Enhanced Distribution Support**: Add support for Gentoo (portage) and Void Linux (xbps) while improving existing distribution detection
3. **Direct Linux Integration**: Strengthen direct interaction with /proc, /sys, netlink, and syscalls
4. **Architecture Preservation**: Maintain Clean Architecture with proper dependency flow (cmd → core → modules)
5. **Performance Optimization**: Achieve startup < 100ms and command execution < 500ms
6. **Backward Compatibility**: Maintain 100% compatibility for all existing CLI commands and flags

### Current State Analysis

The codebase currently has:
- **No TUI dependencies** in go.mod (clean state)
- **internal/ui/** directory exists but appears unused
- **Clean Architecture** with ports/adapters pattern in place
- **ProfileEngine** for distribution detection (supports apt, dnf, yum, pacman, zypper, apk)
- **console.BoxTable** for human-readable table output
- **Global flags** for JSON/YAML output, dry-run, verbose, etc.

### Design Approach

This design follows a **preservation and enhancement** strategy:
- Remove unused TUI infrastructure
- Enhance existing ProfileEngine for better distribution detection
- Add new distribution support (Gentoo, Void Linux)
- Strengthen direct Linux subsystem integration
- Maintain all existing CLI functionality
- Add comprehensive property-based testing

---

## Architecture

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                         CLI Layer                            │
│  (cmd/troncli/commands/)                                     │
│  ┌──────────┬──────────┬──────────┬──────────┬──────────┐  │
│  │  root    │ system   │ service  │ process  │  pkg     │  │
│  │  agent   │ network  │ disk     │ users    │  ...     │  │
│  └──────────┴──────────┴──────────┴──────────┴──────────┘  │
│         │                                                     │
│         ▼ (uses)                                             │
└─────────────────────────────────────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────────────────────────────┐
│                      Core Layer                              │
│  (internal/core/)                                            │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  Ports (Interfaces)                                  │   │
│  │  - ServiceManager                                    │   │
│  │  - PackageManager                                    │   │
│  │  - ProcessManager                                    │   │
│  │  - NetworkManager                                    │   │
│  │  - DiskManager                                       │   │
│  │  - UserManager                                       │   │
│  └─────────────────────────────────────────────────────┘   │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  Domain Models                                       │   │
│  │  - SystemProfile                                     │   │
│  │  - ServiceUnit                                       │   │
│  │  - PackageInfo                                       │   │
│  └─────────────────────────────────────────────────────┘   │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  Services                                            │   │
│  │  - ProfileEngine (distribution detection)           │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
         │
         ▼ (implements)
┌─────────────────────────────────────────────────────────────┐
│                    Modules Layer                             │
│  (internal/modules/)                                         │
│  ┌──────────┬──────────┬──────────┬──────────┬──────────┐  │
│  │ service  │   pkg    │ process  │ network  │  disk    │  │
│  │  users   │ firewall │   lvm    │  ssh     │  ...     │  │
│  └──────────┴──────────┴──────────┴──────────┴──────────┘  │
│         │                                                     │
│         ▼ (uses)                                             │
└─────────────────────────────────────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────────────────────────────┐
│                  Linux Subsystems                            │
│  - /proc filesystem (process info, network stats)           │
│  - /sys filesystem (hardware, devices, network)             │
│  - netlink sockets (network configuration)                  │
│  - syscalls (kill, setpriority, statvfs)                    │
│  - Package managers (apt, dnf, pacman, portage, xbps)       │
│  - Init systems (systemd, openrc, runit, sysvinit)          │
└─────────────────────────────────────────────────────────────┘
```

### Dependency Flow

The architecture follows Clean Architecture principles with strict dependency rules:

1. **cmd → internal/core**: Commands depend on core ports (interfaces)
2. **internal/core → internal/modules**: Core defines interfaces, modules implement them
3. **internal/modules → Linux**: Modules interact with Linux subsystems
4. **No reverse dependencies**: Lower layers never depend on upper layers

### Output Formatting Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Command Execution                         │
└─────────────────────────────────────────────────────────────┘
                        │
                        ▼
         ┌──────────────┴──────────────┐
         │                              │
         ▼                              ▼
┌─────────────────┐          ┌─────────────────┐
│  Structured     │          │  Human-Readable │
│  Output         │          │  Output         │
│  (--json/yaml)  │          │  (default)      │
└─────────────────┘          └─────────────────┘
         │                              │
         ▼                              ▼
┌─────────────────┐          ┌─────────────────┐
│ JSON/YAML       │          │ console.BoxTable│
│ Marshaler       │          │ Formatter       │
└─────────────────┘          └─────────────────┘
```

---

## Components and Interfaces

### 1. TUI Removal Components

#### 1.1 Files to Remove

```
internal/ui/                    # Entire directory
├── components/                 # TUI components
├── console/                    # TUI console
├── themes/                     # TUI themes
└── views/                      # TUI views
```

#### 1.2 Verification Strategy

After removal, verify:
- No imports of TUI libraries in any Go file
- No TUI-related command flags (--tui, --interactive-mode)
- Build succeeds without errors
- All tests pass

### 2. Enhanced ProfileEngine

#### 2.1 Interface

```go
// ProfileEngine detects system characteristics
type ProfileEngine interface {
    DetectProfile() (*SystemProfile, error)
    DetectDistribution() (string, string, error)  // name, version
    DetectPackageManager() (string, error)
    DetectInitSystem() (string, error)
    DetectFirewall() (string, error)
    DetectNetworkStack() (string, error)
    DetectEnvironment() (string, error)
}
```

#### 2.2 Enhanced SystemProfile

```go
type SystemProfile struct {
    Distro         string // ubuntu, fedora, arch, alpine, opensuse, gentoo, void
    Version        string // 22.04, 39, etc.
    InitSystem     string // systemd, openrc, runit, sysvinit
    PackageManager string // apt, dnf, yum, pacman, zypper, apk, portage, xbps
    Firewall       string // nftables, iptables, firewalld, ufw
    NetworkStack   string // netplan, ifcfg, interfaces, NetworkManager, systemd-networkd
    Environment    string // WSL, Docker, Kubernetes, VM, BareMetal
    Kernel         string // Linux kernel version
    Architecture   string // x86_64, aarch64, etc.
}
```

#### 2.3 Distribution Detection Logic

```
Detection Priority:
1. Read /etc/os-release (ID, VERSION_ID, NAME, PRETTY_NAME)
2. Fallback to /etc/lsb-release
3. Fallback to distribution-specific files:
   - /etc/debian_version (Debian/Ubuntu)
   - /etc/redhat-release (RHEL/CentOS/Fedora)
   - /etc/arch-release (Arch)
   - /etc/alpine-release (Alpine)
   - /etc/gentoo-release (Gentoo)
   - /etc/void-release (Void)
```

#### 2.4 Package Manager Detection

```go
// Detection order (first found wins)
var packageManagers = []struct {
    name       string
    binary     string
    testCmd    []string
    distros    []string
}{
    {"apt",     "apt",      []string{"apt", "--version"},           []string{"ubuntu", "debian"}},
    {"dnf",     "dnf",      []string{"dnf", "--version"},           []string{"fedora", "rhel"}},
    {"yum",     "yum",      []string{"yum", "--version"},           []string{"centos", "rhel"}},
    {"pacman",  "pacman",   []string{"pacman", "--version"},        []string{"arch", "manjaro"}},
    {"zypper",  "zypper",   []string{"zypper", "--version"},        []string{"opensuse", "sles"}},
    {"apk",     "apk",      []string{"apk", "--version"},           []string{"alpine"}},
    {"portage", "emerge",   []string{"emerge", "--version"},        []string{"gentoo"}},
    {"xbps",    "xbps-install", []string{"xbps-install", "--version"}, []string{"void"}},
}
```

#### 2.5 Init System Detection

```go
// Detection logic
func detectInitSystem() string {
    // 1. Check PID 1 executable
    target, _ := os.Readlink("/proc/1/exe")
    if strings.Contains(target, "systemd") {
        return "systemd"
    }
    
    // 2. Check for init system markers
    if _, err := os.Stat("/run/openrc"); err == nil {
        return "openrc"
    }
    if _, err := os.Stat("/run/runit"); err == nil {
        return "runit"
    }
    
    // 3. Fallback to sysvinit
    return "sysvinit"
}
```

### 3. Universal Package Manager Adapter

#### 3.1 Interface (existing)

```go
type PackageManager interface {
    DetectManager() (string, error)
    Install(packageName string) error
    Remove(packageName string) error
    Update() error
    Upgrade() error
    Search(query string) ([]PackageInfo, error)
}
```

#### 3.2 Implementation Strategy

```go
type UniversalPackageManager struct {
    profile  *SystemProfile
    executor Executor
}

func (m *UniversalPackageManager) Install(pkg string) error {
    switch m.profile.PackageManager {
    case "apt":
        return m.executor.Exec("apt", "install", "-y", pkg)
    case "dnf":
        return m.executor.Exec("dnf", "install", "-y", pkg)
    case "yum":
        return m.executor.Exec("yum", "install", "-y", pkg)
    case "pacman":
        return m.executor.Exec("pacman", "-S", "--noconfirm", pkg)
    case "zypper":
        return m.executor.Exec("zypper", "install", "-y", pkg)
    case "apk":
        return m.executor.Exec("apk", "add", pkg)
    case "portage":
        return m.executor.Exec("emerge", pkg)
    case "xbps":
        return m.executor.Exec("xbps-install", "-y", pkg)
    default:
        return fmt.Errorf("unsupported package manager: %s", m.profile.PackageManager)
    }
}
```

### 4. Universal Service Manager Adapter

#### 4.1 Interface (existing)

```go
type ServiceManager interface {
    ListServices() ([]ServiceUnit, error)
    StartService(name string) error
    StopService(name string) error
    RestartService(name string) error
    EnableService(name string) error
    DisableService(name string) error
    GetServiceStatus(name string) (string, error)
    GetServiceLogs(name string, lines int) (string, error)
}
```

#### 4.2 Implementation Strategy

```go
type UniversalServiceManager struct {
    profile  *SystemProfile
    executor Executor
}

func (m *UniversalServiceManager) StartService(name string) error {
    switch m.profile.InitSystem {
    case "systemd":
        return m.executor.Exec("systemctl", "start", name)
    case "openrc":
        return m.executor.Exec("rc-service", name, "start")
    case "runit":
        return m.executor.Exec("sv", "start", name)
    case "sysvinit":
        return m.executor.Exec("service", name, "start")
    default:
        return fmt.Errorf("unsupported init system: %s", m.profile.InitSystem)
    }
}
```

---

### 5. Direct Linux Subsystem Integration

#### 5.1 Process Management (/proc filesystem)

```go
// ProcessReader reads from /proc filesystem
type ProcessReader struct{}

// ReadProcessTree reads process hierarchy from /proc
func (r *ProcessReader) ReadProcessTree() ([]Process, error) {
    entries, err := os.ReadDir("/proc")
    if err != nil {
        return nil, err
    }
    
    var processes []Process
    for _, entry := range entries {
        if !entry.IsDir() {
            continue
        }
        
        // Check if directory name is a PID (numeric)
        pid, err := strconv.Atoi(entry.Name())
        if err != nil {
            continue
        }
        
        proc, err := r.readProcess(pid)
        if err != nil {
            continue // Process may have exited
        }
        processes = append(processes, proc)
    }
    return processes, nil
}

// readProcess reads process info from /proc/[pid]/
func (r *ProcessReader) readProcess(pid int) (Process, error) {
    proc := Process{PID: pid}
    
    // Read /proc/[pid]/stat for basic info
    stat, err := os.ReadFile(fmt.Sprintf("/proc/%d/stat", pid))
    if err != nil {
        return proc, err
    }
    
    // Parse stat file (see proc(5) man page)
    fields := strings.Fields(string(stat))
    if len(fields) >= 4 {
        proc.Name = strings.Trim(fields[1], "()")
        proc.State = fields[2]
        proc.PPID, _ = strconv.Atoi(fields[3])
    }
    
    // Read /proc/[pid]/cmdline for full command
    cmdline, _ := os.ReadFile(fmt.Sprintf("/proc/%d/cmdline", pid))
    proc.Cmdline = string(bytes.ReplaceAll(cmdline, []byte{0}, []byte(" ")))
    
    // Read /proc/[pid]/status for additional info
    status, _ := os.ReadFile(fmt.Sprintf("/proc/%d/status", pid))
    proc.Status = parseStatus(status)
    
    return proc, nil
}

// KillProcess sends signal to process using kill(2) syscall
func (r *ProcessReader) KillProcess(pid int, signal syscall.Signal) error {
    return syscall.Kill(pid, signal)
}

// ReniceProcess changes process priority using setpriority(2)
func (r *ProcessReader) ReniceProcess(pid int, priority int) error {
    return syscall.Setpriority(syscall.PRIO_PROCESS, pid, priority)
}

// ReadOpenFiles reads open file descriptors from /proc/[pid]/fd
func (r *ProcessReader) ReadOpenFiles(pid int) ([]string, error) {
    fdDir := fmt.Sprintf("/proc/%d/fd", pid)
    entries, err := os.ReadDir(fdDir)
    if err != nil {
        return nil, err
    }
    
    var files []string
    for _, entry := range entries {
        link, err := os.Readlink(filepath.Join(fdDir, entry.Name()))
        if err != nil {
            continue
        }
        files = append(files, link)
    }
    return files, nil
}
```

#### 5.2 Network Management (/sys and netlink)

```go
// NetworkReader reads network information
type NetworkReader struct{}

// ReadInterfaces reads network interfaces from /sys/class/net
func (r *NetworkReader) ReadInterfaces() ([]NetworkInterface, error) {
    entries, err := os.ReadDir("/sys/class/net")
    if err != nil {
        return nil, err
    }
    
    var interfaces []NetworkInterface
    for _, entry := range entries {
        iface, err := r.readInterface(entry.Name())
        if err != nil {
            continue
        }
        interfaces = append(interfaces, iface)
    }
    return interfaces, nil
}

// readInterface reads interface details from /sys/class/net/[name]
func (r *NetworkReader) readInterface(name string) (NetworkInterface, error) {
    iface := NetworkInterface{Name: name}
    basePath := fmt.Sprintf("/sys/class/net/%s", name)
    
    // Read MAC address
    mac, _ := os.ReadFile(filepath.Join(basePath, "address"))
    iface.MAC = strings.TrimSpace(string(mac))
    
    // Read MTU
    mtu, _ := os.ReadFile(filepath.Join(basePath, "mtu"))
    iface.MTU, _ = strconv.Atoi(strings.TrimSpace(string(mtu)))
    
    // Read operational state
    state, _ := os.ReadFile(filepath.Join(basePath, "operstate"))
    iface.State = strings.TrimSpace(string(state))
    
    // Read speed (if available)
    speed, _ := os.ReadFile(filepath.Join(basePath, "speed"))
    iface.Speed, _ = strconv.Atoi(strings.TrimSpace(string(speed)))
    
    return iface, nil
}

// ReadSocketStats reads socket statistics from /proc/net/tcp and /proc/net/udp
func (r *NetworkReader) ReadSocketStats() (SocketStats, error) {
    stats := SocketStats{}
    
    // Read TCP sockets
    tcpData, err := os.ReadFile("/proc/net/tcp")
    if err == nil {
        stats.TCP = parseTCPSockets(tcpData)
    }
    
    // Read TCP6 sockets
    tcp6Data, _ := os.ReadFile("/proc/net/tcp6")
    stats.TCP6 = parseTCPSockets(tcp6Data)
    
    // Read UDP sockets
    udpData, _ := os.ReadFile("/proc/net/udp")
    stats.UDP = parseUDPSockets(udpData)
    
    // Read UDP6 sockets
    udp6Data, _ := os.ReadFile("/proc/net/udp6")
    stats.UDP6 = parseUDPSockets(udp6Data)
    
    return stats, nil
}

// ReadRouteTable reads routing table from /proc/net/route
func (r *NetworkReader) ReadRouteTable() ([]Route, error) {
    data, err := os.ReadFile("/proc/net/route")
    if err != nil {
        return nil, err
    }
    
    return parseRouteTable(data), nil
}
```

#### 5.3 Disk and Filesystem Management

```go
// DiskReader reads disk and filesystem information
type DiskReader struct{}

// ReadBlockDevices reads block devices from /sys/block
func (r *DiskReader) ReadBlockDevices() ([]BlockDevice, error) {
    entries, err := os.ReadDir("/sys/block")
    if err != nil {
        return nil, err
    }
    
    var devices []BlockDevice
    for _, entry := range entries {
        dev, err := r.readBlockDevice(entry.Name())
        if err != nil {
            continue
        }
        devices = append(devices, dev)
    }
    return devices, nil
}

// readBlockDevice reads device details from /sys/block/[name]
func (r *DiskReader) readBlockDevice(name string) (BlockDevice, error) {
    dev := BlockDevice{Name: name}
    basePath := fmt.Sprintf("/sys/block/%s", name)
    
    // Read size (in 512-byte sectors)
    size, _ := os.ReadFile(filepath.Join(basePath, "size"))
    sectors, _ := strconv.ParseInt(strings.TrimSpace(string(size)), 10, 64)
    dev.Size = sectors * 512
    
    // Read removable flag
    removable, _ := os.ReadFile(filepath.Join(basePath, "removable"))
    dev.Removable = strings.TrimSpace(string(removable)) == "1"
    
    // Read model (if available)
    model, _ := os.ReadFile(filepath.Join(basePath, "device/model"))
    dev.Model = strings.TrimSpace(string(model))
    
    return dev, nil
}

// ReadMountPoints reads mount points from /proc/mounts
func (r *DiskReader) ReadMountPoints() ([]MountPoint, error) {
    data, err := os.ReadFile("/proc/mounts")
    if err != nil {
        return nil, err
    }
    
    var mounts []MountPoint
    lines := strings.Split(string(data), "\n")
    for _, line := range lines {
        fields := strings.Fields(line)
        if len(fields) < 4 {
            continue
        }
        
        mount := MountPoint{
            Device:     fields[0],
            MountPoint: fields[1],
            FSType:     fields[2],
            Options:    fields[3],
        }
        mounts = append(mounts, mount)
    }
    return mounts, nil
}

// GetFilesystemUsage uses statvfs(2) syscall
func (r *DiskReader) GetFilesystemUsage(path string) (FilesystemUsage, error) {
    var stat syscall.Statfs_t
    err := syscall.Statfs(path, &stat)
    if err != nil {
        return FilesystemUsage{}, err
    }
    
    usage := FilesystemUsage{
        Path:      path,
        Total:     stat.Blocks * uint64(stat.Bsize),
        Free:      stat.Bfree * uint64(stat.Bsize),
        Available: stat.Bavail * uint64(stat.Bsize),
        Used:      (stat.Blocks - stat.Bfree) * uint64(stat.Bsize),
    }
    
    if usage.Total > 0 {
        usage.UsedPercent = float64(usage.Used) / float64(usage.Total) * 100
    }
    
    return usage, nil
}
```

#### 5.4 User and Group Management

```go
// UserReader reads user and group information
type UserReader struct{}

// ReadUsers reads users from /etc/passwd
func (r *UserReader) ReadUsers() ([]User, error) {
    data, err := os.ReadFile("/etc/passwd")
    if err != nil {
        return nil, err
    }
    
    var users []User
    lines := strings.Split(string(data), "\n")
    for _, line := range lines {
        if line == "" || strings.HasPrefix(line, "#") {
            continue
        }
        
        fields := strings.Split(line, ":")
        if len(fields) < 7 {
            continue
        }
        
        uid, _ := strconv.Atoi(fields[2])
        gid, _ := strconv.Atoi(fields[3])
        
        user := User{
            Name:    fields[0],
            UID:     uid,
            GID:     gid,
            Comment: fields[4],
            Home:    fields[5],
            Shell:   fields[6],
        }
        users = append(users, user)
    }
    return users, nil
}

// ReadGroups reads groups from /etc/group
func (r *UserReader) ReadGroups() ([]Group, error) {
    data, err := os.ReadFile("/etc/group")
    if err != nil {
        return nil, err
    }
    
    var groups []Group
    lines := strings.Split(string(data), "\n")
    for _, line := range lines {
        if line == "" || strings.HasPrefix(line, "#") {
            continue
        }
        
        fields := strings.Split(line, ":")
        if len(fields) < 4 {
            continue
        }
        
        gid, _ := strconv.Atoi(fields[2])
        
        members := []string{}
        if fields[3] != "" {
            members = strings.Split(fields[3], ",")
        }
        
        group := Group{
            Name:    fields[0],
            GID:     gid,
            Members: members,
        }
        groups = append(groups, group)
    }
    return groups, nil
}
```

---

## Data Models

### Core Domain Models

```go
// SystemProfile represents detected system characteristics
type SystemProfile struct {
    Distro         string
    Version        string
    InitSystem     string
    PackageManager string
    Firewall       string
    NetworkStack   string
    Environment    string
    Kernel         string
    Architecture   string
}

// Process represents a running process
type Process struct {
    PID      int
    PPID     int
    Name     string
    State    string
    Cmdline  string
    Status   map[string]string
    CPU      float64
    Memory   uint64
    Threads  int
}

// NetworkInterface represents a network interface
type NetworkInterface struct {
    Name  string
    MAC   string
    MTU   int
    State string
    Speed int
    IPv4  []string
    IPv6  []string
}

// SocketStats represents socket statistics
type SocketStats struct {
    TCP  []TCPSocket
    TCP6 []TCPSocket
    UDP  []UDPSocket
    UDP6 []UDPSocket
}

// TCPSocket represents a TCP socket
type TCPSocket struct {
    LocalAddr  string
    LocalPort  int
    RemoteAddr string
    RemotePort int
    State      string
    PID        int
}

// Route represents a routing table entry
type Route struct {
    Destination string
    Gateway     string
    Netmask     string
    Interface   string
    Metric      int
}

// BlockDevice represents a block device
type BlockDevice struct {
    Name      string
    Size      int64
    Removable bool
    Model     string
    Serial    string
}

// MountPoint represents a filesystem mount
type MountPoint struct {
    Device     string
    MountPoint string
    FSType     string
    Options    string
}

// FilesystemUsage represents filesystem usage statistics
type FilesystemUsage struct {
    Path        string
    Total       uint64
    Free        uint64
    Available   uint64
    Used        uint64
    UsedPercent float64
}

// User represents a system user
type User struct {
    Name    string
    UID     int
    GID     int
    Comment string
    Home    string
    Shell   string
}

// Group represents a system group
type Group struct {
    Name    string
    GID     int
    Members []string
}

// ServiceUnit represents a system service
type ServiceUnit struct {
    Name        string
    Status      string
    Enabled     bool
    PID         int
    Description string
    LoadState   string
    ActiveState string
    SubState    string
}

// PackageInfo represents a software package
type PackageInfo struct {
    Name        string
    Version     string
    Description string
    Installed   bool
    Manager     string
}
```

### Output Format Models

```go
// CommandOutput represents the output of a command
type CommandOutput struct {
    Success bool        `json:"success" yaml:"success"`
    Data    interface{} `json:"data,omitempty" yaml:"data,omitempty"`
    Error   string      `json:"error,omitempty" yaml:"error,omitempty"`
    Message string      `json:"message,omitempty" yaml:"message,omitempty"`
}

// TableOutput represents tabular data for console.BoxTable
type TableOutput struct {
    Title   string
    Headers []string
    Rows    [][]string
    Footer  string
}
```

---

## Correctness Properties

A property is a characteristic or behavior that should hold true across all valid executions of a system—essentially, a formal statement about what the system should do. Properties serve as the bridge between human-readable specifications and machine-verifiable correctness guarantees.

### Property Reflection

After analyzing all acceptance criteria, I identified the following testable properties. During reflection, I consolidated redundant properties:

- REQ-CLI-001 and REQ-COMPAT-001 both test backward compatibility → Combined into Property 1
- REQ-DISTRO-002 and PROP-IDEM-001 both test package installation idempotence → Combined into Property 4
- REQ-DISTRO-003 and PROP-IDEM-002 both test service start idempotence → Combined into Property 5

### Property 1: CLI Backward Compatibility

For any command that existed in the previous version, executing it with the same arguments should produce identical output and exit codes.

**Validates: Requirements REQ-CLI-001, REQ-COMPAT-001**

### Property 2: No TUI Imports

For any Go source file in the codebase, the file should not contain import statements for TUI libraries (tview, tcell, bubbletea, termui, gocui).

**Validates: Requirements REQ-TUI-002**

### Property 3: JSON Output Validity

For any command that supports the --json flag, the output should be valid, parseable JSON.

**Validates: Requirements REQ-CLI-002**

### Property 4: YAML Output Validity

For any command that supports the --yaml flag, the output should be valid, parseable YAML.

**Validates: Requirements REQ-CLI-002**

### Property 5: Distribution Detection Correctness

For any supported Linux distribution, when the ProfileEngine reads the distribution's /etc/os-release file, it should correctly identify the distribution name and package manager.

**Validates: Requirements REQ-DISTRO-001**

### Property 6: Package Installation Idempotence

For any package and any package manager, installing the package twice should result in the same final system state (package installed, no errors on second install).

**Validates: Requirements REQ-DISTRO-002, PROP-IDEM-001**

### Property 7: Service Start Idempotence

For any service and any init system, starting the service twice should result in the same final state (service running, no errors on second start).

**Validates: Requirements REQ-DISTRO-003, PROP-IDEM-002**

### Property 8: Process Kill Effectiveness

For any running process, sending SIGKILL to the process should eventually result in the process no longer existing in the process table.

**Validates: Requirements REQ-LINUX-001**

### Property 9: User File Parsing Round-Trip

For any valid /etc/passwd file content, parsing it to User structs and then formatting back to /etc/passwd format should produce equivalent content.

**Validates: Requirements REQ-LINUX-004**

### Property 10: Group File Parsing Round-Trip

For any valid /etc/group file content, parsing it to Group structs and then formatting back to /etc/group format should produce equivalent content.

**Validates: Requirements REQ-LINUX-004**

### Property 11: Agent Functionality Preservation

For any agent command that existed in the previous version, executing it should produce identical behavior (same operations, same outputs).

**Validates: Requirements REQ-AGENT-001**

### Property 12: Plugin Operation Idempotence

For any plugin, installing it twice should result in the same final state (plugin installed, no errors on second install). Similarly, removing a non-existent plugin should not cause errors.

**Validates: Requirements REQ-PLUGIN-001**

### Property 13: State-Modifying Operation Idempotence

For any state-modifying operation (install, start, enable, create), executing the operation twice should result in the same final system state.

**Validates: Requirements REQ-REL-002**

### Property 14: JSON Serialization Round-Trip

For any system data structure (SystemProfile, Process, ServiceUnit, etc.), serializing to JSON and deserializing should produce an equivalent value.

**Validates: Requirements PROP-RT-001**

### Property 15: YAML Serialization Round-Trip

For any system data structure (SystemProfile, Process, ServiceUnit, etc.), serializing to YAML and deserializing should produce an equivalent value.

**Validates: Requirements PROP-RT-002**

### Property 16: Dry-Run Safety

For any command executed with the --dry-run flag, the system state before and after execution should be identical (no modifications made).

**Validates: Requirements PROP-SAFE-001**

### Property 17: Policy Engine Enforcement

For any command that the policy engine blocks, the command should not execute and should return an error.

**Validates: Requirements PROP-SAFE-002**

---

## Error Handling

### Error Handling Strategy

All errors in TRONCLI should follow a consistent pattern that provides:
1. **Context**: What operation was being performed
2. **Cause**: Why the operation failed
3. **Remediation**: What the user can do to fix it

### Error Types

```go
// CommandError represents a command execution error
type CommandError struct {
    Command     string   // Command that failed
    Args        []string // Arguments provided
    Cause       error    // Underlying error
    Remediation string   // Suggested fix
    ExitCode    int      // Exit code (1=error, 2=usage error)
}

func (e *CommandError) Error() string {
    msg := fmt.Sprintf("Command '%s' failed: %v", e.Command, e.Cause)
    if e.Remediation != "" {
        msg += fmt.Sprintf("\nSuggestion: %s", e.Remediation)
    }
    return msg
}

// SystemError represents a system-level error
type SystemError struct {
    Operation   string // Operation that failed (e.g., "read /proc/stat")
    Cause       error  // Underlying error
    Recoverable bool   // Whether the operation can be retried
}

// DistributionError represents distribution detection errors
type DistributionError struct {
    Distro string // Distribution that caused the error
    Reason string // Why detection failed
}
```

### Error Handling Patterns

#### 1. File System Errors

```go
func readProcFile(path string) ([]byte, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        if os.IsNotExist(err) {
            return nil, &SystemError{
                Operation:   fmt.Sprintf("read %s", path),
                Cause:       err,
                Recoverable: false,
            }
        }
        if os.IsPermission(err) {
            return nil, &SystemError{
                Operation:   fmt.Sprintf("read %s", path),
                Cause:       fmt.Errorf("permission denied (try running with sudo)"),
                Recoverable: false,
            }
        }
        return nil, err
    }
    return data, nil
}
```

#### 2. Command Execution Errors

```go
func executeCommand(cmd string, args ...string) error {
    output, err := exec.Command(cmd, args...).CombinedOutput()
    if err != nil {
        return &CommandError{
            Command:     cmd,
            Args:        args,
            Cause:       fmt.Errorf("%v: %s", err, string(output)),
            Remediation: suggestRemediation(cmd, err),
            ExitCode:    1,
        }
    }
    return nil
}

func suggestRemediation(cmd string, err error) string {
    if strings.Contains(err.Error(), "not found") {
        return fmt.Sprintf("Install the package containing '%s'", cmd)
    }
    if strings.Contains(err.Error(), "permission denied") {
        return "Try running with sudo or as root"
    }
    return ""
}
```

#### 3. Distribution Detection Errors

```go
func detectDistribution() (*SystemProfile, error) {
    profile := &SystemProfile{}
    
    // Try primary method
    if err := detectFromOSRelease(profile); err == nil {
        return profile, nil
    }
    
    // Try fallback methods
    if err := detectFromLSBRelease(profile); err == nil {
        return profile, nil
    }
    
    // All methods failed
    return nil, &DistributionError{
        Distro: "unknown",
        Reason: "Could not read /etc/os-release or /etc/lsb-release",
    }
}
```

#### 4. Graceful Degradation

When optional features fail, the system should continue with reduced functionality:

```go
func getSystemInfo() SystemInfo {
    info := SystemInfo{}
    
    // Try to get CPU info (optional)
    if cpu, err := getCPUInfo(); err == nil {
        info.CPU = cpu
    } else {
        logger.Warn("Failed to get CPU info: %v", err)
    }
    
    // Try to get memory info (required)
    mem, err := getMemoryInfo()
    if err != nil {
        return SystemInfo{}, fmt.Errorf("failed to get memory info: %w", err)
    }
    info.Memory = mem
    
    return info
}
```

### Exit Codes

```go
const (
    ExitSuccess      = 0  // Command succeeded
    ExitError        = 1  // Command failed
    ExitUsageError   = 2  // Invalid usage (wrong arguments, flags)
    ExitPermission   = 3  // Permission denied
    ExitNotFound     = 4  // Resource not found
    ExitTimeout      = 5  // Operation timed out
)
```

### Error Output Format

Errors should be output in a consistent format based on the output mode:

#### Human-Readable (default)

```
Error: Failed to start service 'nginx'
Cause: Unit nginx.service not found
Suggestion: Check if nginx is installed: apt list --installed | grep nginx
```

#### JSON Format (--json)

```json
{
  "success": false,
  "error": "Failed to start service 'nginx'",
  "details": {
    "command": "systemctl start nginx",
    "cause": "Unit nginx.service not found",
    "remediation": "Check if nginx is installed: apt list --installed | grep nginx"
  }
}
```

#### YAML Format (--yaml)

```yaml
success: false
error: "Failed to start service 'nginx'"
details:
  command: "systemctl start nginx"
  cause: "Unit nginx.service not found"
  remediation: "Check if nginx is installed: apt list --installed | grep nginx"
```

---

## Testing Strategy

### Dual Testing Approach

TRONCLI will use a comprehensive testing strategy combining unit tests and property-based tests:

- **Unit tests**: Verify specific examples, edge cases, and error conditions
- **Property tests**: Verify universal properties across all inputs

Both approaches are complementary and necessary for comprehensive coverage. Unit tests catch concrete bugs in specific scenarios, while property tests verify general correctness across a wide range of inputs.

### Unit Testing

#### Scope

Unit tests should focus on:
- Specific examples that demonstrate correct behavior
- Integration points between components
- Edge cases (empty inputs, boundary values, special characters)
- Error conditions (file not found, permission denied, invalid input)

#### Guidelines

- Avoid writing too many unit tests for scenarios that property tests cover
- Use table-driven tests for multiple similar cases
- Mock external dependencies (filesystem, command execution)
- Test both success and failure paths

#### Example Unit Tests

```go
func TestProfileEngine_DetectDistro(t *testing.T) {
    tests := []struct {
        name           string
        osReleaseData  string
        expectedDistro string
        expectedVersion string
    }{
        {
            name: "Ubuntu 22.04",
            osReleaseData: `ID=ubuntu
VERSION_ID="22.04"`,
            expectedDistro: "ubuntu",
            expectedVersion: "22.04",
        },
        {
            name: "Fedora 39",
            osReleaseData: `ID=fedora
VERSION_ID=39`,
            expectedDistro: "fedora",
            expectedVersion: "39",
        },
        {
            name: "Gentoo (no version)",
            osReleaseData: `ID=gentoo`,
            expectedDistro: "gentoo",
            expectedVersion: "",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}

func TestProcessReader_ReadProcess_NotFound(t *testing.T) {
    reader := NewProcessReader()
    _, err := reader.ReadProcess(999999) // Non-existent PID
    if err == nil {
        t.Error("Expected error for non-existent process")
    }
}

func TestPackageManager_Install_AlreadyInstalled(t *testing.T) {
    // Test idempotence: installing already-installed package should succeed
    mgr := NewMockPackageManager()
    mgr.Install("curl")
    err := mgr.Install("curl") // Second install
    if err != nil {
        t.Errorf("Second install should succeed: %v", err)
    }
}
```

### Property-Based Testing

#### Library Selection

Use **gopter** (https://github.com/leanovate/gopter) for property-based testing in Go:
- Mature library with good documentation
- Supports custom generators
- Integrates well with standard Go testing
- Provides shrinking for minimal failing examples

#### Configuration

Each property test MUST:
- Run minimum 100 iterations (due to randomization)
- Include a comment tag referencing the design property
- Use appropriate generators for the domain

#### Tag Format

```go
// Feature: troncli-cli-enhancement, Property 3: JSON Output Validity
func TestProperty_JSONOutputValidity(t *testing.T) {
    // Test implementation
}
```

#### Property Test Examples

##### Property 3: JSON Output Validity

```go
// Feature: troncli-cli-enhancement, Property 3: JSON Output Validity
func TestProperty_JSONOutputValidity(t *testing.T) {
    properties := gopter.NewProperties(nil)
    
    properties.Property("all commands produce valid JSON with --json flag", 
        prop.ForAll(
            func(cmd Command) bool {
                output := executeCommand(cmd, "--json")
                var result interface{}
                err := json.Unmarshal([]byte(output), &result)
                return err == nil
            },
            genCommand(), // Generator for valid commands
        ),
    )
    
    properties.TestingRun(t, gopter.ConsoleReporter(false))
}
```

##### Property 6: Package Installation Idempotence

```go
// Feature: troncli-cli-enhancement, Property 6: Package Installation Idempotence
func TestProperty_PackageInstallIdempotence(t *testing.T) {
    properties := gopter.NewProperties(nil)
    
    properties.Property("installing package twice produces same state",
        prop.ForAll(
            func(pkg string, mgr PackageManager) bool {
                // First install
                mgr.Install(pkg)
                state1 := getSystemState()
                
                // Second install
                err := mgr.Install(pkg)
                state2 := getSystemState()
                
                // Should succeed and state should be identical
                return err == nil && state1.Equals(state2)
            },
            genPackageName(),
            genPackageManager(),
        ),
    )
    
    properties.TestingRun(t, gopter.ConsoleReporter(false))
}
```

##### Property 14: JSON Serialization Round-Trip

```go
// Feature: troncli-cli-enhancement, Property 14: JSON Serialization Round-Trip
func TestProperty_JSONRoundTrip(t *testing.T) {
    properties := gopter.NewProperties(nil)
    
    properties.Property("JSON round-trip preserves SystemProfile",
        prop.ForAll(
            func(profile SystemProfile) bool {
                // Serialize
                data, err := json.Marshal(profile)
                if err != nil {
                    return false
                }
                
                // Deserialize
                var decoded SystemProfile
                err = json.Unmarshal(data, &decoded)
                if err != nil {
                    return false
                }
                
                // Compare
                return profile.Equals(decoded)
            },
            genSystemProfile(),
        ),
    )
    
    properties.TestingRun(t, gopter.ConsoleReporter(false))
}
```

##### Property 16: Dry-Run Safety

```go
// Feature: troncli-cli-enhancement, Property 16: Dry-Run Safety
func TestProperty_DryRunSafety(t *testing.T) {
    properties := gopter.NewProperties(nil)
    
    properties.Property("dry-run never modifies system state",
        prop.ForAll(
            func(cmd Command) bool {
                stateBefore := captureSystemState()
                
                // Execute with --dry-run
                executeCommand(cmd, "--dry-run")
                
                stateAfter := captureSystemState()
                
                // State should be identical
                return stateBefore.Equals(stateAfter)
            },
            genStateModifyingCommand(),
        ),
    )
    
    properties.TestingRun(t, gopter.ConsoleReporter(false))
}
```

#### Custom Generators

```go
// Generator for valid package names
func genPackageName() gopter.Gen {
    return gen.AlphaString().
        SuchThat(func(s string) bool {
            return len(s) > 0 && len(s) < 100
        })
}

// Generator for SystemProfile
func genSystemProfile() gopter.Gen {
    return gopter.CombineGens(
        gen.OneConstOf("ubuntu", "fedora", "arch", "alpine", "gentoo", "void"),
        gen.AlphaString(),
        gen.OneConstOf("systemd", "openrc", "runit", "sysvinit"),
        gen.OneConstOf("apt", "dnf", "pacman", "apk", "portage", "xbps"),
    ).Map(func(vals []interface{}) SystemProfile {
        return SystemProfile{
            Distro:         vals[0].(string),
            Version:        vals[1].(string),
            InitSystem:     vals[2].(string),
            PackageManager: vals[3].(string),
        }
    })
}

// Generator for valid commands
func genCommand() gopter.Gen {
    return gen.OneConstOf(
        Command{Name: "system", Subcommand: "info"},
        Command{Name: "service", Subcommand: "list"},
        Command{Name: "process", Subcommand: "tree"},
        Command{Name: "network", Subcommand: "interfaces"},
    )
}
```

### Integration Testing

#### Docker-Based Testing

Test on multiple distributions using Docker containers:

```bash
# Test matrix
distributions=(
    "ubuntu:22.04"
    "ubuntu:24.04"
    "fedora:39"
    "fedora:40"
    "archlinux:latest"
    "alpine:latest"
    "opensuse/leap:15"
)

for distro in "${distributions[@]}"; do
    echo "Testing on $distro"
    docker run --rm -v $(pwd):/app $distro /app/test-in-container.sh
done
```

#### Test Script

```bash
#!/bin/bash
# test-in-container.sh

set -e

# Install dependencies
if command -v apt &> /dev/null; then
    apt update && apt install -y golang
elif command -v dnf &> /dev/null; then
    dnf install -y golang
elif command -v pacman &> /dev/null; then
    pacman -Sy --noconfirm go
fi

# Build
cd /app
go build -o troncli cmd/troncli/main.go

# Run tests
./troncli system info
./troncli service list
./troncli process tree

echo "✓ All tests passed on $(cat /etc/os-release | grep PRETTY_NAME)"
```

### Performance Benchmarking

```go
func BenchmarkStartup(b *testing.B) {
    for i := 0; i < b.N; i++ {
        cmd := exec.Command("./troncli", "--version")
        cmd.Run()
    }
}

func BenchmarkSystemInfo(b *testing.B) {
    for i := 0; i < b.N; i++ {
        cmd := exec.Command("./troncli", "system", "info")
        cmd.Run()
    }
}

func BenchmarkProcessTree(b *testing.B) {
    for i := 0; i < b.N; i++ {
        reader := NewProcessReader()
        reader.ReadProcessTree()
    }
}
```

### Test Coverage Goals

- **Unit test coverage**: 80%+ of code
- **Property test coverage**: All 17 correctness properties implemented
- **Integration test coverage**: All 7+ distributions tested
- **Performance benchmarks**: Startup < 100ms, commands < 500ms

---

## Implementation Approach

### Phase 1: TUI Removal

#### Step 1.1: Audit TUI Dependencies

```bash
# Check go.mod for TUI libraries
grep -E "tview|tcell|bubbletea|termui|gocui" go.mod

# Check for TUI imports in code
find . -name "*.go" -exec grep -l "tview\|tcell\|bubbletea" {} \;

# Check for TUI command flags
grep -r "tui\|interactive-mode" cmd/
```

#### Step 1.2: Remove internal/ui Directory

```bash
# Backup first (optional)
tar -czf internal-ui-backup.tar.gz internal/ui/

# Remove directory
rm -rf internal/ui/
```

#### Step 1.3: Remove TUI Command Flags

Search for and remove any TUI-related flags in command files:
- `--tui`
- `--interactive-mode`
- `--dashboard`

#### Step 1.4: Verify Build

```bash
# Clean build
go clean
go mod tidy

# Build
go build -o troncli cmd/troncli/main.go

# Verify no TUI dependencies
go mod graph | grep -E "tview|tcell|bubbletea"
```

### Phase 2: Enhanced Distribution Support

#### Step 2.1: Enhance ProfileEngine

Add detection for Gentoo and Void Linux:

```go
// Add to detectDistro function
func (e *ProfileEngine) detectDistro(ctx context.Context, p *domain.SystemProfile) {
    // Try /etc/os-release first
    out, err := e.executor.Exec(ctx, "cat", "/etc/os-release")
    if err == nil {
        p.Distro, p.Version = parseOSRelease(out.Stdout)
        return
    }
    
    // Fallback to distribution-specific files
    if data, err := os.ReadFile("/etc/gentoo-release"); err == nil {
        p.Distro = "gentoo"
        p.Version = parseGentooRelease(string(data))
        return
    }
    
    if data, err := os.ReadFile("/etc/void-release"); err == nil {
        p.Distro = "void"
        p.Version = parseVoidRelease(string(data))
        return
    }
}
```

#### Step 2.2: Add Package Manager Support

Extend package manager detection and operations:

```go
// Add to detectPackageManager
func (e *ProfileEngine) detectPackageManager(p *domain.SystemProfile) {
    managers := []struct {
        name   string
        binary string
    }{
        {"apt", "apt"},
        {"dnf", "dnf"},
        {"yum", "yum"},
        {"pacman", "pacman"},
        {"zypper", "zypper"},
        {"apk", "apk"},
        {"portage", "emerge"},  // New
        {"xbps", "xbps-install"}, // New
    }
    
    for _, mgr := range managers {
        if _, err := exec.LookPath(mgr.binary); err == nil {
            p.PackageManager = mgr.name
            return
        }
    }
}
```

#### Step 2.3: Implement Package Manager Adapters

Create adapters for new package managers:

```go
// internal/modules/pkg/portage.go
type PortageManager struct {
    executor Executor
}

func (m *PortageManager) Install(pkg string) error {
    return m.executor.Exec("emerge", pkg)
}

func (m *PortageManager) Remove(pkg string) error {
    return m.executor.Exec("emerge", "--unmerge", pkg)
}

func (m *PortageManager) Update() error {
    return m.executor.Exec("emerge", "--sync")
}

func (m *PortageManager) Upgrade() error {
    return m.executor.Exec("emerge", "--update", "--deep", "--newuse", "@world")
}

// internal/modules/pkg/xbps.go
type XBPSManager struct {
    executor Executor
}

func (m *XBPSManager) Install(pkg string) error {
    return m.executor.Exec("xbps-install", "-y", pkg)
}

func (m *XBPSManager) Remove(pkg string) error {
    return m.executor.Exec("xbps-remove", "-y", pkg)
}

func (m *XBPSManager) Update() error {
    return m.executor.Exec("xbps-install", "-S")
}

func (m *XBPSManager) Upgrade() error {
    return m.executor.Exec("xbps-install", "-u")
}
```

### Phase 3: Direct Linux Integration

#### Step 3.1: Implement Process Reader

Create `internal/modules/process/proc_reader.go`:

```go
package process

import (
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "syscall"
)

type ProcReader struct{}

func NewProcReader() *ProcReader {
    return &ProcReader{}
}

func (r *ProcReader) ReadProcessTree() ([]Process, error) {
    // Implementation from Components section
}

func (r *ProcReader) KillProcess(pid int, signal syscall.Signal) error {
    return syscall.Kill(pid, signal)
}

func (r *ProcReader) ReniceProcess(pid int, priority int) error {
    return syscall.Setpriority(syscall.PRIO_PROCESS, pid, priority)
}
```

#### Step 3.2: Implement Network Reader

Create `internal/modules/network/sys_reader.go`:

```go
package network

import (
    "os"
    "path/filepath"
    "strconv"
    "strings"
)

type SysReader struct{}

func NewSysReader() *SysReader {
    return &SysReader{}
}

func (r *SysReader) ReadInterfaces() ([]NetworkInterface, error) {
    // Implementation from Components section
}

func (r *SysReader) ReadSocketStats() (SocketStats, error) {
    // Implementation from Components section
}
```

#### Step 3.3: Implement Disk Reader

Create `internal/modules/disk/sys_reader.go`:

```go
package disk

import (
    "os"
    "path/filepath"
    "syscall"
)

type SysReader struct{}

func NewSysReader() *SysReader {
    return &SysReader{}
}

func (r *SysReader) ReadBlockDevices() ([]BlockDevice, error) {
    // Implementation from Components section
}

func (r *SysReader) GetFilesystemUsage(path string) (FilesystemUsage, error) {
    // Implementation from Components section
}
```

#### Step 3.4: Implement User Reader

Create `internal/modules/users/etc_reader.go`:

```go
package users

import (
    "os"
    "strings"
)

type EtcReader struct{}

func NewEtcReader() *EtcReader {
    return &EtcReader{}
}

func (r *EtcReader) ReadUsers() ([]User, error) {
    // Implementation from Components section
}

func (r *EtcReader) ReadGroups() ([]Group, error) {
    // Implementation from Components section
}
```

### Phase 4: Output Formatting Enhancement

#### Step 4.1: Enhance console.BoxTable

The existing `internal/console/table.go` is already good. Ensure it's used consistently:

```go
// Example usage in commands
func displayServices(services []ServiceUnit) {
    table := console.NewBoxTable(os.Stdout)
    table.SetTitle("System Services")
    table.SetHeaders([]string{"Name", "Status", "Enabled", "PID"})
    
    for _, svc := range services {
        table.AddRow([]string{
            svc.Name,
            svc.Status,
            fmt.Sprintf("%v", svc.Enabled),
            fmt.Sprintf("%d", svc.PID),
        })
    }
    
    table.Render()
}
```

#### Step 4.2: Ensure JSON/YAML Output

Add helper functions for structured output:

```go
// internal/console/output.go
package console

import (
    "encoding/json"
    "fmt"
    "os"
    
    "gopkg.in/yaml.v3"
)

func OutputJSON(data interface{}) error {
    encoder := json.NewEncoder(os.Stdout)
    encoder.SetIndent("", "  ")
    return encoder.Encode(data)
}

func OutputYAML(data interface{}) error {
    encoder := yaml.NewEncoder(os.Stdout)
    defer encoder.Close()
    return encoder.Encode(data)
}

func OutputTable(table *BoxTable) {
    table.Render()
}
```

#### Step 4.3: Update Commands to Support All Formats

```go
// Example command implementation
func runServiceList(cmd *cobra.Command, args []string) error {
    // Get services
    services, err := serviceManager.ListServices()
    if err != nil {
        return err
    }
    
    // Output based on flags
    if flagJSON {
        return console.OutputJSON(services)
    }
    if flagYAML {
        return console.OutputYAML(services)
    }
    
    // Default: table output
    table := console.NewBoxTable(os.Stdout)
    table.SetTitle("System Services")
    table.SetHeaders([]string{"Name", "Status", "Enabled"})
    for _, svc := range services {
        table.AddRow([]string{svc.Name, svc.Status, fmt.Sprintf("%v", svc.Enabled)})
    }
    console.OutputTable(table)
    return nil
}
```

### Phase 5: Testing Implementation

#### Step 5.1: Set Up Property Testing

```bash
# Install gopter
go get github.com/leanovate/gopter
```

#### Step 5.2: Create Test Structure

```
internal/
├── modules/
│   ├── pkg/
│   │   ├── manager.go
│   │   ├── manager_test.go          # Unit tests
│   │   └── manager_property_test.go # Property tests
│   ├── service/
│   │   ├── manager.go
│   │   ├── manager_test.go
│   │   └── manager_property_test.go
│   └── process/
│       ├── reader.go
│       ├── reader_test.go
│       └── reader_property_test.go
```

#### Step 5.3: Implement Property Tests

Create property tests for each of the 17 correctness properties identified in the design.

#### Step 5.4: Set Up Integration Tests

```bash
# Create integration test directory
mkdir -p test/integration

# Create Docker-based test script
cat > test/integration/test-all-distros.sh << 'EOF'
#!/bin/bash
set -e

distributions=(
    "ubuntu:22.04"
    "ubuntu:24.04"
    "fedora:39"
    "archlinux:latest"
    "alpine:latest"
    "opensuse/leap:15"
)

for distro in "${distributions[@]}"; do
    echo "Testing on $distro"
    docker run --rm -v $(pwd):/app $distro /app/test/integration/test-in-container.sh
done
EOF

chmod +x test/integration/test-all-distros.sh
```

### Phase 6: Documentation Updates

#### Step 6.1: Update README.md

- Remove TUI references
- Add new distribution support
- Update architecture diagram
- Add performance benchmarks

#### Step 6.2: Update COMMAND.md

- Verify all commands are documented
- Add examples for new distributions
- Document JSON/YAML output formats

#### Step 6.3: Create Migration Guide

Document changes for existing users:
- TUI removed (use CLI commands instead)
- New distributions supported
- Performance improvements
- No breaking changes to CLI

---

## Migration and Deployment

### Migration Strategy

#### Pre-Migration Checklist

- [ ] Backup current codebase
- [ ] Document all existing CLI commands and their outputs
- [ ] Create baseline performance benchmarks
- [ ] Set up test environments for all target distributions

#### Migration Phases

##### Phase 1: TUI Removal (Week 1)
- Remove internal/ui/ directory
- Remove TUI dependencies from go.mod
- Remove TUI command flags
- Verify build succeeds
- Run existing tests

##### Phase 2: Distribution Enhancement (Week 2)
- Enhance ProfileEngine for Gentoo and Void Linux
- Implement Portage package manager adapter
- Implement XBPS package manager adapter
- Add distribution-specific tests
- Test on Docker containers

##### Phase 3: Linux Integration (Week 3)
- Implement ProcReader for /proc filesystem
- Implement SysReader for /sys filesystem
- Implement EtcReader for /etc files
- Add syscall wrappers (kill, setpriority, statvfs)
- Add unit tests for each reader

##### Phase 4: Testing (Week 4)
- Implement all 17 property-based tests
- Run integration tests on all distributions
- Performance benchmarking
- Fix any issues found

##### Phase 5: Documentation and Release (Week 5)
- Update README.md
- Update COMMAND.md
- Create migration guide
- Create release notes
- Tag version 2.0.0

### Rollback Plan

If critical issues are discovered:

1. **Immediate Rollback**: Revert to previous Git tag
2. **Partial Rollback**: Keep enhancements, revert problematic changes
3. **Fix Forward**: If issues are minor, fix in place

### Deployment Strategy

#### Build Process

```bash
# Clean build
make clean

# Run tests
make test

# Run property tests
make test-properties

# Run integration tests
make test-integration

# Build for all platforms
make build-all
```

#### Release Artifacts

- Binary for Linux x86_64
- Binary for Linux aarch64
- Debian package (.deb)
- RPM package (.rpm)
- AUR package (PKGBUILD)
- Alpine package (apk)
- Source tarball

#### Distribution Channels

1. **GitHub Releases**: Primary distribution channel
2. **Package Repositories**: 
   - Debian/Ubuntu: PPA
   - Fedora/RHEL: COPR
   - Arch: AUR
   - Alpine: Community repository
3. **Container Images**: Docker Hub, GitHub Container Registry

### Monitoring and Validation

#### Post-Deployment Validation

```bash
# Verify version
troncli --version

# Test basic commands
troncli system info
troncli service list
troncli process tree

# Test JSON output
troncli system info --json | jq .

# Test YAML output
troncli system info --yaml

# Performance check
time troncli --version  # Should be < 100ms
time troncli system info  # Should be < 200ms
```

#### Success Metrics

- [ ] Zero TUI dependencies in go.mod
- [ ] All CLI commands functional
- [ ] Support for 7+ distributions verified
- [ ] 100% backward compatibility maintained
- [ ] Startup time < 100ms
- [ ] Command execution < 500ms
- [ ] 80%+ unit test coverage
- [ ] All 17 property tests passing
- [ ] Integration tests passing on all distributions

---

## Risk Assessment and Mitigation

### Technical Risks

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| Breaking existing CLI users | Low | High | Maintain 100% backward compatibility, extensive testing |
| Distribution-specific bugs | Medium | Medium | Test on multiple distributions, use Docker for CI |
| Performance regression | Low | Medium | Benchmark before/after, optimize hot paths |
| Missing edge cases | Medium | Low | Property-based testing, extensive unit tests |
| Incomplete TUI removal | Low | Low | Automated checks in CI, grep for TUI imports |

### Operational Risks

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| User confusion about TUI removal | Medium | Low | Clear documentation, migration guide |
| Package repository delays | Medium | Low | Prepare packages early, test submission process |
| Compatibility issues with old kernels | Low | Medium | Document minimum kernel version (4.0+) |

### Mitigation Strategies

#### For Breaking Changes
- Maintain strict backward compatibility
- Version all APIs and data formats
- Provide deprecation warnings before removal

#### For Distribution Issues
- Test on actual hardware, not just containers
- Engage with distribution maintainers
- Provide distribution-specific documentation

#### For Performance Issues
- Profile before and after changes
- Optimize critical paths
- Use caching where appropriate

---

## Appendix A: File Structure

### Current Structure
```
troncli/
├── cmd/
│   └── troncli/
│       ├── commands/
│       │   ├── root.go
│       │   ├── system.go
│       │   ├── service.go
│       │   ├── process.go
│       │   ├── network.go
│       │   ├── disk.go
│       │   ├── users.go
│       │   ├── pkg.go
│       │   ├── agent.go
│       │   └── ...
│       └── main.go
├── internal/
│   ├── core/
│   │   ├── domain/
│   │   │   └── profile.go
│   │   ├── ports/
│   │   │   ├── service.go
│   │   │   ├── package.go
│   │   │   ├── process.go
│   │   │   └── ...
│   │   ├── services/
│   │   │   └── profile.go
│   │   ├── adapter/
│   │   │   └── executor.go
│   │   └── logger/
│   │       └── logger.go
│   ├── modules/
│   │   ├── service/
│   │   ├── pkg/
│   │   ├── process/
│   │   ├── network/
│   │   ├── disk/
│   │   ├── users/
│   │   └── ...
│   ├── console/
│   │   └── table.go
│   ├── agent/
│   ├── policy/
│   ├── ui/              # TO BE REMOVED
│   └── voice/
├── go.mod
├── go.sum
└── README.md
```

### New Files to Add

```
internal/
├── modules/
│   ├── pkg/
│   │   ├── portage.go          # New: Gentoo support
│   │   └── xbps.go             # New: Void Linux support
│   ├── process/
│   │   └── proc_reader.go      # New: /proc filesystem reader
│   ├── network/
│   │   └── sys_reader.go       # New: /sys filesystem reader
│   ├── disk/
│   │   └── sys_reader.go       # New: /sys/block reader
│   └── users/
│       └── etc_reader.go       # New: /etc/passwd,group reader
├── console/
│   └── output.go               # New: JSON/YAML helpers
test/
├── integration/
│   ├── test-all-distros.sh     # New: Integration test runner
│   └── test-in-container.sh    # New: Container test script
└── property/
    ├── json_test.go            # New: Property tests
    ├── idempotence_test.go     # New: Property tests
    └── ...
```

---

## Appendix B: Command Reference

### All CLI Commands (Preserved)

```
troncli
├── system
│   ├── info
│   ├── uptime
│   └── resources
├── service
│   ├── list
│   ├── start <name>
│   ├── stop <name>
│   ├── restart <name>
│   ├── enable <name>
│   ├── disable <name>
│   ├── status <name>
│   └── logs <name>
├── process
│   ├── tree
│   ├── list
│   ├── kill <pid>
│   ├── renice <pid> <priority>
│   └── ports
├── network
│   ├── interfaces
│   ├── routes
│   ├── sockets
│   └── stats
├── disk
│   ├── list
│   ├── usage
│   ├── mounts
│   └── lvm
├── users
│   ├── list
│   ├── groups
│   ├── add <name>
│   └── remove <name>
├── pkg
│   ├── install <name>
│   ├── remove <name>
│   ├── update
│   ├── upgrade
│   └── search <query>
├── agent
│   ├── setup
│   ├── chat
│   └── capabilities
├── plugin
│   ├── list
│   ├── install <name>
│   └── remove <name>
└── completion
    ├── bash
    ├── zsh
    ├── fish
    └── powershell
```

### Global Flags (Preserved)

```
--json          Output in JSON format
--yaml          Output in YAML format
--quiet         Suppress output
--dry-run       Simulate without executing
--timeout       Timeout in seconds (default: 30)
--verbose       Enable verbose logging
--no-color      Disable color output
--log-file      Log file path
```

---

## Appendix C: Distribution Support Matrix

| Distribution | Package Manager | Init System | Status |
|--------------|----------------|-------------|--------|
| Ubuntu 22.04+ | apt | systemd | ✅ Supported |
| Debian 11+ | apt | systemd | ✅ Supported |
| Fedora 38+ | dnf | systemd | ✅ Supported |
| RHEL 8+ | dnf | systemd | ✅ Supported |
| CentOS 7 | yum | systemd | ✅ Supported |
| Arch Linux | pacman | systemd | ✅ Supported |
| Alpine Linux | apk | openrc | ✅ Supported |
| openSUSE Leap | zypper | systemd | ✅ Supported |
| Gentoo | portage | openrc | 🆕 New |
| Void Linux | xbps | runit | 🆕 New |

---

## Appendix D: Performance Targets

### Startup Performance

| Metric | Target | Measurement |
|--------|--------|-------------|
| Cold start | < 100ms | `time troncli --version` |
| Warm start | < 50ms | Second execution |
| Binary size | < 20MB | `ls -lh troncli` |

### Command Performance

| Command | Target | Notes |
|---------|--------|-------|
| `system info` | < 200ms | Reads /proc, /sys |
| `service list` | < 1000ms | Queries systemd |
| `process tree` | < 500ms | Reads /proc |
| `network interfaces` | < 300ms | Reads /sys/class/net |
| `disk usage` | < 400ms | Reads /proc/mounts, statvfs |

### Resource Usage

| Metric | Target |
|--------|--------|
| Memory (idle) | < 10MB |
| Memory (peak) | < 50MB |
| CPU (idle) | 0% |
| CPU (peak) | < 50% single core |

---

## Summary

This design document provides a comprehensive technical specification for removing TUI components from TRONCLI and enhancing it as a pure CLI tool with improved multi-distribution support and direct Linux integration.

Key design decisions:
1. **Preservation over rewrite**: Keep existing architecture and CLI functionality
2. **Enhancement over replacement**: Add new distribution support without breaking existing support
3. **Direct integration**: Use /proc, /sys, and syscalls for better performance and reliability
4. **Comprehensive testing**: Property-based tests for correctness, integration tests for compatibility
5. **Backward compatibility**: 100% compatibility for all existing CLI commands and flags

The implementation follows a phased approach over 5 weeks, with clear success metrics and rollback plans. All 17 correctness properties will be implemented as property-based tests using gopter, ensuring the system behaves correctly across a wide range of inputs and scenarios.
