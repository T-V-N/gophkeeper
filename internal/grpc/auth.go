package grpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/T-V-N/gophkeeper/internal/auth"
	"github.com/T-V-N/gophkeeper/internal/config"
)

func InitAuthInterceptor(cfg *config.Config) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var token string
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Internal, "cant parse metadata")
		}

		values := md.Get("auth_token")

		if len(values) == 0 {
			return nil, status.Error(codes.Unauthenticated, "not authorized")
		}

		token = values[0]
		uid, err := auth.ParseToken(token, []byte(cfg.SecretKey))

		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		md.Set("uid", uid)
		grpc.SetHeader(ctx, md)

		return handler(metadata.NewIncomingContext(ctx, md), req)
	}
}
