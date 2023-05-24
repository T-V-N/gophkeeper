package app

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/T-V-N/gophkeeper/internal/config"
	"github.com/T-V-N/gophkeeper/internal/storage"
	"github.com/T-V-N/gophkeeper/internal/storage/s3"
	"github.com/T-V-N/gophkeeper/internal/utils"
)

type File interface {
	CreateFile(ctx context.Context, uid, fileName string) (string, error)
	UpdateFile(ctx context.Context, id, fileName, s3Link string, isDeleted bool, committedAt time.Time) error
	ListFilesByUID(ctx context.Context, uid string) (*[]storage.File, error)
	GetFileByID(ctx context.Context, id string) (*storage.File, error)
	Close()
}

type S3Store interface {
	GetUploadLink(ctx context.Context, id string) (string, error)
	GetFileInfo(ctx context.Context, id string) (time.Time, string, error)
}

type FileApp struct {
	File   File
	S3     S3Store
	Cfg    *config.Config
	logger *zap.SugaredLogger
}

func InitFileApp(cfg *config.Config, logger *zap.SugaredLogger) (*FileApp, error) {
	file, err := storage.InitFileStorage(cfg)

	if err != nil {
		return nil, err
	}

	s3Store := s3.InitS3Storage(context.Background(), &cfg.S3Config)

	return &FileApp{File: file, S3: s3Store, Cfg: cfg, logger: logger}, nil
}

func (fa *FileApp) CreateFile(ctx context.Context, uid, fileName string) (string, error) {
	id, err := fa.File.CreateFile(ctx, uid, fileName)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (fa *FileApp) RequestUpdateFile(ctx context.Context, uid, id string) (string, error) {
	file, err := fa.File.GetFileByID(ctx, id)

	if err != nil {
		return "", err
	}

	if file.UID != uid {
		return "", utils.ErrNotAuthorized
	}

	currentTime := time.Now()

	if int(currentTime.Sub(file.CommittedAt).Minutes()) < fa.Cfg.S3Config.FileUpdateTimeWindow {
		return "", utils.ErrConflict
	}

	return fa.S3.GetUploadLink(ctx, id)
}

func (fa *FileApp) CommitUpdateFile(ctx context.Context, uid, id string, previousCommitedAt time.Time, forceUpdate bool) error {
	file, err := fa.File.GetFileByID(ctx, id)

	if err != nil {
		return err
	}

	if file.UID != uid {
		return utils.ErrNotAuthorized
	}

	fileCommitted, url, err := fa.S3.GetFileInfo(ctx, id)

	if err != nil {
		return err
	}

	if (int64(file.CommittedAt.Compare(previousCommitedAt)) != 0) && !forceUpdate && file.CommittedAt.Unix() != 0 {
		return utils.ErrConflict
	}

	return fa.File.UpdateFile(ctx, id, file.FileName, url, file.IsDeleted, fileCommitted)
}

func (fa *FileApp) ListFile(ctx context.Context, uid string, existingFiles []ExistingFiles) ([]storage.File, error) {
	files, err := fa.File.ListFilesByUID(ctx, uid)

	if err != nil {
		return nil, err
	}

	var result []storage.File

	for _, file := range *files {
		include := true

		for _, existingFile := range existingFiles {
			if (file.CommittedAt.Unix() == existingFile.CommittedAt.Unix()) && (file.ID == existingFile.ID) {
				include = false
				break
			}
		}

		if include {
			result = append(result, file)
		}
	}

	return result, nil
}
