package app

import (
	"context"

	"github.com/T-V-N/gophkeeper/internal/config"
	"github.com/T-V-N/gophkeeper/internal/storage"
	"github.com/T-V-N/gophkeeper/internal/utils"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type LogPasswordApp struct {
	LogPassword *storage.LogPasswordStorage
	Cfg         *config.Config
	logger      *zap.SugaredLogger
}

func InitLogPasswordApp(conn *pgxpool.Pool, cfg *config.Config, logger *zap.SugaredLogger) (*LogPasswordApp, error) {
	lp, err := storage.InitLogPassword(conn)

	if err != nil {
		return nil, err
	}

	return &LogPasswordApp{lp, cfg, logger}, nil
}

func (lp *LogPasswordApp) CreateLogPassword(ctx context.Context, uid, loginHash, passwordHash, resourceName string) (string, error) {
	entryHash := utils.PackedCheckSum([]string{loginHash, passwordHash})

	id, err := lp.LogPassword.CreateLogPassword(ctx, uid, loginHash, passwordHash, resourceName, entryHash)

	if err != nil {
		return "", utils.WrapError(err, utils.ErrDBLayer)
	}

	return id, nil
}

func (lp *LogPasswordApp) UpdateLogPassword(ctx context.Context, id, loginHash, passwordHash, previousHash, resourceName string, isDeleted, forceUpdate bool) error {
	logPass, err := lp.LogPassword.GetLogPasswordByID(ctx, id)

	if err != nil {
		return utils.WrapError(err, utils.ErrDBLayer)
	}

	if logPass == nil {
		return utils.ErrNotFound
	}

	if (previousHash != logPass.EntryHash) && !forceUpdate {
		return utils.ErrConflict
	}

	entryHash := utils.PackedCheckSum([]string{loginHash, passwordHash})

	return lp.LogPassword.UpdateLogPassword(ctx, id, loginHash, passwordHash, resourceName, entryHash, isDeleted)
}
