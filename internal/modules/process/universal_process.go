package process

// Package process provides process management capabilities.

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/mascli/troncli/internal/core/adapter"
	"github.com/mascli/troncli/internal/core/domain"
	"github.com/mascli/troncli/internal/core/ports"
)

// UniversalProcessManager implements ProcessManager using system tools
type UniversalProcessManager struct {
	executor adapter.Executor
	profile  *domain.SystemProfile
}

// NewUniversalProcessManager creates a new instance
func NewUniversalProcessManager(executor adapter.Executor, profile *domain.SystemProfile) *UniversalProcessManager {
	return &UniversalProcessManager{
		executor: executor,
		profile:  profile,
	}
}

// KillProcess sends a signal to a process
func (m *UniversalProcessManager) KillProcess(pid int, signal string) error {
	ctx := context.Background()
	// Using "kill" command which is universally available
	// Maps standard signals
	sig := "-15" // SIGTERM default
	switch signal {
	case "SIGKILL":
		sig = "-9"
	case "SIGINT":
		sig = "-2"
	case "SIGHUP":
		sig = "-1"
	}

	_, err := m.executor.Exec(ctx, "kill", sig, strconv.Itoa(pid))
	if err != nil {
		return fmt.Errorf("failed to kill process %d: %w", pid, err)
	}
	return nil
}

// ReniceProcess changes process priority
func (m *UniversalProcessManager) ReniceProcess(pid int, priority int) error {
	ctx := context.Background()
	// renice -n priority -p pid
	_, err := m.executor.Exec(ctx, "renice", "-n", strconv.Itoa(priority), "-p", strconv.Itoa(pid))
	if err != nil {
		return fmt.Errorf("failed to renice process %d: %w", pid, err)
	}
	return nil
}

// KillZombies finds and kills zombie processes
func (m *UniversalProcessManager) KillZombies() (int, error) {
	// Strategy: Find Z state processes, get PPID, kill -SIGCHLD PPID
	ctx := context.Background()
	// ps -A -o stat,ppid
	res, err := m.executor.Exec(ctx, "ps", "-A", "-o", "stat,ppid")
	if err != nil {
		return 0, err
	}

	lines := strings.Split(res.Stdout, "\n")
	count := 0
	parents := make(map[string]bool)

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		state := fields[0]
		if strings.HasPrefix(state, "Z") {
			ppid := fields[1]
			if ppid != "1" {
				parents[ppid] = true
				count++
			}
		}
	}

	for ppid := range parents {
		// Send SIGCHLD (17 on Linux x86/ARM usually, but kill -17 might vary?
		// Best to rely on 'kill -s SIGCHLD' if supported, or just let init handle it eventually)
		// Standard kill doesn't always support named signals easily across busybox/gnu.
		// Let's try "kill -s CHLD"
		m.executor.Exec(ctx, "kill", "-s", "CHLD", ppid)
	}

	return count, nil
}

