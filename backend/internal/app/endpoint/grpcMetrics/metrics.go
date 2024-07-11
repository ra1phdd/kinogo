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
	AvgTime(userTime time.Time)
	ReturningUsers(uuid string)
	PageViews()
	Comments()
	Registrations()
}

type Endpoint struct {
	Metrics Metrics
	pb.UnimplementedMetricsV1Server
}

func (e *Endpoint) NewUser(_ context.Context, req *pb.NewUserRequest) (*pb.MetricResponse, error) {
	e.Metrics.NewUsers()

	return &pb.MetricResponse{}, nil
}

func (e *Endpoint) AvgTimeOnSite(_ context.Context, req *pb.AvgTimeOnSiteRequest) (*pb.MetricResponse, error) {
	e.Metrics.AvgTime(req.Time.AsTime())

	return &pb.MetricResponse{}, nil
}
