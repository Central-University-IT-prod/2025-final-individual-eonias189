package postgres

import (
	"fmt"
	"time"
)

type Config struct {
	Host            string        `env:"POSTGRES_HOST" env-default:"localhost"`
	Port            int           `env:"POSTGRES_PORT" env-default:"5432"`
	DB              string        `env:"POSTGRES_DB" env-default:"postgres"`
	User            string        `env:"POSTGRES_USER" env-default:"postgres"`
	Password        string        `env:"POSTGRES_PASSWORD"`
	MaxOpenConns    int           `env:"POSTGRES_MAX_OPEN_CONNS" env-default:"10"`
	MaxIdleConns    int           `env:"POSTGRES_MAX_IDLE_CONNS" env-default:"5"`
	ConnMaxLifetime time.Duration `env:"POSTGRES_CONN_MAX_LIFETIME" env-default:"3m"`
	ConnMaxIdleTime time.Duration `env:"POSTGRES_CONN_MAX_IDLE_TIME" env-default:"2m"`
}

func (cfg Config) GetConnString() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DB,
	)
}
