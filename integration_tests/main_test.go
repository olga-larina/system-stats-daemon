//go:build integration
// +build integration

package integration

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/olga-larina/system-stats-daemon/internal/logger"
	"github.com/olga-larina/system-stats-daemon/internal/server/grpc/pb"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var globalConfig *Config

type IntegrationTestSuite struct {
	suite.Suite
	cfg        *Config
	logg       *logger.Logger
	grpcClient pb.SystemStatsServiceClient
	grpcConn   *grpc.ClientConn
}

func (s *IntegrationTestSuite) SetupSuite() {
	var err error
	s.cfg = globalConfig
	ctx := context.Background()

	// logger
	s.logg, err = logger.New(s.cfg.Logger.Level)
	if err != nil {
		log.Fatalf("failed building logger %v", err)
		os.Exit(1)
	}

	// grpc client
	s.grpcConn, err = grpc.NewClient(s.cfg.GrpcClient.GrpcURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		s.logg.Error(ctx, err, "failed to connect to grpc")
		os.Exit(1)
	}
	s.grpcClient = pb.NewSystemStatsServiceClient(s.grpcConn)

	s.logg.Info(ctx, "suite started")
}

func (s *IntegrationTestSuite) TearDownSuite() {
	ctx := context.Background()
	if s.grpcConn != nil {
		defer func() {
			if err := s.grpcConn.Close(); err != nil {
				s.logg.Error(ctx, err, "failed to close grpc connection")
			}
		}()
	}

	s.logg.Info(ctx, "suite finished")
}

func TestMain(m *testing.M) {
	var err error

	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		configFile = "/etc/integration_tests/config.yaml"
	}

	globalConfig, err = NewConfig(configFile)
	if err != nil {
		log.Fatalf("failed reading config %v", err)
		os.Exit(1)
	}

	code := m.Run()

	os.Exit(code)
}

func TestIntegrationTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(IntegrationTestSuite))
}
