# Implementation Plan: TRONCLI CLI Enhancement

## Overview

This implementation plan converts the TRONCLI CLI Enhancement design into actionable coding tasks. The project removes all TUI components and enhances TRONCLI as a pure CLI tool with improved multi-distribution Linux support (adding Gentoo and Void Linux), direct Linux subsystem integration (/proc, /sys, syscalls), and comprehensive property-based testing.

The implementation follows 5 phases over 5 weeks, maintaining 100% backward compatibility with existing CLI commands while adding new capabilities.

## Tasks

- [x] 1. Phase 1: TUI Removal and Verification
  - [x] 1.1 Audit codebase for TUI dependencies
    - Run grep commands to find TUI library imports (tview, tcell, bubbletea, termui, gocui)
    - Check go.mod for TUI dependencies
    - Search for TUI-related command flags (--tui, --interactive-mode, --dashboard)
    - Document findings in audit report
    - _Requirements: REQ-TUI-001, REQ-TUI-002_

  - [x] 1.2 Remove internal/ui directory and TUI code
    - Remove internal/ui/ directory completely
    - Remove any TUI-related command flags from cmd/troncli/commands/
    - Remove TUI imports from all Go files
    - _Requirements: REQ-TUI-001, REQ-TUI-002_

  - [x] 1.3 Clean up dependencies and verify build
    - Run `go mod tidy` to clean dependencies
    - Build project with `go build -o troncli cmd/troncli/main.go`
    - Verify no TUI dependencies remain with `go mod graph`
    - Run existing tests to ensure nothing broke
    - _Requirements: REQ-TUI-002, REQ-COMPAT-001_

  - [ ]* 1.4 Write property test for TUI removal
    - **Property 2: No TUI Imports**
    - **Validates: Requirements REQ-TUI-002**
    - Create test/property/tui_test.go
    - Scan all .go files to verify no TUI library imports exist
    - Use gopter to test across all source files

- [x] 2. Checkpoint - Verify TUI removal complete
  - Ensure all tests pass, ask the user if questions arise.

- [x] 3. Phase 2: Enhanced Distribution Support
  - [x] 3.1 Enhance ProfileEngine for Gentoo and Void Linux detection
    - Modify internal/core/services/profile.go
    - Add detection for /etc/gentoo-release
    - Add detection for /etc/void-release
    - Add parseGentooRelease() and parseVoidRelease() functions
    - Update detectDistro() to include new distributions
    - _Requirements: REQ-DISTRO-001_

  - [x] 3.2 Add package manager detection for Portage and XBPS
    - Modify detectPackageManager() in internal/core/services/profile.go
    - Add detection for "emerge" binary (Portage)
    - Add detection for "xbps-install" binary (XBPS)
    - Update package manager list to include portage and xbps
    - _Requirements: REQ-DISTRO-001_

  - [ ]* 3.3 Write unit tests for distribution detection
    - Create internal/core/services/profile_test.go
    - Test Ubuntu, Fedora, Arch, Alpine, Gentoo, Void detection
    - Test package manager detection for all 8 managers
    - Test fallback behavior when files don't exist
    - _Requirements: REQ-DISTRO-001_

  - [ ]* 3.4 Write property test for distribution detection
    - **Property 5: Distribution Detection Correctness**
    - **Validates: Requirements REQ-DISTRO-001**
    - Create test/property/distro_test.go
    - Generate mock /etc/os-release content for all distributions
    - Verify correct distribution and package manager identification

