package lvm

import (
	"bytes"
	"fmt"
	"os/exec"
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

func (m *LinuxLVMManager) ListPhysicalVolumes() ([]ports.PhysicalVolume, error) {
	out, err := m.runCommand("pvs", "--noheadings", "--separator", ":", "--units", "g", "-o", "pv_name,vg_name,pv_size,pv_free")
	if err != nil {
		return nil, err
	}
	
	var pvs []ports.PhysicalVolume
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		parts := strings.Split(strings.TrimSpace(line), ":")
		if len(parts) >= 4 {
			pvs = append(pvs, ports.PhysicalVolume{
				Name:   parts[0],
				VGName: parts[1],
				Size:   parts[2],
				Free:   parts[3],
			})
		}
	}
	return pvs, nil
}

func (m *LinuxLVMManager) ListVolumeGroups() ([]ports.VolumeGroup, error) {
	out, err := m.runCommand("vgs", "--noheadings", "--separator", ":", "--units", "g", "-o", "vg_name,vg_size,vg_free,pv_count,lv_count")
	if err != nil {
		return nil, err
	}
	
	var vgs []ports.VolumeGroup
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		parts := strings.Split(strings.TrimSpace(line), ":")
		if len(parts) >= 5 {
			var pvCount, lvCount int
			fmt.Sscanf(parts[3], "%d", &pvCount)
			fmt.Sscanf(parts[4], "%d", &lvCount)
			
			vgs = append(vgs, ports.VolumeGroup{
				Name:    parts[0],
				Size:    parts[1],
				Free:    parts[2],
				PVCount: pvCount,
				LVCount: lvCount,
			})
		}
	}
	return vgs, nil
}

func (m *LinuxLVMManager) ListLogicalVolumes() ([]ports.LogicalVolume, error) {
	out, err := m.runCommand("lvs", "--noheadings", "--separator", ":", "--units", "g", "-o", "lv_name,vg_name,lv_path,lv_size,lv_attr")
	if err != nil {
		return nil, err
	}
	
	var lvs []ports.LogicalVolume
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		parts := strings.Split(strings.TrimSpace(line), ":")
		if len(parts) >= 5 {
			lvs = append(lvs, ports.LogicalVolume{
				Name:   parts[0],
				VGName: parts[1],
				Path:   parts[2],
				Size:   parts[3],
				Status: parts[4], // Simplified status parsing
			})
		}
	}
	return lvs, nil
}

func (m *LinuxLVMManager) CreateLogicalVolume(vgName string, lvName string, size string) error {
	_, err := m.runCommand("lvcreate", "-L", size, "-n", lvName, vgName)
	return err
}

func (m *LinuxLVMManager) ExtendLogicalVolume(lvPath string, size string) error {
	_, err := m.runCommand("lvextend", "-L", "+"+size, lvPath, "-r") // -r to resize filesystem
	return err
}

func (m *LinuxLVMManager) ReduceLogicalVolume(lvPath string, size string) error {
	_, err := m.runCommand("lvreduce", "-L", "-"+size, lvPath, "-r")
	return err
}

func (m *LinuxLVMManager) RemoveLogicalVolume(lvPath string) error {
	_, err := m.runCommand("lvremove", "-f", lvPath)
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
