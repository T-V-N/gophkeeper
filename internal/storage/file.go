package storage

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type FileStorage struct {
	Conn *pgxpool.Pool
}

type File struct {
	ID          string
	UID         string
	FileName    string
	S3Link      string
	CommittedAt time.Time
	IsDeleted   bool
}

func InitFile(conn *pgxpool.Pool) (*FileStorage, error) {
	return &FileStorage{conn}, nil
}

func (f *FileStorage) CreateFile(ctx context.Context, uid, fileName string) (string, error) {
	sqlStatement := `
	INSERT INTO CARDS (uid, file_name)
	VALUES ($1, $2)
	RETURNING id;`

	var id string
	err := f.Conn.QueryRow(ctx, sqlStatement, uid, fileName).Scan(&id)

	if err != nil {
		return id, err
	}

	return id, nil
}

func (f *FileStorage) UpdateFile(ctx context.Context, id, fileName, S3Link string, isDeleted bool, CommitedAt time.Time) error {
	updateFileSQL := `
	UPDATE FILE SET 
	file_name = $2,
	s3_link = $3,
	committed_at = $4,
	is_deleted = $5,
	WHERE id = $1
	`

	_, err := f.Conn.Exec(ctx, updateFileSQL, id, fileName, S3Link, CommitedAt, isDeleted)

	if err != nil {
		return err
	}

	return nil
}

func (f *FileStorage) GetFileByID(ctx context.Context, id string) (*File, error) {
	sqlStatement := `
	SELECT id, uid, file_name, s3_link, committed_at, is_deleted FROM files WHERE ID = $1 
	`

	row, err := f.Conn.Query(ctx, sqlStatement, id)
	if err != nil {
		return nil, err
	}

	defer row.Close()

	file := File{}

	err = row.Scan(&file.ID, &file.UID, &file.FileName, &file.S3Link, &file.CommittedAt, &file.IsDeleted)

	if err != nil {
		return nil, err
	}

	return &file, nil
}

func (f *FileStorage) ListFilesByUID(ctx context.Context, uid string) (*[]File, error) {
	sqlStatement := `
	SELECT id, uid, file_name, s3_link, committed_at, is_deleted FROM files WHERE UID = $1 
	`

	rows, err := f.Conn.Query(ctx, sqlStatement, uid)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	Files := []File{}

	for rows.Next() {
		entry := File{}
		err = rows.Scan(&entry.ID, &entry.UID, &entry.FileName, &entry.S3Link, &entry.CommittedAt, &entry.IsDeleted)

		if err != nil {
			return nil, err
		}

		Files = append(Files, entry)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return &Files, nil
}
