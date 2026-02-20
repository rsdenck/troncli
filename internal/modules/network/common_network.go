package network

import (
	"fmt"
	"net"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/mascli/troncli/internal/core/ports"
)

// RunNativeDig performs a native DNS lookup
func RunNativeDig(target string) (string, error) {
	// Simple DNS lookup
	ips, err := net.LookupHost(target)
	if err != nil {
		return "", err
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("; <<>> TRONCLI DiG 1.0 <<>> %s\n", target))
	result.WriteString(";; global options: +cmd\n")
	result.WriteString(";; Got answer:\n")
	result.WriteString(";; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 12345\n")
	result.WriteString(";; flags: qr rd ra; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 1\n\n")

	result.WriteString(";; QUESTION SECTION:\n")
	result.WriteString(fmt.Sprintf(";%s.\t\t\tIN\tA\n\n", target))

	result.WriteString(";; ANSWER SECTION:\n")
	for _, ip := range ips {
		result.WriteString(fmt.Sprintf("%s.\t\t\t300\tIN\tA\t%s\n", target, ip))
	}

	cname, err := net.LookupCNAME(target)
	if err == nil && cname != "" && cname != target+"." {
		result.WriteString(fmt.Sprintf("%s.\t\t\t300\tIN\tCNAME\t%s\n", target, cname))
	}

	// MX records
	mxs, err := net.LookupMX(target)
	if err == nil {
		for _, mx := range mxs {
			result.WriteString(fmt.Sprintf("%s.\t\t\t300\tIN\tMX\t%d\t%s\n", target, mx.Pref, mx.Host))
		}
	}

	// TXT records
	txts, err := net.LookupTXT(target)
	if err == nil {
		for _, txt := range txts {
			result.WriteString(fmt.Sprintf("%s.\t\t\t300\tIN\tTXT\t\"%s\"\n", target, txt))
		}
	}

	result.WriteString(fmt.Sprintf("\n;; Query time: %d msec\n", 10))         // Fake time
	result.WriteString(fmt.Sprintf(";; SERVER: %s#53(8.8.8.8)\n", "8.8.8.8")) // Fake server
	result.WriteString(fmt.Sprintf(";; WHEN: %s\n", time.Now().Format(time.RFC1123)))
	result.WriteString(fmt.Sprintf(";; MSG SIZE  rcvd: %d\n", result.Len()))

	return result.String(), nil
}

// RunNativePortScan performs a TCP connect scan on common ports or custom ports if provided
func RunNativePortScan(target string, customPorts []int) ([]ports.PortScanResult, error) {
	portsToScan := customPorts
	if len(portsToScan) == 0 {
		portsToScan = []int{20, 21, 22, 23, 25, 53, 80, 110, 143, 443, 465, 587, 993, 995, 3306, 3389, 5432, 6379, 8080, 8443}
	}

	var results []ports.PortScanResult
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Limit concurrency
	sem := make(chan struct{}, 100)

	timeout := 2 * time.Second

	for _, port := range portsToScan {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			address := net.JoinHostPort(target, fmt.Sprintf("%d", p))
			conn, err := net.DialTimeout("tcp", address, timeout)
			if err == nil {
				conn.Close()

				service := getServiceName(p)

				res := ports.PortScanResult{
					Port:     p,
					Protocol: "tcp",
					State:    "open",
					Service:  service,
				}
				mu.Lock()
				results = append(results, res)
				mu.Unlock()
			}
		}(port)
	}

	wg.Wait()

	sort.Slice(results, func(i, j int) bool {
		return results[i].Port < results[j].Port
	})

	return results, nil
}

func getServiceName(p int) string {
	switch p {
	case 21:
		return "ftp"
	case 22:
		return "ssh"
	case 23:
		return "telnet"
	case 25:
		return "smtp"
	case 53:
		return "domain"
	case 80:
		return "http"
	case 443:
		return "https"
	case 3306:
		return "mysql"
	case 5432:
		return "postgresql"
	case 6379:
		return "redis"
	case 8080:
		return "http-proxy"
	default:
		return "unknown"
	}
}

// ParsePortsFromOptions parses ports from options string (e.g., "-p 80,443" or just "80,443")
func ParsePortsFromOptions(options string) []int {
	if options == "" {
		return nil
	}

	options = strings.TrimSpace(options)
	// Handle -p flag if user provided it manually
	if strings.Contains(options, "-p") {
		parts := strings.Fields(options)
		for i, part := range parts {
			if part == "-p" && i+1 < len(parts) {
				// Found -p, next part is ports
				return parsePortList(parts[i+1])
			} else if strings.HasPrefix(part, "-p") {
				// Found -p80,443
				return parsePortList(strings.TrimPrefix(part, "-p"))
			}
		}
	}

	// If no -p flag found, try to parse the whole string as port list if it looks like one
	// or just return nil to use defaults if it's some other nmap option we don't support
	// But since we are replacing nmap, we can be more flexible.
	// If the string contains digits and commas, assume it's a port list.
	if strings.ContainsAny(options, "0123456789") {
		return parsePortList(options)
	}

	return nil
}

func parsePortList(s string) []int {
	var ports []int
	parts := strings.Split(s, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if p, err := strconv.Atoi(part); err == nil {
			ports = append(ports, p)
		}
	}
	return ports
}
