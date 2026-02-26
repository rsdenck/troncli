package network

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mascli/troncli/internal/core/ports"
)

// SysReader reads network information directly from /sys/class/net
type SysReader struct{}

// NewSysReader creates a new SysReader instance
func NewSysReader() *SysReader {
	return &SysReader{}
}

// ReadInterfaces reads network interfaces from /sys/class/net
func (r *SysReader) ReadInterfaces() ([]ports.NetworkInterface, error) {
	entries, err := os.ReadDir("/sys/class/net")
	if err != nil {
		return nil, fmt.Errorf("failed to read /sys/class/net: %w", err)
	}

	var interfaces []ports.NetworkInterface
	for _, entry := range entries {
		iface, err := r.readInterface(entry.Name())
		if err != nil {
			// Log error but continue with other interfaces
			continue
		}
		interfaces = append(interfaces, iface)
	}

	return interfaces, nil
}

// readInterface reads interface details from /sys/class/net/[name]
func (r *SysReader) readInterface(name string) (ports.NetworkInterface, error) {
	iface := ports.NetworkInterface{Name: name}
	basePath := fmt.Sprintf("/sys/class/net/%s", name)

	// Read MAC address
	mac, err := os.ReadFile(filepath.Join(basePath, "address"))
	if err == nil {
		iface.HardwareAddr = strings.TrimSpace(string(mac))
	}

	// Read MTU
	mtuData, err := os.ReadFile(filepath.Join(basePath, "mtu"))
	if err == nil {
		mtu, err := strconv.Atoi(strings.TrimSpace(string(mtuData)))
		if err == nil {
			iface.MTU = mtu
		}
	}

	// Read operational state
	stateData, err := os.ReadFile(filepath.Join(basePath, "operstate"))
	if err == nil {
		state := strings.TrimSpace(string(stateData))
		// Map operstate to UP/DOWN
		if state == "up" {
			iface.State = "UP"
		} else {
			iface.State = "DOWN"
		}
	}

	// Read interface index
	ifindexData, err := os.ReadFile(filepath.Join(basePath, "ifindex"))
	if err == nil {
		ifindex, err := strconv.Atoi(strings.TrimSpace(string(ifindexData)))
		if err == nil {
			iface.Index = ifindex
		}
	}

	// Read speed (if available) - may not exist for all interfaces
	speedData, err := os.ReadFile(filepath.Join(basePath, "speed"))
	if err == nil {
		// Speed is in Mbps, we don't have a field for it in the current struct
		// but we read it as part of the implementation
		_ = strings.TrimSpace(string(speedData))
	}

	// Read flags from the net package to maintain compatibility
	// This provides additional information like FlagUp, FlagBroadcast, etc.
	netIface, err := net.InterfaceByName(name)
	if err == nil {
		iface.Flags = netIface.Flags
		
		// Get IP addresses using net package
		addrs, err := netIface.Addrs()
		if err == nil {
			var ipAddrs []string
			for _, addr := range addrs {
				ipAddrs = append(ipAddrs, addr.String())
			}
			iface.IPAddresses = ipAddrs
		}
	}

	return iface, nil
}

// ReadSocketStats reads socket statistics from /proc/net/tcp, /proc/net/udp
func (r *SysReader) ReadSocketStats() (ports.SocketStats, error) {
	stats := ports.SocketStats{}

	// Read TCP sockets (IPv4)
	tcpData, err := os.ReadFile("/proc/net/tcp")
	if err == nil {
		stats.TCP = parseTCPSockets(tcpData)
	}

	// Read TCP6 sockets (IPv6)
	tcp6Data, err := os.ReadFile("/proc/net/tcp6")
	if err == nil {
		stats.TCP6 = parseTCPSockets(tcp6Data)
	}

	// Read UDP sockets (IPv4)
	udpData, err := os.ReadFile("/proc/net/udp")
	if err == nil {
		stats.UDP = parseUDPSockets(udpData)
	}

	// Read UDP6 sockets (IPv6)
	udp6Data, err := os.ReadFile("/proc/net/udp6")
	if err == nil {
		stats.UDP6 = parseUDPSockets(udp6Data)
	}

	return stats, nil
}

// ReadRouteTable reads routing table from /proc/net/route
func (r *SysReader) ReadRouteTable() ([]ports.Route, error) {
	data, err := os.ReadFile("/proc/net/route")
	if err != nil {
		return nil, fmt.Errorf("failed to read /proc/net/route: %w", err)
	}

	return parseRouteTable(data), nil
}

// parseTCPSockets parses TCP socket data from /proc/net/tcp or /proc/net/tcp6
func parseTCPSockets(data []byte) []ports.TCPSocket {
	var sockets []ports.TCPSocket
	lines := strings.Split(string(data), "\n")

	// Skip header line
	for i, line := range lines {
		if i == 0 || strings.TrimSpace(line) == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 10 {
			continue
		}

		// Parse local address and port
		localAddr, localPort := parseAddress(fields[1])
		// Parse remote address and port
		remoteAddr, remotePort := parseAddress(fields[2])
		// Parse state
		state := parseTCPState(fields[3])
		// Parse UID
		uid, _ := strconv.Atoi(fields[7])
		// Parse inode
		inode, _ := strconv.Atoi(fields[9])

		socket := ports.TCPSocket{
			LocalAddr:  localAddr,
			LocalPort:  localPort,
			RemoteAddr: remoteAddr,
			RemotePort: remotePort,
			State:      state,
			UID:        uid,
			Inode:      inode,
		}
		sockets = append(sockets, socket)
	}

	return sockets
}

