# TRONCLI CLI Enhancement - Requirements Document

**Feature:** Remove TUI and Enhance Multi-Distribution CLI Support  
**Spec Type:** Feature  
**Workflow:** Requirements-First  
**Created:** 2026-02-25

---

## 1. Executive Summary

This specification defines the requirements for removing the Terminal User Interface (TUI) components from TRONCLI and enhancing it as a pure Command-Line Interface (CLI) tool with improved multi-distribution Linux support and direct interaction with all Linux subsystems.

### 1.1 Goals

1. **Remove TUI Completely**: Eliminate all TUI-related code, dependencies, and references
2. **Preserve CLI Functionality**: Maintain all existing CLI commands and features
3. **Enhance Multi-Distribution Support**: Improve compatibility across Linux distributions
4. **Direct Linux Integration**: Strengthen direct interaction with Linux subsystems
5. **Maintain Architecture**: Preserve Clean Architecture principles

### 1.2 Success Criteria

- Zero TUI dependencies in go.mod
- All CLI commands functional and tested
- Support for 7+ major Linux distributions
- 100% backward compatibility for CLI commands
- Performance improvements (startup < 100ms, commands < 500ms)

---

## 2. Functional Requirements

### 2.1 TUI Removal Requirements

#### REQ-TUI-001: Remove TUI Dependencies
**Type:** Functional  
**Priority:** Critical  

**Description:**  
WHEN the system builds the project,  
THEN it SHALL NOT include any TUI-related dependencies (tview, tcell, bubbletea, etc.) in go.mod.

**Acceptance Criteria:**
- [ ] go.mod contains no TUI library dependencies
- [ ] go.sum contains no TUI library checksums
- [ ] Build completes without TUI-related imports

**Correctness Property:**
```
Property: NoTUIDependencies
∀ dependency ∈ go.mod:
  dependency.name ∉ {tview, tcell, bubbletea, termui, gocui}
```

---

#### REQ-TUI-002: Remove TUI Source Code
**Type:** Functional  
**Priority:** Critical  

**Description:**  
WHEN the system is analyzed for TUI code,  
THEN it SHALL NOT contain any TUI implementation files or directories.

**Acceptance Criteria:**
- [ ] internal/ui/ directory is removed
- [ ] No files import TUI libraries
- [ ] No TUI-related command flags (--tui, --interactive-mode)

**Correctness Property:**
```
Property: NoTUICode
∀ file ∈ codebase:
  ¬∃ import ∈ file.imports: import.path.contains("tview") ∨
                             import.path.contains("tcell") ∨
                             import.path.contains("bubbletea")
```

---

#### REQ-TUI-003: Remove TUI Commands
**Type:** Functional  
**Priority:** High  

**Description:**  
WHEN the user lists available commands,  
THEN the system SHALL NOT display any TUI-specific commands or subcommands.

**Acceptance Criteria:**
- [ ] No 'troncli tui' command exists
- [ ] No 'troncli dashboard' command exists
- [ ] Help output shows only CLI commands

---

### 2.2 CLI Preservation Requirements

#### REQ-CLI-001: Preserve All CLI Commands
**Type:** Functional  
**Priority:** Critical  

**Description:**  
WHEN the user executes any existing CLI command,  
THEN the system SHALL execute it with identical behavior to the previous version.

**Acceptance Criteria:**
- [ ] All commands in COMMAND.md are functional
- [ ] Command syntax remains unchanged
- [ ] Output format remains consistent

**Correctness Property:**
```
Property: CLIBackwardCompatibility
∀ command ∈ legacy_commands:
  execute(command, args) = legacy_execute(command, args)
```

---

#### REQ-CLI-002: Preserve Global Flags
**Type:** Functional  
**Priority:** Critical  

**Description:**  
WHEN the user provides global flags,  
THEN the system SHALL honor all flags: --json, --yaml, --quiet, --dry-run, --timeout, --verbose, --no-color, --log-file.

**Acceptance Criteria:**
- [ ] --json produces valid JSON output
- [ ] --yaml produces valid YAML output
- [ ] --quiet suppresses non-error output
- [ ] --dry-run simulates without executing
- [ ] --verbose shows detailed logs
- [ ] --no-color removes ANSI codes
- [ ] --timeout enforces time limits

**Correctness Property:**
```
Property: JSONOutputValidity
∀ command ∈ commands:
  output = execute(command, ["--json"])
  ⇒ isValidJSON(output)

Property: YAMLOutputValidity
∀ command ∈ commands:
  output = execute(command, ["--yaml"])
  ⇒ isValidYAML(output)
```

---

#### REQ-CLI-003: Preserve Output Formatting
**Type:** Functional  
**Priority:** High  

