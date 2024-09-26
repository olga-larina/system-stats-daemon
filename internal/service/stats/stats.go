package stats

import (
	"context"

	"github.com/olga-larina/system-stats-daemon/internal/model"
	"github.com/robfig/cron/v3"
)

type Repo interface {
	Update(stats *model.SystemStats, periodToKeep uint32) error
	GetAvg(period uint32) (*model.SystemStats, error)
}

type Collector interface {
	Collect(ctx context.Context) (*model.SystemStats, error)
}

type SettingsService interface {
	GetMax() (uint32, bool)
}

type Logger interface {
	Debug(ctx context.Context, msg string, args ...any)
	Info(ctx context.Context, msg string, args ...any)
	Warn(ctx context.Context, msg string, args ...any)
	Error(ctx context.Context, err error, msg string, args ...any)
}

type Service struct {
	logger          Logger
	repo            Repo
	collector       Collector
	settings        SettingsService
	collectCron     *cron.Cron
	collectCronSpec string
	ctx             context.Context
}

func NewService(
	ctx context.Context,
	logger Logger,
	repo Repo,
	collector Collector,
	settings SettingsService,
	collectCronSpec string,
) *Service {
	return &Service{
		repo:            repo,
		logger:          logger,
		collector:       collector,
		settings:        settings,
		collectCron:     cron.New(cron.WithSeconds()),
		collectCronSpec: collectCronSpec,
		ctx:             ctx,
	}
}

func (s *Service) Start(ctx context.Context) error {
	s.logger.Info(ctx, "starting stats service")

	s.collectCron.AddFunc(s.collectCronSpec, func() {
		s.collect()
	})
	s.collectCron.Start()

	s.logger.Info(ctx, "started stats service")
	return nil
}

func (s *Service) Stop(ctx context.Context) error {
	s.logger.Info(ctx, "stopping stats service")

	ctx = s.collectCron.Stop()
	<-ctx.Done()

	s.logger.Info(ctx, "stopped stats service")
	return nil
}

func (s *Service) GetSnapshot(calcPeriod uint32) (*model.SystemStats, error) {
	return s.repo.GetAvg(calcPeriod)
}

func (s *Service) collect() {
	s.logger.Debug(s.ctx, "start collecting stats")

	periodToKeep, exists := s.settings.GetMax()
	if !exists {
		s.logger.Debug(s.ctx, "finished collecting stats - no clients")
		return
	}

	stats, err := s.collector.Collect(s.ctx)
	if err != nil {
		s.logger.Error(s.ctx, err, "failed collecting stats")
		return
	}

	if err := s.repo.Update(stats, periodToKeep); err != nil {
		s.logger.Error(s.ctx, err, "failed updating stats", "periodToKeep", periodToKeep)
	} else {
		s.logger.Debug(s.ctx, "succeeded collecting stats")
	}
}