- [x] 4. Phase 2: Package Manager Adapters
  - [x] 4.1 Implement Portage package manager adapter
    - Create internal/modules/pkg/portage.go
    - Implement PortageManager struct with Executor
    - Implement Install(), Remove(), Update(), Upgrade() methods
    - Map emerge commands correctly (emerge, emerge --unmerge, emerge --sync)
    - _Requirements: REQ-DISTRO-002_

  - [x] 4.2 Implement XBPS package manager adapter
    - Create internal/modules/pkg/xbps.go
    - Implement XBPSManager struct with Executor
    - Implement Install(), Remove(), Update(), Upgrade() methods
    - Map xbps commands correctly (xbps-install, xbps-remove)
    - _Requirements: REQ-DISTRO-002_

  - [x] 4.3 Update universal package manager to use new adapters
    - Modify internal/modules/pkg/manager.go
    - Add cases for "portage" and "xbps" in Install(), Remove(), Update(), Upgrade()
    - Ensure proper error handling for unsupported managers
    - _Requirements: REQ-DISTRO-002_

  - [ ]* 4.4 Write unit tests for package manager adapters
    - Test Portage adapter with mock executor
    - Test XBPS adapter with mock executor
    - Test error conditions and edge cases
    - _Requirements: REQ-DISTRO-002_

  - [ ]* 4.5 Write property test for package installation idempotence
    - **Property 6: Package Installation Idempotence**
    - **Validates: Requirements REQ-DISTRO-002, PROP-IDEM-001**
    - Create test/property/package_test.go
    - Test installing package twice produces same state
    - Use gopter with custom package name generator

- [x] 5. Phase 2: Service Manager Adapters
  - [x] 5.1 Enhance universal service manager for multiple init systems
    - Modify internal/modules/service/manager.go
    - Add support for openrc (rc-service commands)
    - Add support for runit (sv commands)
    - Add support for sysvinit (service commands)
    - Update StartService(), StopService(), RestartService() for all init systems
    - _Requirements: REQ-DISTRO-003_

  - [ ]* 5.2 Write unit tests for service manager adapters
    - Test systemd, openrc, runit, sysvinit adapters
    - Test all service operations (start, stop, restart, enable, disable)
    - Test error conditions
    - _Requirements: REQ-DISTRO-003_

  - [ ]* 5.3 Write property test for service start idempotence
    - **Property 7: Service Start Idempotence**
    - **Validates: Requirements REQ-DISTRO-003, PROP-IDEM-002**
    - Create test/property/service_test.go
    - Test starting service twice produces same state
    - Use gopter with custom service name generator

- [x] 6. Checkpoint - Verify distribution support complete
  - Ensure all tests pass, ask the user if questions arise.

- [x] 7. Phase 3: Direct Linux Integration - Process Management
  - [x] 7.1 Implement ProcReader for /proc filesystem
    - Create internal/modules/process/proc_reader.go
    - Implement ProcReader struct
    - Implement ReadProcessTree() to read from /proc
    - Parse /proc/[pid]/stat, /proc/[pid]/cmdline, /proc/[pid]/status
    - Handle process enumeration and parsing errors gracefully
    - _Requirements: REQ-LINUX-001_

  - [x] 7.2 Implement process control syscalls
    - Add KillProcess(pid, signal) using syscall.Kill
    - Add ReniceProcess(pid, priority) using syscall.Setpriority
    - Add ReadOpenFiles(pid) to read /proc/[pid]/fd
    - Implement proper error handling for permission denied
    - _Requirements: REQ-LINUX-001_

  - [ ]* 7.3 Write unit tests for ProcReader
    - Test ReadProcessTree() with mock /proc data
    - Test process parsing edge cases (missing fields, invalid PIDs)
    - Test error conditions (non-existent PID, permission denied)
    - _Requirements: REQ-LINUX-001_

  - [ ]* 7.4 Write property test for process kill effectiveness
    - **Property 8: Process Kill Effectiveness**
    - **Validates: Requirements REQ-LINUX-001**
    - Create test/property/process_test.go
    - Test that SIGKILL removes process from process table
    - Use gopter with process lifecycle testing

- [x] 8. Phase 3: Direct Linux Integration - Network Management
  - [x] 8.1 Implement SysReader for network interfaces
    - Create internal/modules/network/sys_reader.go
    - Implement SysReader struct
    - Implement ReadInterfaces() to read from /sys/class/net
    - Parse interface details (MAC, MTU, state, speed)
    - _Requirements: REQ-LINUX-002_

  - [x] 8.2 Implement network statistics readers
    - Implement ReadSocketStats() to parse /proc/net/tcp, /proc/net/udp
    - Implement ReadRouteTable() to parse /proc/net/route
    - Create parseTCPSockets(), parseUDPSockets(), parseRouteTable() helpers
    - Handle IPv4 and IPv6 sockets
    - _Requirements: REQ-LINUX-002_

  - [ ]* 8.3 Write unit tests for network readers
    - Test ReadInterfaces() with mock /sys data
    - Test socket statistics parsing
    - Test route table parsing
    - Test error conditions
    - _Requirements: REQ-LINUX-002_

