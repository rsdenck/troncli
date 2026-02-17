//go:build !linux

package disk

import (
	"errors"

	"github.com/mascli/troncli/internal/core/ports"
)

type OtherOSDiskManager struct{}

func NewLinuxDiskManager() ports.DiskManager {
	return &OtherOSDiskManager{}
}

func (m *OtherOSDiskManager) ListBlockDevices() ([]ports.BlockDevice, error) {
	return nil, errors.New("disk management not supported on this OS")
}

func (m *OtherOSDiskManager) GetFilesystemUsage(path string) (ports.FilesystemUsage, error) {
	return ports.FilesystemUsage{}, errors.New("filesystem usage not supported on this OS")
}

func (m *OtherOSDiskManager) GetMounts() ([]string, error) {
	return nil, errors.New("mounts not supported on this OS")
}
