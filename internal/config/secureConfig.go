package config

type SecureConfig struct {
	SecretKey       string `env:"SECRET_KEY" envDefault:"secret"`
	JWTExpireTiming int64  `env:"JWT_EXPIRE_TIMING" envDefault:"10000"`
}
