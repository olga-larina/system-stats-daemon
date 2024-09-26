package fs

import (
	"github.com/olga-larina/system-stats-daemon/internal/service/collector"
)

type InodeCollector struct {
	executor collector.CommandExecutor
}

func NewInodeCollector(executor collector.CommandExecutor) collector.MetricCollector {
	return &InodeCollector{executor: executor}
}

func (c *InodeCollector) Name() string {
	return "FsInode"
}

func WithCollectorFsInodeStats(enabled bool) func(collector.CommandExecutor) collector.MetricCollector {
	if enabled {
		return NewInodeCollector
	}
	return nil
}
