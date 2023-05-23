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
	tPB "github.com/T-V-N/gophkeeper/internal/grpc/textNote"
	mock "github.com/T-V-N/gophkeeper/internal/mocks"
	service "github.com/T-V-N/gophkeeper/internal/service/textNoteService"
	"github.com/T-V-N/gophkeeper/internal/storage"
	"github.com/T-V-N/gophkeeper/internal/utils"
)

func TestTextNoteService_CreateTextNote(t *testing.T) {
	type want struct {
		errCode  codes.Code
		response tPB.CreateTextNoteResponse
	}

	tests := []struct {
		name         string
		request      tPB.CreateTextNoteRequest
		want         want
		ctx          context.Context
		omitCreateTN bool
	}{
		{
			name: "create text note normally",
			request: tPB.CreateTextNoteRequest{
				NoteTextHash: "sdjasdsadasdasdasdasd",
				NoteName:     "some name",
			},
			want: want{
				errCode: codes.OK,
				response: tPB.CreateTextNoteResponse{
					Id: "1",
				},
			},
			ctx: metadata.NewIncomingContext(context.Background(), metadata.Pairs("uid", "1")),
		},
		{
			name: "create card unauthorized",
			request: tPB.CreateTextNoteRequest{
				NoteTextHash: "sdjasdsadasdasdasdasd",
				NoteName:     "some name",
			},
			want: want{
				errCode: codes.Unauthenticated,
				response: tPB.CreateTextNoteResponse{
					Id: "1",
				},
			},
			ctx:          context.Background(),
			omitCreateTN: true,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{}
	env.Parse(cfg)

	tn := mock.NewMockTextNote(ctrl)

	tnApp := app.TextNoteApp{TextNote: tn, Cfg: cfg}
	tnService := service.TextNoteService{tPB.UnimplementedTextNoteServer{}, &tnApp}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.omitCreateTN {
				tn.EXPECT().CreateTextNote(tt.ctx, "1", tt.request.NoteName, tt.request.NoteTextHash, gomock.Any())
			}
			_, err := tnService.CreateTextNote(tt.ctx, &tt.request)

			if unwrapError, ok := status.FromError(err); ok {
				assert.Equal(t, unwrapError.Code(), tt.want.errCode)
			} else {
				t.Error("can't parse error, unexpected test result")
			}
		})
	}
}

func TestTextNoteService_UpdateTextNote(t *testing.T) {
	type want struct {
		errCode  codes.Code
		response tPB.UpdateTextNoteResponse
	}

	tests := []struct {
		name                string
		request             tPB.UpdateTextNoteRequest
		existingTextNote    storage.TextNote
		existingTextNoteErr error
		want                want
		ctx                 context.Context
		omitUpdateTN        bool
	}{
		{
			name: "update text note normally",
			request: tPB.UpdateTextNoteRequest{
				Id:           "1",
				NoteTextHash: "asvdsf",
				NoteName:     "someNewName",
				ForceUpdate:  false,
				PreviousHash: "sameHash",
				IsDeleted:    false,
			},
			existingTextNote: storage.TextNote{
				ID:           "1",
				UID:          "1",
				NoteName:     "someName",
				NoteTextHash: "asdasd",
				EntryHash:    "sameHash",
			},
			want: want{
				errCode:  codes.OK,
				response: tPB.UpdateTextNoteResponse{},
			},
			ctx: metadata.NewIncomingContext(context.Background(), metadata.Pairs("uid", "1")),
		},
		{
			name: "update text note which doesn't exist",
			request: tPB.UpdateTextNoteRequest{
				Id:           "1",
				NoteTextHash: "asvdsf",
				NoteName:     "someNewName",
				ForceUpdate:  false,
				PreviousHash: "sameHash",
				IsDeleted:    false,
			},
			existingTextNote: storage.TextNote{
				ID:           "1",
				UID:          "1",
				NoteName:     "someName",
				NoteTextHash: "asdasd",
				EntryHash:    "sameHash",
			},
			existingTextNoteErr: utils.ErrNotFound,
			want: want{
				errCode:  codes.NotFound,
				response: tPB.UpdateTextNoteResponse{},
			},
			omitUpdateTN: true,
			ctx:          metadata.NewIncomingContext(context.Background(), metadata.Pairs("uid", "1")),
		},
		{
			name: "try update text note which already updated",
			request: tPB.UpdateTextNoteRequest{
				Id:           "1",
				NoteTextHash: "asvdsf",
				NoteName:     "someNewName",
				ForceUpdate:  false,
				PreviousHash: "differentHash",
				IsDeleted:    false,
			},
			existingTextNote: storage.TextNote{
				ID:           "1",
				UID:          "1",
				NoteName:     "someName",
				NoteTextHash: "asdasd",
				EntryHash:    "sameHash",
			},
			existingTextNoteErr: nil,
			want: want{
				errCode:  codes.Unavailable,
				response: tPB.UpdateTextNoteResponse{},
			},
			omitUpdateTN: true,
			ctx:          metadata.NewIncomingContext(context.Background(), metadata.Pairs("uid", "1")),
		},
		{
			name: "try update text note which already updated but force",
			request: tPB.UpdateTextNoteRequest{
				Id:           "1",
				NoteTextHash: "asvdsf",
				NoteName:     "someNewName",
				ForceUpdate:  true,
				PreviousHash: "differentHash",
				IsDeleted:    false,
			},
			existingTextNote: storage.TextNote{
				ID:           "1",
				UID:          "1",
				NoteName:     "someName",
				NoteTextHash: "asdasd",
				EntryHash:    "sameHash",
			},
			existingTextNoteErr: nil,
			want: want{
				errCode:  codes.OK,
				response: tPB.UpdateTextNoteResponse{},
			},
			ctx: metadata.NewIncomingContext(context.Background(), metadata.Pairs("uid", "1")),
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{}
	env.Parse(cfg)

	tn := mock.NewMockTextNote(ctrl)

	tnApp := app.TextNoteApp{TextNote: tn, Cfg: cfg}
	tnService := service.TextNoteService{tPB.UnimplementedTextNoteServer{}, &tnApp}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tn.EXPECT().GetTextNoteByID(tt.ctx, tt.request.Id).Return(&tt.existingTextNote, tt.existingTextNoteErr)
			if !tt.omitUpdateTN {
				tn.EXPECT().UpdateTextNote(tt.ctx, tt.request.Id, tt.request.NoteName, tt.request.NoteTextHash, gomock.Any(), tt.request.IsDeleted)
			}
			_, err := tnService.UpdateTextNote(tt.ctx, &tt.request)

			if unwrapError, ok := status.FromError(err); ok {
				assert.Equal(t, unwrapError.Code(), tt.want.errCode)
			} else {
				t.Error("can't parse error, unexpected test result")
			}
		})
	}
}