- [x] 9. Phase 3: Direct Linux Integration - Disk Management
  - [x] 9.1 Implement SysReader for block devices
    - Create internal/modules/disk/sys_reader.go
    - Implement SysReader struct
    - Implement ReadBlockDevices() to read from /sys/block
    - Parse device details (size, removable flag, model)
    - _Requirements: REQ-LINUX-003_

  - [x] 9.2 Implement filesystem readers
    - Implement ReadMountPoints() to parse /proc/mounts
    - Implement GetFilesystemUsage(path) using syscall.Statfs
    - Calculate usage statistics (total, free, available, used, percent)
    - _Requirements: REQ-LINUX-003_

  - [ ]* 9.3 Write unit tests for disk readers
    - Test ReadBlockDevices() with mock /sys data
    - Test ReadMountPoints() with mock /proc/mounts
    - Test GetFilesystemUsage() with various paths
    - Test error conditions
    - _Requirements: REQ-LINUX-003_

- [x] 10. Phase 3: Direct Linux Integration - User Management
  - [x] 10.1 Implement EtcReader for users and groups
    - Create internal/modules/users/etc_reader.go
    - Implement EtcReader struct
    - Implement ReadUsers() to parse /etc/passwd
    - Implement ReadGroups() to parse /etc/group
    - Parse user fields (name, UID, GID, comment, home, shell)
    - Parse group fields (name, GID, members)
    - _Requirements: REQ-LINUX-004_

  - [ ]* 10.2 Write unit tests for user readers
    - Test ReadUsers() with mock /etc/passwd data
    - Test ReadGroups() with mock /etc/group data
    - Test parsing edge cases (empty fields, special characters)
    - Test error conditions (file not found, permission denied)
    - _Requirements: REQ-LINUX-004_

  - [ ]* 10.3 Write property tests for user/group file parsing
    - **Property 9: User File Parsing Round-Trip**
    - **Property 10: Group File Parsing Round-Trip**
    - **Validates: Requirements REQ-LINUX-004**
    - Create test/property/users_test.go
    - Test parsing and formatting /etc/passwd produces equivalent content
    - Test parsing and formatting /etc/group produces equivalent content

- [x] 11. Checkpoint - Verify Linux integration complete
  - Ensure all tests pass, ask the user if questions arise.

- [~] 12. Phase 4: Output Formatting Enhancement
  - [x] 12.1 Create JSON/YAML output helpers
    - Create internal/console/output.go
    - Implement OutputJSON(data) using json.Encoder
    - Implement OutputYAML(data) using yaml.Encoder
    - Implement OutputTable(table) wrapper for BoxTable
    - Add proper indentation for JSON (2 spaces)
    - _Requirements: REQ-CLI-002_

  - [~] 12.2 Update all commands to support JSON/YAML output
    - Modify cmd/troncli/commands/system.go for --json/--yaml
    - Modify cmd/troncli/commands/service.go for --json/--yaml
    - Modify cmd/troncli/commands/process.go for --json/--yaml
    - Modify cmd/troncli/commands/network.go for --json/--yaml
    - Modify cmd/troncli/commands/disk.go for --json/--yaml
    - Modify cmd/troncli/commands/users.go for --json/--yaml
    - Modify cmd/troncli/commands/pkg.go for --json/--yaml
    - Ensure all commands check flagJSON and flagYAML before output
    - _Requirements: REQ-CLI-002_

  - [ ]* 12.3 Write property tests for output format validity
    - **Property 3: JSON Output Validity**
    - **Property 4: YAML Output Validity**
    - **Validates: Requirements REQ-CLI-002**
    - Create test/property/output_test.go
    - Test all commands produce valid JSON with --json flag
    - Test all commands produce valid YAML with --yaml flag
    - Use gopter with command generator

