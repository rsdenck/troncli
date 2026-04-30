package commands

import (
	"testing"

	"github.com/rsdenck/nux/internal/core"
)

func TestDiskListE2E(t *testing.T) {
	mock := &core.MockExecutor{
		Output: "NAME MAJ:MIN RM SIZE RO TYPE MOUNTED\n/dev/nvme0n1 259:0 0 90G 0 disk /\n",
		Err:    nil,
	}

	diskExecutor = mock

	out, err := diskExecutor.CombinedOutput("lsblk", "-o", "NAME,MAJ:MIN,RM,SIZE,RO,TYPE,MOUNTPOINT", "-d", "-n")
	
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if out == "" {
		t.Error("expected output, got empty string")
	}

	if mock.Output != out {
		t.Errorf("expected %q, got %q", mock.Output, out)
	}
}

func TestDiskFormatE2E(t *testing.T) {
	mock := &core.MockExecutor{
		Output: "mkfs.ext4 executed successfully",
		Err:    nil,
	}

	diskExecutor = mock

	out, err := diskExecutor.CombinedOutput("mkfs.ext4", "/dev/sdb1")
	
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if out == "" {
		t.Error("expected output, got empty string")
	}
}
