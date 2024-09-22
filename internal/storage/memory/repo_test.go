package memory

import (
	"testing"

	"github.com/olga-larina/system-stats-daemon/internal/model"
	"github.com/stretchr/testify/require"
)

const epsilon = 1e-6

func TestStatsRepoCalcAvgBasic(t *testing.T) {
	t.Parallel()

	summator := NewStatsSummator(
		WithSummatorLoadAvgStats(true),
		WithSummatorCPUStats(true),
		WithSummatorDisksLoadStats(true),
		WithSummatorFilesystemStats(true),
	)

	averager := NewStatsAverager(
		WithAveragerLoadAvgStats(true),
		WithAveragerCPUStats(true),
		WithAveragerDisksLoadStats(true),
		WithAveragerFilesystemStats(true),
	)

	t.Run("get avg with zero period", func(t *testing.T) {
		t.Parallel()

		repo := NewStatsRepo(summator, averager)
		stats, err := repo.GetAvg(0)
		require.Nil(t, stats)
		require.ErrorIs(t, err, model.ErrPeriodNotValid)
	})

	t.Run("get avg with no data", func(t *testing.T) {
		t.Parallel()

		repo := NewStatsRepo(summator, averager)
		stats, err := repo.GetAvg(1)
		require.NoError(t, err)
		require.Equal(t, &model.SystemStats{
			LoadAvg:          &model.LoadAvgStats{},
			CPU:              &model.CPUStats{},
			DisksLoad:        &model.DisksLoadStats{},
			FilesystemsMb:    &model.FilesystemsMbStats{},
			FilesystemsInode: &model.FilesystemsInodeStats{},
		}, stats)
	})
}

func TestStatsRepoCalcAvg(t *testing.T) {
	t.Parallel()

	summator := NewStatsSummator(
		WithSummatorLoadAvgStats(true),
		WithSummatorCPUStats(true),
		WithSummatorDisksLoadStats(true),
		WithSummatorFilesystemStats(true),
	)

	averager := NewStatsAverager(
		WithAveragerLoadAvgStats(true),
		WithAveragerCPUStats(true),
		WithAveragerDisksLoadStats(true),
		WithAveragerFilesystemStats(true),
	)

	tests := prepareTestStatsRepoCalcAvg()

	for _, tc := range tests {
		tc := tc

		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			repo := prepareDataTestStatsRepoCalcAvg(summator, averager)
			stats, err := repo.GetAvg(tc.period)
			require.NoError(t, err)
			require.InDelta(t, tc.expected.LoadAvg.LoadAvg1, stats.LoadAvg.LoadAvg1, epsilon)
			require.InDelta(t, tc.expected.LoadAvg.LoadAvg5, stats.LoadAvg.LoadAvg5, epsilon)
			require.InDelta(t, tc.expected.LoadAvg.LoadAvg15, stats.LoadAvg.LoadAvg15, epsilon)
			require.InDelta(t, tc.expected.CPU.UserMode, stats.CPU.UserMode, epsilon)
			require.InDelta(t, tc.expected.CPU.SystemMode, stats.CPU.SystemMode, epsilon)
			require.InDelta(t, tc.expected.CPU.Idle, stats.CPU.Idle, epsilon)
			disks := []string{"disk1", "disk2", "disk3"}
			for _, disk := range disks {
				expectedDisk := tc.expected.DisksLoad.Disks[disk]
				actualDisk := stats.DisksLoad.Disks[disk]
				if expectedDisk == nil {
					require.Nil(t, actualDisk)
				} else {
					require.InDelta(t, expectedDisk.Kbs, actualDisk.Kbs, epsilon)
					require.InDelta(t, expectedDisk.Tps, actualDisk.Tps, epsilon)
				}
			}
			filesystems := []model.Filesystem{
				{Name: "tmpfs", MountedOn: "/dev"},
				{Name: "/dev/vda", MountedOn: "/etc/hosts"},
				{Name: "tmpfs", MountedOn: "/sys/firmware"},
			}
			for _, fs := range filesystems {
				expectedFsSpace := tc.expected.FilesystemsMb.Fs[fs]
				actuaFsSpace := stats.FilesystemsMb.Fs[fs]
				if expectedFsSpace == nil {
					require.Nil(t, actuaFsSpace)
				} else {
					require.InDelta(t, expectedFsSpace.Used, actuaFsSpace.Used, epsilon)
					require.InDelta(t, expectedFsSpace.Total, actuaFsSpace.Total, epsilon)
				}
				expectedFsInode := tc.expected.FilesystemsInode.Fs[fs]
				actuaFsInode := stats.FilesystemsInode.Fs[fs]
				if expectedFsInode == nil {
					require.Nil(t, actuaFsInode)
				} else {
					require.InDelta(t, expectedFsInode.Used, actuaFsInode.Used, epsilon)
					require.InDelta(t, expectedFsInode.Total, actuaFsInode.Total, epsilon)
				}
			}
		})
	}
}

