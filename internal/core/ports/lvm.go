package ports

// PhysicalVolume represents an LVM PV
type PhysicalVolume struct {
	Name   string
	VGName string
	Size   string
	Free   string
}

// VolumeGroup represents an LVM VG
type VolumeGroup struct {
	Name    string
	Size    string
	Free    string
	PVCount int
	LVCount int
}

// LogicalVolume represents an LVM LV
type LogicalVolume struct {
	Name   string
	VGName string
	Path   string
	Size   string
	Status string
}

// LVMManager defines operations for Logical Volume Management
type LVMManager interface {
	// List operations
	ListPhysicalVolumes() ([]PhysicalVolume, error)
	ListVolumeGroups() ([]VolumeGroup, error)
	ListLogicalVolumes() ([]LogicalVolume, error)

	// Logical Volume Operations
	CreateLogicalVolume(vgName string, lvName string, size string) error
	ExtendLogicalVolume(lvPath string, size string) error
	ReduceLogicalVolume(lvPath string, size string) error
	RemoveLogicalVolume(lvPath string) error
	ResizeFileSystem(lvPath string) error // Resize FS after LV resize

	// Physical Volume Operations
	CreatePhysicalVolume(device string) error
	RemovePhysicalVolume(device string) error
	ResizePhysicalVolume(device string) error

	// Volume Group Operations
	CreateVolumeGroup(vgName string, pvs []string) error
	ExtendVolumeGroup(vgName string, pvName string) error
	ReduceVolumeGroup(vgName string, pvName string) error
	RemoveVolumeGroup(vgName string) error

	// System Operations
	ScanDevices() error // Trigger LVM scan
	RescanSCSI() error  // Trigger SCSI bus rescan for new disks
}
