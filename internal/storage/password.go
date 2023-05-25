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

type LogPasswordStorage struct {
	Conn *pgxpool.Pool
}

type LogPassword struct {
	ID           string
	UID          string
	LoginHash    string
	PasswordHash string
	ResourceName string
	EntryHash    string
	IsDeleted    bool
}

func InitLogPasswordStorage(cfg *config.Config) (*LogPasswordStorage, error) {
	conn, err := pgxpool.New(context.Background(), cfg.DatabaseURI)

	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err.Error())
		return nil, err
	}

	return &LogPasswordStorage{conn}, nil
}

func (lp *LogPasswordStorage) CreateLogPassword(ctx context.Context, uid, loginHash, passwordHash, resourceName, entryHash string) (string, error) {
	sqlStatement := `
	INSERT INTO log_passwords (uid, login_hash, password_hash, resource_name, entry_hash)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id;`

	var id string
	err := lp.Conn.QueryRow(ctx, sqlStatement, uid, loginHash, passwordHash, resourceName, entryHash).Scan(&id)

	if err != nil {
		return id, utils.WrapError(err, utils.ErrDBLayer)
	}

	return id, nil
}

func (lp *LogPasswordStorage) UpdateLogPassword(ctx context.Context, id, loginHash, passwordHash, resourceName, entryHash string, isDeleted bool) error {
	updateBalanceSQL := `
	UPDATE LOG_PASSWORDS SET 
	login_hash = $2,
	password_hash = $3,
	resource_name = $4,
	entry_hash = $5,
	is_deleted = $6

	WHERE id = $1
	`

	_, err := lp.Conn.Exec(ctx, updateBalanceSQL, id, loginHash, passwordHash, resourceName, entryHash, isDeleted)

	if err != nil {
		return utils.WrapError(err, utils.ErrDBLayer)
	}

	return nil
}

func (lp *LogPasswordStorage) ListLogPasswordByUID(ctx context.Context, uid string) ([]LogPassword, error) {
	sqlStatement := `
	SELECT id, uid, login_hash, password_hash, resource_name, entry_hash, is_deleted FROM log_passwords WHERE UID = $1 
	`

	rows, err := lp.Conn.Query(ctx, sqlStatement, uid)
	if err != nil {
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, utils.ErrNotFound
			}

			return nil, utils.WrapError(err, utils.ErrDBLayer)
		}
	}

	defer rows.Close()

	logPasses := []LogPassword{}

	for rows.Next() {
		entry := LogPassword{}
		err = rows.Scan(&entry.ID, &entry.UID, &entry.LoginHash, &entry.PasswordHash, &entry.ResourceName, &entry.EntryHash, &entry.IsDeleted)

		if err != nil {
			return nil, err
		}

		logPasses = append(logPasses, entry)
	}

	err = rows.Err()
	if err != nil {
		return nil, utils.WrapError(err, utils.ErrDBLayer)
	}

	return logPasses, nil
}

func (lp *LogPasswordStorage) GetLogPasswordByID(ctx context.Context, id string) (*LogPassword, error) {
	sqlStatement := `
	SELECT id, uid, login_hash, password_hash, resource_name, entry_hash, is_deleted FROM log_passwords WHERE ID = $1 
	`
	logPass := LogPassword{}

	err := lp.Conn.QueryRow(ctx, sqlStatement, id).Scan(&logPass.ID, &logPass.UID, &logPass.LoginHash, &logPass.PasswordHash, &logPass.ResourceName, &logPass.EntryHash, &logPass.IsDeleted)
	if err != nil {
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, utils.ErrNotFound
			}

			return nil, utils.WrapError(err, utils.ErrDBLayer)
		}
	}

	if err != nil {
		return nil, utils.WrapError(err, utils.ErrDBLayer)
	}

	return &logPass, nil
}

func (lp *LogPasswordStorage) Close() {
	lp.Conn.Close()
}
