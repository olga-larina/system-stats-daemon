package memory

import (
	"github.com/olga-larina/system-stats-daemon/internal/model"
)

func NewStatsAverager(averagerOpts ...func() StatsAverager) StatsAverager {
	averagers := make([]StatsAverager, 0)
	for _, opt := range averagerOpts {
		if averager := opt(); averager != nil {
			averagers = append(averagers, averager)
		}
	}

	return func(result *model.SystemStats, period float64) *model.SystemStats {
		for _, averager := range averagers {
			result = averager(result, period)
		}
		return result
	}
}

func WithAveragerLoadAvgStats(enabled bool) func() StatsAverager {
	if enabled {
		return averagerLoadAvgStats
	}
	return nil
}

func WithAveragerCPUStats(enabled bool) func() StatsAverager {
	if enabled {
		return averagerCPUStats
	}
	return nil
}

func WithAveragerDisksLoadStats(enabled bool) func() StatsAverager {
	if enabled {
		return averagerDisksLoadStats
	}
	return nil
}

func WithAveragerFilesystemStats(enabled bool) func() StatsAverager {
	if enabled {
		return averagerFilesystemStats
	}
	return nil
}

func averagerLoadAvgStats() StatsAverager {
	return func(result *model.SystemStats, period float64) *model.SystemStats {
		if result.LoadAvg == nil {
			return result
		}
		result.LoadAvg.LoadAvg1 /= period
		result.LoadAvg.LoadAvg5 /= period
		result.LoadAvg.LoadAvg15 /= period
		return result
	}
}

func averagerCPUStats() StatsAverager {
	return func(result *model.SystemStats, period float64) *model.SystemStats {
		if result.CPU == nil {
			return result
		}
		result.CPU.UserMode /= period
		result.CPU.SystemMode /= period
		result.CPU.Idle /= period
		return result
	}
}

func averagerDisksLoadStats() StatsAverager {
	return func(result *model.SystemStats, period float64) *model.SystemStats {
		if result.DisksLoad == nil {
			return result
		}
		for _, diskLoad := range result.DisksLoad.Disks {
			diskLoad.Tps /= period
			diskLoad.Kbs /= period
		}
		return result
	}
}

func averagerFilesystemStats() StatsAverager {
	return func(result *model.SystemStats, period float64) *model.SystemStats {
		if result.FilesystemsMb == nil || result.FilesystemsInode == nil {
			return result
		}
		result = averagerFilesystemSpaceStats(result, period)
		result = averagerFilesystemInodeStats(result, period)
		return result
	}
}

func averagerFilesystemSpaceStats(result *model.SystemStats, period float64) *model.SystemStats {
	if result.FilesystemsMb == nil {
		return result
	}
	for _, fsStats := range result.FilesystemsMb.Fs {
		fsStats.Used /= period
		fsStats.Total /= period
	}
	return result
}

func averagerFilesystemInodeStats(result *model.SystemStats, period float64) *model.SystemStats {
	if result.FilesystemsInode == nil {
		return result
	}
	for _, fsStats := range result.FilesystemsInode.Fs {
		fsStats.Used /= period
		fsStats.Total /= period
	}
	return result
}
