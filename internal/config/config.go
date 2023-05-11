package config

import (
	"flag"

	"github.com/caarlos0/env/v8"
)

type Config struct {
	S3URL                string `env:"S3URL"`
	S3AccessKey          string `env:"S3_ACCESS_KEY"`
	S3Secret             string `env:"S3_SECRET"`
	S3Bucket             string `env:"S3_BUCKET"`
	DatabaseURI          string `env:"DATABASE_URI" envDefault:"postgresql://root:root@localhost:5433/keeper"`
	MigrationsPath       string `env:"MIGRATIONS_PATH" envDefault:"migrations"`
	RunAddress           string `env:"RUN_ADDRESS" envDefault:":8080"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS" envDefault:"http://127.0.0.1:8888"`
	JWTExpireTiming      int64  `env:"JWT_EXPIRE_TIMING" envDefault:"10000"`
	SecretKey            string `env:"SECRET_KEY" envDefault:"secret"`
	CompressLevel        int    `env:"COMPRESS_LEVEL" envDefault:"5"`
	CheckOrderDelay      uint   `env:"CHECK_ORDER_DELAY" envDefault:"1000"`
	CheckOrderInterval   uint   `env:"CHECK_ORDER_INTERVAL" envDefault:"2"`
	WorkerLimit          int    `env:"WORKER_LIMIT" envDefault:"1"`
	ContextCancelTimeout int    `env:"CONTEXT_CANCEL_AMOUNT" envDefault:"10"`
	FileUpdateTimeWindow int    `env:"FILE_UPDATE_TIME_WINDOW" envDefault:"5"`

	RegisterTemplateID string `env:"REGISTER_TEMPLATE_ID"`
	SengridAPIKey      string `env:"SENDGRID_API_KEY"`
}

func Init() (*Config, error) {
	cfg := &Config{}
	err := env.Parse(cfg)

	if err != nil {
		return nil, err
	}

	flag.StringVar(&cfg.RunAddress, "a", cfg.RunAddress, "server address")
	flag.StringVar(&cfg.DatabaseURI, "d", cfg.DatabaseURI, "db address")
	flag.StringVar(&cfg.AccrualSystemAddress, "r", cfg.AccrualSystemAddress, "accrual system address")
	flag.Parse()

	return cfg, nil
}
