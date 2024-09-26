package memory

import (
	"sync"

	"github.com/olga-larina/system-stats-daemon/internal/model"
)

type StatsRepo struct {
	statsCache []*model.SystemStats
	summator   StatsSummator
	averager   StatsAverager
	mx         sync.RWMutex
}

type (
	StatsSummator func(result *model.SystemStats, stats *model.SystemStats) *model.SystemStats
	StatsAverager func(result *model.SystemStats, period float64) *model.SystemStats
)

func NewStatsRepo(summator StatsSummator, averager StatsAverager) *StatsRepo {
	return &StatsRepo{
		statsCache: make([]*model.SystemStats, 0),
		summator:   summator,
		averager:   averager,
	}
}

func (c *StatsRepo) Update(stats *model.SystemStats, periodToKeep uint32) error {
	c.mx.Lock()
	defer c.mx.Unlock()

	c.statsCache = append(c.statsCache, stats)
	if len(c.statsCache) > int(periodToKeep) {
		firstIndex := len(c.statsCache) - int(periodToKeep)
		c.statsCache = c.statsCache[firstIndex:]
	}

	return nil
}

func (c *StatsRepo) GetAvg(period uint32) (*model.SystemStats, error) {
	c.mx.RLock()
	defer c.mx.RUnlock()

	if period == 0 {
		return nil, model.ErrPeriodNotValid
	}

	if len(c.statsCache) == 0 {
		return &model.SystemStats{
			LoadAvg:          &model.LoadAvgStats{},
			CPU:              &model.CPUStats{},
			DisksLoad:        &model.DisksLoadStats{},
			FilesystemsMb:    &model.FilesystemsMbStats{},
			FilesystemsInode: &model.FilesystemsInodeStats{},
		}, nil
	}

	firstIndex := max(0, len(c.statsCache)-int(period))

	return c.calcAvg(firstIndex, len(c.statsCache))
}

func (c *StatsRepo) calcAvg(firstIndex int, lastIndex int) (*model.SystemStats, error) {
	period := float64(lastIndex - firstIndex)
	result := &model.SystemStats{}

	for i := firstIndex; i < lastIndex; i++ {
		result = c.summator(result, c.statsCache[i])
	}

	return c.averager(result, period), nil
}
