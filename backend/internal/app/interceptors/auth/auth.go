package icAuth

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"kinogo/internal/app/interceptors"
)

func AuthCheckerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Применяем только к методу AddComments
	if info.FullMethod == interceptors.AddComment || info.FullMethod == interceptors.UpdateComment || info.FullMethod == interceptors.DeleteComment {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
		}

		tokens := md.Get("token")
		if len(tokens) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "missing auth token")
		}

		// token := tokens[0]
		// проверка
		fmt.Println("проверка аутентификации")
	}

	return handler(ctx, req)
}
