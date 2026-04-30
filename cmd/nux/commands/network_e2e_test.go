package commands

import (
	"testing"

	"github.com/rsdenck/nux/internal/core"
)

func TestNetworkListE2E(t *testing.T) {
	mock := &core.MockExecutor{
		Output: "INTERFACE  TYPE  STATE  CONNECTION\neth0      ether  connected  Wired Connection 1\n",
		Err:    nil,
	}

	networkExecutor = mock

	out, err := networkExecutor.CombinedOutput("nmcli", "-t", "-f", "NAME,TYPE,STATE,CONNECTION", "device", "status")
	
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if out == "" {
		t.Error("expected output, got empty string")
	}
}
