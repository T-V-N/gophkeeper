package app

import (
	"context"

	"github.com/T-V-N/gophkeeper/internal/config"
	"github.com/T-V-N/gophkeeper/internal/storage"
	"github.com/T-V-N/gophkeeper/internal/utils"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type CardApp struct {
	Card   *storage.CardStorage
	Cfg    *config.Config
	logger *zap.SugaredLogger
}

func InitCardApp(conn *pgxpool.Pool, cfg *config.Config, logger *zap.SugaredLogger) (*CardApp, error) {
	c, err := storage.InitCard(conn)

	if err != nil {
		return nil, err
	}

	return &CardApp{c, cfg, logger}, nil
}

func (c *CardApp) CreateCard(ctx context.Context, uid, cardNumberHash, validUntilHash, CVVHash, lastFourDigits string) (string, error) {
	entryHash := utils.PackedCheckSum([]string{cardNumberHash, validUntilHash, CVVHash})

	id, err := c.Card.CreateCard(ctx, uid, cardNumberHash, validUntilHash, CVVHash, lastFourDigits, entryHash)

	if err != nil {
		return "", utils.WrapError(err, utils.ErrDBLayer)
	}

	return id, nil
}

func (c *CardApp) UpdateCard(ctx context.Context, id, cardNumberHash, validUntilHash, CVVHash, lastFourDigits, previousHash string, isDeleted, forceUpdate bool) error {
	card, err := c.Card.GetCardByID(ctx, id)

	if err != nil {
		return utils.WrapError(err, utils.ErrDBLayer)
	}

	if card == nil {
		return utils.ErrNotFound
	}

	if (previousHash != card.EntryHash) && !forceUpdate {
		return utils.ErrConflict
	}

	entryHash := utils.PackedCheckSum([]string{cardNumberHash, validUntilHash, CVVHash})

	return c.Card.UpdateCard(ctx, id, cardNumberHash, validUntilHash, CVVHash, lastFourDigits, entryHash, isDeleted)
}
