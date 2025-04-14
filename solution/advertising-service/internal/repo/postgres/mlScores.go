package postgres

import (
	"advertising/advertising-service/internal/models"
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type MLScoresRepo struct {
	db *sqlx.DB
	sq sq.StatementBuilderType
}

func NewMlScoresRepo(db *sqlx.DB) *MLScoresRepo {
	return &MLScoresRepo{
		db: db,
		sq: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (msr *MLScoresRepo) UpsertMLScore(ctx context.Context, mlscore models.MLScore) error {
	op := "MLScoresRepo.UpsertMLScore"

	query, args, err := msr.sq.
		Insert("ml_scores").
		Columns("client_id", "advertiser_id", "score").
		Values(mlscore.ClientId, mlscore.AdvertiserId, mlscore.Score).
		Suffix(
			`ON CONFLICT (client_id, advertiser_id)
			DO UPDATE SET
			score = EXCLUDED.score
		`,
		).ToSql()
	if err != nil {
		return fmt.Errorf("%s: build query: %s", op, err)
	}

	if _, err := msr.db.ExecContext(ctx, query, args...); err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23503":
				switch pqErr.Constraint {
				case "ml_scores_client_id_fkey":
					return models.ErrClientNotFound
				case "ml_scores_advertiser_id_fkey":
					return models.ErrAdvertiserNotFound
				}
			}
		}
		return fmt.Errorf("%s: db.ExecContext: %w", op, err)
	}

	return nil
}