//nolint:dupl
func prepareTestStatsRepoCalcAvg() []struct {
	period   uint32
	expected model.SystemStats
	testName string
} {
	return []struct {
		period   uint32
		expected model.SystemStats
		testName string
	}{
		{
			period: 4,
			expected: model.SystemStats{
				LoadAvg: &model.LoadAvgStats{LoadAvg1: 2.20333333, LoadAvg5: 2.21, LoadAvg15: 2.64333333},
				CPU:     &model.CPUStats{UserMode: 6.96, SystemMode: 8.55333333, Idle: 84.48},
				DisksLoad: &model.DisksLoadStats{
					Disks: map[string]*model.DiskLoad{
						"disk1": {Tps: 0.75, Kbs: 14.82333333},
						"disk2": {Tps: 2.25666667, Kbs: 33.77666667},
						"disk3": {Tps: 2.63, Kbs: 30.04},
					},
				},
				FilesystemsMb: &model.FilesystemsMbStats{
					Fs: map[model.Filesystem]*model.FilesystemStats{
						{Name: "tmpfs", MountedOn: "/dev"}:          {Used: 3.66666667, Total: 121},
						{Name: "/dev/vda", MountedOn: "/etc/hosts"}: {Used: 4, Total: 194.33333333},
						{Name: "tmpfs", MountedOn: "/sys/firmware"}: {Used: 27.66666667, Total: 165.666666667},
					},
				},
				FilesystemsInode: &model.FilesystemsInodeStats{
					Fs: map[model.Filesystem]*model.FilesystemStats{
						{Name: "tmpfs", MountedOn: "/dev"}:          {Used: 7.33333333, Total: 285},
						{Name: "/dev/vda", MountedOn: "/etc/hosts"}: {Used: 20.66666667, Total: 26},
						{Name: "tmpfs", MountedOn: "/sys/firmware"}: {Used: 1.33333333, Total: 4.666666667},
					},
				},
			},
			testName: "get avg with period >= number of data in repo",
		},
		{
			period: 2,
			expected: model.SystemStats{
				LoadAvg: &model.LoadAvgStats{LoadAvg1: 2.105, LoadAvg5: 1.86, LoadAvg15: 2.565},
				CPU:     &model.CPUStats{UserMode: 5.775, SystemMode: 7.39, Idle: 86.835},
				DisksLoad: &model.DisksLoadStats{
					Disks: map[string]*model.DiskLoad{
						"disk1": {Tps: 0.615, Kbs: 6.17},
						"disk2": {Tps: 3.385, Kbs: 50.665},
						"disk3": {Tps: 3.945, Kbs: 45.06},
					},
				},
				FilesystemsMb: &model.FilesystemsMbStats{
					Fs: map[model.Filesystem]*model.FilesystemStats{
						{Name: "tmpfs", MountedOn: "/dev"}:          {Used: 5, Total: 178},
						{Name: "/dev/vda", MountedOn: "/etc/hosts"}: {Used: 6, Total: 291.5},
						{Name: "tmpfs", MountedOn: "/sys/firmware"}: {Used: 41.5, Total: 248.5},
					},
				},
				FilesystemsInode: &model.FilesystemsInodeStats{
					Fs: map[model.Filesystem]*model.FilesystemStats{
						{Name: "tmpfs", MountedOn: "/dev"}:          {Used: 8.5, Total: 417},
						{Name: "/dev/vda", MountedOn: "/etc/hosts"}: {Used: 31, Total: 39},
						{Name: "tmpfs", MountedOn: "/sys/firmware"}: {Used: 2, Total: 7},
					},
				},
			},
			testName: "get avg with period < number of data in repo",
		},
		{
			period: 1,
			expected: model.SystemStats{
				LoadAvg: &model.LoadAvgStats{LoadAvg1: 1.11, LoadAvg5: 2.22, LoadAvg15: 3.33},
				CPU:     &model.CPUStats{UserMode: 4.44, SystemMode: 5.55, Idle: 90.01},
				DisksLoad: &model.DisksLoadStats{
					Disks: map[string]*model.DiskLoad{
						"disk2": {Tps: 2.21, Kbs: 44.55},
					},
				},
				FilesystemsMb: &model.FilesystemsMbStats{
					Fs: map[model.Filesystem]*model.FilesystemStats{
						{Name: "tmpfs", MountedOn: "/sys/firmware"}: {Used: 33, Total: 45},
					},
				},
				FilesystemsInode: &model.FilesystemsInodeStats{
					Fs: map[model.Filesystem]*model.FilesystemStats{
						{Name: "tmpfs", MountedOn: "/sys/firmware"}: {Used: 3, Total: 8},
					},
				},
			},
			testName: "get avg with period == 1",
		},
	}
}