**Description:**  
WHEN the user executes commands without format flags,  
THEN the system SHALL output human-readable formatted tables using the console.BoxTable component.

**Acceptance Criteria:**
- [ ] Table output is properly aligned
- [ ] Headers and footers are displayed
- [ ] Long text is truncated with "..."
- [ ] Colors work unless --no-color is set

---

### 2.3 Multi-Distribution Support Requirements

#### REQ-DISTRO-001: Detect Linux Distribution
**Type:** Functional  
**Priority:** Critical  

**Description:**  
WHEN the system starts,  
THEN it SHALL automatically detect the Linux distribution and version.

**Acceptance Criteria:**
- [ ] Detects Ubuntu/Debian (apt)
- [ ] Detects RHEL/CentOS/Fedora (dnf/yum)
- [ ] Detects Arch Linux (pacman)
- [ ] Detects Alpine Linux (apk)
- [ ] Detects openSUSE (zypper)
- [ ] Detects Gentoo (portage)
- [ ] Detects Void Linux (xbps)

**Correctness Property:**
```
Property: DistributionDetection
∀ system ∈ supported_distributions:
  profile = detectProfile(system)
  ⇒ profile.distribution = system.actual_distribution ∧
    profile.packageManager = system.actual_package_manager
```

---

#### REQ-DISTRO-002: Universal Package Management
**Type:** Functional  
**Priority:** Critical  

**Description:**  
WHEN the user installs a package,  
THEN the system SHALL use the appropriate package manager for the detected distribution.

**Acceptance Criteria:**
- [ ] apt install works on Debian/Ubuntu
- [ ] dnf install works on Fedora/RHEL 8+
- [ ] yum install works on CentOS/RHEL 7
- [ ] pacman -S works on Arch
- [ ] apk add works on Alpine
- [ ] zypper install works on openSUSE

**Correctness Property:**
```
Property: PackageInstallIdempotence
∀ package ∈ packages:
  install(package); install(package)
  ⇒ packageInstalled(package) ∧ systemState = systemState_after_first_install
```

---

#### REQ-DISTRO-003: Universal Service Management
**Type:** Functional  
**Priority:** Critical  

**Description:**  
WHEN the user manages services,  
THEN the system SHALL use the appropriate init system (systemd, openrc, sysvinit, runit).

**Acceptance Criteria:**
- [ ] systemctl commands work on systemd systems
- [ ] rc-service commands work on OpenRC systems
- [ ] service commands work on SysVinit systems
- [ ] sv commands work on runit systems

**Correctness Property:**
```
Property: ServiceStartIdempotence
∀ service ∈ services:
  start(service); start(service)
  ⇒ serviceRunning(service) ∧ ¬error_occurred
```

---

### 2.4 Direct Linux Subsystem Integration

#### REQ-LINUX-001: Process Management
**Type:** Functional  
**Priority:** High  

**Description:**  
WHEN the user manages processes,  
THEN the system SHALL interact directly with /proc filesystem and process signals.

**Acceptance Criteria:**
- [ ] Process tree reads from /proc
- [ ] Kill sends signals via kill(2) syscall
- [ ] Renice uses setpriority(2) syscall
- [ ] Open files read from /proc/[pid]/fd
- [ ] Ports read from /proc/net/tcp and /proc/net/udp

**Correctness Property:**
```
Property: ProcessKillEffectiveness
∀ pid ∈ running_processes:
  killProcess(pid, "SIGKILL")
  ⇒ eventually(¬processExists(pid))
```

---

#### REQ-LINUX-002: Network Management
**Type:** Functional  
**Priority:** High  

**Description:**  
WHEN the user manages network interfaces,  
THEN the system SHALL interact directly with netlink sockets and /sys/class/net.

**Acceptance Criteria:**
- [ ] Interface info reads from /sys/class/net
- [ ] IP configuration uses ip command or netlink
- [ ] Socket states read from /proc/net/tcp*
- [ ] Route table reads from /proc/net/route

---

#### REQ-LINUX-003: Disk and Filesystem Management
**Type:** Functional  
**Priority:** High  

**Description:**  
WHEN the user manages disks and filesystems,  
THEN the system SHALL interact directly with /sys/block, /proc/mounts, and LVM tools.

**Acceptance Criteria:**
- [ ] Block devices read from /sys/block
- [ ] Mount points read from /proc/mounts
- [ ] LVM uses pvs, vgs, lvs commands
- [ ] Filesystem usage uses statvfs(2) syscall

---

#### REQ-LINUX-004: User and Group Management
**Type:** Functional  
**Priority:** High  

**Description:**  
WHEN the user manages users and groups,  
THEN the system SHALL interact with /etc/passwd, /etc/group, and shadow files.

**Acceptance Criteria:**
- [ ] User list reads from /etc/passwd
- [ ] Group list reads from /etc/group
- [ ] User creation uses useradd command
- [ ] Group creation uses groupadd command