// parseUDPSockets parses UDP socket data from /proc/net/udp or /proc/net/udp6
func parseUDPSockets(data []byte) []ports.UDPSocket {
	var sockets []ports.UDPSocket
	lines := strings.Split(string(data), "\n")

	// Skip header line
	for i, line := range lines {
		if i == 0 || strings.TrimSpace(line) == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 10 {
			continue
		}

		// Parse local address and port
		localAddr, localPort := parseAddress(fields[1])
		// Parse remote address and port
		remoteAddr, remotePort := parseAddress(fields[2])
		// Parse state (UDP states are simpler)
		state := parseUDPState(fields[3])
		// Parse UID
		uid, _ := strconv.Atoi(fields[7])
		// Parse inode
		inode, _ := strconv.Atoi(fields[9])

		socket := ports.UDPSocket{
			LocalAddr:  localAddr,
			LocalPort:  localPort,
			RemoteAddr: remoteAddr,
			RemotePort: remotePort,
			State:      state,
			UID:        uid,
			Inode:      inode,
		}
		sockets = append(sockets, socket)
	}

	return sockets
}

// parseRouteTable parses routing table data from /proc/net/route
func parseRouteTable(data []byte) []ports.Route {
	var routes []ports.Route
	lines := strings.Split(string(data), "\n")

	// Skip header line
	for i, line := range lines {
		if i == 0 || strings.TrimSpace(line) == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 8 {
			continue
		}

		// Parse fields
		iface := fields[0]
		destination := parseHexIP(fields[1])
		gateway := parseHexIP(fields[2])
		flags, _ := strconv.ParseInt(fields[3], 10, 32)
		refcnt, _ := strconv.Atoi(fields[4])
		use, _ := strconv.Atoi(fields[5])
		metric, _ := strconv.Atoi(fields[6])
		mask := parseHexIP(fields[7])

		route := ports.Route{
			Interface:   iface,
			Destination: destination,
			Gateway:     gateway,
			Flags:       int(flags),
			RefCnt:      refcnt,
			Use:         use,
			Metric:      metric,
			Mask:        mask,
		}
		routes = append(routes, route)
	}

	return routes
}

// parseAddress parses address:port from hex format (e.g., "0100007F:1F90")
func parseAddress(addrPort string) (string, int) {
	parts := strings.Split(addrPort, ":")
	if len(parts) != 2 {
		return "", 0
	}

	// Parse IP address from hex
	addr := parseHexIP(parts[0])

	// Parse port from hex
	port, err := strconv.ParseInt(parts[1], 16, 32)
	if err != nil {
		return addr, 0
	}

	return addr, int(port)
}

// parseHexIP converts hex IP address to dotted decimal notation
// Handles both IPv4 (8 hex chars) and IPv6 (32 hex chars)
func parseHexIP(hexIP string) string {
	if len(hexIP) == 8 {
		// IPv4 address (little-endian)
		ip, err := strconv.ParseUint(hexIP, 16, 32)
		if err != nil {
			return "0.0.0.0"
		}
		return fmt.Sprintf("%d.%d.%d.%d",
			byte(ip), byte(ip>>8), byte(ip>>16), byte(ip>>24))
	} else if len(hexIP) == 32 {
		// IPv6 address
		var parts []string
		for i := 0; i < 32; i += 4 {
			parts = append(parts, hexIP[i:i+4])
		}
		return strings.Join(parts, ":")
	}
	return hexIP
}

// parseTCPState converts hex state code to human-readable string
func parseTCPState(hexState string) string {
	state, err := strconv.ParseInt(hexState, 16, 32)
	if err != nil {
		return "UNKNOWN"
	}

	states := map[int]string{
		0x01: "ESTABLISHED",
		0x02: "SYN_SENT",
		0x03: "SYN_RECV",
		0x04: "FIN_WAIT1",
		0x05: "FIN_WAIT2",
		0x06: "TIME_WAIT",
		0x07: "CLOSE",
		0x08: "CLOSE_WAIT",
		0x09: "LAST_ACK",
		0x0A: "LISTEN",
		0x0B: "CLOSING",
	}

	if stateName, ok := states[int(state)]; ok {
		return stateName
	}
	return "UNKNOWN"
}

// parseUDPState converts hex state code to human-readable string
func parseUDPState(hexState string) string {
	state, err := strconv.ParseInt(hexState, 16, 32)
	if err != nil {
		return "UNKNOWN"
	}

	states := map[int]string{
		0x07: "CLOSE",
		0x0A: "LISTEN",
	}

	if stateName, ok := states[int(state)]; ok {
		return stateName
	}
	return "UNKNOWN"
}
