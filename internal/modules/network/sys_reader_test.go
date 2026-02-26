package network

import (
	"os"
	"testing"
)

func TestSysReader_ReadInterfaces(t *testing.T) {
	// Skip test if /sys/class/net doesn't exist (non-Linux systems)
	if _, err := os.Stat("/sys/class/net"); os.IsNotExist(err) {
		t.Skip("Skipping test: /sys/class/net not available (non-Linux system)")
	}

	reader := NewSysReader()
	interfaces, err := reader.ReadInterfaces()
	
	if err != nil {
		t.Fatalf("ReadInterfaces() failed: %v", err)
	}

	if len(interfaces) == 0 {
		t.Error("Expected at least one network interface, got none")
	}

	// Verify that at least one interface has expected fields populated
	foundValidInterface := false
	for _, iface := range interfaces {
		if iface.Name != "" && iface.MTU > 0 {
			foundValidInterface = true
			t.Logf("Found interface: %s (MTU: %d, State: %s, MAC: %s)", 
				iface.Name, iface.MTU, iface.State, iface.HardwareAddr)
			break
		}
	}

	if !foundValidInterface {
		t.Error("No valid interface found with name and MTU")
	}
}

func TestSysReader_readInterface(t *testing.T) {
	// Skip test if /sys/class/net doesn't exist (non-Linux systems)
	if _, err := os.Stat("/sys/class/net"); os.IsNotExist(err) {
		t.Skip("Skipping test: /sys/class/net not available (non-Linux system)")
	}

	reader := NewSysReader()
	
	// Try to read the loopback interface which should exist on all Linux systems
	iface, err := reader.readInterface("lo")
	if err != nil {
		t.Fatalf("readInterface('lo') failed: %v", err)
	}

	if iface.Name != "lo" {
		t.Errorf("Expected interface name 'lo', got '%s'", iface.Name)
	}

	if iface.MTU == 0 {
		t.Error("Expected MTU to be set, got 0")
	}

	if iface.State == "" {
		t.Error("Expected State to be set, got empty string")
	}

	t.Logf("Loopback interface: MTU=%d, State=%s, Index=%d", 
		iface.MTU, iface.State, iface.Index)
}

func TestSysReader_readInterface_NonExistent(t *testing.T) {
	// Skip test if /sys/class/net doesn't exist (non-Linux systems)
	if _, err := os.Stat("/sys/class/net"); os.IsNotExist(err) {
		t.Skip("Skipping test: /sys/class/net not available (non-Linux system)")
	}

	reader := NewSysReader()
	
	// Try to read a non-existent interface
	_, err := reader.readInterface("nonexistent999")
	if err == nil {
		t.Error("Expected error for non-existent interface, got nil")
	}
}

func TestSysReader_ReadSocketStats(t *testing.T) {
	// Skip test if /proc/net/tcp doesn't exist (non-Linux systems)
	if _, err := os.Stat("/proc/net/tcp"); os.IsNotExist(err) {
		t.Skip("Skipping test: /proc/net/tcp not available (non-Linux system)")
	}

	reader := NewSysReader()
	stats, err := reader.ReadSocketStats()

	if err != nil {
		t.Fatalf("ReadSocketStats() failed: %v", err)
	}

	// We should have at least some sockets (loopback, etc.)
	totalSockets := len(stats.TCP) + len(stats.TCP6) + len(stats.UDP) + len(stats.UDP6)
	if totalSockets == 0 {
		t.Error("Expected at least one socket, got none")
	}

	t.Logf("Found sockets: TCP=%d, TCP6=%d, UDP=%d, UDP6=%d",
		len(stats.TCP), len(stats.TCP6), len(stats.UDP), len(stats.UDP6))

	// Verify at least one socket has valid data
	if len(stats.TCP) > 0 {
		socket := stats.TCP[0]
		if socket.LocalAddr == "" {
			t.Error("Expected LocalAddr to be set")
		}
		if socket.State == "" {
			t.Error("Expected State to be set")
		}
		t.Logf("Sample TCP socket: %s:%d -> %s:%d [%s]",
			socket.LocalAddr, socket.LocalPort,
			socket.RemoteAddr, socket.RemotePort,
			socket.State)
	}
}

func TestSysReader_ReadRouteTable(t *testing.T) {
	// Skip test if /proc/net/route doesn't exist (non-Linux systems)
	if _, err := os.Stat("/proc/net/route"); os.IsNotExist(err) {
		t.Skip("Skipping test: /proc/net/route not available (non-Linux system)")
	}

	reader := NewSysReader()
	routes, err := reader.ReadRouteTable()

	if err != nil {
		t.Fatalf("ReadRouteTable() failed: %v", err)
	}

	// We should have at least one route (default route or loopback)
	if len(routes) == 0 {
		t.Error("Expected at least one route, got none")
	}

	t.Logf("Found %d routes", len(routes))

	// Verify at least one route has valid data
	if len(routes) > 0 {
		route := routes[0]
		if route.Interface == "" {
			t.Error("Expected Interface to be set")
		}
		t.Logf("Sample route: Interface=%s, Destination=%s, Gateway=%s, Mask=%s, Metric=%d",
			route.Interface, route.Destination, route.Gateway, route.Mask, route.Metric)
	}
}

