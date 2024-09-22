package memory

import (
	"github.com/olga-larina/system-stats-daemon/internal/model"
)

func NewStatsSummator(summatorOpts ...func() StatsSummator) StatsSummator {
	summators := make([]StatsSummator, 0)
	for _, opt := range summatorOpts {
		if summator := opt(); summator != nil {
			summators = append(summators, summator)
		}
	}

	return func(result *model.SystemStats, stats *model.SystemStats) *model.SystemStats {
		for _, summator := range summators {
			result = summator(result, stats)
		}
		return result
	}
}

func WithSummatorLoadAvgStats(enabled bool) func() StatsSummator {
	if enabled {
		return summatorLoadAvgStats
	}
	return nil
}

func WithSummatorCPUStats(enabled bool) func() StatsSummator {
	if enabled {
		return summatorCPUStats
	}
	return nil
}

func WithSummatorDisksLoadStats(enabled bool) func() StatsSummator {
	if enabled {
		return summatorDisksLoadStats
	}
	return nil
}

func WithSummatorFilesystemStats(enabled bool) func() StatsSummator {
	if enabled {
		return summatorFilesystemStats
	}
	return nil
}

func summatorLoadAvgStats() StatsSummator {
	return func(result *model.SystemStats, stats *model.SystemStats) *model.SystemStats {
		if stats.LoadAvg == nil {
			return result
		}
		if result.LoadAvg == nil {
			result.LoadAvg = &model.LoadAvgStats{}
		}
		result.LoadAvg.LoadAvg1 += stats.LoadAvg.LoadAvg1
		result.LoadAvg.LoadAvg5 += stats.LoadAvg.LoadAvg5
		result.LoadAvg.LoadAvg15 += stats.LoadAvg.LoadAvg15
		return result
	}
}

func summatorCPUStats() StatsSummator {
	return func(result *model.SystemStats, stats *model.SystemStats) *model.SystemStats {
		if stats.CPU == nil {
			return result
		}
		if result.CPU == nil {
			result.CPU = &model.CPUStats{}
		}
		result.CPU.UserMode += stats.CPU.UserMode
		result.CPU.SystemMode += stats.CPU.SystemMode
		result.CPU.Idle += stats.CPU.Idle
		return result
	}
}

func summatorDisksLoadStats() StatsSummator {
	return func(result *model.SystemStats, stats *model.SystemStats) *model.SystemStats {
		if stats.DisksLoad == nil {
			return result
		}
		if result.DisksLoad == nil {
			result.DisksLoad = &model.DisksLoadStats{
				Disks: make(map[string]*model.DiskLoad),
			}
		}
		for disk, diskLoadCached := range stats.DisksLoad.Disks {
			diskLoad, exists := result.DisksLoad.Disks[disk]
			if !exists {
				diskLoad = &model.DiskLoad{}
				result.DisksLoad.Disks[disk] = diskLoad
			}
			diskLoad.Tps += diskLoadCached.Tps
			diskLoad.Kbs += diskLoadCached.Kbs
		}
		return result
	}
}

func summatorFilesystemStats() StatsSummator {
	return func(result *model.SystemStats, stats *model.SystemStats) *model.SystemStats {
		if stats.FilesystemsMb == nil || stats.FilesystemsInode == nil {
			return result
		}
		result = summatorFilesystemSpaceStats(result, stats)
		result = summatorFilesystemInodeStats(result, stats)
		return result
	}
}

func summatorFilesystemSpaceStats(result *model.SystemStats, stats *model.SystemStats) *model.SystemStats {
	if stats.FilesystemsMb == nil {
		return result
	}
	if result.FilesystemsMb == nil {
		result.FilesystemsMb = &model.FilesystemsMbStats{
			Fs: make(map[model.Filesystem]*model.FilesystemStats),
		}
	}
	for fs, fsStatsCached := range stats.FilesystemsMb.Fs {
		fsStats, exists := result.FilesystemsMb.Fs[fs]
		if !exists {
			fsStats = &model.FilesystemStats{}
			result.FilesystemsMb.Fs[fs] = fsStats
		}
		fsStats.Used += fsStatsCached.Used
		fsStats.Total += fsStatsCached.Total
	}
	return result
}

func summatorFilesystemInodeStats(result *model.SystemStats, stats *model.SystemStats) *model.SystemStats {
	if stats.FilesystemsInode == nil {
		return result
	}
	if result.FilesystemsInode == nil {
		result.FilesystemsInode = &model.FilesystemsInodeStats{
			Fs: make(map[model.Filesystem]*model.FilesystemStats),
		}
	}
	for fs, fsStatsCached := range stats.FilesystemsInode.Fs {
		fsStats, exists := result.FilesystemsInode.Fs[fs]
		if !exists {
			fsStats = &model.FilesystemStats{}
			result.FilesystemsInode.Fs[fs] = fsStats
		}
		fsStats.Used += fsStatsCached.Used
		fsStats.Total += fsStatsCached.Total
	}
	return result
}
