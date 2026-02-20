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

	// Phase 4: Universal Disk Features
	GetTopFiles(path string, count int) ([]FileNode, error)
	GetInodeUsage(path string) (int, int, error) // used, total
	Cleanup() error
	GetDiskHealth() (string, error)

	// New v1.0 commands
	Mount(source, target, fstype, options string) error
	Unmount(target string) error
	Format(device, fstype string) error
}

// FileNode represents a file or directory
type FileNode struct {
	Name  string
	Path  string
	Size  uint64
	IsDir bool
}
