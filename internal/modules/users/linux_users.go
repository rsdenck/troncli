//go:build linux

package users

import (
	"bufio"
	"os"
	"strings"

	"github.com/mascli/troncli/internal/core/ports"
)

type LinuxUserManager struct{}

func NewLinuxUserManager() ports.UserManager {
	return &LinuxUserManager{}
}

func (m *LinuxUserManager) ListUsers() ([]ports.User, error) {
	file, err := os.Open("/etc/passwd")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var users []ports.User
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) >= 7 {
			users = append(users, ports.User{
				Username: parts[0],
				UID:      parts[2],
				GID:      parts[3],
				Info:     parts[4],
				HomeDir:  parts[5],
				Shell:    parts[6],
			})
		}
	}
	return users, nil
}

func (m *LinuxUserManager) ListGroups() ([]ports.Group, error) {
	file, err := os.Open("/etc/group")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var groups []ports.Group
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) >= 4 {
			members := strings.Split(parts[3], ",")
			if len(members) == 1 && members[0] == "" {
				members = nil
			}
			groups = append(groups, ports.Group{
				Groupname: parts[0],
				GID:       parts[2],
				Members:   members,
			})
		}
	}
	return groups, nil
}
