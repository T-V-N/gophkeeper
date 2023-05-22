package config

type EmailSenderConfig struct {
	RegisterSender       string `env:"REGISTER_SENDER" envDefault:"test@test.com"`
	SenderName           string `env:"SENDER_NAME" envDefault:"Admin"`
	RegisterTemplateID   string `env:"REGISTER_TEMPLATE_ID"`
	SengridAPIKey        string `env:"SENDGRID_API_KEY"`
	FileUpdateTimeWindow int    `env:"FILE_UPDATE_TIME_WINDOW" envDefault:"5"`
}
