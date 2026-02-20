package security

import (
	"fmt"
	"strings"

	"github.com/mascli/troncli/internal/core/ports"
)

// VulnDB simulates a CVE database lookup
// In a full implementation, this would query a local SQLite database (like cve-bin-tool does)
// or an API (OSV). Here we embed a significant set of known vulnerabilities for demonstration.
type VulnDB struct {
	// Map of Product -> Version -> []CVE
	// Simplified for this implementation
	Entries map[string][]VulnerabilityEntry
}

type VulnerabilityEntry struct {
	Product     string
	Version     string // Exact version or range (simplified to exact for now)
	CVEID       string
	Severity    string
	Description string
}

func NewVulnDB() *VulnDB {
	db := &VulnDB{
		Entries: make(map[string][]VulnerabilityEntry),
	}
	db.loadData()
	return db
}

func (db *VulnDB) Check(product, version string) []ports.CVEVulnerability {
	var results []ports.CVEVulnerability
	
	// Normalize
	product = strings.ToLower(product)
	version = strings.TrimSpace(version)

	if entries, ok := db.Entries[product]; ok {
		for _, entry := range entries {
			// Basic version matching
			// "Deep Audit" would parse semantic versions properly
			if entry.Version == version || entry.Version == "*" {
				results = append(results, ports.CVEVulnerability{
					CVEID:       entry.CVEID,
					Severity:    entry.Severity,
					Product:     entry.Product,
					Version:     version,
					Description: entry.Description,
				})
			} else if strings.HasPrefix(entry.Version, "<") {
				// Simple range check logic
				limit := strings.TrimSpace(strings.TrimPrefix(entry.Version, "<"))
				if compareVersions(version, limit) < 0 {
					results = append(results, ports.CVEVulnerability{
						CVEID:       entry.CVEID,
						Severity:    entry.Severity,
						Product:     entry.Product,
						Version:     version,
						Description: entry.Description,
					})
				}
			}
		}
	}
	return results
}

// Simple semantic version comparison: -1 if v1 < v2, 0 if v1 == v2, 1 if v1 > v2
func compareVersions(v1, v2 string) int {
	// Split by dots
	p1 := strings.Split(v1, ".")
	p2 := strings.Split(v2, ".")
	
	max := len(p1)
	if len(p2) > max {
		max = len(p2)
	}

	for i := 0; i < max; i++ {
		n1 := 0
		if i < len(p1) {
			fmt.Sscanf(p1[i], "%d", &n1)
		}
		n2 := 0
		if i < len(p2) {
			fmt.Sscanf(p2[i], "%d", &n2)
		}

		if n1 < n2 {
			return -1
		}
		if n1 > n2 {
			return 1
		}
	}
	return 0
}

func (db *VulnDB) loadData() {
	// Populate with some real-world examples
	
	// OpenSSL
	db.addEntry("openssl", "< 1.1.1t", "CVE-2023-0286", "HIGH", "Type confusion in X.400 address processing")
	db.addEntry("openssl", "< 3.0.8", "CVE-2023-0215", "HIGH", "Use-after-free in BIO_new_NDEF")
	db.addEntry("openssl", "1.0.2g", "CVE-2016-0800", "CRITICAL", "DROWN attack")

	// Python
	db.addEntry("python", "< 3.9.16", "CVE-2022-45061", "MEDIUM", "CPU denial of service via IDNA")
	db.addEntry("python", "< 3.10.9", "CVE-2022-45061", "MEDIUM", "CPU denial of service via IDNA")
	
	// SQLite
	db.addEntry("sqlite", "< 3.40.1", "CVE-2022-46908", "HIGH", "Specific SQL queries can cause a denial of service")

	// Curl
	db.addEntry("curl", "< 7.88.0", "CVE-2023-23916", "MEDIUM", "HTTP multi-header compression denial of service")
	
	// Tcpdump
	db.addEntry("tcpdump", "< 4.99.3", "CVE-2023-1801", "HIGH", "Heap-based buffer overflow in SMB printer")
}

func (db *VulnDB) addEntry(product, version, cve, severity, desc string) {
	entry := VulnerabilityEntry{
		Product:     product,
		Version:     version,
		CVEID:       cve,
		Severity:    severity,
		Description: desc,
	}
	db.Entries[product] = append(db.Entries[product], entry)
}
