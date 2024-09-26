package internalgrpc

import (
	"context"

	"github.com/olga-larina/system-stats-daemon/internal/model"
	"github.com/olga-larina/system-stats-daemon/internal/server/grpc/pb"
)

func (s *Server) ObserveSystemStats(
	req *pb.SystemStatsRequest,
	stream pb.SystemStatsService_ObserveSystemStatsServer,
) error {
	contextFunc := func() context.Context {
		return stream.Context()
	}
	sendFunc := func(stats *model.SystemStats) error {
		return stream.Send(s.converterToProto(stats, &pb.SystemStatsPb{}))
	}

	return s.app.ObserveSystemStats(contextFunc, sendFunc, req.SendPeriod, req.CalcPeriod)
}
