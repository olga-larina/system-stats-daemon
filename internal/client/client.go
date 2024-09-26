package client

import (
	"context"
	"fmt"

	"github.com/olga-larina/system-stats-daemon/internal/server/grpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Logger interface {
	Debug(ctx context.Context, msg string, args ...any)
	Info(ctx context.Context, msg string, args ...any)
	Warn(ctx context.Context, msg string, args ...any)
	Error(ctx context.Context, err error, msg string, args ...any)
}

type Client struct {
	grpcConn   *grpc.ClientConn
	logger     Logger
	url        string
	sendPeriod uint32
	calcPeriod uint32
}

func NewClient(logger Logger, url string, sendPeriod uint32, calcPeriod uint32) *Client {
	return &Client{
		logger:     logger,
		url:        url,
		sendPeriod: sendPeriod,
		calcPeriod: calcPeriod,
	}
}

func (s *Client) Start(ctx context.Context) error {
	s.logger.Info(ctx, "starting client")

	var err error

	s.grpcConn, err = grpc.NewClient(s.url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		s.logger.Error(ctx, err, "failed to connect to grpc")
		return err
	}

	return nil
}

func (s *Client) Stop(ctx context.Context) error {
	s.logger.Info(ctx, "stopping client")
	return s.grpcConn.Close()
}

func (s *Client) ListenStats(ctx context.Context) (<-chan struct{}, error) {
	grpcClient := pb.NewSystemStatsServiceClient(s.grpcConn)

	client, err := grpcClient.ObserveSystemStats(
		ctx,
		&pb.SystemStatsRequest{
			SendPeriod: s.sendPeriod,
			CalcPeriod: s.calcPeriod,
		},
	)
	if err != nil {
		s.logger.Error(ctx, err, "failed to make request")
		return nil, err
	}

	done := make(chan struct{}, 1)
	go func() {
		defer close(done)
		for {
			msg, err := client.Recv()
			if err != nil {
				s.logger.Error(ctx, err, "receive error")
				return
			}
			// s.logger.Info(ctx, "received message", "stats", msg.String())
			s.printSystemStats(msg)
		}
	}()

	return done, nil
}

func (s *Client) printSystemStats(stats *pb.SystemStatsPb) {
	fmt.Println("System stats:")
	if stats.LoadAvgStats != nil {
		fmt.Println("  Load Average:")
		fmt.Printf("     1 minute:     %.2f\n", stats.LoadAvgStats.LoadAvg1)
		fmt.Printf("     5 minutes:    %.2f\n", stats.LoadAvgStats.LoadAvg5)
		fmt.Printf("    15 minutes:    %.2f\n", stats.LoadAvgStats.LoadAvg15)
	}
	if stats.CpuStats != nil {
		fmt.Println("  CPU:")
		fmt.Printf("    User Mode:    %.2f%%\n", stats.CpuStats.UserMode)
		fmt.Printf("    System Mode:  %.2f%%\n", stats.CpuStats.SystemMode)
		fmt.Printf("    Idle:         %.2f%%\n", stats.CpuStats.Idle)
	}
	if stats.DisksLoadStats != nil {
		fmt.Println("  DisksLoad:")
		fmt.Printf("    %-20s %-10s %-10s\n", "Disk", "Tps", "KB/s")
		for _, diskLoad := range stats.DisksLoadStats.Disks {
			fmt.Printf("    %-20s %-10.2f %-10.2f\n", diskLoad.Disk, diskLoad.Tps, diskLoad.Kbs)
		}
	}
	if stats.FilesystemsStats != nil {
		fmt.Println("  Filesystems:")
		fmt.Printf(
			"    %-20s %-40s %-10s %-10s %-15s %-15s\n",
			"Filesystem",
			"Mounted on",
			"Used MB",
			"Used %",
			"Inodes Used",
			"Inodes %",
		)
		for _, fs := range stats.FilesystemsStats.Filesystems {
			fmt.Printf(
				"    %-20s %-40s %-10.2f %-10.2f %-15.2f %-15.2f\n",
				fs.Filesystem,
				fs.MountedOn,
				fs.UsedMb,
				fs.UsedPercent,
				fs.UsedInode,
				fs.UsedInodePercent,
			)
		}
	}
}
