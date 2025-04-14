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

type AdvertiserRepo struct {
	db *sqlx.DB
	sq sq.StatementBuilderType
}

func NewAdvertiserRepo(db *sqlx.DB) *AdvertiserRepo {
	return &AdvertiserRepo{
		db: db,
		sq: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (ar *AdvertiserRepo) GetAdvertiserById(ctx context.Context, id uuid.UUID) (models.Advertiser, error) {
	op := "AdvertiserRepo.GetAdvertiserById"

	query, args, err := ar.sq.
		Select("id", "name").
		From("advertisers").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return models.Advertiser{}, fmt.Errorf("%s: build query: %w", op, err)
	}

	var advertiser models.Advertiser
	if err := ar.db.GetContext(ctx, &advertiser, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Advertiser{}, models.ErrAdvertiserNotFound
		}

		return models.Advertiser{}, fmt.Errorf("%s: db.GetContext: %w", op, err)
	}

	return advertiser, nil
}

func (ar *AdvertiserRepo) UpsertAdvertisers(ctx context.Context, advertisers []models.Advertiser) ([]models.Advertiser, error) {
	op := "AdvertiserRepo.UpsertAdvertisers"

	toInsert := map[uuid.UUID]models.Advertiser{}
	for _, advertiser := range advertisers {
		toInsert[advertiser.Id] = advertiser
	}

	qb := ar.sq.Insert("advertisers").Columns("id", "name")
	for _, advertiser := range toInsert {
		qb = qb.Values(advertiser.Id, advertiser.Name)
	}

	query, args, err := qb.Suffix(
		`ON CONFLICT(id)
			DO UPDATE SET
			name = EXCLUDED.name
		`,
	).ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: build query: %w", op, err)
	}

	if _, err := ar.db.ExecContext(ctx, query, args...); err != nil {
		return nil, fmt.Errorf("%s: db.ExecContext: %w", op, err)
	}

	inserted := make([]models.Advertiser, 0, len(toInsert))
	for _, advertiser := range toInsert {
		inserted = append(inserted, advertiser)
	}

	return inserted, nil
}
