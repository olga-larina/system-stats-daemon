package collector

import (
	"context"
	"fmt"
	"time"

	"github.com/olga-larina/system-stats-daemon/internal/model"
)

type MetricCollector interface {
	ExecuteCommand() ([]byte, error)
	ParseCommandOutput(output string) (any, error)
	Name() string
}

type Logger interface {
	Debug(ctx context.Context, msg string, args ...any)
	Info(ctx context.Context, msg string, args ...any)
	Warn(ctx context.Context, msg string, args ...any)
	Error(ctx context.Context, err error, msg string, args ...any)
}

type StatsCollector struct {
	logger         Logger
	collectors     []MetricCollector
	collectTimeout time.Duration
}

func NewStatsCollector(
	logger Logger,
	collectTimeout time.Duration,
	commandExecutor CommandExecutor,
	collectorOpts ...func(CommandExecutor) MetricCollector,
) *StatsCollector {
	collectors := make([]MetricCollector, 0)
	for _, opt := range collectorOpts {
		if collector := opt(commandExecutor); collector != nil {
			collectors = append(collectors, collector)
		}
	}

	return &StatsCollector{
		logger:         logger,
		collectors:     collectors,
		collectTimeout: collectTimeout,
	}
}

func (s *StatsCollector) Collect(ctx context.Context) (*model.SystemStats, error) {
	var err error
	ctx, cancel := context.WithTimeout(ctx, s.collectTimeout)
	defer cancel()

	resultChans := make([]<-chan any, 0)
	errorChans := make([]<-chan error, 0)

	// стартуем асинхронную обработку
	for _, collector := range s.collectors {
		resultChan, errorChan := s.collectAsync(ctx, collector)
		resultChans = append(resultChans, resultChan)
		errorChans = append(errorChans, errorChan)
	}

	// дожидаемся ответа всех коллекторов (ошибка или результат)
	if err = waitForErrors(ctx, errorChans...); err != nil {
		return nil, err
	}

	return waitForStats(ctx, resultChans...)
}

// обработка вызова одного коллектора.
func (s *StatsCollector) collectAsync(ctx context.Context, collector MetricCollector) (<-chan any, <-chan error) {
	resultChan := make(chan any, 1)
	errorChan := make(chan error, 1)

	go func() {
		var err error
		var output string

		defer close(errorChan)
		defer close(resultChan)
		defer func() {
			if err != nil {
				s.logger.Error(ctx, err, fmt.Sprintf("%s output", collector.Name()), output)
			}
		}()

		var out []byte
		out, err = collector.ExecuteCommand()
		if err != nil {
			errorChan <- err
			return
		}
		output = string(out)

		var result any
		result, err = collector.ParseCommandOutput(output)
		if err != nil {
			errorChan <- err
			return
		}

		resultChan <- result
	}()

	return resultChan, errorChan
}
