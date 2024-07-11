package grpcMetrics

import (
	"context"
	pb "kinogo/pkg/metrics_v1"
	"time"
)

type Metrics interface {
	WriteFromDB()
	Reset()
	UniqueUsers(uuid string)
	NewUsers()
	SpentTime(userTime time.Time, uuid string)
	ReturningUsers(uuid string)
	NewComments()
	NewRegistrations()
}

type Endpoint struct {
	Metrics Metrics
	pb.UnimplementedMetricsV1Server
}

func (e *Endpoint) NewUser(_ context.Context, req *pb.NewUserRequest) (*pb.MetricResponse, error) {
	e.Metrics.NewUsers()

	return &pb.MetricResponse{}, nil
}

func (e *Endpoint) SpentTime(_ context.Context, req *pb.SpentTimeRequest) (*pb.MetricResponse, error) {
	e.Metrics.SpentTime(req.Time.AsTime(), req.Uuid)

	return &pb.MetricResponse{}, nil
}
