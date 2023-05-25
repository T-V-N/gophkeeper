package service

import (
	"context"

	"google.golang.org/grpc/metadata"

	"github.com/T-V-N/gophkeeper/internal/utils"
)

func ExtractUIDFromCtx(ctx context.Context) (string, error) {
	var uid string

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		values := md.Get("uid")
		if len(values) > 0 {
			uid = values[0]
			return uid, nil
		}
	}

	return "", utils.ErrAppLayer
}

// func InitUserRPCService(cfg *config.Config, a *app.UserApp) *UserService {
// 	return &UserService{userPB.UnimplementedUserServer{}, a}
// }
