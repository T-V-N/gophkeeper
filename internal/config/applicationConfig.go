package config

type ApplicationConfig struct {
	RPCPort              string `env:"RPC_PORT" envDefault:":8081"`
	RunAddress           string `env:"RUN_ADDRESS" envDefault:":8080"`
	ContextCancelTimeout int    `env:"CONTEXT_CANCEL_TIMEOUT" envDefault:"10"`
}
