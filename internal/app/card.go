package app

import (
	"context"

	"go.uber.org/zap"

	"github.com/T-V-N/gophkeeper/internal/config"
	"github.com/T-V-N/gophkeeper/internal/storage"
	"github.com/T-V-N/gophkeeper/internal/utils"
)

type Card interface {
	CreateCard(ctx context.Context, uid, cardNumberHash, validUntilHash, CVVHash, lastFourDigits, entryHash string) (string, error)
	UpdateCard(ctx context.Context, id, cardNumberHash, validUntilHash, CVVHash, lastFourDigits, entryHash string, isDeleted bool) error
	ListCardByUID(ctx context.Context, uid string) ([]storage.Card, error)
	GetCardByID(ctx context.Context, id string) (*storage.Card, error)
	Close()
}

type CardApp struct {
	Card   Card
	Cfg    *config.Config
	logger *zap.SugaredLogger
}

func InitCardApp(cfg *config.Config, logger *zap.SugaredLogger) (*CardApp, error) {
	c, err := storage.InitCardStorage(cfg)

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

func (c *CardApp) UpdateCard(ctx context.Context, uid, id, cardNumberHash, validUntilHash, CVVHash, lastFourDigits, previousHash string, isDeleted, forceUpdate bool) error {
	card, err := c.Card.GetCardByID(ctx, id)

	if err != nil {
		return err
	}

	if card.UID != uid {
		return utils.ErrNotAuthorized
	}

	if (previousHash != card.EntryHash) && !forceUpdate {
		return utils.ErrConflict
	}

	entryHash := utils.PackedCheckSum([]string{cardNumberHash, validUntilHash, CVVHash})

	return c.Card.UpdateCard(ctx, id, cardNumberHash, validUntilHash, CVVHash, lastFourDigits, entryHash, isDeleted)
}

func (c *CardApp) ListCard(ctx context.Context, uid string, existingHashes []ExistingHash) ([]storage.Card, error) {
	cards, err := c.Card.ListCardByUID(ctx, uid)

	if err != nil {
		return nil, err
	}

	var result []storage.Card

	for _, card := range cards {
		include := true

		for _, existingHash := range existingHashes {
			if card.EntryHash == existingHash.EntryHash {
				include = false
				break
			}
		}

		if include {
			result = append(result, card)
		}
	}

	return result, nil
}
