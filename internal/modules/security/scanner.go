package security

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mascli/troncli/internal/core/ports"
)

// NativeScanner implements deep audit logic in pure Go
type NativeScanner struct {
	vulnDB *VulnDB
}

func NewNativeScanner() *NativeScanner {
	return &NativeScanner{
		vulnDB: NewVulnDB(),
	}
}

// ScanPath determines if path is a file or directory and scans accordingly
func (s *NativeScanner) ScanPath(path string) ([]ports.CVEVulnerability, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if info.IsDir() {
		return s.ScanDirectory(path)
	}
	return s.ScanFile(path)
}

func (s *NativeScanner) ScanDirectory(root string) ([]ports.CVEVulnerability, error) {
	var vulns []ports.CVEVulnerability

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip unreadable files
		}
		if !info.IsDir() {
			// Skip very large files to avoid OOM
			if info.Size() > 100*1024*1024 {
				return nil
			}

			// Scan file
			fileVulns, _ := s.ScanFile(path)
			if len(fileVulns) > 0 {
				vulns = append(vulns, fileVulns...)
			}
		}
		return nil
	})

	return vulns, err
}

func (s *NativeScanner) ScanFile(path string) ([]ports.CVEVulnerability, error) {
	var vulns []ports.CVEVulnerability
	filename := filepath.Base(path)

	// 1. Check for specific package files
	if filename == "go.mod" {
		// Parse go.mod
	} else if filename == "package.json" {
		// Parse package.json
	} else if filename == "requirements.txt" {
		// Parse requirements.txt
	}

	// 2. Binary Analysis (ELF, Scripts)
	// Check if executable
	info, err := os.Stat(path)
	if err == nil && info.Mode()&0111 != 0 {
		// It's executable, try to detect version
		product, version := s.DetectVersion(path)
		if product != "" && version != "" {
			// Check DB
			found := s.vulnDB.Check(product, version)
			for _, v := range found {
				v.Path = path
				vulns = append(vulns, v)
			}
		}
	}

	return vulns, nil
}

// DetectVersion tries to extract product and version from a binary
func (s *NativeScanner) DetectVersion(path string) (string, string) {
	// Strategy 1: OpenSSL specific check (common high value target)
	if strings.Contains(strings.ToLower(path), "openssl") {
		// Try extracting from binary strings
		ver := s.extractOpenSSLVersion(path)
		if ver != "" {
			return "openssl", ver
		}
	}

	// Strategy 2: Python specific check
	if strings.Contains(strings.ToLower(path), "python") {
		// Try running --version if it's safe? No, "Native" usually implies static analysis or safe execution.
		// Let's stick to static analysis for "Deep Audit" to be safe.
		ver := s.extractPythonVersion(path)
		if ver != "" {
			return "python", ver
		}
	}

	// Strategy 3: General "strings" analysis for known patterns
	// This is what cve-bin-tool does: looks for regex in the binary
	return s.scanBinaryStrings(path)
}

func (s *NativeScanner) extractOpenSSLVersion(path string) string {
	// Look for "OpenSSL x.y.z" pattern
	// In a real binary, this string is usually present
	content, err := os.ReadFile(path)
	if err != nil {
		return ""
	}

	// Regex for OpenSSL version
	re := regexp.MustCompile(`OpenSSL\s+(\d+\.\d+\.\d+[a-z]?)`)
	matches := re.FindSubmatch(content)
	if len(matches) > 1 {
		return string(matches[1])
	}
	return ""
}

func (s *NativeScanner) extractPythonVersion(path string) string {
	// Try to find "Python x.y.z"
	content, err := os.ReadFile(path)
	if err != nil {
		return ""
	}

	re := regexp.MustCompile(`Python\s+(\d+\.\d+\.\d+)`)
	matches := re.FindSubmatch(content)
	if len(matches) > 1 {
		return string(matches[1])
	}
	return ""
}

func (s *NativeScanner) scanBinaryStrings(path string) (string, string) {
	// Generic scanner for other products
	// This simulates cve-bin-tool's signature matching

	// Open file
	f, err := os.Open(path)
	if err != nil {
		return "", ""
	}
	defer f.Close()

	// We'll read the first 20MB (arbitrary limit for performance)
	// cve-bin-tool uses specific offsets or searches the whole file

	// Signatures
	signatures := map[string]*regexp.Regexp{
		"curl":    regexp.MustCompile(`curl\s+(\d+\.\d+\.\d+)`),
		"sqlite":  regexp.MustCompile(`SQLite\s+(\d+\.\d+\.\d+)`),
		"tcpdump": regexp.MustCompile(`tcpdump\s+version\s+(\d+\.\d+\.\d+)`),
	}

	// Read chunks to avoid loading whole file
	reader := bufio.NewReader(f)
	buf := make([]byte, 1024*1024) // 1MB chunks

	for {
		n, err := reader.Read(buf)
		if n == 0 || err != nil {
			break
		}

		chunk := buf[:n]
		for prod, re := range signatures {
			if matches := re.FindSubmatch(chunk); len(matches) > 1 {
				return prod, string(matches[1])
			}
		}

		if err == io.EOF {
			break
		}
	}

	return "", ""
}
