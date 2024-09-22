package internalgrpc

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/olga-larina/system-stats-daemon/internal/model"
	"github.com/olga-larina/system-stats-daemon/internal/server/grpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

const (
	serviceNameDefault = "system-stats-app"
)

//go:generate protoc -I ../../../api SystemStatsService.proto --go_out=. --go-grpc_out=.
type Server struct {
	logger           Logger
	app              Application
	grpcPort         string
	converterToProto ConverterModelToProto
	srv              *grpc.Server
	pb.UnimplementedSystemStatsServiceServer
}

type Logger interface {
	Debug(ctx context.Context, msg string, args ...any)
	Info(ctx context.Context, msg string, args ...any)
	Warn(ctx context.Context, msg string, args ...any)
	Error(ctx context.Context, err error, msg string, args ...any)
}

type Application interface {
	ObserveSystemStats(
		context func() context.Context,
		send func(*model.SystemStats) error,
		sendPeriod uint32,
		calcPeriod uint32,
	) error
}

func NewServer(logger Logger, app Application, grpcPort string, converterToProto ConverterModelToProto) *Server {
	return &Server{
		logger:           logger,
		app:              app,
		grpcPort:         grpcPort,
		converterToProto: converterToProto,
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.logger.Info(ctx, "starting grpc server", "port", s.grpcPort)

	lsn, err := net.Listen("tcp", fmt.Sprintf(":%s", s.grpcPort))
	if err != nil {
		s.logger.Error(ctx, err, "failed to create grpc server")
		return err
	}

	s.srv = grpc.NewServer(
		grpc.ChainStreamInterceptor(
			LoggerInterceptor(s.logger),
		),
	)
	reflection.Register(s.srv)
	pb.RegisterSystemStatsServiceServer(s.srv, s)

	serviceName, found := os.LookupEnv("SERVICE_NAME")
	if !found {
		serviceName = serviceNameDefault
	}
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s.srv, healthServer)
	healthServer.SetServingStatus(serviceName, grpc_health_v1.HealthCheckResponse_SERVING)

	return s.srv.Serve(lsn)
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info(ctx, "stopping grpc server")
	s.srv.GracefulStop()
	return nil
}
