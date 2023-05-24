package storage

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/T-V-N/gophkeeper/internal/config"
	"github.com/T-V-N/gophkeeper/internal/utils"
)

type TextNotesStorage struct {
	Conn *pgxpool.Pool
}

type TextNote struct {
	ID           string
	UID          string
	NoteName     string
	NoteTextHash string
	EntryHash    string
	IsDeleted    bool
}

func InitTextNoteStorage(cfg *config.Config) (*TextNotesStorage, error) {
	conn, err := pgxpool.New(context.Background(), cfg.DatabaseURI)

	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err.Error())
		return nil, err
	}

	return &TextNotesStorage{conn}, nil
}

func (t *TextNotesStorage) CreateTextNote(ctx context.Context, uid, noteName, noteTextHash, entryHash string) (string, error) {
	sqlStatement := `
	INSERT INTO text_notes (uid, note_name, note_text_hash, entry_hash)
	VALUES ($1, $2, $3, $4)
	RETURNING id;`

	var id string
	err := t.Conn.QueryRow(ctx, sqlStatement, uid, noteName, noteTextHash, entryHash).Scan(&id)

	if err != nil {
		return id, utils.WrapError(err, utils.ErrDBLayer)
	}

	return id, nil
}

func (t *TextNotesStorage) UpdateTextNote(ctx context.Context, id, noteName, noteTextHash, entryHash string, isDeleted bool) error {
	updateBalanceSQL := `
	UPDATE text_notes SET 
	note_name = $2,
	note_text_hash = $3,
	entry_hash = $4,
	is_deleted = $5
	WHERE id = $1
	`

	_, err := t.Conn.Exec(ctx, updateBalanceSQL, id, noteName, noteTextHash, entryHash, isDeleted)

	if err != nil {
		return utils.WrapError(err, utils.ErrDBLayer)
	}

	return nil
}

func (t *TextNotesStorage) ListTextNoteByUID(ctx context.Context, uid string) ([]TextNote, error) {
	sqlStatement := `
	SELECT id, uid, note_name, note_text_hash, entry_hash, is_deleted FROM text_notes WHERE UID = $1 
	`

	rows, err := t.Conn.Query(ctx, sqlStatement, uid)
	if err != nil {
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, utils.ErrNotFound
			}

			return nil, utils.WrapError(err, utils.ErrDBLayer)
		}
	}

	defer rows.Close()

	TextNotes := []TextNote{}

	for rows.Next() {
		entry := TextNote{}
		err = rows.Scan(&entry.ID, &entry.UID, &entry.NoteTextHash, &entry.EntryHash, &entry.IsDeleted)

		if err != nil {
			return nil, err
		}

		TextNotes = append(TextNotes, entry)
	}

	err = rows.Err()
	if err != nil {
		return nil, utils.WrapError(err, utils.ErrDBLayer)
	}

	return TextNotes, nil
}

func (t *TextNotesStorage) GetTextNoteByID(ctx context.Context, id string) (*TextNote, error) {
	sqlStatement := `
	SELECT id, uid, note_name, note_text_hash, entry_hash, is_deleted FROM text_notes WHERE ID = $1
	`

	row, err := t.Conn.Query(ctx, sqlStatement, id)
	if err != nil {
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, utils.ErrNotFound
			}
			return nil, utils.WrapError(err, utils.ErrDBLayer)
		}
	}

	defer row.Close()

	note := TextNote{}

	err = row.Scan(&note.ID, &note.UID, &note.NoteName, &note.NoteTextHash, &note.EntryHash, &note.IsDeleted)

	if err != nil {
		return nil, utils.WrapError(err, utils.ErrDBLayer)
	}

	return &note, nil
}

func (t *TextNotesStorage) Close() {
	t.Conn.Close()
}