- [~] 13. Phase 4: Error Handling Enhancement
  - [~] 13.1 Implement error types
    - Create internal/core/errors.go
    - Implement CommandError struct with context and remediation
    - Implement SystemError struct for system-level errors
    - Implement DistributionError struct for detection errors
    - Add Error() methods for all error types
    - _Requirements: REQ-REL-001_

  - [~] 13.2 Add error handling patterns to all modules
    - Update process module with proper error handling
    - Update network module with proper error handling
    - Update disk module with proper error handling
    - Update users module with proper error handling
    - Add suggestRemediation() helper function
    - Implement graceful degradation for optional features
    - _Requirements: REQ-REL-001_

  - [~] 13.3 Implement structured error output
    - Update console/output.go to handle errors in JSON/YAML
    - Format errors consistently in human-readable mode
    - Include error details, cause, and remediation in structured output
    - _Requirements: REQ-REL-001_

- [ ] 14. Phase 4: Property-Based Testing - Core Properties
  - [ ]* 14.1 Write property test for CLI backward compatibility
    - **Property 1: CLI Backward Compatibility**
    - **Validates: Requirements REQ-CLI-001, REQ-COMPAT-001**
    - Create test/property/compatibility_test.go
    - Test existing commands produce identical output and exit codes
    - Compare against baseline outputs

  - [ ]* 14.2 Write property tests for serialization round-trips
    - **Property 14: JSON Serialization Round-Trip**
    - **Property 15: YAML Serialization Round-Trip**
    - **Validates: Requirements PROP-RT-001, PROP-RT-002**
    - Create test/property/serialization_test.go
    - Test SystemProfile, Process, ServiceUnit, etc. round-trip through JSON
    - Test same structures round-trip through YAML
    - Use gopter with custom generators for all data structures

  - [ ]* 14.3 Write property test for dry-run safety
    - **Property 16: Dry-Run Safety**
    - **Validates: Requirements PROP-SAFE-001**
    - Create test/property/dryrun_test.go
    - Capture system state before and after --dry-run commands
    - Verify state is identical (no modifications made)
    - Test across all state-modifying commands

  - [ ]* 14.4 Write property test for policy engine enforcement
    - **Property 17: Policy Engine Enforcement**
    - **Validates: Requirements PROP-SAFE-002**
    - Create test/property/policy_test.go
    - Test blocked commands don't execute and return errors
    - Use gopter to test various policy configurations

- [ ] 15. Phase 4: Property-Based Testing - Agent and Plugin Properties
  - [ ]* 15.1 Write property test for agent functionality preservation
    - **Property 11: Agent Functionality Preservation**
    - **Validates: Requirements REQ-AGENT-001**
    - Create test/property/agent_test.go
    - Test agent commands produce identical behavior to previous version
    - Test agent setup, chat, and capabilities commands

  - [ ]* 15.2 Write property test for plugin operation idempotence
    - **Property 12: Plugin Operation Idempotence**
    - **Validates: Requirements REQ-PLUGIN-001**
    - Create test/property/plugin_test.go
    - Test installing plugin twice produces same state
    - Test removing non-existent plugin doesn't cause errors

  - [ ]* 15.3 Write property test for state-modifying operation idempotence
    - **Property 13: State-Modifying Operation Idempotence**
    - **Validates: Requirements REQ-REL-002**
    - Create test/property/idempotence_test.go
    - Test install, start, enable, create operations are idempotent
    - Verify executing twice produces same final state

- [ ] 16. Checkpoint - Verify all property tests pass
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 17. Phase 4: Integration Testing Setup
  - [ ] 17.1 Create Docker-based integration test infrastructure
    - Create test/integration/ directory
    - Create test/integration/test-all-distros.sh script
    - Create test/integration/test-in-container.sh script
    - Define test matrix for 7+ distributions (Ubuntu, Fedora, Arch, Alpine, openSUSE, Gentoo, Void)
    - _Requirements: REQ-DISTRO-001, REQ-DISTRO-002, REQ-DISTRO-003_

  - [ ]* 17.2 Write integration tests for all distributions
    - Test system info command on each distribution
    - Test service list command on each distribution
    - Test process tree command on each distribution
    - Test package manager detection on each distribution
    - Test init system detection on each distribution
    - Verify correct distribution identification
    - _Requirements: REQ-DISTRO-001, REQ-DISTRO-002, REQ-DISTRO-003_

  - [ ]* 17.3 Set up CI pipeline for integration tests
    - Create .github/workflows/integration.yml (if using GitHub Actions)
    - Configure Docker-based testing in CI
    - Run integration tests on all distributions
    - Set up test result reporting
    - _Requirements: REQ-DISTRO-001_

