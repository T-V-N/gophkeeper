package service_test

import (
	"context"
	"testing"

	"github.com/caarlos0/env/v8"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/T-V-N/gophkeeper/internal/app"
	"github.com/T-V-N/gophkeeper/internal/config"
	lpPB "github.com/T-V-N/gophkeeper/internal/grpc/logPassword"
	mock "github.com/T-V-N/gophkeeper/internal/mocks"
	service "github.com/T-V-N/gophkeeper/internal/service/logPasswordService"
	"github.com/T-V-N/gophkeeper/internal/storage"
	"github.com/T-V-N/gophkeeper/internal/utils"
)

func TestLogPasswordService_CreateLogPassword(t *testing.T) {
	type want struct {
		errCode  codes.Code
		response lpPB.CreateLogPasswordResponse
	}

	tests := []struct {
		name         string
		request      lpPB.CreateLogPasswordRequest
		want         want
		ctx          context.Context
		omitCreateLP bool
	}{
		{
			name: "create card normally",
			request: lpPB.CreateLogPasswordRequest{
				LoginHash:    "jsadasjdjs",
				PasswordHash: "shdahsdhahd",
				ResourceName: "yandex.ru",
			},
			want: want{
				errCode:  codes.OK,
				response: lpPB.CreateLogPasswordResponse{},
			},
			ctx: metadata.NewIncomingContext(context.Background(), metadata.Pairs("uid", "1")),
		},
		{
			name: "create card unauthorized",
			request: lpPB.CreateLogPasswordRequest{
				LoginHash:    "jsadasjdjs",
				PasswordHash: "shdahsdhahd",
				ResourceName: "yandex.ru",
			},
			want: want{
				errCode: codes.Unauthenticated,
				response: lpPB.CreateLogPasswordResponse{
					Id: "1",
				},
			},
			ctx:          context.Background(),
			omitCreateLP: true,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{}
	env.Parse(cfg)

	lp := mock.NewMockLogPassword(ctrl)

	lpApp := app.LogPasswordApp{LogPassword: lp, Cfg: cfg}
	lpService := service.LogPasswordService{lpPB.UnimplementedLogPasswordServer{}, &lpApp}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.omitCreateLP {
				lp.EXPECT().CreateLogPassword(tt.ctx, "1", tt.request.LoginHash, tt.request.PasswordHash, tt.request.ResourceName, gomock.Any()).Return(tt.want.response.Id, nil)
			}
			_, err := lpService.CreateLogPassword(tt.ctx, &tt.request)

			if unwrapError, ok := status.FromError(err); ok {
				assert.Equal(t, unwrapError.Code(), tt.want.errCode)
			} else {
				t.Error("can't parse error, unexpected test result")
			}
		})
	}
}

