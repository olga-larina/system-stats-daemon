package cpu

import (
	"github.com/olga-larina/system-stats-daemon/internal/service/collector"
)

type Collector struct {
	executor collector.CommandExecutor
}

func NewCollector(executor collector.CommandExecutor) collector.MetricCollector {
	return &Collector{executor: executor}
}

func (c *Collector) Name() string {
	return "CPU"
}

func WithCollectorCPUStats(enabled bool) func(collector.CommandExecutor) collector.MetricCollector {
	if enabled {
		return NewCollector
	}
	return nil
}
