package storage

import (
	"context"
	"log"

	"github.com/T-V-N/gophkeeper/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	Conn *pgxpool.Pool
	cfg  config.Config
}

func InitStorage(cfg config.Config) (*Storage, error) {
	conn, err := pgxpool.New(context.Background(), cfg.DatabaseURI)

	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err.Error())
		return nil, err
	}

	return &Storage{conn, cfg}, nil
}
