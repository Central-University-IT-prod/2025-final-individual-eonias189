package config

import (
	"advertising/pkg/minio"
	"advertising/pkg/openai"
	"advertising/pkg/postgres"
	"advertising/pkg/redis"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	ServerPort     int    `env:"SERVER_PORT" env-default:"8080"`
	LogLevel       string `env:"LOG_LEVEL" env-default:"info"`
	StaticBucket   string `env:"MINIO_STATIC_BUCKET" env-default:"static"`
	StaticBaseUrl  string `env:"STATIC_BASE_URL" env-default:"http://localhost:8080/static"`
	PostgresConfig postgres.Config
	RedisConfig    redis.Config
	MinioConfig    minio.Config
	OpenAIConfig   openai.Config
}

func Get() (Config, error) {
	var cfg Config
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}
