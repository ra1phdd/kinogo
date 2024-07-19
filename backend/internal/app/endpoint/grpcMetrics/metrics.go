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
	StreamingPerformance(uuid string, movieId, bufferingCount, bufferingTime int32, playbackError string, viewsTime, duration int32)
}

type Endpoint struct {
	Metrics Metrics
	pb.UnimplementedMetricsV1Server
}

func (e *Endpoint) NewUser(_ context.Context, _ *pb.NewUserRequest) (*pb.MetricResponse, error) {
	e.Metrics.NewUsers()

	return &pb.MetricResponse{}, nil
}

func (e *Endpoint) SpentTime(_ context.Context, req *pb.SpentTimeRequest) (*pb.MetricResponse, error) {
	e.Metrics.SpentTime(req.Time.AsTime(), req.Uuid)

	return &pb.MetricResponse{}, nil
}

func (e *Endpoint) StreamingPerformance(_ context.Context, req *pb.StreamingPerformanceRequest) (*pb.MetricResponse, error) {
	e.Metrics.StreamingPerformance(req.Uuid, req.MovieId, req.BufferingCount, req.BufferingTime, req.PlaybackError, req.ViewsTime, req.Duration)

	return &pb.MetricResponse{}, nil
}
