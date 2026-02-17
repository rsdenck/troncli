package ports

// BlockDevice represents a storage device
type BlockDevice struct {
	Name       string
	Size       string
	Type       string
	MountPoint string
	Children   []BlockDevice
}

// FilesystemUsage represents usage of a mount point
type FilesystemUsage struct {
	Path      string
	Total     uint64
	Used      uint64
	Free      uint64
	Files     uint64 // Inodes
	FilesFree uint64
}

// DiskManager defines operations for storage management
type DiskManager interface {
	ListBlockDevices() ([]BlockDevice, error)
	GetFilesystemUsage(path string) (FilesystemUsage, error)
	GetMounts() ([]string, error)
}
