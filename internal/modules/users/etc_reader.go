//go:build linux

package users

import (
	"bufio"
	"os"
	"strings"

	"github.com/mascli/troncli/internal/core/ports"
)

// EtcReader reads user and group information directly from /etc/passwd and /etc/group files
type EtcReader struct{}

// NewEtcReader creates a new EtcReader instance
func NewEtcReader() *EtcReader {
	return &EtcReader{}
}

// ReadUsers parses /etc/passwd and returns a list of users
// Format: name:password:UID:GID:comment:home:shell
func (r *EtcReader) ReadUsers() ([]ports.User, error) {
	file, err := os.Open("/etc/passwd")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var users []ports.User
	scanner := bufio.NewScanner(file)
	
	for scanner.Scan() {
		line := scanner.Text()
		
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		
		// Parse the line: name:password:UID:GID:comment:home:shell
		fields := strings.Split(line, ":")
		if len(fields) < 7 {
			continue
		}
		
		user := ports.User{
			Username: fields[0],
			UID:      fields[2],
			GID:      fields[3],
			Info:     fields[4],
			HomeDir:  fields[5],
			Shell:    fields[6],
		}
		
		users = append(users, user)
	}
	
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	
	return users, nil
}

// ReadGroups parses /etc/group and returns a list of groups
// Format: name:password:GID:members
func (r *EtcReader) ReadGroups() ([]ports.Group, error) {
	file, err := os.Open("/etc/group")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var groups []ports.Group
	scanner := bufio.NewScanner(file)
	
	for scanner.Scan() {
		line := scanner.Text()
		
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		
		// Parse the line: name:password:GID:members
		fields := strings.Split(line, ":")
		if len(fields) < 4 {
			continue
		}
		
		// Parse members (comma-separated list)
		var members []string
		if fields[3] != "" {
			members = strings.Split(fields[3], ",")
		}
		
		group := ports.Group{
			Groupname: fields[0],
			GID:       fields[2],
			Members:   members,
		}
		
		groups = append(groups, group)
	}
	
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	
	return groups, nil
}
