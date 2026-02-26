//go:build linux

package users

import (
	"os"
	"testing"
)

func TestEtcReader_ReadUsers(t *testing.T) {
	// Skip test if /etc/passwd doesn't exist (non-Linux systems)
	if _, err := os.Stat("/etc/passwd"); os.IsNotExist(err) {
		t.Skip("Skipping test: /etc/passwd not available (non-Linux system)")
	}

	reader := NewEtcReader()
	users, err := reader.ReadUsers()

	if err != nil {
		t.Fatalf("ReadUsers() failed: %v", err)
	}

	if len(users) == 0 {
		t.Error("Expected at least one user, got none")
	}

	// Verify that at least one user has expected fields populated
	foundValidUser := false
	for _, user := range users {
		if user.Username != "" && user.UID != "" && user.GID != "" {
			foundValidUser = true
			t.Logf("Found user: %s (UID: %s, GID: %s, Home: %s, Shell: %s)",
				user.Username, user.UID, user.GID, user.HomeDir, user.Shell)
			break
		}
	}

	if !foundValidUser {
		t.Error("No valid user found with username, UID, and GID")
	}

	// Verify root user exists
	foundRoot := false
	for _, user := range users {
		if user.Username == "root" {
			foundRoot = true
			if user.UID != "0" {
				t.Errorf("Expected root UID to be '0', got '%s'", user.UID)
			}
			t.Logf("Root user: UID=%s, GID=%s, Home=%s, Shell=%s",
				user.UID, user.GID, user.HomeDir, user.Shell)
			break
		}
	}

	if !foundRoot {
		t.Error("Expected to find root user, but didn't")
	}
}

func TestEtcReader_ReadGroups(t *testing.T) {
	// Skip test if /etc/group doesn't exist (non-Linux systems)
	if _, err := os.Stat("/etc/group"); os.IsNotExist(err) {
		t.Skip("Skipping test: /etc/group not available (non-Linux system)")
	}

	reader := NewEtcReader()
	groups, err := reader.ReadGroups()

	if err != nil {
		t.Fatalf("ReadGroups() failed: %v", err)
	}

	if len(groups) == 0 {
		t.Error("Expected at least one group, got none")
	}

	// Verify that at least one group has expected fields populated
	foundValidGroup := false
	for _, group := range groups {
		if group.Groupname != "" && group.GID != "" {
			foundValidGroup = true
			t.Logf("Found group: %s (GID: %s, Members: %v)",
				group.Groupname, group.GID, group.Members)
			break
		}
	}

	if !foundValidGroup {
		t.Error("No valid group found with groupname and GID")
	}

	// Verify root group exists
	foundRoot := false
	for _, group := range groups {
		if group.Groupname == "root" {
			foundRoot = true
			if group.GID != "0" {
				t.Errorf("Expected root GID to be '0', got '%s'", group.GID)
			}
			t.Logf("Root group: GID=%s, Members=%v", group.GID, group.Members)
			break
		}
	}

	if !foundRoot {
		t.Error("Expected to find root group, but didn't")
	}
}

func TestEtcReader_ReadUsers_FieldParsing(t *testing.T) {
	// Skip test if /etc/passwd doesn't exist (non-Linux systems)
	if _, err := os.Stat("/etc/passwd"); os.IsNotExist(err) {
		t.Skip("Skipping test: /etc/passwd not available (non-Linux system)")
	}

	reader := NewEtcReader()
	users, err := reader.ReadUsers()

	if err != nil {
		t.Fatalf("ReadUsers() failed: %v", err)
	}

	// Check that all users have required fields
	for _, user := range users {
		if user.Username == "" {
			t.Error("Found user with empty username")
		}
		if user.UID == "" {
			t.Errorf("User %s has empty UID", user.Username)
		}
		if user.GID == "" {
			t.Errorf("User %s has empty GID", user.Username)
		}
		// HomeDir and Shell should be present but can be empty for some system users
		// Info (comment) can be empty
	}

	t.Logf("Successfully parsed %d users with all required fields", len(users))
}

func TestEtcReader_ReadGroups_FieldParsing(t *testing.T) {
	// Skip test if /etc/group doesn't exist (non-Linux systems)
	if _, err := os.Stat("/etc/group"); os.IsNotExist(err) {
		t.Skip("Skipping test: /etc/group not available (non-Linux system)")
	}

	reader := NewEtcReader()
	groups, err := reader.ReadGroups()

	if err != nil {
		t.Fatalf("ReadGroups() failed: %v", err)
	}

	// Check that all groups have required fields
	for _, group := range groups {
		if group.Groupname == "" {
			t.Error("Found group with empty groupname")
		}
		if group.GID == "" {
			t.Errorf("Group %s has empty GID", group.Groupname)
		}
		// Members can be empty (nil or empty slice)
	}

	t.Logf("Successfully parsed %d groups with all required fields", len(groups))
}

func TestEtcReader_ReadGroups_MembersParsing(t *testing.T) {
	// Skip test if /etc/group doesn't exist (non-Linux systems)
	if _, err := os.Stat("/etc/group"); os.IsNotExist(err) {
		t.Skip("Skipping test: /etc/group not available (non-Linux system)")
	}

	reader := NewEtcReader()
	groups, err := reader.ReadGroups()

	if err != nil {
		t.Fatalf("ReadGroups() failed: %v", err)
	}

	// Find a group with members
	foundGroupWithMembers := false
	for _, group := range groups {
		if len(group.Members) > 0 {
			foundGroupWithMembers = true
			t.Logf("Group %s has %d members: %v",
				group.Groupname, len(group.Members), group.Members)
			
			// Verify members are not empty strings
			for _, member := range group.Members {
				if member == "" {
					t.Errorf("Group %s has empty member in members list", group.Groupname)
				}
			}
			break
		}
	}

	if !foundGroupWithMembers {
		t.Log("No groups with members found (this may be normal on some systems)")
	}
}
