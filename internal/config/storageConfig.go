package config

type StorageConfig struct {
	DatabaseURI    string `env:"DATABASE_URI" envDefault:"postgresql://root:root@localhost:5433/keeper"`
	MigrationsPath string `env:"MIGRATIONS_PATH" envDefault:"migrations"`
}
