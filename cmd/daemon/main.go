package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"

	"github.com/olga-larina/system-stats-daemon/internal/app"
	"github.com/olga-larina/system-stats-daemon/internal/logger"
	internalgrpc "github.com/olga-larina/system-stats-daemon/internal/server/grpc"
	"github.com/olga-larina/system-stats-daemon/internal/service/collector"
	"github.com/olga-larina/system-stats-daemon/internal/service/collector/cpu"
	"github.com/olga-larina/system-stats-daemon/internal/service/collector/diskload"
	"github.com/olga-larina/system-stats-daemon/internal/service/collector/fs"
	"github.com/olga-larina/system-stats-daemon/internal/service/collector/la"
	"github.com/olga-larina/system-stats-daemon/internal/service/settings"
	"github.com/olga-larina/system-stats-daemon/internal/service/stats"
	"github.com/olga-larina/system-stats-daemon/internal/storage/memory"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/stats/daemon/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	config, err := NewConfig(configFile)
	if err != nil {
		log.Fatalf("failed reading config %v", err)
		return
	}

	logg, err := logger.New(config.Logger.Level)
	if err != nil {
		log.Fatalf("failed building logger %v", err)
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	// storage
	summator := memory.NewStatsSummator(
		memory.WithSummatorLoadAvgStats(config.App.Metrics.La),
		memory.WithSummatorCPUStats(config.App.Metrics.CPU),
		memory.WithSummatorDisksLoadStats(config.App.Metrics.DisksLoad),
		memory.WithSummatorFilesystemStats(config.App.Metrics.Filesystem),
	)
	averager := memory.NewStatsAverager(
		memory.WithAveragerLoadAvgStats(config.App.Metrics.La),
		memory.WithAveragerCPUStats(config.App.Metrics.CPU),
		memory.WithAveragerDisksLoadStats(config.App.Metrics.DisksLoad),
		memory.WithAveragerFilesystemStats(config.App.Metrics.Filesystem),
	)
	storage := memory.NewStatsRepo(summator, averager)

	// settings service
	settingsService := settings.NewService()

	// os command executor
	commandExecutor := collector.NewOsCommandExecutor()

	// collector
	statsCollector := collector.NewStatsCollector(
		logg,
		config.App.CollectTimeout,
		commandExecutor,
		la.WithCollectorLoadAvgStats(config.App.Metrics.La),
		cpu.WithCollectorCPUStats(config.App.Metrics.CPU),
		diskload.WithCollectorDisksLoadStats(config.App.Metrics.DisksLoad),
		fs.WithCollectorFsSpaceStats(config.App.Metrics.Filesystem),
		fs.WithCollectorFsInodeStats(config.App.Metrics.Filesystem),
	)

	// stats service
	statsService := stats.NewService(ctx, logg, storage, statsCollector, settingsService, config.App.CollectCronSpec)

	if err := statsService.Start(ctx); err != nil {
		logg.Error(ctx, err, "stats service failed to start")
		return
	}

	// app
	app := app.NewApplication(ctx, statsService, settingsService)

	// grpc server
	converter := internalgrpc.NewConverterModelToProto(
		internalgrpc.WithConverterLoadAvgStats(config.App.Metrics.La),
		internalgrpc.WithConverterCPUStats(config.App.Metrics.CPU),
		internalgrpc.WithConverterDisksLoadStats(config.App.Metrics.DisksLoad),
		internalgrpc.WithConverterFilesystemsStats(config.App.Metrics.Filesystem),
	)
	grpcServer := internalgrpc.NewServer(logg, app, config.GrpcServer.Port, converter)

	go func() {
		if err = grpcServer.Start(ctx); err != nil {
			logg.Error(ctx, err, "grpc failed to serve")
			cancel()
		}
	}()

	logg.Info(ctx, "system-stats daemon is running...")

	<-ctx.Done()

	logg.Info(ctx, "system-stats daemon is stopping...")

	if err := statsService.Stop(ctx); err != nil {
		logg.Error(ctx, err, "failed to stop stats service")
	}

	if err := grpcServer.Stop(ctx); err != nil {
		logg.Error(ctx, err, "failed to stop grpc server")
	}
}
