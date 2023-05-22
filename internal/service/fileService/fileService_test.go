package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/T-V-N/gophkeeper/internal/app"
	"github.com/T-V-N/gophkeeper/internal/config"
	filePB "github.com/T-V-N/gophkeeper/internal/grpc/file"
	mock "github.com/T-V-N/gophkeeper/internal/mocks"
	service "github.com/T-V-N/gophkeeper/internal/service/fileService"
	"github.com/T-V-N/gophkeeper/internal/storage"
	"github.com/caarlos0/env/v8"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestFileService_CreateFile(t *testing.T) {
	type want struct {
		errCode  codes.Code
		response filePB.CreateFileResponse
	}

	tests := []struct {
		name       string
		request    filePB.CreateFileRequest
		want       want
		ctx        context.Context
		omitCreate bool
	}{
		{
			name: "create file normally",
			request: filePB.CreateFileRequest{
				Filename: "jsadasjdjs",
			},
			want: want{
				errCode:  codes.OK,
				response: filePB.CreateFileResponse{},
			},
			ctx: metadata.NewIncomingContext(context.Background(), metadata.Pairs("uid", "1")),
		},
		{
			name: "create file unauthorized",
			request: filePB.CreateFileRequest{
				Filename: "jsadasjdjs",
			},
			want: want{
				errCode: codes.Unauthenticated,
				response: filePB.CreateFileResponse{
					Id: "1",
				},
			},
			ctx:        context.Background(),
			omitCreate: true,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{}
	env.Parse(cfg)

	f := mock.NewMockFile(ctrl)

	fApp := app.FileApp{File: f, Cfg: cfg}
	fileService := service.FileService{filePB.UnimplementedFileServer{}, &fApp}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.omitCreate {
				f.EXPECT().CreateFile(tt.ctx, "1", tt.request.Filename).Return(tt.want.response.Id, nil)
			}
			_, err := fileService.CreateFile(tt.ctx, &tt.request)

			if unwrapError, ok := status.FromError(err); ok {
				assert.Equal(t, unwrapError.Code(), tt.want.errCode)
			} else {
				t.Error("can't parse error, unexpected test result")
			}
		})
	}
}

func TestFileService_RequestUpdateFile(t *testing.T) {
	type want struct {
		errCode  codes.Code
		response filePB.UpdateFileResponse
	}

	tests := []struct {
		name         string
		request      filePB.UpdateFileRequest
		want         want
		ctx          context.Context
		existingFile storage.File
		omitS3Call   bool
	}{
		{
			name: "update file normally",
			request: filePB.UpdateFileRequest{
				Id: "1",
			},
			existingFile: storage.File{
				ID:          "1",
				UID:         "1",
				FileName:    "test",
				S3Link:      "test",
				IsDeleted:   false,
				CommittedAt: time.Unix(0, 0),
			},
			want: want{
				errCode:  codes.OK,
				response: filePB.UpdateFileResponse{},
			},
			ctx: metadata.NewIncomingContext(context.Background(), metadata.Pairs("uid", "1")),
		},
		{
			name: "update file while update on cooldown",
			request: filePB.UpdateFileRequest{
				Id: "1",
			},
			existingFile: storage.File{
				ID:          "1",
				UID:         "1",
				FileName:    "test",
				S3Link:      "test",
				IsDeleted:   false,
				CommittedAt: time.Now(),
			},
			omitS3Call: true,
			want: want{
				errCode:  codes.Unavailable,
				response: filePB.UpdateFileResponse{},
			},
			ctx: metadata.NewIncomingContext(context.Background(), metadata.Pairs("uid", "1")),
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{}
	env.Parse(cfg)

	f := mock.NewMockFile(ctrl)
	s3 := mock.NewMockS3Store(ctrl)

	fApp := app.FileApp{File: f, Cfg: cfg, S3: s3}
	fileService := service.FileService{filePB.UnimplementedFileServer{}, &fApp}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f.EXPECT().GetFileByID(tt.ctx, tt.request.Id).Return(&tt.existingFile, nil)
			if !tt.omitS3Call {
				s3.EXPECT().GetUploadLink(tt.ctx, tt.request.Id).Return("link", nil)
			}

			_, err := fileService.RequestUpdateFile(tt.ctx, &tt.request)

			if unwrapError, ok := status.FromError(err); ok {
				assert.Equal(t, unwrapError.Code(), tt.want.errCode)
			} else {
				t.Error("can't parse error, unexpected test result")
			}
		})
	}
}

