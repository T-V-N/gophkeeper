package service

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/T-V-N/gophkeeper/internal/app"
	"github.com/T-V-N/gophkeeper/internal/config"
	logPassPB "github.com/T-V-N/gophkeeper/internal/grpc/logPassword"
	"github.com/T-V-N/gophkeeper/internal/service"
	"github.com/T-V-N/gophkeeper/internal/utils"
)

type LogPasswordService struct {
	logPassPB.UnimplementedLogPasswordServer
	LogPasswordApp *app.LogPasswordApp
}

func (lps *LogPasswordService) CreateLogPassword(ctx context.Context, in *logPassPB.CreateLogPasswordRequest) (*logPassPB.CreateLogPasswordResponse, error) {
	response := logPassPB.CreateLogPasswordResponse{}

	uid, err := service.ExtractUIDFromCtx(ctx)

	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unathorized")
	}

	id, err := lps.LogPasswordApp.CreateLogPassword(ctx, uid, in.LoginHash, in.PasswordHash, in.ResourceName)

	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}

	response.Id = id

	return &response, nil
}

func (lps *LogPasswordService) UpdateLogPassword(ctx context.Context, in *logPassPB.UpdateLogPasswordRequest) (*logPassPB.UpdateLogPasswordResponse, error) {
	response := logPassPB.UpdateLogPasswordResponse{}

	uid, err := service.ExtractUIDFromCtx(ctx)

	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unathorized")
	}

	err = lps.LogPasswordApp.UpdateLogPassword(ctx, uid, in.Id, in.LoginHash, in.PasswordHash, in.PreviousHash, in.ResourceName, in.IsDeleted, in.ForceUpdate)

	if err != nil {
		switch err {
		case utils.ErrNotFound:
			return nil, status.Error(codes.NotFound, "logpass not found")
		case utils.ErrConflict:
			return nil, status.Error(codes.Unavailable, "cannot update logpass which is already updated, sync first")
		case utils.ErrNotAuthorized:
			return nil, status.Error(codes.Unauthenticated, "unauthorized")
		default:
			return nil, status.Error(codes.Internal, "internal server error ;(")
		}
	}

	return &response, nil
}

func (lps *LogPasswordService) ListLogPassword(ctx context.Context, in *logPassPB.ListLogPasswordRequest) (*logPassPB.ListLogPasswordResponse, error) {
	response := logPassPB.ListLogPasswordResponse{}

	uid, err := service.ExtractUIDFromCtx(ctx)

	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	existingHashes := []app.ExistingHash{}
	for _, hash := range in.ExistingHashes {
		existingHashes = append(existingHashes, app.ExistingHash{ID: hash.Id, EntryHash: hash.EntryHash})
	}

	logPasswords, err := lps.LogPasswordApp.ListLogPassword(ctx, uid, existingHashes)

	if err != nil {
		switch err {
		case utils.ErrNotFound:
			return nil, status.Error(codes.NotFound, "no log passes available")
		default:
			return nil, status.Error(codes.Internal, "internal server error ;(")
		}
	}

	logPasswordPB := []*logPassPB.LogPasswordEntry{}
	for _, logPassword := range logPasswords {
		logPasswordPB = append(logPasswordPB, &logPassPB.LogPasswordEntry{
			Id:           logPassword.ID,
			LoginHash:    logPassword.LoginHash,
			PasswordHash: logPassword.PasswordHash,
			ResourceName: logPassword.ResourceName,
			EntryHash:    logPassword.EntryHash,
		})
	}

	response.LogPasswords = logPasswordPB

	return &response, nil
}

func InitLogPasswordService(cfg *config.Config, a *app.LogPasswordApp) *LogPasswordService {
	return &LogPasswordService{logPassPB.UnimplementedLogPasswordServer{}, a}
}
