package users

import (
	"bufio"
	"os"
	"os/exec"
	"strings"

	"github.com/mascli/troncli/internal/core/adapter"
	"github.com/mascli/troncli/internal/core/ports"
)

type LinuxUserManager struct {
	executor adapter.Executor
}

func NewLinuxUserManager() ports.UserManager {
	return &LinuxUserManager{
		executor: adapter.NewExecutor(),
	}
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
		if len(parts) < 7 {
			continue
		}
		users = append(users, ports.User{
			Username: parts[0],
			UID:      parts[2],
			GID:      parts[3],
			Info:     parts[4],
			HomeDir:  parts[5],
			Shell:    parts[6],
		})
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
		if len(parts) < 4 {
			continue
		}

		members := []string{}
		if parts[3] != "" {
			members = strings.Split(parts[3], ",")
		}

		groups = append(groups, ports.Group{
			Groupname: parts[0],
			GID:       parts[2],
			Members:   members,
		})
	}
	return groups, nil
}

func (m *LinuxUserManager) AddUser(username string, options ports.UserOptions) error {
	args := []string{"useradd"}
	if options.UID != "" {
		args = append(args, "-u", options.UID)
	}
	if options.GID != "" {
		args = append(args, "-g", options.GID)
	}
	if options.Shell != "" {
		args = append(args, "-s", options.Shell)
	}
	if options.HomeDir != "" {
		args = append(args, "-d", options.HomeDir)
	}
	if options.Comment != "" {
		args = append(args, "-c", options.Comment)
	}
	args = append(args, username)

	return exec.Command("sudo", args...).Run()
}

func (m *LinuxUserManager) DeleteUser(username string, removeHome bool) error {
	args := []string{"userdel"}
	if removeHome {
		args = append(args, "-r")
	}
	args = append(args, username)
	return exec.Command("sudo", args...).Run()
}

func (m *LinuxUserManager) ModifyUser(username string, options ports.UserOptions) error {
	args := []string{"usermod"}
	if options.UID != "" {
		args = append(args, "-u", options.UID)
	}
	if options.GID != "" {
		args = append(args, "-g", options.GID)
	}
	if options.Shell != "" {
		args = append(args, "-s", options.Shell)
	}
	if options.HomeDir != "" {
		args = append(args, "-d", options.HomeDir)
	}
	if options.Comment != "" {
		args = append(args, "-c", options.Comment)
	}
	// Add other options as needed, e.g. groups
	args = append(args, username)

	return exec.Command("sudo", args...).Run()
}

func (m *LinuxUserManager) AddGroup(groupname string, gid string) error {
	args := []string{"groupadd"}
	if gid != "" {
		args = append(args, "-g", gid)
	}
	args = append(args, groupname)
	return exec.Command("sudo", args...).Run()
}

func (m *LinuxUserManager) DeleteGroup(groupname string) error {
	return exec.Command("sudo", "groupdel", groupname).Run()
}
