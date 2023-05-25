package app

import (
	"context"

	"go.uber.org/zap"

	"github.com/T-V-N/gophkeeper/internal/config"
	"github.com/T-V-N/gophkeeper/internal/storage"
	"github.com/T-V-N/gophkeeper/internal/utils"
)

type TextNote interface {
	CreateTextNote(ctx context.Context, uid, noteName, noteTextHash, entryHash string) (string, error)
	UpdateTextNote(ctx context.Context, id, noteName, noteTextHash, entryHash string, isDeleted bool) error
	ListTextNoteByUID(ctx context.Context, uid string) ([]storage.TextNote, error)
	GetTextNoteByID(ctx context.Context, id string) (*storage.TextNote, error)
	Close()
}

type TextNoteApp struct {
	TextNote TextNote
	Cfg      *config.Config
	logger   *zap.SugaredLogger
}

func InitTextNoteApp(cfg *config.Config, logger *zap.SugaredLogger) (*TextNoteApp, error) {
	tn, err := storage.InitTextNoteStorage(cfg)

	if err != nil {
		return nil, err
	}

	return &TextNoteApp{tn, cfg, logger}, nil
}

func (t *TextNoteApp) CreateTextNote(ctx context.Context, uid, noteTextHash, noteName string) (string, error) {
	entryHash := utils.PackedCheckSum([]string{noteTextHash})

	id, err := t.TextNote.CreateTextNote(ctx, uid, noteName, noteTextHash, entryHash)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (t *TextNoteApp) UpdateTextNote(ctx context.Context, uid, id, noteTextHash, noteName, previousHash string, isDeleted, forceUpdate bool) error {
	textNote, err := t.TextNote.GetTextNoteByID(ctx, id)

	if err != nil {
		return err
	}

	if textNote.UID != uid {
		return utils.ErrNotAuthorized
	}

	if (previousHash != textNote.EntryHash) && !forceUpdate {
		return utils.ErrConflict
	}

	entryHash := utils.PackedCheckSum([]string{noteTextHash})

	return t.TextNote.UpdateTextNote(ctx, id, noteName, noteTextHash, entryHash, isDeleted)
}

func (t *TextNoteApp) ListTextNote(ctx context.Context, uid string, existingHashes []ExistingHash) ([]storage.TextNote, error) {
	notes, err := t.TextNote.ListTextNoteByUID(ctx, uid)

	if err != nil {
		return nil, err
	}

	var result []storage.TextNote

	for _, note := range notes {
		include := true

		for _, existingHash := range existingHashes {
			if note.EntryHash == existingHash.EntryHash {
				include = false
				break
			}
		}

		if include {
			result = append(result, note)
		}
	}

	return result, nil
}
