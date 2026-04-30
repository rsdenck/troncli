package commands

import (
	"testing"

	"github.com/rsdenck/nux/internal/core"
)

func TestServiceListE2E(t *testing.T) {
	mock := &core.MockExecutor{
		Output: "UNIT LOAD ACTIVE SUB DESCRIPTION\nssh.service loaded active running SSH server\n",
		Err:    nil,
	}

	serviceExecutor = mock

	out, err := serviceExecutor.CombinedOutput("systemctl", "list-units", "--type=service", "--all", "--no-pager")
	
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if out == "" {
		t.Error("expected output, got empty string")
	}
}

func TestServiceStartE2E(t *testing.T) {
	mock := &core.MockExecutor{
		Output: "service started",
		Err:    nil,
	}

	serviceExecutor = mock

	_, err := serviceExecutor.CombinedOutput("systemctl", "start", "nginx")
	
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