func TestParseTCPSockets(t *testing.T) {
	// Sample data from /proc/net/tcp
	sampleData := []byte(`  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode
   0: 0100007F:1F90 00000000:0000 0A 00000000:00000000 00:00000000 00000000  1000        0 12345 1 0000000000000000 100 0 0 10 0
   1: 00000000:0050 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 67890 1 0000000000000000 100 0 0 10 0`)

	sockets := parseTCPSockets(sampleData)

	if len(sockets) != 2 {
		t.Errorf("Expected 2 sockets, got %d", len(sockets))
	}

	// Check first socket (127.0.0.1:8080)
	if len(sockets) > 0 {
		socket := sockets[0]
		if socket.LocalAddr != "127.0.0.1" {
			t.Errorf("Expected LocalAddr '127.0.0.1', got '%s'", socket.LocalAddr)
		}
		if socket.LocalPort != 8080 {
			t.Errorf("Expected LocalPort 8080, got %d", socket.LocalPort)
		}
		if socket.State != "LISTEN" {
			t.Errorf("Expected State 'LISTEN', got '%s'", socket.State)
		}
		if socket.UID != 1000 {
			t.Errorf("Expected UID 1000, got %d", socket.UID)
		}
	}

	// Check second socket (0.0.0.0:80)
	if len(sockets) > 1 {
		socket := sockets[1]
		if socket.LocalAddr != "0.0.0.0" {
			t.Errorf("Expected LocalAddr '0.0.0.0', got '%s'", socket.LocalAddr)
		}
		if socket.LocalPort != 80 {
			t.Errorf("Expected LocalPort 80, got %d", socket.LocalPort)
		}
	}
}

func TestParseUDPSockets(t *testing.T) {
	// Sample data from /proc/net/udp
	sampleData := []byte(`  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode ref pointer drops
   0: 00000000:0035 00000000:0000 07 00000000:00000000 00:00000000 00000000     0        0 11111 2 0000000000000000 0`)

	sockets := parseUDPSockets(sampleData)

	if len(sockets) != 1 {
		t.Errorf("Expected 1 socket, got %d", len(sockets))
	}

	if len(sockets) > 0 {
		socket := sockets[0]
		if socket.LocalAddr != "0.0.0.0" {
			t.Errorf("Expected LocalAddr '0.0.0.0', got '%s'", socket.LocalAddr)
		}
		if socket.LocalPort != 53 {
			t.Errorf("Expected LocalPort 53, got %d", socket.LocalPort)
		}
		if socket.State != "CLOSE" {
			t.Errorf("Expected State 'CLOSE', got '%s'", socket.State)
		}
	}
}

func TestParseRouteTable(t *testing.T) {
	// Sample data from /proc/net/route
	sampleData := []byte(`Iface	Destination	Gateway 	Flags	RefCnt	Use	Metric	Mask		MTU	Window	IRTT
eth0	00000000	0101A8C0	0003	0	0	100	00000000	0	0	0
eth0	0001A8C0	00000000	0001	0	0	0	00FFFFFF	0	0	0`)

	routes := parseRouteTable(sampleData)

	if len(routes) != 2 {
		t.Errorf("Expected 2 routes, got %d", len(routes))
	}

	// Check first route (default gateway)
	if len(routes) > 0 {
		route := routes[0]
		if route.Interface != "eth0" {
			t.Errorf("Expected Interface 'eth0', got '%s'", route.Interface)
		}
		if route.Destination != "0.0.0.0" {
			t.Errorf("Expected Destination '0.0.0.0', got '%s'", route.Destination)
		}
		if route.Gateway != "192.168.1.1" {
			t.Errorf("Expected Gateway '192.168.1.1', got '%s'", route.Gateway)
		}
		if route.Metric != 100 {
			t.Errorf("Expected Metric 100, got %d", route.Metric)
		}
	}

	// Check second route (local network)
	if len(routes) > 1 {
		route := routes[1]
		if route.Destination != "192.168.1.0" {
			t.Errorf("Expected Destination '192.168.1.0', got '%s'", route.Destination)
		}
		if route.Mask != "255.255.255.0" {
			t.Errorf("Expected Mask '255.255.255.0', got '%s'", route.Mask)
		}
	}
}

func TestParseHexIP(t *testing.T) {
	tests := []struct {
		name     string
		hexIP    string
		expected string
	}{
		{"Loopback", "0100007F", "127.0.0.1"},
		{"Zero", "00000000", "0.0.0.0"},
		{"Gateway", "0101A8C0", "192.168.1.1"},
		{"Network", "0001A8C0", "192.168.1.0"},
		{"Broadcast", "FFFFFFFF", "255.255.255.255"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseHexIP(tt.hexIP)
			if result != tt.expected {
				t.Errorf("parseHexIP(%s) = %s, expected %s", tt.hexIP, result, tt.expected)
			}
		})
	}
}

func TestParseTCPState(t *testing.T) {
	tests := []struct {
		hexState string
		expected string
	}{
		{"01", "ESTABLISHED"},
		{"02", "SYN_SENT"},
		{"03", "SYN_RECV"},
		{"0A", "LISTEN"},
		{"06", "TIME_WAIT"},
		{"FF", "UNKNOWN"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := parseTCPState(tt.hexState)
			if result != tt.expected {
				t.Errorf("parseTCPState(%s) = %s, expected %s", tt.hexState, result, tt.expected)
			}
		})
	}
}

func TestParseAddress(t *testing.T) {
	tests := []struct {
		name         string
		addrPort     string
		expectedAddr string
		expectedPort int
	}{
		{"Loopback 8080", "0100007F:1F90", "127.0.0.1", 8080},
		{"Any address port 80", "00000000:0050", "0.0.0.0", 80},
		{"Any address port 443", "00000000:01BB", "0.0.0.0", 443},
		{"DNS port 53", "00000000:0035", "0.0.0.0", 53},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addr, port := parseAddress(tt.addrPort)
			if addr != tt.expectedAddr {
				t.Errorf("parseAddress(%s) addr = %s, expected %s", tt.addrPort, addr, tt.expectedAddr)
			}
			if port != tt.expectedPort {
				t.Errorf("parseAddress(%s) port = %d, expected %d", tt.addrPort, port, tt.expectedPort)
			}
		})
	}
}
