package collector

import (
	"context"
	"testing"
	"time"

	"github.com/olga-larina/system-stats-daemon/internal/model"
	"github.com/olga-larina/system-stats-daemon/internal/service/collector/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

/*
Генерация моков:

	go install github.com/vektra/mockery/v2@v2.43.2

	mockery --all --case underscore --keeptree --dir internal/service/collector \
		--output internal/service/collector/mocks --with-expecter --log-level warn.
*/
func TestCollector(t *testing.T) {
	t.Parallel()

	cpu := model.CPUStats{
		UserMode:   9.33,
		SystemMode: 10.88,
		Idle:       79.77,
	}
	la := model.LoadAvgStats{
		LoadAvg1:  2.4,
		LoadAvg5:  2.91,
		LoadAvg15: 2.8,
	}
	diskload := model.DisksLoadStats{
		Disks: map[string]*model.DiskLoad{
			"disk1": {Tps: 1.23, Kbs: 12.34},
			"disk2": {Tps: 4.56, Kbs: 56.78},
			"disk3": {Tps: 7.89, Kbs: 90.12},
		},
	}
	filesystemsMb := model.FilesystemsMbStats{
		Fs: map[model.Filesystem]*model.FilesystemStats{
			{Name: "tmpfs", MountedOn: "/dev"}:          {Used: 10, Total: 356},
			{Name: "/dev/vda", MountedOn: "/etc/hosts"}: {Used: 12, Total: 583},
			{Name: "tmpfs", MountedOn: "/sys/firmware"}: {Used: 50, Total: 452},
		},
	}
	filesystemsInode := model.FilesystemsInodeStats{
		Fs: map[model.Filesystem]*model.FilesystemStats{
			{Name: "tmpfs", MountedOn: "/dev"}:          {Used: 17, Total: 834},
			{Name: "/dev/vda", MountedOn: "/etc/hosts"}: {Used: 62, Total: 78},
			{Name: "tmpfs", MountedOn: "/sys/firmware"}: {Used: 1, Total: 6},
		},
	}

	ctx := context.Background()

	t.Run("collect successful stats", func(t *testing.T) {
		t.Parallel()
		collectTimeout := time.Duration(10) * time.Second
		cpu := cpu
		la := la
		diskload := diskload
		filesystemsMb := filesystemsMb
		filesystemsInode := filesystemsInode

		mockedLogger := mocks.NewLogger(t)
		mockedCollector1 := mocks.NewMetricCollector(t)
		mockedCollector2 := mocks.NewMetricCollector(t)
		mockedCollector3 := mocks.NewMetricCollector(t)
		mockedCollector4 := mocks.NewMetricCollector(t)
		mockedCollector5 := mocks.NewMetricCollector(t)
		mockedExecutor := mocks.NewCommandExecutor(t)
		metricCollectors := []func(CommandExecutor) MetricCollector{
			func(CommandExecutor) MetricCollector { return mockedCollector1 },
			func(CommandExecutor) MetricCollector { return mockedCollector2 },
			func(CommandExecutor) MetricCollector { return mockedCollector3 },
			func(CommandExecutor) MetricCollector { return mockedCollector4 },
			func(CommandExecutor) MetricCollector { return mockedCollector5 },
		}

		collector := NewStatsCollector(mockedLogger, collectTimeout, mockedExecutor, metricCollectors...)

		mockedCollector1.EXPECT().ExecuteCommand().Return([]byte("test1"), nil)
		mockedCollector2.EXPECT().ExecuteCommand().Return([]byte("test2"), nil)
		mockedCollector3.EXPECT().ExecuteCommand().Return([]byte("test3"), nil)
		mockedCollector4.EXPECT().ExecuteCommand().Return([]byte("test4"), nil)
		mockedCollector5.EXPECT().ExecuteCommand().Return([]byte("test5"), nil)

		mockedCollector1.EXPECT().ParseCommandOutput("test1").Return(&cpu, nil)
		mockedCollector2.EXPECT().ParseCommandOutput("test2").Return(&la, nil)
		mockedCollector3.EXPECT().ParseCommandOutput("test3").Return(&diskload, nil)
		mockedCollector4.EXPECT().ParseCommandOutput("test4").Return(&filesystemsMb, nil)
		mockedCollector5.EXPECT().ParseCommandOutput("test5").Return(&filesystemsInode, nil)

		stats, err := collector.Collect(ctx)
		require.NoError(t, err)
		require.Equal(t, &model.SystemStats{
			CPU:              &cpu,
			LoadAvg:          &la,
			DisksLoad:        &diskload,
			FilesystemsMb:    &filesystemsMb,
			FilesystemsInode: &filesystemsInode,
		}, stats)
	})

	t.Run("error if one of collectors fails", func(t *testing.T) {
		t.Parallel()
		collectTimeout := time.Duration(10) * time.Second
		cpu := cpu

		mockedLogger := mocks.NewLogger(t)
		mockedCollector1 := mocks.NewMetricCollector(t)
		mockedCollector2 := mocks.NewMetricCollector(t)
		mockedExecutor := mocks.NewCommandExecutor(t)
		metricCollectors := []func(CommandExecutor) MetricCollector{
			func(CommandExecutor) MetricCollector { return mockedCollector1 },
			func(CommandExecutor) MetricCollector { return mockedCollector2 },
		}

		collector := NewStatsCollector(mockedLogger, collectTimeout, mockedExecutor, metricCollectors...)

		mockedCollector1.EXPECT().ExecuteCommand().Return([]byte("test1"), nil).Maybe()
		mockedCollector2.EXPECT().ExecuteCommand().Return([]byte("test2"), nil)

		mockedCollector1.EXPECT().ParseCommandOutput("test1").Return(&cpu, nil).Maybe()
		mockedCollector2.EXPECT().ParseCommandOutput("test2").Return(nil, model.ErrStatsNotValid)

		mockedLogger.EXPECT().Error(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return().Maybe()

		mockedCollector1.EXPECT().Name().Return("a1").Maybe()
		mockedCollector2.EXPECT().Name().Return("a2").Maybe()

		stats, err := collector.Collect(ctx)
		require.Nil(t, stats)
		require.ErrorIs(t, err, model.ErrStatsNotValid)
	})

	t.Run("error if timeout received", func(t *testing.T) {
		t.Parallel()
		collectTimeout := time.Duration(0) * time.Second
		cpu := cpu
		la := la

		mockedLogger := mocks.NewLogger(t)
		mockedCollector1 := mocks.NewMetricCollector(t)
		mockedCollector2 := mocks.NewMetricCollector(t)
		mockedExecutor := mocks.NewCommandExecutor(t)
		metricCollectors := []func(CommandExecutor) MetricCollector{
			func(CommandExecutor) MetricCollector { return mockedCollector1 },
			func(CommandExecutor) MetricCollector { return mockedCollector2 },
		}

		collector := NewStatsCollector(mockedLogger, collectTimeout, mockedExecutor, metricCollectors...)

		mockedCollector1.EXPECT().ExecuteCommand().Return([]byte("test1"), nil).Maybe()
		mockedCollector2.EXPECT().ExecuteCommand().Return([]byte("test2"), nil).Maybe()

		mockedCollector1.EXPECT().ParseCommandOutput("test1").Return(&cpu, nil).Maybe()
		mockedCollector2.EXPECT().ParseCommandOutput("test2").Return(&la, nil).Maybe()

		mockedLogger.EXPECT().Error(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return().Maybe()

		mockedCollector1.EXPECT().Name().Return("a1").Maybe()
		mockedCollector2.EXPECT().Name().Return("a2").Maybe()

		stats, err := collector.Collect(ctx)
		require.Nil(t, stats)
		require.ErrorIs(t, err, context.DeadlineExceeded)
	})
}
