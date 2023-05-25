package app

import (
	"context"

	"go.uber.org/zap"

	"github.com/T-V-N/gophkeeper/internal/config"
	"github.com/T-V-N/gophkeeper/internal/storage"
	"github.com/T-V-N/gophkeeper/internal/utils"
)

type LogPassword interface {
	CreateLogPassword(ctx context.Context, uid, loginHash, passwordHash, resourceName, entryHash string) (string, error)
	UpdateLogPassword(ctx context.Context, id, loginHash, passwordHash, resourceName, entryHash string, isDeleted bool) error
	ListLogPasswordByUID(ctx context.Context, uid string) ([]storage.LogPassword, error)
	GetLogPasswordByID(ctx context.Context, id string) (*storage.LogPassword, error)
	Close()
}

type LogPasswordApp struct {
	LogPassword LogPassword
	Cfg         *config.Config
	logger      *zap.SugaredLogger
}

func InitLogPasswordApp(cfg *config.Config, logger *zap.SugaredLogger) (*LogPasswordApp, error) {
	lp, err := storage.InitLogPasswordStorage(cfg)

	if err != nil {
		return nil, err
	}

	return &LogPasswordApp{lp, cfg, logger}, nil
}

func (lp *LogPasswordApp) CreateLogPassword(ctx context.Context, uid, loginHash, passwordHash, resourceName string) (string, error) {
	entryHash := utils.PackedCheckSum([]string{loginHash, passwordHash})

	id, err := lp.LogPassword.CreateLogPassword(ctx, uid, loginHash, passwordHash, resourceName, entryHash)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (lp *LogPasswordApp) UpdateLogPassword(ctx context.Context, uid, id, loginHash, passwordHash, previousHash, resourceName string, isDeleted, forceUpdate bool) error {
	logPass, err := lp.LogPassword.GetLogPasswordByID(ctx, id)

	if err != nil {
		return err
	}

	if logPass.UID != uid {
		return utils.ErrNotAuthorized
	}

	if (previousHash != logPass.EntryHash) && !forceUpdate {
		return utils.ErrConflict
	}

	entryHash := utils.PackedCheckSum([]string{loginHash, passwordHash})

	return lp.LogPassword.UpdateLogPassword(ctx, id, loginHash, passwordHash, resourceName, entryHash, isDeleted)
}

func (lp *LogPasswordApp) ListLogPassword(ctx context.Context, uid string, existingHashes []ExistingHash) ([]storage.LogPassword, error) {
	logPasswords, err := lp.LogPassword.ListLogPasswordByUID(ctx, uid)

	if err != nil {
		return nil, err
	}

	var result []storage.LogPassword

	for _, logPassword := range logPasswords {
		include := true

		for _, existingHash := range existingHashes {
			if logPassword.EntryHash == existingHash.EntryHash {
				include = false
				break
			}
		}

		if include {
			result = append(result, logPassword)
		}
	}

	return result, nil
}
