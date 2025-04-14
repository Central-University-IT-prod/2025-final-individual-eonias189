package redis

type Config struct {
	Host string `env:"REDIS_HOST" env-default:"localhost"`
	Port int    `env:"REDIS_PORT" env-default:"6379"`
	Db   int    `env:"REDIS_DB" env-default:"0"`
}
