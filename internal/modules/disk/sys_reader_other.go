//go:build !linux

package disk

import (
	"errors"

	"github.com/mascli/troncli/internal/core/ports"
)

// ReadMountPoints is not supported on non-Linux systems
func (r *SysReader) ReadMountPoints() ([]ports.MountPoint, error) {
	return nil, errors.New("ReadMountPoints not supported on this OS")
}

// GetFilesystemUsage is not supported on non-Linux systems
func (r *SysReader) GetFilesystemUsage(path string) (ports.FilesystemUsage, error) {
	return ports.FilesystemUsage{}, errors.New("GetFilesystemUsage not supported on this OS")
}
