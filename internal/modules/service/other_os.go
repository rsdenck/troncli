//go:build !linux

package service

import (
	"errors"

	"github.com/mascli/troncli/internal/core/ports"
)

type OtherOSServiceManager struct{}

func NewLinuxServiceManager() ports.ServiceManager {
	return &OtherOSServiceManager{}
}

func (m *OtherOSServiceManager) ListServices() ([]ports.ServiceUnit, error) {
	return nil, errors.New("service management not supported on this OS")
}

func (m *OtherOSServiceManager) StartService(name string) error {
	return errors.New("service management not supported on this OS")
}

func (m *OtherOSServiceManager) StopService(name string) error {
	return errors.New("service management not supported on this OS")
}

func (m *OtherOSServiceManager) RestartService(name string) error {
	return errors.New("service management not supported on this OS")
}

func (m *OtherOSServiceManager) EnableService(name string) error {
	return errors.New("service management not supported on this OS")
}

func (m *OtherOSServiceManager) DisableService(name string) error {
	return errors.New("service management not supported on this OS")
}

func (m *OtherOSServiceManager) GetServiceStatus(name string) (string, error) {
	return "", errors.New("service management not supported on this OS")
}

func (m *OtherOSServiceManager) GetServiceLogs(name string, lines int) (string, error) {
	return "", errors.New("service management not supported on this OS")
}
