package commands

import (
	"testing"
)

func TestDiskCommand(t *testing.T) {
	if diskCmd == nil {
		t.Error("diskCmd should not be nil")
	}
	
	if diskCmd.Use != "disk" {
		t.Errorf("Expected disk command Use to be 'disk', got %s", diskCmd.Use)
	}
}

func TestDiskListCommand(t *testing.T) {
	if diskListCmd == nil {
		t.Error("diskListCmd should not be nil")
	}
	
	if diskListCmd.Use != "list" {
		t.Errorf("Expected Use to be 'list', got %s", diskListCmd.Use)
	}
}

func TestDiskUsageCommand(t *testing.T) {
	if diskUsageCmd == nil {
		t.Error("diskUsageCmd should not be nil")
	}
}

func TestDiskLVMCommand(t *testing.T) {
	if diskLvmCmd == nil {
		t.Error("diskLvmCmd should not be nil")
	}
}

func TestLVMCreateCommand(t *testing.T) {
	if lvmCreateCmd == nil {
		t.Error("lvmCreateCmd should not be nil")
	}
}

func TestLVMDisplayCommand(t *testing.T) {
	if lvmDisplayCmd == nil {
		t.Error("lvmDisplayCmd should not be nil")
	}
}
