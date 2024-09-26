package collector

import (
	"context"
	"sync"

	"github.com/olga-larina/system-stats-daemon/internal/model"
)

func waitForStats(ctx context.Context, results ...<-chan any) (*model.SystemStats, error) {
	stats := &model.SystemStats{}
	resultsChan := merge(results...)
	finished := 0
	for finished < len(results) {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case result, ok := <-resultsChan:
			if !ok {
				finished++
			} else if err := populateStats(stats, result); err != nil {
				return nil, err
			}
		}
	}
	return stats, nil
}

func populateStats(stats *model.SystemStats, value any) error {
	switch v := value.(type) {
	case *model.LoadAvgStats:
		stats.LoadAvg = v
	case model.LoadAvgStats:
		stats.LoadAvg = &v
	case *model.CPUStats:
		stats.CPU = v
	case model.CPUStats:
		stats.CPU = &v
	case *model.DisksLoadStats:
		stats.DisksLoad = v
	case model.DisksLoadStats:
		stats.DisksLoad = &v
	case *model.FilesystemsMbStats:
		stats.FilesystemsMb = v
	case model.FilesystemsMbStats:
		stats.FilesystemsMb = &v
	case *model.FilesystemsInodeStats:
		stats.FilesystemsInode = v
	case model.FilesystemsInodeStats:
		stats.FilesystemsInode = &v
	default:
		return model.ErrStatsNotValid
	}
	return nil
}

func waitForErrors(ctx context.Context, errs ...<-chan error) error {
	errorsChan := merge(errs...)
	finished := 0
	for finished < len(errs) {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err, ok := <-errorsChan:
			if err != nil {
				return err
			}
			if !ok {
				finished++
			}
		}
	}
	return nil
}

func merge[T any](channels ...<-chan T) <-chan T {
	outChan := make(chan T)
	wg := sync.WaitGroup{}

	output := func(channel <-chan T) {
		defer wg.Done()

		for val := range channel {
			outChan <- val
		}
	}

	for _, channel := range channels {
		wg.Add(1)
		go output(channel)
	}

	go func() {
		wg.Wait()
		close(outChan)
	}()

	return outChan
}
