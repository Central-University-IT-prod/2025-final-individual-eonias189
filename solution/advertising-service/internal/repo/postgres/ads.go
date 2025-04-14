package postgres

import (
	"advertising/advertising-service/internal/models"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type AdsRepo struct {
	db *sqlx.DB
}

func NewAdsRepo(db *sqlx.DB) *AdsRepo {
	return &AdsRepo{
		db: db,
	}
}

func (ar *AdsRepo) GetAdForClient(ctx context.Context, client models.Client, currentDay int) (models.Ad, error) {
	op := "AdsRepo.GetAdForClient"

	// $1 - client id
	// $2 - current day
	// $3 - client gender
	// $4 - client location
	// $5 - client age

	args := []any{
		client.Id, currentDay, client.Gender, client.Location, client.Age,
	}

	query := `
	WITH
		ml_scores_max_score AS
		(
			SELECT
				CASE
					WHEN max(score) != 0
						THEN max(score)
					ELSE 1
				END AS max_score
			FROM ml_scores
		),
		impressions_counted AS
		(
			SELECT campaign_id, count(*) AS impressions_count
			FROM impressions
			GROUP BY campaign_id
		),
		impressions_by_client AS
		(
			SELECT campaign_id, true AS impressed_by_client
			FROM impressions
			WHERE client_id = $1
		),
    	campaigns_filtered AS
		(
			SELECT *
			FROM campaigns
			WHERE
				$2 BETWEEN campaigns.start_date AND campaigns.end_date AND
				(campaigns.gender IS NULL OR campaigns.gender = 'ALL' OR campaigns.gender = $3) AND
				(campaigns.location IS NULL OR campaigns.location = $4) AND
				$5 BETWEEN COALESCE(campaigns.age_from, -1) AND COALESCE(campaigns.age_to, 999)
		),
		candidates AS 
		(
			SELECT
				campaigns.id AS campaign_id,
				campaigns.advertiser_id AS advertiser_id,
				campaigns.ad_title AS ad_title,
				campaigns.ad_text AS ad_text,
				campaigns.ad_image_url AS ad_image_url,
          		cost_per_impression + (COALESCE(ml_scores.score, 0)::double precision / max_score::double precision) * 0.5 * cost_per_click AS profit,
				COALESCE(ml_scores.score, 0) AS score,
				ABS(campaigns.impressions_limit - COALESCE(impressions_counted.impressions_count, 0)) AS limits_diff
			FROM campaigns_filtered campaigns
			LEFT JOIN ml_scores ON
				ml_scores.client_id = $1 AND
				ml_scores.advertiser_id = campaigns.advertiser_id
			LEFT JOIN impressions_counted ON impressions_counted.campaign_id = campaigns.id
			LEFT JOIN impressions_by_client on impressions_by_client.campaign_id = campaigns.id
			JOIN ml_scores_max_score ON true
            WHERE
				COALESCE(impressions_by_client.impressed_by_client, false) = false AND
				COALESCE(impressions_counted.impressions_count, 0) < ROUND(campaigns.impressions_limit::double precision * 1.05)
 		),
		candidates_max_values AS
        (
        	SELECT
				CASE
					WHEN max(profit) != 0
						THEN max(profit)
					ELSE 0.1
				END AS max_profit,
				CASE
					WHEN max(limits_diff) != 0
						THEN max(limits_diff)
					ELSE 1
				END AS max_limits_diff,
				CASE
					WHEN max(score) != 0
						THEN max(score)
					ELSE 1
				END AS max_score
          	FROM candidates
        ),
        candidates_normalized_profit AS
        (
        	SELECT
          		*,
          		profit / max_profit AS profit_normalized,
				limits_diff::double precision / max_limits_diff::double precision AS limits_compilance,
				score::double precision / max_score::double precision AS score_normalized
          	FROM candidates
          	JOIN candidates_max_values ON true
        )
	SELECT campaign_id, advertiser_id, ad_title, ad_text, ad_image_url
	FROM candidates_normalized_profit
	ORDER BY profit_normalized + score_normalized * 0.25 + limits_compilance * 0.1 DESC
	LIMIT 1
	`
	var ad models.Ad
	if err := ar.db.GetContext(ctx, &ad, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Ad{}, models.ErrNoAdsForClient
		}
		return models.Ad{}, fmt.Errorf("%s: execute query: %w", op, err)
	}

	return ad, nil
}
