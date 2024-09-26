package internalgrpc

import (
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

const (
	timeLayout = "02/Jan/2006:15:04:05 -0700"
)

func LoggerInterceptor(logger Logger) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		startTime := time.Now()

		var addr string
		if p, exists := peer.FromContext(ss.Context()); exists {
			addr = p.Addr.String()
		}
		logger.Info(ss.Context(), "grpc request started",
			"ip", addr,
			"startTime", startTime.Format(timeLayout),
			"method", info.FullMethod,
		)

		err := handler(srv, ss)

		elapsed := time.Since(startTime)
		respStatus := status.Code(err).String()
		logger.Info(ss.Context(), "grpc request finished",
			"ip", addr,
			"startTime", startTime.Format(timeLayout),
			"method", info.FullMethod,
			"statusCode", respStatus,
			"latency", elapsed.Milliseconds(),
		)

		return err
	}
}