//nolint:dupl
func prepareDataTestStatsRepoCalcAvg(summator StatsSummator, averager StatsAverager) *StatsRepo {
	repo := NewStatsRepo(summator, averager)
	repo.statsCache = append(repo.statsCache, &model.SystemStats{
		LoadAvg: &model.LoadAvgStats{LoadAvg1: 2.4, LoadAvg5: 2.91, LoadAvg15: 2.8},
		CPU:     &model.CPUStats{UserMode: 9.33, SystemMode: 10.88, Idle: 79.77},
		DisksLoad: &model.DisksLoadStats{
			Disks: map[string]*model.DiskLoad{
				"disk1": {Tps: 1.02, Kbs: 32.13},
			},
		},
		FilesystemsMb: &model.FilesystemsMbStats{
			Fs: map[model.Filesystem]*model.FilesystemStats{
				{Name: "tmpfs", MountedOn: "/dev"}: {Used: 1, Total: 7},
			},
		},
		FilesystemsInode: &model.FilesystemsInodeStats{
			Fs: map[model.Filesystem]*model.FilesystemStats{
				{Name: "tmpfs", MountedOn: "/dev"}: {Used: 5, Total: 21},
			},
		},
	})
	repo.statsCache = append(repo.statsCache, &model.SystemStats{
		LoadAvg: &model.LoadAvgStats{LoadAvg1: 3.1, LoadAvg5: 1.5, LoadAvg15: 1.8},
		CPU:     &model.CPUStats{UserMode: 7.11, SystemMode: 9.23, Idle: 83.66},
		DisksLoad: &model.DisksLoadStats{
			Disks: map[string]*model.DiskLoad{
				"disk1": {Tps: 1.23, Kbs: 12.34},
				"disk2": {Tps: 4.56, Kbs: 56.78},
				"disk3": {Tps: 7.89, Kbs: 90.12},
			},
		},
		FilesystemsMb: &model.FilesystemsMbStats{
			Fs: map[model.Filesystem]*model.FilesystemStats{
				{Name: "tmpfs", MountedOn: "/dev"}:          {Used: 10, Total: 356},
				{Name: "/dev/vda", MountedOn: "/etc/hosts"}: {Used: 12, Total: 583},
				{Name: "tmpfs", MountedOn: "/sys/firmware"}: {Used: 50, Total: 452},
			},
		},
		FilesystemsInode: &model.FilesystemsInodeStats{
			Fs: map[model.Filesystem]*model.FilesystemStats{
				{Name: "tmpfs", MountedOn: "/dev"}:          {Used: 17, Total: 834},
				{Name: "/dev/vda", MountedOn: "/etc/hosts"}: {Used: 62, Total: 78},
				{Name: "tmpfs", MountedOn: "/sys/firmware"}: {Used: 1, Total: 6},
			},
		},
	})
	repo.statsCache = append(repo.statsCache, &model.SystemStats{
		LoadAvg: &model.LoadAvgStats{LoadAvg1: 1.11, LoadAvg5: 2.22, LoadAvg15: 3.33},
		CPU:     &model.CPUStats{UserMode: 4.44, SystemMode: 5.55, Idle: 90.01},
		DisksLoad: &model.DisksLoadStats{
			Disks: map[string]*model.DiskLoad{
				"disk2": {Tps: 2.21, Kbs: 44.55},
			},
		},
		FilesystemsMb: &model.FilesystemsMbStats{
			Fs: map[model.Filesystem]*model.FilesystemStats{
				{Name: "tmpfs", MountedOn: "/sys/firmware"}: {Used: 33, Total: 45},
			},
		},
		FilesystemsInode: &model.FilesystemsInodeStats{
			Fs: map[model.Filesystem]*model.FilesystemStats{
				{Name: "tmpfs", MountedOn: "/sys/firmware"}: {Used: 3, Total: 8},
			},
		},
	})
	return repo
}

func TestStatsRepoUpdate(t *testing.T) {
	t.Parallel()

	summator := NewStatsSummator(
		WithSummatorLoadAvgStats(true),
		WithSummatorCPUStats(true),
		WithSummatorDisksLoadStats(true),
		WithSummatorFilesystemStats(true),
	)

	averager := NewStatsAverager(
		WithAveragerLoadAvgStats(true),
		WithAveragerCPUStats(true),
		WithAveragerDisksLoadStats(true),
		WithAveragerFilesystemStats(true),
	)

	t.Run("update with period < number of data in repo", func(t *testing.T) {
		t.Parallel()

		repo := NewStatsRepo(summator, averager)
		repo.statsCache = append(repo.statsCache, &model.SystemStats{})
		repo.statsCache = append(repo.statsCache, &model.SystemStats{})
		repo.statsCache = append(repo.statsCache, &model.SystemStats{})

		err := repo.Update(&model.SystemStats{}, 2)
		require.NoError(t, err)
		require.Equal(t, 2, len(repo.statsCache))
	})

	t.Run("update with period > number of data in repo", func(t *testing.T) {
		t.Parallel()

		repo := NewStatsRepo(summator, averager)
		repo.statsCache = append(repo.statsCache, &model.SystemStats{})
		repo.statsCache = append(repo.statsCache, &model.SystemStats{})

		err := repo.Update(&model.SystemStats{}, 10)
		require.NoError(t, err)
		require.Equal(t, 3, len(repo.statsCache))
	})

	t.Run("update with zero period", func(t *testing.T) {
		t.Parallel()

		repo := NewStatsRepo(summator, averager)
		repo.statsCache = append(repo.statsCache, &model.SystemStats{})

		err := repo.Update(&model.SystemStats{}, 0)
		require.NoError(t, err)
		require.Equal(t, 0, len(repo.statsCache))
	})
}
