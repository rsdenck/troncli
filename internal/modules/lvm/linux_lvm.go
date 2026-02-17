package lvm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/mascli/troncli/internal/core/ports"
)

// LinuxLVMManager implements ports.LVMManager using Linux LVM tools
type LinuxLVMManager struct {
	// Options like sudo can be injected
	Sudo bool
}

func NewLinuxLVMManager(sudo bool) ports.LVMManager {
	return &LinuxLVMManager{Sudo: sudo}
}

func (m *LinuxLVMManager) runCommand(args ...string) (string, error) {
	var cmd *exec.Cmd
	if m.Sudo {
		cmdArgs := append([]string{"sudo"}, args...)
		cmd = exec.Command(cmdArgs[0], cmdArgs[1:]...)
	} else {
		cmd = exec.Command(args[0], args[1:]...)
	}
	
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("lvm command failed: %w, stderr: %s", err, stderr.String())
	}
	
	return strings.TrimSpace(out.String()), nil
}

// JSON structs for LVM reports
type lvmReport struct {
	Report []struct {
		Pv []map[string]string `json:"pv"`
		Vg []map[string]string `json:"vg"`
		Lv []map[string]string `json:"lv"`
	} `json:"report"`
}

func (m *LinuxLVMManager) ListPhysicalVolumes() ([]ports.PhysicalVolume, error) {
	// Use JSON output for robustness
	out, err := m.runCommand("pvs", "--reportformat", "json", "--units", "g", "-o", "pv_name,vg_name,pv_size,pv_free")
	if err != nil {
		return nil, err
	}

	var report lvmReport
	if err := json.Unmarshal([]byte(out), &report); err != nil {
		return nil, fmt.Errorf("failed to parse pvs output: %w", err)
	}

	var pvs []ports.PhysicalVolume
	for _, r := range report.Report {
		for _, item := range r.Pv {
			pvs = append(pvs, ports.PhysicalVolume{
				Name:   item["pv_name"],
				VGName: item["vg_name"],
				Size:   item["pv_size"],
				Free:   item["pv_free"],
			})
		}
	}
	return pvs, nil
}

func (m *LinuxLVMManager) ListVolumeGroups() ([]ports.VolumeGroup, error) {
	out, err := m.runCommand("vgs", "--reportformat", "json", "--units", "g", "-o", "vg_name,vg_size,vg_free,pv_count,lv_count")
	if err != nil {
		return nil, err
	}

	var report lvmReport
	if err := json.Unmarshal([]byte(out), &report); err != nil {
		return nil, fmt.Errorf("failed to parse vgs output: %w", err)
	}

	var vgs []ports.VolumeGroup
	for _, r := range report.Report {
		for _, item := range r.Vg {
			pvCount, _ := strconv.Atoi(item["pv_count"])
			lvCount, _ := strconv.Atoi(item["lv_count"])
			
			vgs = append(vgs, ports.VolumeGroup{
				Name:    item["vg_name"],
				Size:    item["vg_size"],
				Free:    item["vg_free"],
				PVCount: pvCount,
				LVCount: lvCount,
			})
		}
	}
	return vgs, nil
}

func (m *LinuxLVMManager) ListLogicalVolumes() ([]ports.LogicalVolume, error) {
	out, err := m.runCommand("lvs", "--reportformat", "json", "--units", "g", "-o", "lv_name,vg_name,lv_path,lv_size,lv_attr")
	if err != nil {
		return nil, err
	}

	var report lvmReport
	if err := json.Unmarshal([]byte(out), &report); err != nil {
		return nil, fmt.Errorf("failed to parse lvs output: %w", err)
	}

	var lvs []ports.LogicalVolume
	for _, r := range report.Report {
		for _, item := range r.Lv {
			lvs = append(lvs, ports.LogicalVolume{
				Name:   item["lv_name"],
				VGName: item["vg_name"],
				Path:   item["lv_path"],
				Size:   item["lv_size"],
				Status: item["lv_attr"], // Raw attributes for now
			})
		}
	}
	return lvs, nil
}

func (m *LinuxLVMManager) CreateLogicalVolume(vgName string, lvName string, size string) error {
	// Validate inputs (basic)
	if vgName == "" || lvName == "" || size == "" {
		return fmt.Errorf("invalid arguments: vgName, lvName and size are required")
	}
	_, err := m.runCommand("lvcreate", "-L", size, "-n", lvName, vgName)
	return err
}

func (m *LinuxLVMManager) ExtendLogicalVolume(lvPath string, size string) error {
	if lvPath == "" || size == "" {
		return fmt.Errorf("invalid arguments: lvPath and size are required")
	}
	// -r to resize filesystem automatically
	_, err := m.runCommand("lvextend", "-L", "+"+size, lvPath, "-r") 
	return err
}

func (m *LinuxLVMManager) ReduceLogicalVolume(lvPath string, size string) error {
	if lvPath == "" || size == "" {
		return fmt.Errorf("invalid arguments: lvPath and size are required")
	}
	// -r to resize filesystem automatically (safely)
	_, err := m.runCommand("lvreduce", "-L", "-"+size, lvPath, "-r")
	return err
}

func (m *LinuxLVMManager) RemoveLogicalVolume(lvPath string) error {
	if lvPath == "" {
		return fmt.Errorf("invalid arguments: lvPath is required")
	}
	_, err := m.runCommand("lvremove", "-f", lvPath)
	return err
}

func (m *LinuxLVMManager) CreatePhysicalVolume(device string) error {
	if device == "" {
		return fmt.Errorf("device path is required")
	}
	_, err := m.runCommand("pvcreate", device)
	return err
}

func (m *LinuxLVMManager) CreateVolumeGroup(vgName string, pvs []string) error {
	if vgName == "" || len(pvs) == 0 {
		return fmt.Errorf("vgName and at least one PV are required")
	}
	args := append([]string{"vgcreate", vgName}, pvs...)
	_, err := m.runCommand(args...)
	return err
}

func (m *LinuxLVMManager) RemoveVolumeGroup(vgName string) error {
	if vgName == "" {
		return fmt.Errorf("vgName is required")
	}
	_, err := m.runCommand("vgremove", "-f", vgName)
	return err
}

func (m *LinuxLVMManager) ScanDevices() error {
	_, err := m.runCommand("pvscan")
	if err != nil {
		return err
	}
	_, err = m.runCommand("vgscan")
	return err
}
