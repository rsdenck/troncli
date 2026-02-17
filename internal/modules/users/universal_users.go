package users

import (
	"bufio"
	"context"
	"fmt"
	"strings"

	"github.com/mascli/troncli/internal/core/adapter"
	"github.com/mascli/troncli/internal/core/domain"
	"github.com/mascli/troncli/internal/core/ports"
)

// UniversalUserManager implements UserManager using system tools
type UniversalUserManager struct {
	executor adapter.Executor
	profile  *domain.SystemProfile
}

// NewUniversalUserManager creates a new instance
func NewUniversalUserManager(executor adapter.Executor, profile *domain.SystemProfile) *UniversalUserManager {
	return &UniversalUserManager{
		executor: executor,
		profile:  profile,
	}
}

func (m *UniversalUserManager) ListUsers() ([]ports.User, error) {
	ctx := context.Background()
	// cat /etc/passwd
	res, err := m.executor.Exec(ctx, "cat", "/etc/passwd")
	if err != nil {
		return nil, err
	}

	var users []ports.User
	scanner := bufio.NewScanner(strings.NewReader(res.Stdout))
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

func (m *UniversalUserManager) ListGroups() ([]ports.Group, error) {
	ctx := context.Background()
	// cat /etc/group
	res, err := m.executor.Exec(ctx, "cat", "/etc/group")
	if err != nil {
		return nil, err
	}

	var groups []ports.Group
	scanner := bufio.NewScanner(strings.NewReader(res.Stdout))
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

func (m *UniversalUserManager) AddUser(username string, options ports.UserOptions) error {
	ctx := context.Background()
	// useradd [options] username
	args := []string{}

	if options.UID != "" {
		args = append(args, "-u", options.UID)
	}
	if options.GID != "" {
		args = append(args, "-g", options.GID)
	}
	if len(options.Groups) > 0 {
		args = append(args, "-G", strings.Join(options.Groups, ","))
	}
	if options.Shell != "" {
		args = append(args, "-s", options.Shell)
	}
	if options.HomeDir != "" {
		args = append(args, "-d", options.HomeDir)
		args = append(args, "-m") // Create home
	}
	if options.Comment != "" {
		args = append(args, "-c", options.Comment)
	}

	args = append(args, username)

	_, err := m.executor.Exec(ctx, "useradd", args...)
	if err != nil {
		return fmt.Errorf("failed to add user: %w", err)
	}
	return nil
}

func (m *UniversalUserManager) DeleteUser(username string, removeHome bool) error {
	ctx := context.Background()
	args := []string{}
	if removeHome {
		args = append(args, "-r")
	}
	args = append(args, username)

	_, err := m.executor.Exec(ctx, "userdel", args...)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

func (m *UniversalUserManager) ModifyUser(username string, options ports.UserOptions) error {
	ctx := context.Background()
	args := []string{}

	if options.UID != "" {
		args = append(args, "-u", options.UID)
	}
	if options.GID != "" {
		args = append(args, "-g", options.GID)
	}
	if len(options.Groups) > 0 {
		args = append(args, "-G", strings.Join(options.Groups, ","))
	}
	if options.Shell != "" {
		args = append(args, "-s", options.Shell)
	}
	if options.HomeDir != "" {
		args = append(args, "-d", options.HomeDir)
		args = append(args, "-m") // Move home content if possible, usually -m used with -d
	}
	if options.Comment != "" {
		args = append(args, "-c", options.Comment)
	}

	args = append(args, username)

	_, err := m.executor.Exec(ctx, "usermod", args...)
	if err != nil {
		return fmt.Errorf("failed to modify user: %w", err)
	}
	return nil
}

func (m *UniversalUserManager) AddGroup(groupname string, gid string) error {
	ctx := context.Background()
	args := []string{}

	if gid != "" {
		args = append(args, "-g", gid)
	}
	args = append(args, groupname)

	_, err := m.executor.Exec(ctx, "groupadd", args...)
	if err != nil {
		return fmt.Errorf("failed to add group: %w", err)
	}
	return nil
}

func (m *UniversalUserManager) DeleteGroup(groupname string) error {
	ctx := context.Background()

	_, err := m.executor.Exec(ctx, "groupdel", groupname)
	if err != nil {
		return fmt.Errorf("failed to delete group: %w", err)
	}
	return nil
}
