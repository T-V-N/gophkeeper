package config

import (
	"github.com/caarlos0/env/v8"
)

type Config struct {
	S3Config
	StorageConfig
	SecureConfig
	EmailSenderConfig
	ApplicationConfig
}

func Init() (*Config, error) {
	cfg := &Config{}
	err := env.Parse(cfg)

	if err != nil {
		return nil, err
	}

	return cfg, nil
}
