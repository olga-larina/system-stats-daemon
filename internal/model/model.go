package model

type SystemStats struct {
	LoadAvg          *LoadAvgStats
	CPU              *CPUStats
	DisksLoad        *DisksLoadStats
	FilesystemsMb    *FilesystemsMbStats
	FilesystemsInode *FilesystemsInodeStats
}

type LoadAvgStats struct {
	LoadAvg1  float64 // 1 minute
	LoadAvg5  float64 // 5 minutes
	LoadAvg15 float64 // 15 minutes
}

type CPUStats struct {
	UserMode   float64 // percent
	SystemMode float64 // percent
	Idle       float64 // percent
}

type DisksLoadStats struct {
	Disks map[string]*DiskLoad
}

type DiskLoad struct {
	Tps float64 // transfers per second
	Kbs float64 // kilobytes (read+write) per second
}

type FilesystemsMbStats struct {
	Fs map[Filesystem]*FilesystemStats
}

type FilesystemsInodeStats struct {
	Fs map[Filesystem]*FilesystemStats
}

type Filesystem struct {
	Name      string
	MountedOn string
}

type FilesystemStats struct {
	Used  float64
	Total float64
}
