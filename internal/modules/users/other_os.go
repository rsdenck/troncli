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
