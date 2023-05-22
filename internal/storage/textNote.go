package storage

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/T-V-N/gophkeeper/internal/config"
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
	INSERT INTO CARDS (uid, note_name, note_text_hash, entry_hash)
	VALUES ($1, $2, $3, $4)
	RETURNING id;`

	var id string
	err := t.Conn.QueryRow(ctx, sqlStatement, uid, noteName, noteTextHash, entryHash).Scan(&id)

	if err != nil {
		return id, err
	}

	return id, nil
}

func (t *TextNotesStorage) UpdateTextNote(ctx context.Context, id, noteName, noteTextHash, entryHash string, isDeleted bool) error {
	updateBalanceSQL := `
	UPDATE USERS SET 
	card_number_hash = $2,
	valid_until_hash = $3,
	cvv_hash = $4,
	last_four_digits = $5,
	entry_hash = $6,
	is_deleted = $7
	WHERE id = $1
	`

	_, err := t.Conn.Exec(ctx, updateBalanceSQL, id, id, noteName, noteTextHash, entryHash, isDeleted)

	if err != nil {
		return err
	}

	return nil
}

func (t *TextNotesStorage) ListTextNoteByUID(ctx context.Context, uid string) ([]TextNote, error) {
	sqlStatement := `
	SELECT id, uid, note_name, note_text_hash, entry_hash, is_deleted FROM text_notes WHERE UID = $1 
	`

	rows, err := t.Conn.Query(ctx, sqlStatement, uid)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	return TextNotes, nil
}

func (t *TextNotesStorage) GetTextNoteByID(ctx context.Context, id string) (*TextNote, error) {
	sqlStatement := `
	SELECT id, uid, note_name, note_text_hash, entry_hash, is_deleted FROM text_notes WHERE ID = $1
	`

	row, err := t.Conn.Query(ctx, sqlStatement, id)
	if err != nil {
		return nil, err
	}

	defer row.Close()

	note := TextNote{}

	err = row.Scan(&note.ID, &note.UID, &note.NoteName, &note.NoteTextHash, &note.EntryHash, &note.IsDeleted)

	if err != nil {
		return nil, err
	}

	return &note, nil
}

func (t *TextNotesStorage) Close() {
	t.Conn.Close()
}
