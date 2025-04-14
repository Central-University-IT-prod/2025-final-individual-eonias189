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
	"github.com/lib/pq"
)

type ClientActionsRepo struct {
	db *sqlx.DB
	sq sq.StatementBuilderType
}

func NewClientActionsRepo(db *sqlx.DB) *ClientActionsRepo {
	return &ClientActionsRepo{
		db: db,
		sq: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (car *ClientActionsRepo) RecordImpression(ctx context.Context, impression models.Impression) error {
	op := "ClientActionsRepo.RecordImpression"

	query, args, err := car.sq.
		Insert("impressions").
		Columns("client_id", "campaign_id", "date", "profit").
		Values(impression.ClientId, impression.CampaignId, impression.Date, impression.Profit).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: build query: %w", op, err)
	}

	if _, err := car.db.ExecContext(ctx, query, args...); err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505":
				return models.ErrAlreadyImpressed
			case "23503":
				switch pqErr.Constraint {
				case "impressions_client_id_fkey":
					return models.ErrClientNotFound
				case "impressions_campaign_id_fkey":
					return models.ErrCampaignNotFound
				}
			}
		}
		return fmt.Errorf("%s: db.ExecContext: %w", op, err)
	}

	return nil
}

func (car *ClientActionsRepo) RecordClick(ctx context.Context, click models.Click) error {
	op := "ClientActionsRepo.RecordClick"

	query, args, err := car.sq.
		Insert("clicks").
		Columns("client_id", "campaign_id", "date", "profit").
		Values(click.ClientId, click.CampaignId, click.Date, click.Profit).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: build query: %w", op, err)
	}

	if _, err := car.db.ExecContext(ctx, query, args...); err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505":
				return models.ErrAlreadyClicked
			case "23503":
				switch pqErr.Constraint {
				case "clicks_client_id_fkey":
					return models.ErrClientNotFound
				case "clicks_campaign_id_fkey":
					return models.ErrCampaignNotFound
				}
			}
		}
		return fmt.Errorf("%s: db.ExecContext: %w", op, err)
	}

	return nil
}

func (car *ClientActionsRepo) CheckImpressed(ctx context.Context, clientId, campaignId uuid.UUID) (bool, error) {
	op := "ClientActionsRepo.CheckImpressed"

	query, args, err := car.sq.
		Select("true").
		From("impressions").
		Where(sq.Eq{
			"client_id":   clientId,
			"campaign_id": campaignId,
		}).
		ToSql()
	if err != nil {
		return false, fmt.Errorf("%s: build query: %w", op, err)
	}

	var impressed bool
	if err := car.db.QueryRowContext(ctx, query, args...).Scan(&impressed); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}

		return false, fmt.Errorf("%s: db.QueryRowContext: %w", op, err)
	}

	return true, nil
}