func TestFileService_CommitUpdateFile(t *testing.T) {
	type want struct {
		errCode  codes.Code
		response filePB.CommitUpdateResponse
	}

	tests := []struct {
		name         string
		request      filePB.CommitUpdateRequest
		want         want
		ctx          context.Context
		existingFile storage.File
		S3CommitedAt time.Time
		omitS3Call   bool
		omitUpdate   bool
	}{
		{
			name: "commit file normally",
			request: filePB.CommitUpdateRequest{
				Id:          "1",
				CommittedAt: timestamppb.New(time.Unix(0, 0)),
				ForceUpdate: false,
			},
			existingFile: storage.File{
				ID:          "1",
				UID:         "1",
				FileName:    "test",
				S3Link:      "test",
				IsDeleted:   false,
				CommittedAt: time.Unix(0, 0),
			},
			want: want{
				errCode:  codes.OK,
				response: filePB.CommitUpdateResponse{},
			},
			ctx: metadata.NewIncomingContext(context.Background(), metadata.Pairs("uid", "1")),
		},
		{
			name: "commit file which already updated",
			request: filePB.CommitUpdateRequest{
				Id:          "1",
				CommittedAt: timestamppb.New(time.Unix(1, 0)),
				ForceUpdate: false,
			},
			existingFile: storage.File{
				ID:          "1",
				UID:         "1",
				FileName:    "test",
				S3Link:      "test",
				IsDeleted:   false,
				CommittedAt: time.Unix(2, 0),
			},
			S3CommitedAt: time.Unix(3, 0),
			want: want{
				errCode:  codes.Unavailable,
				response: filePB.CommitUpdateResponse{},
			},
			omitUpdate: true,
			ctx:        metadata.NewIncomingContext(context.Background(), metadata.Pairs("uid", "1")),
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{}
	env.Parse(cfg)

	f := mock.NewMockFile(ctrl)
	s3 := mock.NewMockS3Store(ctrl)

	fApp := app.FileApp{File: f, Cfg: cfg, S3: s3}
	fileService := service.FileService{filePB.UnimplementedFileServer{}, &fApp}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f.EXPECT().GetFileByID(tt.ctx, tt.request.Id).Return(&tt.existingFile, nil)
			if !tt.omitS3Call {
				s3.EXPECT().GetFileUpdatedAt(tt.ctx, tt.request.Id).Return(tt.S3CommitedAt, nil)
			}
			if !tt.omitUpdate {
				f.EXPECT().UpdateFile(tt.ctx, tt.request.Id, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			}
			_, err := fileService.CommitUpdateFile(tt.ctx, &tt.request)

			if unwrapError, ok := status.FromError(err); ok {
				assert.Equal(t, unwrapError.Code(), tt.want.errCode)
			} else {
				t.Error("can't parse error, unexpected test result")
			}
		})
	}
}

func TestFileService_ListFile(t *testing.T) {
	type want struct {
		errCode  codes.Code
		response filePB.ListFileResponse
	}

	tests := []struct {
		name            string
		request         filePB.ListFileRequest
		existingFiles   []storage.File
		existingFileErr error
		want            want
		ctx             context.Context
		uid             string
	}{
		{
			name:          "list file normally",
			request:       filePB.ListFileRequest{},
			existingFiles: []storage.File{{CommittedAt: time.Unix(0, 0), UID: "1"}, {CommittedAt: time.Unix(0, 0), UID: "1"}, {CommittedAt: time.Unix(0, 0), UID: "1"}},
			want: want{
				errCode: codes.OK,
				response: filePB.ListFileResponse{
					Files: []*filePB.FileEntry{{CommitedAt: timestamppb.New(time.Unix(0, 0))}, {CommitedAt: timestamppb.New(time.Unix(0, 0))}, {CommitedAt: timestamppb.New(time.Unix(0, 0))}},
				},
			},
			ctx: metadata.NewIncomingContext(context.Background(), metadata.Pairs("uid", "1")),
			uid: "1",
		},
		{
			name: "list file but remove some existing files",
			request: filePB.ListFileRequest{
				ExistingFiles: []*filePB.ExistingFile{
					{Id: "1", CommittedAt: timestamppb.New(time.Unix(0, 0))},
					{Id: "3", CommittedAt: timestamppb.New(time.Unix(0, 0))},
				},
			},
			existingFiles: []storage.File{{ID: "1", CommittedAt: time.Unix(0, 0), UID: "1"}, {ID: "2", CommittedAt: time.Unix(0, 0), UID: "1"}, {ID: "3", CommittedAt: time.Unix(0, 0), UID: "1"}, {ID: "4", CommittedAt: time.Unix(0, 0), UID: "4"}},
			want: want{
				errCode: codes.OK,
				response: filePB.ListFileResponse{
					Files: []*filePB.FileEntry{{CommitedAt: timestamppb.New(time.Unix(0, 0))}, {CommitedAt: timestamppb.New(time.Unix(0, 0))}},
				},
			},
			ctx: metadata.NewIncomingContext(context.Background(), metadata.Pairs("uid", "1")),
			uid: "1",
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{}
	env.Parse(cfg)

	tn := mock.NewMockFile(ctrl)

	tnApp := app.FileApp{File: tn, Cfg: cfg}
	tnService := service.FileService{filePB.UnimplementedFileServer{}, &tnApp}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tn.EXPECT().ListFilesByUID(tt.ctx, tt.uid).Return(&tt.existingFiles, nil)
			resp, err := tnService.ListFile(tt.ctx, &tt.request)

			if unwrapError, ok := status.FromError(err); ok {
				assert.Equal(t, unwrapError.Code(), tt.want.errCode)
				assert.Equal(t, len(resp.Files), len(tt.want.response.Files))

			} else {
				t.Error("can't parse error, unexpected test result")
			}
		})
	}
}