func TestLogPasswordService_UpdateLogPassword(t *testing.T) {
	type want struct {
		errCode  codes.Code
		response lpPB.UpdateLogPasswordResponse
	}

	tests := []struct {
		name                     string
		request                  lpPB.UpdateLogPasswordRequest
		existingLogPassword      storage.LogPassword
		existingLogPasswordError error
		want                     want
		ctx                      context.Context
		omitUpdateTN             bool
	}{
		{
			name: "update lp log password normally",
			request: lpPB.UpdateLogPasswordRequest{
				Id:           "1",
				LoginHash:    "asvdsf",
				PasswordHash: "asvdsf",
				ResourceName: "someNewName",
				ForceUpdate:  false,
				PreviousHash: "sameHash",
				IsDeleted:    false,
			},
			existingLogPassword: storage.LogPassword{
				ID:           "1",
				UID:          "1",
				LoginHash:    "asdasd",
				PasswordHash: "asdasd",
				ResourceName: "someName",
				EntryHash:    "sameHash",
			},
			want: want{
				errCode:  codes.OK,
				response: lpPB.UpdateLogPasswordResponse{},
			},
			ctx: metadata.NewIncomingContext(context.Background(), metadata.Pairs("uid", "1")),
		},
		{
			name: "update lp which doesn't exist",
			request: lpPB.UpdateLogPasswordRequest{
				Id:           "10",
				LoginHash:    "asvdsf",
				PasswordHash: "asvdsf",
				ResourceName: "someNewName",
				ForceUpdate:  false,
				PreviousHash: "sameHash",
				IsDeleted:    false,
			},
			existingLogPassword: storage.LogPassword{
				ID:           "1",
				UID:          "1",
				LoginHash:    "asdasd",
				PasswordHash: "asdasd",
				ResourceName: "someName",
				EntryHash:    "sameHash",
			},
			existingLogPasswordError: utils.ErrNotFound,
			want: want{
				errCode:  codes.NotFound,
				response: lpPB.UpdateLogPasswordResponse{},
			},
			omitUpdateTN: true,
			ctx:          metadata.NewIncomingContext(context.Background(), metadata.Pairs("uid", "1")),
		},
		{
			name: "try update lp which already updated",
			request: lpPB.UpdateLogPasswordRequest{
				Id:           "1",
				LoginHash:    "asvdsf",
				PasswordHash: "asvdsf",
				ResourceName: "someNewName",
				ForceUpdate:  false,
				PreviousHash: "differentHash",
				IsDeleted:    false,
			},
			existingLogPassword: storage.LogPassword{
				ID:           "1",
				UID:          "1",
				LoginHash:    "asdasd",
				PasswordHash: "asdasd",
				ResourceName: "someName",
				EntryHash:    "sameHash",
			},
			existingLogPasswordError: nil,
			want: want{
				errCode:  codes.Unavailable,
				response: lpPB.UpdateLogPasswordResponse{},
			},
			omitUpdateTN: true,
			ctx:          metadata.NewIncomingContext(context.Background(), metadata.Pairs("uid", "1")),
		},
		{
			name: "try update lp which already updated but force",
			request: lpPB.UpdateLogPasswordRequest{
				Id:           "1",
				LoginHash:    "asvdsf",
				PasswordHash: "asvdsf",
				ResourceName: "someNewName",
				ForceUpdate:  true,
				PreviousHash: "differentHash",
				IsDeleted:    false,
			},
			existingLogPassword: storage.LogPassword{
				ID:           "1",
				UID:          "1",
				LoginHash:    "asdasd",
				PasswordHash: "asdasd",
				ResourceName: "someName",
				EntryHash:    "sameHash",
			},
			existingLogPasswordError: nil,
			want: want{
				errCode:  codes.OK,
				response: lpPB.UpdateLogPasswordResponse{},
			},
			ctx: metadata.NewIncomingContext(context.Background(), metadata.Pairs("uid", "1")),
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{}
	env.Parse(cfg)

	lp := mock.NewMockLogPassword(ctrl)

	lpApp := app.LogPasswordApp{LogPassword: lp, Cfg: cfg}
	tnService := service.LogPasswordService{lpPB.UnimplementedLogPasswordServer{}, &lpApp}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lp.EXPECT().GetLogPasswordByID(tt.ctx, tt.request.Id).Return(&tt.existingLogPassword, tt.existingLogPasswordError)
			if !tt.omitUpdateTN {
				lp.EXPECT().UpdateLogPassword(tt.ctx, tt.request.Id, tt.request.LoginHash, tt.request.PasswordHash, tt.request.ResourceName, gomock.Any(), tt.request.IsDeleted)
			}
			_, err := tnService.UpdateLogPassword(tt.ctx, &tt.request)

			if unwrapError, ok := status.FromError(err); ok {
				assert.Equal(t, unwrapError.Code(), tt.want.errCode)
			} else {
				t.Error("can't parse error, unexpected test result")
			}
		})
	}
}

func TestLogPasswordService_ListLogPassword(t *testing.T) {
	type want struct {
		errCode  codes.Code
		response lpPB.ListLogPasswordResponse
	}

	tests := []struct {
		name                   string
		request                lpPB.ListLogPasswordRequest
		existingLogPasswords   []storage.LogPassword
		existingLogPasswordErr error
		want                   want
		ctx                    context.Context
		uid                    string
	}{
		{
			name:                 "list log password normally",
			request:              lpPB.ListLogPasswordRequest{},
			existingLogPasswords: []storage.LogPassword{{EntryHash: "1", UID: "1"}, {EntryHash: "2", UID: "1"}, {EntryHash: "4", UID: "1"}},
			want: want{
				errCode: codes.OK,
				response: lpPB.ListLogPasswordResponse{
					LogPasswords: []*lpPB.LogPasswordEntry{{EntryHash: "1"}, {EntryHash: "2"}, {EntryHash: "3"}},
				},
			},
			ctx: metadata.NewIncomingContext(context.Background(), metadata.Pairs("uid", "1")),
			uid: "1",
		},
		{
			name: "list log password but remove some existing log passwords",
			request: lpPB.ListLogPasswordRequest{
				ExistingHashes: []*lpPB.ExistingLogPasswordHash{{Id: "1", EntryHash: "1"}, {Id: "3", EntryHash: "3"}},
			},
			existingLogPasswords: []storage.LogPassword{{ID: "1", EntryHash: "1", UID: "1"}, {ID: "2", EntryHash: "2"}, {ID: "3", EntryHash: "3", UID: "1"}, {ID: "1", EntryHash: "4"}},
			want: want{
				errCode: codes.OK,
				response: lpPB.ListLogPasswordResponse{
					LogPasswords: []*lpPB.LogPasswordEntry{{EntryHash: "1"}, {EntryHash: "3"}},
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

	tn := mock.NewMockLogPassword(ctrl)

	tnApp := app.LogPasswordApp{LogPassword: tn, Cfg: cfg}
	tnService := service.LogPasswordService{lpPB.UnimplementedLogPasswordServer{}, &tnApp}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tn.EXPECT().ListLogPasswordByUID(tt.ctx, tt.uid).Return(tt.existingLogPasswords, nil)
			resp, err := tnService.ListLogPassword(tt.ctx, &tt.request)

			if unwrapError, ok := status.FromError(err); ok {
				assert.Equal(t, unwrapError.Code(), tt.want.errCode)
				assert.Equal(t, len(resp.LogPasswords), len(tt.want.response.LogPasswords))

			} else {
				t.Error("can't parse error, unexpected test result")
			}
		})
	}
}
