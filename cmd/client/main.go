package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"

	"github.com/olga-larina/system-stats-daemon/internal/client"
	"github.com/olga-larina/system-stats-daemon/internal/logger"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/stats/client/config.yaml", "Path to configuration file")
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

	// client
	client := client.NewClient(logg, config.GrpcClient.GrpcURL, config.App.SendPeriodSeconds, config.App.CalcPeriodSeconds)

	if err := client.Start(ctx); err != nil {
		logg.Error(ctx, err, "client failed to start")
		return
	}

	done, err := client.ListenStats(ctx)
	if err != nil {
		logg.Error(ctx, err, "client failed to make request")
		return
	}

	logg.Info(ctx, "system-stats client is running...")

	select {
	case <-ctx.Done():
	case <-done:
	}

	if err := client.Stop(ctx); err != nil {
		logg.Error(ctx, err, "failed to stop client")
	}

	<-done
}
