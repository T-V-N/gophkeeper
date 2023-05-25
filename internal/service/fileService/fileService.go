package service

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/T-V-N/gophkeeper/internal/app"
	"github.com/T-V-N/gophkeeper/internal/config"
	filePB "github.com/T-V-N/gophkeeper/internal/grpc/file"
	"github.com/T-V-N/gophkeeper/internal/service"
	"github.com/T-V-N/gophkeeper/internal/utils"
)

type FileService struct {
	filePB.UnimplementedFileServer
	FileApp *app.FileApp
}

func (fs *FileService) CreateFile(ctx context.Context, in *filePB.CreateFileRequest) (*filePB.CreateFileResponse, error) {
	response := filePB.CreateFileResponse{}

	uid, err := service.ExtractUIDFromCtx(ctx)

	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unathorized")
	}

	id, err := fs.FileApp.CreateFile(ctx, uid, in.Filename)

	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}

	response.Id = id

	return &response, nil
}

func (fs *FileService) RequestUpdateFile(ctx context.Context, in *filePB.UpdateFileRequest) (*filePB.UpdateFileResponse, error) {
	response := filePB.UpdateFileResponse{}

	uid, err := service.ExtractUIDFromCtx(ctx)

	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unathorized")
	}

	link, err := fs.FileApp.RequestUpdateFile(ctx, uid, in.Id)

	if err != nil {
		switch err {
		case utils.ErrNotFound:
			return nil, status.Error(codes.NotFound, "file not found")
		case utils.ErrConflict:
			return nil, status.Error(codes.Unavailable, "cannot update file note which is already being updated, sync first")
		case utils.ErrNotAuthorized:
			return nil, status.Error(codes.Unauthenticated, "unauthorized")
		case utils.ErrThirdParty:
			return nil, status.Error(codes.Unavailable, "cannot connect to third party service")
		default:
			return nil, status.Error(codes.Internal, "internal server error ;(")
		}
	}

	response.UploadLink = link

	return &response, nil
}

func (fs *FileService) CommitUpdateFile(ctx context.Context, in *filePB.CommitUpdateRequest) (*filePB.CommitUpdateResponse, error) {
	response := filePB.CommitUpdateResponse{}

	uid, err := service.ExtractUIDFromCtx(ctx)

	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unathorized")
	}

	err = fs.FileApp.CommitUpdateFile(ctx, uid, in.Id, in.CommittedAt.AsTime(), in.ForceUpdate)

	if err != nil {
		switch err {
		case utils.ErrNotFound:
			return nil, status.Error(codes.NotFound, "file not found")
		case utils.ErrConflict:
			return nil, status.Error(codes.Unavailable, "cannot update file note which is already being updated, sync first")
		case utils.ErrNotAuthorized:
			return nil, status.Error(codes.Unauthenticated, "unauthorized")
		case utils.ErrThirdParty:
			return nil, status.Error(codes.Unavailable, "cannot connect to third party service")
		default:
			return nil, status.Error(codes.Internal, "internal server error ;(")
		}
	}

	return &response, nil
}

func (fs *FileService) ListFile(ctx context.Context, in *filePB.ListFileRequest) (*filePB.ListFileResponse, error) {
	response := filePB.ListFileResponse{}

	uid, err := service.ExtractUIDFromCtx(ctx)

	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unathorized")
	}

	existingFiles := []app.ExistingFiles{}
	for _, hash := range in.ExistingFiles {
		existingFiles = append(existingFiles, app.ExistingFiles{ID: hash.Id, CommittedAt: hash.CommittedAt.AsTime()})
	}

	files, err := fs.FileApp.ListFile(ctx, uid, existingFiles)

	if err != nil {
		response.Error = err.Error()
		return &response, nil
	}

	filePBResponse := []*filePB.FileEntry{}
	for _, file := range files {
		filePBResponse = append(filePBResponse, &filePB.FileEntry{
			Id:         file.ID,
			FileName:   file.FileName,
			S3Link:     file.S3Link,
			CommitedAt: timestamppb.New(file.CommittedAt),
		})
	}

	response.Files = filePBResponse

	return &response, nil
}

func InitFileService(cfg *config.Config, a *app.FileApp) *FileService {
	return &FileService{filePB.UnimplementedFileServer{}, a}
}
