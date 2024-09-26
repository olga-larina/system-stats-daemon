package collector

import (
	"context"
	"errors"
	"testing"

	"github.com/olga-larina/system-stats-daemon/internal/model"
	"github.com/stretchr/testify/require"
)

//nolint:funlen
func TestStatsPipeline(t *testing.T) {
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

	t.Run("populate CPU stats", func(t *testing.T) {
		t.Parallel()
		stats := &model.SystemStats{}
		cpu := cpu

		err := populateStats(stats, &cpu)
		require.NoError(t, err)
		require.NotNil(t, stats.CPU)
		require.Nil(t, stats.LoadAvg)
		require.Nil(t, stats.DisksLoad)
		require.Nil(t, stats.FilesystemsMb)
		require.Nil(t, stats.FilesystemsInode)
		require.Equal(t, cpu, *stats.CPU)
	})

	t.Run("populate LA stats", func(t *testing.T) {
		t.Parallel()
		stats := &model.SystemStats{}
		la := la

		err := populateStats(stats, &la)
		require.NoError(t, err)
		require.NotNil(t, stats.LoadAvg)
		require.Nil(t, stats.CPU)
		require.Nil(t, stats.DisksLoad)
		require.Nil(t, stats.FilesystemsMb)
		require.Nil(t, stats.FilesystemsInode)
		require.Equal(t, la, *stats.LoadAvg)
	})

	t.Run("populate DisksLoad stats", func(t *testing.T) {
		t.Parallel()
		stats := &model.SystemStats{}
		diskload := diskload

		err := populateStats(stats, &diskload)
		require.NoError(t, err)
		require.NotNil(t, stats.DisksLoad)
		require.Nil(t, stats.LoadAvg)
		require.Nil(t, stats.CPU)
		require.Nil(t, stats.FilesystemsMb)
		require.Nil(t, stats.FilesystemsInode)
		require.Equal(t, diskload, *stats.DisksLoad)
	})

	t.Run("populate Filesystems space stats", func(t *testing.T) {
		t.Parallel()
		stats := &model.SystemStats{}
		filesystemsMb := filesystemsMb

		err := populateStats(stats, &filesystemsMb)
		require.NoError(t, err)
		require.NotNil(t, stats.FilesystemsMb)
		require.Nil(t, stats.LoadAvg)
		require.Nil(t, stats.CPU)
		require.Nil(t, stats.DisksLoad)
		require.Nil(t, stats.FilesystemsInode)
		require.Equal(t, filesystemsMb, *stats.FilesystemsMb)
	})

	t.Run("populate Filesystems inode stats", func(t *testing.T) {
		t.Parallel()
		stats := &model.SystemStats{}
		filesystemsInode := filesystemsInode

		err := populateStats(stats, &filesystemsInode)
		require.NoError(t, err)
		require.NotNil(t, stats.FilesystemsInode)
		require.Nil(t, stats.LoadAvg)
		require.Nil(t, stats.CPU)
		require.Nil(t, stats.DisksLoad)
		require.Nil(t, stats.FilesystemsMb)
		require.Equal(t, filesystemsInode, *stats.FilesystemsInode)
	})

	t.Run("wait for successful results", func(t *testing.T) {
		t.Parallel()
		resultChan1 := make(chan any, 1)
		resultChan2 := make(chan any, 1)
		resultChan3 := make(chan any, 1)
		resultChan4 := make(chan any, 1)
		resultChan5 := make(chan any, 1)
		cpu := cpu
		la := la
		diskload := diskload
		filesystemsMb := filesystemsMb
		filesystemsInode := filesystemsInode

		resultChan1 <- &cpu
		resultChan2 <- &la
		resultChan3 <- &diskload
		resultChan4 <- &filesystemsMb
		resultChan5 <- &filesystemsInode

		close(resultChan1)
		close(resultChan2)
		close(resultChan3)
		close(resultChan4)
		close(resultChan5)

		stats, err := waitForStats(ctx, resultChan1, resultChan2, resultChan3, resultChan4, resultChan5)
		require.NoError(t, err)
		require.Equal(t, &model.SystemStats{
			CPU:              &cpu,
			LoadAvg:          &la,
			DisksLoad:        &diskload,
			FilesystemsMb:    &filesystemsMb,
			FilesystemsInode: &filesystemsInode,
		}, stats)
	})

	t.Run("wait for single result", func(t *testing.T) {
		t.Parallel()
		resultChan1 := make(chan any, 1)
		resultChan2 := make(chan any, 1)
		cpu := cpu

		resultChan1 <- &cpu

		close(resultChan1)
		close(resultChan2)

		stats, err := waitForStats(ctx, resultChan1, resultChan2)
		require.NoError(t, err)
		require.Equal(t, &model.SystemStats{
			CPU: &cpu,
		}, stats)
	})
}

func TestStatsPipelineError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("error if populating unknown stats", func(t *testing.T) {
		t.Parallel()
		stats := &model.SystemStats{}
		value := "abc"

		err := populateStats(stats, value)
		require.ErrorIs(t, err, model.ErrStatsNotValid)
		require.Nil(t, stats.LoadAvg)
		require.Nil(t, stats.CPU)
		require.Nil(t, stats.DisksLoad)
		require.Nil(t, stats.FilesystemsMb)
		require.Nil(t, stats.FilesystemsInode)
	})

	t.Run("wait for result if error", func(t *testing.T) {
		t.Parallel()
		resultChan1 := make(chan any, 1)
		resultChan2 := make(chan any, 1)

		resultChan1 <- "abc"

		close(resultChan1)
		close(resultChan2)

		stats, err := waitForStats(ctx, resultChan1, resultChan2)
		require.Nil(t, stats)
		require.ErrorIs(t, err, model.ErrStatsNotValid)
	})

	t.Run("wait for result if context is done", func(t *testing.T) {
		t.Parallel()
		resultChan1 := make(chan any, 1)
		resultChan2 := make(chan any, 1)

		ctx, cancel := context.WithCancel(ctx)
		cancel()

		stats, err := waitForStats(ctx, resultChan1, resultChan2)
		require.Nil(t, stats)
		require.ErrorIs(t, err, context.Canceled)

		close(resultChan1)
		close(resultChan2)
	})
}

func TestStatsErrors(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("wait for no errors", func(t *testing.T) {
		t.Parallel()
		errorChan1 := make(chan error, 1)
		errorChan2 := make(chan error, 1)

		close(errorChan1)
		close(errorChan2)

		err := waitForErrors(ctx, errorChan1, errorChan2)
		require.NoError(t, err)
	})

	t.Run("wait for error", func(t *testing.T) {
		t.Parallel()
		errorChan1 := make(chan error, 1)
		errorChan2 := make(chan error, 1)
		errExp := errors.New("test error")

		errorChan1 <- errExp

		close(errorChan1)
		close(errorChan2)

		err := waitForErrors(ctx, errorChan1, errorChan2)
		require.ErrorIs(t, err, errExp)
	})

	t.Run("wait for error if context is done", func(t *testing.T) {
		t.Parallel()
		errorChan1 := make(chan error, 1)
		errorChan2 := make(chan error, 1)

		ctx, cancel := context.WithCancel(ctx)
		cancel()

		err := waitForErrors(ctx, errorChan1, errorChan2)
		require.ErrorIs(t, err, context.Canceled)

		close(errorChan1)
		close(errorChan2)
	})
}
