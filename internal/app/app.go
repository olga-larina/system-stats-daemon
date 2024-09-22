package app

import (
	"context"
	"time"

	"github.com/olga-larina/system-stats-daemon/internal/model"
)

type StatsService interface {
	GetSnapshot(calcPeriod uint32) (*model.SystemStats, error)
}

type SettingsService interface {
	Add(calcPeriod uint32)
	Remove(calcPeriod uint32) bool
}

type StatsApplication struct {
	ctx      context.Context
	stats    StatsService
	settings SettingsService
}

func NewApplication(ctx context.Context, stats StatsService, settings SettingsService) *StatsApplication {
	return &StatsApplication{
		ctx:      ctx,
		stats:    stats,
		settings: settings,
	}
}

func (app *StatsApplication) ObserveSystemStats(
	context func() context.Context,
	send func(*model.SystemStats) error,
	sendPeriodSeconds uint32,
	calcPeriodSeconds uint32,
) error {
	if sendPeriodSeconds > calcPeriodSeconds {
		return model.ErrPeriodNotValid
	}

	app.settings.Add(calcPeriodSeconds)
	defer app.settings.Remove(calcPeriodSeconds)

	firstDelay := time.NewTimer(time.Duration(calcPeriodSeconds-sendPeriodSeconds) * time.Second)

	select {
	case <-context().Done():
		firstDelay.Stop()
		return context().Err()
	case <-firstDelay.C:
	}

	ticker := time.NewTicker(time.Duration(sendPeriodSeconds) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-context().Done():
			return context().Err()
		case <-app.ctx.Done():
			return app.ctx.Err()
		case <-ticker.C:
			snapshot, err := app.stats.GetSnapshot(calcPeriodSeconds)
			if err != nil {
				return err
			}

			if err := send(snapshot); err != nil {
				return err
			}
		}
	}
}
