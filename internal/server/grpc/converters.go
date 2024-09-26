package internalgrpc

import (
	"math"

	"github.com/olga-larina/system-stats-daemon/internal/model"
	"github.com/olga-larina/system-stats-daemon/internal/server/grpc/pb"
)

type ConverterModelToProto func(stats *model.SystemStats, proto *pb.SystemStatsPb) *pb.SystemStatsPb

func NewConverterModelToProto(converterOpts ...func() ConverterModelToProto) ConverterModelToProto {
	converters := make([]ConverterModelToProto, 0)
	for _, opt := range converterOpts {
		if converter := opt(); converter != nil {
			converters = append(converters, converter)
		}
	}

	return func(stats *model.SystemStats, proto *pb.SystemStatsPb) *pb.SystemStatsPb {
		for _, converter := range converters {
			proto = converter(stats, proto)
		}
		return proto
	}
}

func WithConverterLoadAvgStats(enabled bool) func() ConverterModelToProto {
	if enabled {
		return converterLoadAvgStats
	}
	return nil
}

func WithConverterCPUStats(enabled bool) func() ConverterModelToProto {
	if enabled {
		return converterCPUStats
	}
	return nil
}

func WithConverterDisksLoadStats(enabled bool) func() ConverterModelToProto {
	if enabled {
		return converterDisksLoadStats
	}
	return nil
}

func WithConverterFilesystemsStats(enabled bool) func() ConverterModelToProto {
	if enabled {
		return converterFilesystemsStats
	}
	return nil
}

func converterLoadAvgStats() ConverterModelToProto {
	return func(stats *model.SystemStats, proto *pb.SystemStatsPb) *pb.SystemStatsPb {
		if stats.LoadAvg == nil {
			return proto
		}
		proto.LoadAvgStats = &pb.LoadAvgStatsPb{
			LoadAvg1:  stats.LoadAvg.LoadAvg1,
			LoadAvg5:  stats.LoadAvg.LoadAvg5,
			LoadAvg15: stats.LoadAvg.LoadAvg15,
		}
		return proto
	}
}

func converterCPUStats() ConverterModelToProto {
	return func(stats *model.SystemStats, proto *pb.SystemStatsPb) *pb.SystemStatsPb {
		if stats.CPU == nil {
			return proto
		}
		proto.CpuStats = &pb.CpuStatsPb{
			UserMode:   stats.CPU.UserMode,
			SystemMode: stats.CPU.SystemMode,
			Idle:       stats.CPU.Idle,
		}
		return proto
	}
}

func converterDisksLoadStats() ConverterModelToProto {
	return func(stats *model.SystemStats, proto *pb.SystemStatsPb) *pb.SystemStatsPb {
		if stats.DisksLoad == nil {
			return proto
		}
		disksLoad := make([]*pb.DiskLoadStatsPb, 0)
		for disk, diskLoad := range stats.DisksLoad.Disks {
			disksLoad = append(disksLoad, &pb.DiskLoadStatsPb{
				Disk: disk,
				Tps:  diskLoad.Tps,
				Kbs:  diskLoad.Kbs,
			})
		}
		proto.DisksLoadStats = &pb.DisksLoadStatsPb{
			Disks: disksLoad,
		}
		return proto
	}
}

func converterFilesystemsStats() ConverterModelToProto {
	return func(stats *model.SystemStats, proto *pb.SystemStatsPb) *pb.SystemStatsPb {
		if stats.FilesystemsMb == nil || stats.FilesystemsInode == nil {
			return proto
		}
		filesystemsMap := make(map[model.Filesystem]*pb.FilesystemStatsPb)
		for fs, fsStats := range stats.FilesystemsMb.Fs {
			var usedPercent float64
			if !isZero(fsStats.Total) {
				usedPercent = fsStats.Used / fsStats.Total * 100
			}
			filesystemsMap[fs] = &pb.FilesystemStatsPb{
				Filesystem:  fs.Name,
				MountedOn:   fs.MountedOn,
				UsedMb:      fsStats.Used,
				UsedPercent: usedPercent,
			}
		}
		for fs, fsStats := range stats.FilesystemsInode.Fs {
			var usedInodePercent float64
			if !isZero(fsStats.Total) {
				usedInodePercent = fsStats.Used / fsStats.Total * 100
			}
			filesystemsCached, exists := filesystemsMap[fs]
			if !exists {
				filesystemsCached = &pb.FilesystemStatsPb{
					Filesystem: fs.Name,
					MountedOn:  fs.MountedOn,
				}
				filesystemsMap[fs] = filesystemsCached
			}
			filesystemsCached.UsedInode = fsStats.Used
			filesystemsCached.UsedInodePercent = usedInodePercent
		}
		filesystems := make([]*pb.FilesystemStatsPb, 0, len(filesystemsMap))
		for _, fs := range filesystemsMap {
			filesystems = append(filesystems, fs)
		}
		proto.FilesystemsStats = &pb.FilesystemsStatsPb{
			Filesystems: filesystems,
		}
		return proto
	}
}

func isZero(value float64) bool {
	return math.Abs(value) < 1e-5
}
