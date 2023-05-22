package service

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/T-V-N/gophkeeper/internal/app"
	"github.com/T-V-N/gophkeeper/internal/config"
	textNotePB "github.com/T-V-N/gophkeeper/internal/grpc/textNote"
	"github.com/T-V-N/gophkeeper/internal/service"
	"github.com/T-V-N/gophkeeper/internal/utils"
)

type TextNoteService struct {
	textNotePB.UnimplementedTextNoteServer
	TextNoteApp *app.TextNoteApp
}

func (tns *TextNoteService) CreateTextNote(ctx context.Context, in *textNotePB.CreateTextNoteRequest) (*textNotePB.CreateTextNoteResponse, error) {
	response := textNotePB.CreateTextNoteResponse{}

	uid, err := service.ExtractUIDFromCtx(ctx)

	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unathorized")
	}

	id, err := tns.TextNoteApp.CreateTextNote(ctx, uid, in.NoteTextHash, in.NoteName)

	if err != nil {
		switch errors.Unwrap(err) {
		case utils.ErrDBLayer:
			return nil, status.Error(codes.Internal, err.(utils.WrappedAPIError).Message())
		default:
			return nil, status.Error(codes.Internal, err.(utils.WrappedAPIError).Message())
		}
	}

	response.Id = id

	return &response, nil
}

func (tns *TextNoteService) UpdateTextNote(ctx context.Context, in *textNotePB.UpdateTextNoteRequest) (*textNotePB.UpdateTextNoteResponse, error) {
	response := textNotePB.UpdateTextNoteResponse{}

	uid, err := service.ExtractUIDFromCtx(ctx)

	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unathorized")
	}

	err = tns.TextNoteApp.UpdateTextNote(ctx, uid, in.Id, in.NoteTextHash, in.NoteName, in.PreviousHash, in.IsDeleted, in.ForceUpdate)

	if err != nil {
		switch errors.Unwrap(err) {
		case utils.ErrNoData:
			return nil, status.Error(codes.NotFound, err.(utils.WrappedAPIError).Message())
		case utils.ErrConflict:
			return nil, status.Error(codes.Unavailable, "conflict")
		case utils.ErrNotAuthorized:
			return nil, status.Error(codes.Unauthenticated, "unathorized")
		case utils.ErrDBLayer:
			return nil, status.Error(codes.Internal, err.(utils.WrappedAPIError).Message())
		default:
			return nil, status.Error(codes.Internal, err.(utils.WrappedAPIError).Message())
		}
	}

	return &response, nil
}

func (tns *TextNoteService) ListTextNote(ctx context.Context, in *textNotePB.ListTextNoteRequest) (*textNotePB.ListTextNoteResponse, error) {
	response := textNotePB.ListTextNoteResponse{}

	uid, err := service.ExtractUIDFromCtx(ctx)

	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unathorized")
	}

	existingHashes := []app.ExistingHash{}
	for _, hash := range in.ExistingHashes {
		existingHashes = append(existingHashes, app.ExistingHash{ID: hash.Id, EntryHash: hash.EntryHash})
	}

	textNotes, err := tns.TextNoteApp.ListTextNote(ctx, uid, existingHashes)

	if err != nil {
		switch errors.Unwrap(err) {
		case utils.ErrNoData:
			return nil, status.Error(codes.NotFound, err.(utils.WrappedAPIError).Message())
		case utils.ErrDBLayer:
			return nil, status.Error(codes.Internal, err.(utils.WrappedAPIError).Message())
		default:
			return nil, status.Error(codes.Internal, err.(utils.WrappedAPIError).Message())
		}
	}

	textNotePBResponse := []*textNotePB.TextNoteEntry{}
	for _, textNote := range textNotes {
		textNotePBResponse = append(textNotePBResponse, &textNotePB.TextNoteEntry{
			Id:           textNote.ID,
			NoteTextHash: textNote.NoteTextHash,
			NoteName:     textNote.NoteName,
			EntryHash:    textNote.EntryHash,
		})
	}

	response.Notes = textNotePBResponse

	return &response, nil
}

func InittextNoteService(cfg *config.Config, a *app.TextNoteApp) *TextNoteService {
	return &TextNoteService{textNotePB.UnimplementedTextNoteServer{}, a}
}
