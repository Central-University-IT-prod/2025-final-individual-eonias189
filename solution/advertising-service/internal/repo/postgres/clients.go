package postgres

import (
	"advertising/advertising-service/internal/models"
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ClientsRepo struct {
	db *sqlx.DB
	sq sq.StatementBuilderType
}

func NewClientRepo(db *sqlx.DB) *ClientsRepo {
	return &ClientsRepo{
		db: db,
		sq: sq.StatementBuilderType{}.PlaceholderFormat(sq.Dollar),
	}
}

func (cr *ClientsRepo) GetClientById(ctx context.Context, id uuid.UUID) (models.Client, error) {
	op := "ClientsRepo.GetClientById"

	query, args, err := cr.sq.
		Select("id", "login", "age", "location", "gender").
		From("clients").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return models.Client{}, fmt.Errorf("%s: build query: %w", op, err)
	}

	var client models.Client
	if err := cr.db.GetContext(ctx, &client, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Client{}, models.ErrClientNotFound
		}

		return models.Client{}, fmt.Errorf("%s: db.GetContext: %w", op, err)
	}

	return client, nil
}

func (cr *ClientsRepo) UpsertClients(ctx context.Context, clients []models.Client) ([]models.Client, error) {
	op := "ClientsRepo.UpsertClients"

	qb := cr.sq.
		Insert("clients").
		Columns("id", "login", "age", "location", "gender")

	toInsert := map[uuid.UUID]models.Client{}
	for _, client := range clients {
		toInsert[client.Id] = client
	}

	for _, client := range toInsert {
		qb = qb.Values(client.Id, client.Login, client.Age, client.Location, client.Gender)
	}

	query, args, err := qb.
		Suffix(`ON CONFLICT (id) DO UPDATE SET
			login = EXCLUDED.login,
			age = EXCLUDED.age,
			location = EXCLUDED.location,
			gender = EXCLUDED.gender`,
		).ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: build query: %w", op, err)
	}

	if _, err := cr.db.ExecContext(ctx, query, args...); err != nil {
		return nil, fmt.Errorf("%s: db.ExecContext: %w", op, err)
	}

	inserted := make([]models.Client, 0, len(toInsert))
	for _, client := range toInsert {
		inserted = append(inserted, client)
	}

	return inserted, nil
}
