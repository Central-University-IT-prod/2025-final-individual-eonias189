package postgres

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Migrate(migrationsDir string, databaseUrl string) error {
	m, err := migrate.New(fmt.Sprintf("file://%s", migrationsDir), databaseUrl)
	if err != nil {
		return err
	}

	if err := m.Up(); err != migrate.ErrNoChange {
		return err
	}

	return nil
}
