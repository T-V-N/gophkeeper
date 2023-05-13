package app

import (
	"context"
	"time"

	"github.com/T-V-N/gophkeeper/internal/config"
	"github.com/T-V-N/gophkeeper/internal/storage"
	"github.com/T-V-N/gophkeeper/internal/storage/s3"

	"github.com/T-V-N/gophkeeper/internal/utils"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type FileApp struct {
	File   *storage.FileStorage
	s3     *s3.S3Store
	Cfg    *config.Config
	logger *zap.SugaredLogger
}

func InitFileApp(conn *pgxpool.Pool, cfg *config.Config, logger *zap.SugaredLogger) (*FileApp, error) {
	file, err := storage.InitFile(conn)

	if err != nil {
		return nil, err
	}

	s3Store := s3.InitS3Storage(context.Background(), cfg)

	return &FileApp{file, s3Store, cfg, logger}, nil
}

func (fa *FileApp) CreateFile(ctx context.Context, uid, fileName string) (string, error) {
	id, err := fa.File.CreateFile(ctx, uid, fileName)

	if err != nil {
		return "", utils.WrapError(err, utils.ErrDBLayer)
	}

	return id, nil
}

func (fa *FileApp) RequestUpdateFile(ctx context.Context, id string) (string, error) {
	file, err := fa.File.GetFileByID(ctx, id)

	if err != nil {
		return "", utils.WrapError(err, utils.ErrDBLayer)
	}

	if file == nil {
		return "", utils.ErrNotFound
	}

	currentTime := time.Now()

	if int(currentTime.Sub(file.CommittedAt).Minutes()) < fa.Cfg.FileUpdateTimeWindow {
		return "", utils.ErrConflict
	}

	return fa.s3.GetUploadLink(ctx, id)
}

func (fa *FileApp) CommitUpdateFile(ctx context.Context, id string, committedAt time.Time, forceUpdate bool) error {
	file, err := fa.File.GetFileByID(ctx, id)

	if err != nil {
		return utils.WrapError(err, utils.ErrDBLayer)
	}

	if file == nil {
		return utils.ErrNotFound
	}

	if (committedAt != file.CommittedAt) && !forceUpdate {
		return utils.ErrConflict
	}

	return fa.File.UpdateFile(ctx, id, file.FileName, file.S3Link, file.IsDeleted, committedAt)
}
