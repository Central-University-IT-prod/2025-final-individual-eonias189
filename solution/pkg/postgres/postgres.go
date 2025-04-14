package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func Connect(ctx context.Context, cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", cfg.GetConnString())
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