- [ ] 18. Phase 4: Performance Benchmarking
  - [ ]* 18.1 Create performance benchmarks
    - Create test/benchmark/startup_test.go
    - Implement BenchmarkStartup for cold/warm start
    - Implement BenchmarkSystemInfo
    - Implement BenchmarkServiceList
    - Implement BenchmarkProcessTree
    - Implement BenchmarkNetworkInterfaces
    - Implement BenchmarkDiskUsage
    - _Requirements: REQ-PERF-001_

  - [ ]* 18.2 Run benchmarks and verify performance targets
    - Run benchmarks with `go test -bench=.`
    - Verify startup < 100ms
    - Verify system info < 200ms
    - Verify service list < 1000ms
    - Verify process tree < 500ms
    - Document results
    - _Requirements: REQ-PERF-001_

- [ ] 19. Checkpoint - Verify all tests and benchmarks pass
  - Ensure all tests pass, ask the user if questions arise.

- [~] 20. Phase 5: Documentation Updates
  - [~] 20.1 Update README.md
    - Remove all TUI references
    - Add Gentoo and Void Linux to supported distributions
    - Update architecture diagram
    - Add performance benchmark results
    - Update installation instructions
    - _Requirements: REQ-DOC-001_

  - [~] 20.2 Update COMMAND.md
    - Verify all commands are documented
    - Add examples for Gentoo and Void Linux
    - Document JSON output format for all commands
    - Document YAML output format for all commands
    - Add troubleshooting section
    - _Requirements: REQ-DOC-001_

  - [~] 20.3 Create migration guide
    - Create MIGRATION.md document
    - Document TUI removal (use CLI commands instead)
    - Document new distributions supported
    - Document performance improvements
    - Emphasize no breaking changes to CLI
    - Provide upgrade instructions
    - _Requirements: REQ-DOC-001_

  - [~] 20.4 Update API documentation
    - Generate godoc for all packages
    - Document all public interfaces
    - Add code examples for key functions
    - Document error handling patterns
    - _Requirements: REQ-DOC-001_

- [~] 21. Phase 5: Release Preparation
  - [~] 21.1 Create release artifacts
    - Build binary for Linux x86_64
    - Build binary for Linux aarch64
    - Create Debian package (.deb)
    - Create RPM package (.rpm)
    - Create AUR PKGBUILD
    - Create Alpine apk
    - Create source tarball
    - _Requirements: REQ-REL-003_

  - [~] 21.2 Create release notes
    - Create CHANGELOG.md entry for v2.0.0
    - Document all new features (Gentoo, Void Linux support)
    - Document TUI removal
    - Document performance improvements
    - Document all bug fixes
    - List all breaking changes (none expected)
    - _Requirements: REQ-REL-003_

  - [~] 21.3 Prepare distribution channels
    - Upload binaries to GitHub Releases
    - Submit to Debian PPA
    - Submit to Fedora COPR
    - Submit to AUR
    - Submit to Alpine community repository
    - Create Docker images for Docker Hub
    - _Requirements: REQ-REL-003_

  - [~] 21.4 Final validation and release
    - Run full test suite one final time
    - Verify all 17 property tests pass
    - Verify integration tests pass on all distributions
    - Verify performance benchmarks meet targets
    - Tag version v2.0.0 in Git
    - Push release to GitHub
    - Announce release
    - _Requirements: REQ-REL-003_

- [~] 22. Final Checkpoint - Release complete
  - Ensure all tests pass, ask the user if questions arise.

## Notes

- Tasks marked with `*` are optional testing tasks and can be skipped for faster MVP
- Each task references specific requirements for traceability
- Property tests validate universal correctness properties using gopter
- Integration tests ensure compatibility across 7+ Linux distributions
- All implementation uses Go (project language)
- Checkpoints ensure incremental validation at phase boundaries
- Zero breaking changes to existing CLI functionality
- Performance targets: startup < 100ms, commands < 500ms
- Test coverage goal: 80%+ unit tests, 100% property tests
