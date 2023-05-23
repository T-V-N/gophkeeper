package storage

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/T-V-N/gophkeeper/internal/config"
	"github.com/T-V-N/gophkeeper/internal/utils"
)

type UserStorage struct {
	Conn *pgxpool.Pool
}

type User struct {
	UID              string
	Email            string
	PasswordHash     string
	VerificationCode string
	TOTPSecret       string
	ConfirmedAt      time.Time
	CreatedAt        string
	TOTPEnabled      bool
}

func InitUserStorage(cfg *config.Config) (*UserStorage, error) {
	conn, err := pgxpool.New(context.Background(), cfg.DatabaseURI)

	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err.Error())
		return nil, err
	}

	return &UserStorage{conn}, nil
}

func (user *UserStorage) CreateUser(ctx context.Context, email, passwordHash, verificationCode string) (string, error) {
	sqlStatement := `
	INSERT INTO users (email, password_hash, verification_code)
	VALUES ($1, $2, $3)
	RETURNING uid;`

	var uid string
	err := user.Conn.QueryRow(ctx, sqlStatement, email, passwordHash, verificationCode).Scan(&uid)

	var pgErr *pgconn.PgError
	if err != nil {
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return "", utils.ErrDuplicate
		}

		return "", utils.WrapError(err, utils.ErrDBLayer)
	}

	return uid, nil
}

func (user *UserStorage) UpdateUser(ctx context.Context, uid, email, passwordHash, TOTPSecret string, TOTPEnabled bool, confirmedAt time.Time) error {
	updateBalanceSQL := `
	UPDATE USERS 
	SET 
	email = $2,
	password_hash = $3,
	totp_secret = $4,
	totp_enabled = $5,
	confirmed_at = $6
	
	WHERE uid = $1
	`

	_, err := user.Conn.Exec(ctx, updateBalanceSQL, uid, email, passwordHash, TOTPSecret, TOTPEnabled, confirmedAt)

	if err != nil {
		return utils.WrapError(err, utils.ErrDBLayer)
	}

	return nil
}

func (user *UserStorage) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	sqlStatement := `
	SELECT uid, email, password_hash, totp_secret, totp_enabled, confirmed_at, verification_code FROM USERS
	WHERE email = $1
	`

	var u User

	var totpSecret sql.NullString
	var confirmedAt sql.NullTime

	err := user.Conn.QueryRow(ctx, sqlStatement, email).Scan(&u.UID, &u.Email, &u.PasswordHash, &totpSecret, &u.TOTPEnabled, &confirmedAt, &u.VerificationCode)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrNotFound
		}

		return nil, utils.WrapError(err, utils.ErrDBLayer)
	}

	if totpSecret.Valid {
		u.TOTPSecret = totpSecret.String
	}

	if confirmedAt.Valid {
		u.ConfirmedAt = confirmedAt.Time
	}

	return &u, nil
}

func (user *UserStorage) GetUserByID(ctx context.Context, uid string) (*User, error) {
	sqlStatement := `
	SELECT uid, email, password_hash, totp_secret, totp_enabled, confirmed_at, verification_code FROM USERS
	WHERE uid = $1
	`

	var u User
	var totpSecret sql.NullString
	var confirmedAt sql.NullTime
	err := user.Conn.QueryRow(ctx, sqlStatement, uid).Scan(&u.UID, &u.Email, &u.PasswordHash, &totpSecret, &u.TOTPEnabled, &confirmedAt, &u.VerificationCode)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrNotFound
		}

		return nil, utils.WrapError(err, utils.ErrDBLayer)
	}

	if totpSecret.Valid {
		u.TOTPSecret = totpSecret.String
	}

	if confirmedAt.Valid {
		u.ConfirmedAt = confirmedAt.Time
	}

	return &u, nil
}

func (user *UserStorage) Close() {
	user.Conn.Close()
}
