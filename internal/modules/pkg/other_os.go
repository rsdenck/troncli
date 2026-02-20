//go:build !linux

package pkg

import (
	"errors"
	"github.com/mascli/troncli/internal/core/ports"
)

type OtherOSPackageManager struct{}

func NewLinuxPackageManager() ports.PackageManager {
	return &OtherOSPackageManager{}
}

func (m *OtherOSPackageManager) DetectManager() (string, error) {
	return "", errors.New("package management not supported on this OS")
}

func (m *OtherOSPackageManager) Install(packageName string) error {
	return errors.New("package management not supported on this OS")
}

func (m *OtherOSPackageManager) Remove(packageName string) error {
	return errors.New("package management not supported on this OS")
}

func (m *OtherOSPackageManager) Update() error {
	return errors.New("package management not supported on this OS")
}

func (m *OtherOSPackageManager) Upgrade() error {
	return errors.New("package management not supported on this OS")
}

func (m *OtherOSPackageManager) Search(query string) ([]ports.PackageInfo, error) {
	return nil, errors.New("package management not supported on this OS")
}
