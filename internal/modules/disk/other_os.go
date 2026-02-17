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

func (m *OtherOSDiskManager) Cleanup() error {
	return nil
}

func (m *OtherOSDiskManager) GetFilesystemUsage(path string) (ports.FilesystemUsage, error) {
	return ports.FilesystemUsage{}, errors.New("filesystem usage not supported on this OS")
}

func (m *OtherOSDiskManager) GetDiskHealth() (string, error) {
	return "", errors.New("disk health not supported on this OS")
}

func (m *OtherOSDiskManager) GetTopFiles(path string, count int) ([]ports.FileNode, error) {
	return nil, errors.New("top files not supported on this OS")
}

func (m *OtherOSDiskManager) GetInodeUsage(path string) (int, int, error) {
	return 0, 0, errors.New("inode usage not supported on this OS")
}

func (m *OtherOSDiskManager) GetMounts() ([]string, error) {
	return nil, errors.New("mounts not supported on this OS")
}

func (m *OtherOSDiskManager) Mount(source, target, fstype, options string) error {
	return errors.New("mount not supported on this OS")
}

func (m *OtherOSDiskManager) Unmount(target string) error {
	return errors.New("unmount not supported on this OS")
}

func (m *OtherOSDiskManager) Format(device, fstype string) error {
	return errors.New("format not supported on this OS")
}