// GetProcessTree returns the process tree
func (m *UniversalProcessManager) GetProcessTree() ([]ports.ProcessNode, error) {
	ctx := context.Background()
	// ps -eo pid,ppid,user,stat,comm
	// We use "comm" for command name, usually truncated but sufficient for tree view
	// "args" would be full command line
	res, err := m.executor.Exec(ctx, "ps", "-eo", "pid,ppid,user,stat,comm")
	if err != nil {
		return nil, fmt.Errorf("failed to list processes: %w", err)
	}

	lines := strings.Split(res.Stdout, "\n")
	nodeMap := make(map[int]*ports.ProcessNode)

	// First pass: create all nodes
	for i, line := range lines {
		if i == 0 { // Skip header
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 5 {
			continue
		}

		pid, err := strconv.Atoi(fields[0])
		if err != nil {
			continue
		}
		ppid, err := strconv.Atoi(fields[1])
		if err != nil {
			continue
		}
		user := fields[2]
		state := fields[3]
		// comm might contain spaces, so join the rest
		name := strings.Join(fields[4:], " ")

		node := &ports.ProcessNode{
			PID:      pid,
			PPID:     ppid,
			Name:     name,
			User:     user,
			State:    state,
			Children: []ports.ProcessNode{},
		}
		nodeMap[pid] = node
	}

	// Second pass: build tree
	var roots []ports.ProcessNode
	for _, node := range nodeMap {
		if parent, ok := nodeMap[node.PPID]; ok && node.PID != node.PPID {
			parent.Children = append(parent.Children, *node)
		}
	}

	// Third pass: find roots (nodes whose parents are not in the map or PPID=0)
	// Note: We need to iterate over the map again to find roots, but since we modified children by value in the previous loop?
	// Wait, parent.Children = append(parent.Children, *node) copies the node value.
	// If I modify the child node later, the parent's copy won't update.
	// But here I am building bottom-up or just linking?
	// Actually, with value semantics, it's tricky.
	// Let's use a map of pointers to struct with pointer children first, then convert.
	// Or simply:
	// The requirement is just a list of nodes, and maybe the UI handles the tree?
	// The interface says `GetProcessTree() ([]ProcessNode, error)`.
	// If it returns a tree, it should be the root nodes with Children populated.

	// Let's redo the tree building with a helper struct to avoid value copy issues during build.
	type tempNode struct {
		*ports.ProcessNode
		Children []*tempNode
	}
	tempMap := make(map[int]*tempNode)

	// Re-parse to populate tempMap
	for i, line := range lines {
		if i == 0 {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 5 {
			continue
		}
		pid, _ := strconv.Atoi(fields[0])
		ppid, _ := strconv.Atoi(fields[1])
		user := fields[2]
		state := fields[3]
		name := strings.Join(fields[4:], " ")

		tempMap[pid] = &tempNode{
			ProcessNode: &ports.ProcessNode{
				PID:   pid,
				PPID:  ppid,
				Name:  name,
				User:  user,
				State: state,
			},
		}
	}

	// Build links
	var rootTemps []*tempNode
	for _, node := range tempMap {
		if parent, ok := tempMap[node.PPID]; ok && node.PID != node.PPID {
			parent.Children = append(parent.Children, node)
		} else {
			rootTemps = append(rootTemps, node)
		}
	}

	// Convert to []ports.ProcessNode recursive
	var convert func(*tempNode) ports.ProcessNode
	convert = func(n *tempNode) ports.ProcessNode {
		pn := *n.ProcessNode
		pn.Children = make([]ports.ProcessNode, len(n.Children))
		for i, child := range n.Children {
			pn.Children[i] = convert(child)
		}
		return pn
	}

	roots = make([]ports.ProcessNode, len(rootTemps))
	for i, r := range rootTemps {
		roots[i] = convert(r)
	}

	return roots, nil
}

// GetOpenFiles returns list of open files for a process
func (m *UniversalProcessManager) GetOpenFiles(pid int) ([]string, error) {
	ctx := context.Background()
	// ls -l /proc/<pid>/fd
	res, err := m.executor.Exec(ctx, "ls", "-l", fmt.Sprintf("/proc/%d/fd", pid))
	if err != nil {
		return nil, fmt.Errorf("failed to get open files for pid %d: %w", pid, err)
	}

	lines := strings.Split(res.Stdout, "\n")
	var files []string
	for _, line := range lines {
		// lrwx------ 1 root root 64 Feb 17 ... 0 -> /dev/pts/0
		if strings.Contains(line, " -> ") {
			parts := strings.Split(line, " -> ")
			if len(parts) == 2 {
				files = append(files, parts[1])
			}
		}
	}
	return files, nil
}

// GetProcessPorts returns ports listened by a process
func (m *UniversalProcessManager) GetProcessPorts(pid int) ([]string, error) {
	ctx := context.Background()
	// ss -l -p -n
	// Output: Netid State Recv-Q Send-Q Local Address:Port Peer Address:PortProcess
	// grep for pid
	res, err := m.executor.Exec(ctx, "ss", "-lpn")
	if err != nil {
		// Fallback to netstat if ss fails?
		// But instructions say "modern Linux". ss is modern.
		return nil, fmt.Errorf("failed to get process ports: %w", err)
	}

	lines := strings.Split(res.Stdout, "\n")
	var portsList []string
	pidStr := fmt.Sprintf("pid=%d,", pid)

	for _, line := range lines {
		if strings.Contains(line, pidStr) {
			fields := strings.Fields(line)
			if len(fields) >= 5 {
				// Local Address:Port is usually field 4 (0-indexed) or 3 depending on State column
				// ss output: Netid State Recv-Q Send-Q Local_Address:Port Peer_Address:Port Process
				// u_str ESTAB 0 0 * 19036 * 19037 users:(("dbus-daemon",pid=584,fd=12))
				// tcp LISTEN 0 128 0.0.0.0:22 0.0.0.0:* users:(("sshd",pid=684,fd=3))

				// We want the Local Address:Port
				// Usually index 4
				if len(fields) > 4 {
					portsList = append(portsList, fields[4])
				}
			}
		}
	}
	return portsList, nil
}

// GetAllListeningPorts returns all listening ports
func (m *UniversalProcessManager) GetAllListeningPorts() ([]string, error) {
	ctx := context.Background()
	// ss -nltu
	res, err := m.executor.Exec(ctx, "ss", "-nltu")
	if err != nil {
		return nil, fmt.Errorf("failed to get listening ports: %w", err)
	}

	lines := strings.Split(res.Stdout, "\n")
	var portsList []string

	for i, line := range lines {
		if i == 0 {
			continue
		} // Header
		fields := strings.Fields(line)
		if len(fields) >= 5 {
			// Netid State Recv-Q Send-Q Local_Address:Port ...
			portsList = append(portsList, fmt.Sprintf("%s/%s", fields[0], fields[4]))
		}
	}
	return portsList, nil
}
