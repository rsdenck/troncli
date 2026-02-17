//go:build !linux

package users

import (
	"errors"
	"github.com/mascli/troncli/internal/core/ports"
)

type OtherOSUserManager struct{}

func NewLinuxUserManager() ports.UserManager {
	return &OtherOSUserManager{}
}

func (m *OtherOSUserManager) ListUsers() ([]ports.User, error) {
	return nil, errors.New("users management not supported on this OS")
}

func (m *OtherOSUserManager) ListGroups() ([]ports.Group, error) {
	return nil, errors.New("groups management not supported on this OS")
}

func (m *OtherOSUserManager) AddUser(username string, options ports.UserOptions) error {
	return errors.New("user management not supported on this OS")
}

func (m *OtherOSUserManager) DeleteUser(username string, removeHome bool) error {
	return errors.New("user management not supported on this OS")
}

func (m *OtherOSUserManager) ModifyUser(username string, options ports.UserOptions) error {
	return errors.New("user management not supported on this OS")
}

func (m *OtherOSUserManager) AddGroup(groupname string, gid string) error {
	return errors.New("group management not supported on this OS")
}

func (m *OtherOSUserManager) DeleteGroup(groupname string) error {
	return errors.New("group management not supported on this OS")
}
