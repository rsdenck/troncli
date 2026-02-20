package security

// Package security provides security scanning capabilities.

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/mascli/troncli/internal/core/ports"
)

// CveToolManager implements ports.SecurityManager wrapping cve-bin-tool
type CveToolManager struct{}

// NewCveToolManager creates a new security manager
func NewCveToolManager() *CveToolManager {
	return &CveToolManager{}
}

// ScanDirectory runs cve-bin-tool on a directory
func (m *CveToolManager) ScanDirectory(path string) ([]ports.CVEVulnerability, error) {
	return m.runScan(path)
}

// ScanBinary runs cve-bin-tool on a specific binary
func (m *CveToolManager) ScanBinary(path string) ([]ports.CVEVulnerability, error) {
	return m.runScan(path)
}

func (m *CveToolManager) runScan(input string) ([]ports.CVEVulnerability, error) {
	// cve-bin-tool --format json <input>
	cmd := exec.Command("cve-bin-tool", "--format", "json", input)
	out, err := cmd.Output()
	if err != nil {
		// cve-bin-tool might return exit code 1 if vulns found? 
		// Need to check documentation or behavior. 
		// Usually tools return non-zero on error OR vulns found.
		// If output is valid JSON, we can proceed.
		if len(out) == 0 {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
	}

	var results []struct {
		CVE         string `json:"cve_number"`
		Severity    string `json:"severity"`
		Product     string `json:"product"`
		Version     string `json:"version"`
		Description string `json:"description"` // might not be in standard output, checking schema
		Path        string `json:"path"`
	}

	// cve-bin-tool JSON output is usually a list of objects
	if err := json.Unmarshal(out, &results); err != nil {
		// Try parsing as object if it returns summary
		return nil, fmt.Errorf("failed to parse scan results: %w", err)
	}

	var vulns []ports.CVEVulnerability
	for _, r := range results {
		vulns = append(vulns, ports.CVEVulnerability{
			CVEID:    r.CVE,
			Severity: r.Severity,
			Product:  r.Product,
			Version:  r.Version,
			Path:     r.Path,
		})
	}

	return vulns, nil
}

// IsToolInstalled checks if cve-bin-tool is available
func (m *CveToolManager) IsToolInstalled() bool {
	_, err := exec.LookPath("cve-bin-tool")
	return err == nil
}

// InstallTool attempts to install cve-bin-tool via pip
func (m *CveToolManager) InstallTool() error {
	// pip install cve-bin-tool
	cmd := exec.Command("pip", "install", "cve-bin-tool")
	return cmd.Run()
}
