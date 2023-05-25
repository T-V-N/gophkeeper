package config

type S3Config struct {
	S3URL                string `env:"S3URL"`
	S3AccessKey          string `env:"S3_ACCESS_KEY"`
	S3Secret             string `env:"S3_SECRET"`
	S3Bucket             string `env:"S3_BUCKET"`
	FileUpdateTimeWindow int    `env:"FILE_UPDATE_TIME_WINDOW" envDefault:"5"`
}
