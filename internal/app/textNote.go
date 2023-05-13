package app

import (
	"context"

	"github.com/T-V-N/gophkeeper/internal/config"
	"github.com/T-V-N/gophkeeper/internal/storage"
	"github.com/T-V-N/gophkeeper/internal/utils"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type TextNoteApp struct {
	TextNote *storage.TextNotesStorage
	Cfg      *config.Config
	logger   *zap.SugaredLogger
}

func InitTextNoteApp(conn *pgxpool.Pool, cfg *config.Config, logger *zap.SugaredLogger) (*TextNoteApp, error) {
	tn, err := storage.InitTextNote(conn)

	if err != nil {
		return nil, err
	}

	return &TextNoteApp{tn, cfg, logger}, nil
}

func (t *TextNoteApp) CreateTextNote(ctx context.Context, uid, noteTextHash, noteName string) (string, error) {
	entryHash := utils.PackedCheckSum([]string{noteTextHash})

	id, err := t.TextNote.CreateTextNote(ctx, uid, noteName, noteTextHash, entryHash)

	if err != nil {
		return "", utils.WrapError(err, utils.ErrDBLayer)
	}

	return id, nil
}

func (t *TextNoteApp) UpdateTextNote(ctx context.Context, id, noteTextHash, noteName, previousHash string, isDeleted, forceUpdate bool) error {
	logPass, err := t.TextNote.GetTextNoteByID(ctx, id)

	if err != nil {
		return utils.WrapError(err, utils.ErrDBLayer)
	}

	if logPass == nil {
		return utils.ErrNotFound
	}

	if (previousHash != logPass.EntryHash) && !forceUpdate {
		return utils.ErrConflict
	}

	entryHash := utils.PackedCheckSum([]string{noteTextHash})

	return t.TextNote.UpdateTextNote(ctx, id, noteName, noteTextHash, entryHash, isDeleted)
}