---

### 2.5 AI Agent Integration

#### REQ-AGENT-001: Preserve AI Agent Functionality
**Type:** Functional  
**Priority:** Medium  

**Description:**  
WHEN the user interacts with AI agents,  
THEN the system SHALL maintain all agent capabilities (Ollama, Claude, OpenAI, LlamaCpp, Local).

**Acceptance Criteria:**
- [ ] Agent configuration loads from ~/.troncli/agent_config.yaml
- [ ] Capabilities registry enforces allowed/blocked intents
- [ ] All agent adapters remain functional
- [ ] Agent commands execute system operations

---

### 2.6 Plugin System

#### REQ-PLUGIN-001: Preserve Plugin System
**Type:** Functional  
**Priority:** Medium  

**Description:**  
WHEN the user manages plugins,  
THEN the system SHALL support plugin installation, removal, and listing.

**Acceptance Criteria:**
- [ ] Plugin list shows installed plugins
- [ ] Plugin install downloads and activates plugins
- [ ] Plugin remove deactivates and deletes plugins

---

## 3. Non-Functional Requirements

### 3.1 Performance

#### REQ-PERF-001: Startup Performance
**Type:** Non-Functional  
**Priority:** High  

**Description:**  
The system SHALL start in less than 100 milliseconds on modern hardware.

**Acceptance Criteria:**
- [ ] Cold start < 100ms (measured with `time troncli --version`)
- [ ] Warm start < 50ms

---

#### REQ-PERF-002: Command Execution Performance
**Type:** Non-Functional  
**Priority:** High  

**Description:**  
The system SHALL execute commands in less than 500 milliseconds (excluding network/disk I/O).

**Acceptance Criteria:**
- [ ] `troncli system info` < 200ms
- [ ] `troncli process tree` < 500ms
- [ ] `troncli service list` < 1000ms (acceptable due to systemd query)

---

### 3.2 Reliability

#### REQ-REL-001: Error Handling
**Type:** Non-Functional  
**Priority:** Critical  

**Description:**  
The system SHALL provide clear, actionable error messages for all failure scenarios.

**Acceptance Criteria:**
- [ ] Errors include context (command, arguments, system state)
- [ ] Errors suggest remediation steps
- [ ] Errors use exit codes (0=success, 1=error, 2=usage error)

---

#### REQ-REL-002: Idempotence
**Type:** Non-Functional  
**Priority:** Critical  

**Description:**  
The system SHALL ensure idempotent operations for all state-modifying commands.

**Acceptance Criteria:**
- [ ] Installing an installed package succeeds without changes
- [ ] Starting a running service succeeds without changes
- [ ] Creating an existing user fails gracefully

**Correctness Property:**
```
Property: OperationIdempotence
∀ operation ∈ state_modifying_operations:
  state1 = execute(operation)
  state2 = execute(operation)
  ⇒ state1 = state2
```

---

### 3.3 Maintainability

#### REQ-MAINT-001: Clean Architecture
**Type:** Non-Functional  
**Priority:** High  

**Description:**  
The system SHALL maintain Clean Architecture with proper dependency flow: cmd → internal/core → internal/modules.

**Acceptance Criteria:**
- [ ] No circular dependencies
- [ ] Ports (interfaces) defined in internal/core/ports
- [ ] Implementations in internal/modules
- [ ] Commands in cmd/troncli/commands

---

#### REQ-MAINT-002: Code Quality
**Type:** Non-Functional  
**Priority:** High  

**Description:**  
The system SHALL pass all linting and formatting checks.

**Acceptance Criteria:**
- [ ] `golangci-lint run` passes with zero errors
- [ ] `go fmt` produces no changes
- [ ] `go vet` reports no issues

---

### 3.4 Compatibility

#### REQ-COMPAT-001: Backward Compatibility
**Type:** Non-Functional  
**Priority:** Critical  

**Description:**  
The system SHALL maintain 100% backward compatibility for all CLI commands and flags.

**Acceptance Criteria:**
- [ ] All commands in COMMAND.md work identically
- [ ] All flags produce identical output
- [ ] Exit codes remain unchanged

---

#### REQ-COMPAT-002: Shell Completion
**Type:** Non-Functional  
**Priority:** Medium  

**Description:**  
The system SHALL provide shell completion for bash, zsh, fish, and powershell.

**Acceptance Criteria:**
- [ ] `troncli completion bash` generates valid bash completion
- [ ] `troncli completion zsh` generates valid zsh completion
- [ ] `troncli completion fish` generates valid fish completion
- [ ] `troncli completion powershell` generates valid powershell completion

---

## 4. Correctness Properties

### 4.1 Round-Trip Properties

