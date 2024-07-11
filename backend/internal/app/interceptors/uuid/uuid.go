package icUuid

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"kinogo/internal/app/interceptors"
	metrics "kinogo/internal/app/services/metrics"
)

func UUIDCheckerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if info.FullMethod != interceptors.GetCommentsById {
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			uuids := md.Get("uuid")
			if len(uuids) != 0 {
				m := metrics.New()
				m.UniqueUsers(uuids[0])
				m.ReturningUsers(uuids[0])
				fmt.Println("проверка uuid")
			}
		}
	}

	return handler(ctx, req)
}
