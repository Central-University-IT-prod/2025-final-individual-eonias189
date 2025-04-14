package minio

type Config struct {
	Endpoint string `env:"MINIO_ENDPOINT" env-default:"localhost:9000"`
	User     string `env:"MINIO_USER"`
	Password string `env:"MINIO_PASSWORD"`
}
