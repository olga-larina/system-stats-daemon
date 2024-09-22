package fs

import (
	"github.com/olga-larina/system-stats-daemon/internal/service/collector"
)

type SpaceCollector struct {
	executor collector.CommandExecutor
}

func NewSpaceCollector(executor collector.CommandExecutor) collector.MetricCollector {
	return &SpaceCollector{executor: executor}
}

func (c *SpaceCollector) Name() string {
	return "FsSpace"
}

func WithCollectorFsSpaceStats(enabled bool) func(collector.CommandExecutor) collector.MetricCollector {
	if enabled {
		return NewSpaceCollector
	}
	return nil
}
