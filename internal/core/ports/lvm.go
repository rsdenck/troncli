package ports

// PhysicalVolume represents an LVM PV
type PhysicalVolume struct {
	Name string
	VGName string
	Size string
	Free string
}

// VolumeGroup represents an LVM VG
type VolumeGroup struct {
	Name string
	Size string
	Free string
	PVCount int
	LVCount int
}

// LogicalVolume represents an LVM LV
type LogicalVolume struct {
	Name string
	VGName string
	Path string
	Size string
	Status string
}

// LVMManager defines operations for Logical Volume Management
type LVMManager interface {
	ListPhysicalVolumes() ([]PhysicalVolume, error)
	ListVolumeGroups() ([]VolumeGroup, error)
	ListLogicalVolumes() ([]LogicalVolume, error)
	
	CreateLogicalVolume(vgName string, lvName string, size string) error
	ExtendLogicalVolume(lvPath string, size string) error
	ReduceLogicalVolume(lvPath string, size string) error
	RemoveLogicalVolume(lvPath string) error
	
	CreatePhysicalVolume(device string) error
	CreateVolumeGroup(vgName string, pvs []string) error
	RemoveVolumeGroup(vgName string) error
	
	ScanDevices() error
}