func TestTextNoteService_ListTextNote(t *testing.T) {
	type want struct {
		errCode  codes.Code
		response tPB.ListTextNoteResponse
	}

	tests := []struct {
		name                string
		request             tPB.ListTextNoteRequest
		existingTextNotes   []storage.TextNote
		existingTextNoteErr error
		want                want
		ctx                 context.Context
		uid                 string
	}{
		{
			name:              "list text note normally",
			request:           tPB.ListTextNoteRequest{},
			existingTextNotes: []storage.TextNote{{EntryHash: "1", UID: "1"}, {EntryHash: "2", UID: "1"}, {EntryHash: "4", UID: "1"}},
			want: want{
				errCode: codes.OK,
				response: tPB.ListTextNoteResponse{
					Notes: []*tPB.TextNoteEntry{{EntryHash: "1"}, {EntryHash: "2"}, {EntryHash: "3"}},
				},
			},
			ctx: metadata.NewIncomingContext(context.Background(), metadata.Pairs("uid", "1")),
			uid: "1",
		},
		{
			name: "list text note but remove some existing text notes",
			request: tPB.ListTextNoteRequest{
				ExistingHashes: []*tPB.ExistingTextNoteHash{{Id: "1", EntryHash: "1"}, {Id: "3", EntryHash: "3"}},
			},
			existingTextNotes: []storage.TextNote{{ID: "1", EntryHash: "1", UID: "1"}, {ID: "2", EntryHash: "2"}, {ID: "3", EntryHash: "3", UID: "1"}, {ID: "1", EntryHash: "4"}},
			want: want{
				errCode: codes.OK,
				response: tPB.ListTextNoteResponse{
					Notes: []*tPB.TextNoteEntry{{EntryHash: "1"}, {EntryHash: "3"}},
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

	tn := mock.NewMockTextNote(ctrl)

	tnApp := app.TextNoteApp{TextNote: tn, Cfg: cfg}
	tnService := service.TextNoteService{tPB.UnimplementedTextNoteServer{}, &tnApp}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tn.EXPECT().ListTextNoteByUID(tt.ctx, tt.uid).Return(tt.existingTextNotes, nil)
			resp, err := tnService.ListTextNote(tt.ctx, &tt.request)

			if unwrapError, ok := status.FromError(err); ok {
				assert.Equal(t, unwrapError.Code(), tt.want.errCode)
				assert.Equal(t, len(resp.Notes), len(tt.want.response.Notes))

			} else {
				t.Error("can't parse error, unexpected test result")
			}
		})
	}
}