#### PROP-RT-001: JSON Round-Trip
```
Property: JSONRoundTrip
∀ data ∈ system_data:
  json_str = toJSON(data)
  data' = fromJSON(json_str)
  ⇒ data = data'
```

#### PROP-RT-002: YAML Round-Trip
```
Property: YAMLRoundTrip
∀ data ∈ system_data:
  yaml_str = toYAML(data)
  data' = fromYAML(yaml_str)
  ⇒ data = data'
```

---

### 4.2 Idempotence Properties

#### PROP-IDEM-001: Package Installation Idempotence
```
Property: PackageInstallIdempotence
∀ package ∈ packages:
  install(package)
  state1 = getSystemState()
  install(package)
  state2 = getSystemState()
  ⇒ state1 = state2
```

#### PROP-IDEM-002: Service Start Idempotence
```
Property: ServiceStartIdempotence
∀ service ∈ services:
  start(service)
  state1 = getServiceState(service)
  start(service)
  state2 = getServiceState(service)
  ⇒ state1 = state2 ∧ state1.running = true
```

---

### 4.3 Safety Properties

#### PROP-SAFE-001: No Data Loss on Dry-Run
```
Property: DryRunSafety
∀ command ∈ commands:
  state_before = getSystemState()
  execute(command, ["--dry-run"])
  state_after = getSystemState()
  ⇒ state_before = state_after
```

#### PROP-SAFE-002: Policy Engine Enforcement
```
Property: PolicyEnforcement
∀ command ∈ dangerous_commands:
  ¬policyEngine.allows(command)
  ⇒ ¬execute(command)
```

---

## 5. Out of Scope

The following items are explicitly OUT OF SCOPE for this specification:

1. **TUI Functionality**: No terminal user interface, dashboards, or interactive menus
2. **GUI Development**: No graphical user interface
3. **Web Interface**: No web-based management interface
4. **Windows/macOS Support**: Linux-only tool
5. **New Features**: Focus is on removal and enhancement, not new capabilities
6. **Breaking Changes**: No changes to existing CLI command syntax

---

## 6. Dependencies

### 6.1 External Dependencies

- Go 1.24+
- Linux kernel 4.0+ (for /proc, /sys interfaces)
- Standard Linux utilities (ps, kill, ip, ss, etc.)

### 6.2 Go Module Dependencies

**Keep:**
- github.com/spf13/cobra (CLI framework)
- gopkg.in/yaml.v3 (YAML parsing)
- github.com/kevinburke/ssh_config (SSH configuration)

**Remove:**
- Any TUI libraries (tview, tcell, bubbletea, etc.)

---

## 7. Testing Strategy

### 7.1 Unit Tests

- Test all modules independently
- Mock system calls and command execution
- Achieve 80%+ code coverage

### 7.2 Integration Tests

- Test on multiple Linux distributions (Docker containers)
- Test all CLI commands end-to-end
- Verify output formats (JSON, YAML, table)

### 7.3 Property-Based Tests

- Implement all correctness properties using rapid/gopter
- Test idempotence properties
- Test round-trip properties for parsers

### 7.4 Manual Testing

- Test on physical machines with different distributions
- Verify performance benchmarks
- Test shell completion in real shells

---

## 8. Migration Plan

### 8.1 Phase 1: TUI Removal
1. Identify all TUI dependencies
2. Remove internal/ui/ directory
3. Remove TUI imports from commands
4. Update go.mod and go.sum
5. Verify build succeeds

### 8.2 Phase 2: CLI Enhancement
1. Audit all CLI commands
2. Add missing distribution support
3. Enhance error messages
4. Improve performance

### 8.3 Phase 3: Testing
1. Write property-based tests
2. Test on multiple distributions
3. Performance benchmarking
4. Documentation updates

### 8.4 Phase 4: Release
1. Update README.md
2. Update COMMAND.md
3. Create release notes
4. Tag version 2.0.0

---

## 9. Success Metrics

- **TUI Removal**: Zero TUI dependencies, zero TUI code
- **CLI Preservation**: 100% backward compatibility
- **Multi-Distro Support**: Tested on 7+ distributions
- **Performance**: Startup < 100ms, commands < 500ms
- **Test Coverage**: 80%+ unit test coverage
- **Property Tests**: All correctness properties pass

---

## 10. Risks and Mitigations

| Risk | Impact | Probability | Mitigation |
|------|--------|-------------|------------|
| Breaking existing CLI users | High | Low | Maintain 100% backward compatibility |
| Distribution-specific bugs | Medium | Medium | Test on multiple distributions |
| Performance regression | Medium | Low | Benchmark before/after |
| Missing TUI dependencies | Low | Low | Thorough dependency audit |

---

## Document History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2026-02-25 | Kiro | Initial requirements document |

