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

type CardStorage struct {
	Conn *pgxpool.Pool
}

type Card struct {
	ID             string
	UID            string
	CardNumberHash string
	ValidUntilHash string
	CVVHash        string
	LastFourDigits string
	EntryHash      string
	IsDeleted      bool
}

func InitCardStorage(cfg *config.Config) (*CardStorage, error) {
	conn, err := pgxpool.New(context.Background(), cfg.DatabaseURI)

	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err.Error())
		return nil, err
	}

	return &CardStorage{conn}, nil
}

func (c *CardStorage) CreateCard(ctx context.Context, uid, cardNumberHash, validUntilHash, CVVHash, lastFourDigits, entryHash string) (string, error) {
	sqlStatement := `
	INSERT INTO cards (uid, card_number_hash, valid_until_hash, CVV_hash, last_four_digits, entry_hash)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id;`

	var id string
	err := c.Conn.QueryRow(ctx, sqlStatement, uid, cardNumberHash, validUntilHash, CVVHash, lastFourDigits, entryHash).Scan(&id)

	if err != nil {
		return id, utils.WrapError(err, utils.ErrDBLayer)
	}

	return id, nil
}

func (c *CardStorage) UpdateCard(ctx context.Context, id, cardNumberHash, validUntilHash, CVVHash, lastFourDigits, entryHash string, isDeleted bool) error {
	updateBalanceSQL := `
	UPDATE CARDS SET 
	card_number_hash = $2,
	valid_until_hash = $3,
	CVV_hash = $4,
	last_four_digits = $5,
	entry_hash = $6,
	is_deleted = $7

	WHERE id = $1
	`

	_, err := c.Conn.Exec(ctx, updateBalanceSQL, id, cardNumberHash, validUntilHash, CVVHash, lastFourDigits, entryHash, isDeleted)

	if err != nil {
		return utils.WrapError(err, utils.ErrDBLayer)
	}

	return nil
}

func (c *CardStorage) ListCardByUID(ctx context.Context, uid string) ([]Card, error) {
	sqlStatement := `
	SELECT id, uid, card_number_hash, valid_until_hash, CVV_hash, last_four_digits, entry_hash, is_deleted FROM cards WHERE UID = $1 
	`

	rows, err := c.Conn.Query(ctx, sqlStatement, uid)
	if err != nil {
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, utils.ErrNotFound
			}

			return nil, utils.WrapError(err, utils.ErrDBLayer)
		}
	}

	defer rows.Close()

	cards := []Card{}

	for rows.Next() {
		entry := Card{}
		err = rows.Scan(&entry.ID, &entry.UID, &entry.CardNumberHash, &entry.ValidUntilHash, &entry.CVVHash, &entry.LastFourDigits, &entry.EntryHash, &entry.IsDeleted)

		if err != nil {
			return nil, err
		}

		cards = append(cards, entry)
	}

	err = rows.Err()
	if err != nil {
		return nil, utils.WrapError(err, utils.ErrDBLayer)
	}

	return cards, nil
}

func (c *CardStorage) GetCardByID(ctx context.Context, id string) (*Card, error) {
	sqlStatement := `
	SELECT id, uid, card_number_hash, valid_until_hash, CVV_hash, last_four_digits, entry_hash, is_deleted FROM cards WHERE ID = $1 
	`

	card := Card{}

	err := c.Conn.QueryRow(ctx, sqlStatement, id).Scan(&card.ID, &card.UID, &card.CardNumberHash, &card.ValidUntilHash, &card.CVVHash, &card.LastFourDigits, &card.EntryHash, &card.IsDeleted)
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

	return &card, nil
}

func (c *CardStorage) Close() {
	c.Conn.Close()
}
