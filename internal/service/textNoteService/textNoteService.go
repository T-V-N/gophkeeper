package service

import (
	"context"

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
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	id, err := tns.TextNoteApp.CreateTextNote(ctx, uid, in.NoteTextHash, in.NoteName)

	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}

	response.Id = id

	return &response, nil
}

func (tns *TextNoteService) UpdateTextNote(ctx context.Context, in *textNotePB.UpdateTextNoteRequest) (*textNotePB.UpdateTextNoteResponse, error) {
	response := textNotePB.UpdateTextNoteResponse{}

	uid, err := service.ExtractUIDFromCtx(ctx)

	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	err = tns.TextNoteApp.UpdateTextNote(ctx, uid, in.Id, in.NoteTextHash, in.NoteName, in.PreviousHash, in.IsDeleted, in.ForceUpdate)

	if err != nil {
		switch err {
		case utils.ErrNotFound:
			return nil, status.Error(codes.NotFound, "note not found")
		case utils.ErrConflict:
			return nil, status.Error(codes.Unavailable, "cannot update text note which is already updated, sync first")
		case utils.ErrNotAuthorized:
			return nil, status.Error(codes.Unauthenticated, "unauthorized")
		default:
			return nil, status.Error(codes.Internal, "internal server error ;(")
		}
	}

	return &response, nil
}

func (tns *TextNoteService) ListTextNote(ctx context.Context, in *textNotePB.ListTextNoteRequest) (*textNotePB.ListTextNoteResponse, error) {
	response := textNotePB.ListTextNoteResponse{}

	uid, err := service.ExtractUIDFromCtx(ctx)

	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	existingHashes := []app.ExistingHash{}
	for _, hash := range in.ExistingHashes {
		existingHashes = append(existingHashes, app.ExistingHash{ID: hash.Id, EntryHash: hash.EntryHash})
	}

	textNotes, err := tns.TextNoteApp.ListTextNote(ctx, uid, existingHashes)

	if err != nil {
		switch err {
		case utils.ErrNotFound:
			return nil, status.Error(codes.NotFound, "no notes available")
		default:
			return nil, status.Error(codes.Internal, "internal server error ;(")
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
