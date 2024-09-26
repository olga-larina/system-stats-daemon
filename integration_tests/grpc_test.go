//go:build integration
// +build integration

package integration

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/olga-larina/system-stats-daemon/internal/server/grpc/pb"
	"github.com/stretchr/testify/require"
)

func (s *IntegrationTestSuite) TestGrpcObserveSystemStats() {
	t := s.T()

	var err error

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	client, err := s.grpcClient.ObserveSystemStats(
		ctx,
		&pb.SystemStatsRequest{
			SendPeriod: s.cfg.App.SendPeriodSeconds,
			CalcPeriod: s.cfg.App.CalcPeriodSeconds,
		},
	)
	require.NoError(t, err)

	expectedMinEvents := int64(5)
	var receivedEvents int64

	done := make(chan struct{}, 1)
	go func() {
		defer close(done)
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			msg, err := client.Recv()
			if err != nil {
				s.logg.Error(ctx, err, "client receiving stopped")
				return
			}
			s.logg.Info(ctx, "received message", "stats", msg.String())
			require.NotNil(t, msg)
			require.NotNil(t, msg.LoadAvgStats)
			require.NotNil(t, msg.CpuStats)
			require.NotNil(t, msg.DisksLoadStats)
			require.NotNil(t, msg.FilesystemsStats)
			atomic.AddInt64(&receivedEvents, 1)
		}
	}()

	// ожидание получения необходимого количества событий
	require.Eventually(
		t,
		func() bool {
			select {
			case <-ctx.Done():
				return false
			default:
				return receivedEvents >= expectedMinEvents
			}
		},
		time.Duration(s.cfg.App.CalcPeriodSeconds+uint32(expectedMinEvents+1)*s.cfg.App.SendPeriodSeconds)*time.Second,
		time.Second,
		"events were received?",
	)
	cancel()

	<-done
}
