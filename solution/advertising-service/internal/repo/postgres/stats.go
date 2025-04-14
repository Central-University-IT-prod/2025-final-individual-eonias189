package postgres

import (
	"advertising/advertising-service/internal/models"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type StatsRepo struct {
	db *sqlx.DB
}

func NewStatsRepo(db *sqlx.DB) *StatsRepo {
	return &StatsRepo{
		db: db,
	}
}

func (sr *StatsRepo) GetStatsForCampaign(ctx context.Context, campaignId uuid.UUID) (models.Stats, error) {
	op := "StatsRepo.GetStatsForCampaign"

	// $1 - campaign id
	query := `
	WITH
		impressions_stats AS
		(
			SELECT
				count(*) AS impressions_count,
				COALESCE(sum(profit), 0) AS spent_impressions
			FROM impressions
			WHERE campaign_id = $1
		),
		clicks_stats AS
		(
			SELECT
				count(*) AS clicks_count,
				COALESCE(sum(profit), 0) AS spent_clicks
			FROM clicks
			WHERE campaign_id = $1
		)
	SELECT *
	FROM impressions_stats
	JOIN clicks_stats ON true
	`

	var stats models.Stats
	if err := sr.db.GetContext(ctx, &stats, query, campaignId); err != nil {
		return models.Stats{}, fmt.Errorf("%s: db.GetContext: %w", op, err)
	}

	if stats.ImpressionsCount == 0 {
		stats.Conversion = 0
	} else {
		stats.Conversion = float64(stats.ClicksCount) / float64(stats.ImpressionsCount) * 100
	}

	stats.SpentTotal = stats.SpentImpressions + stats.SpentClicks

	return stats, nil
}

func (sr *StatsRepo) GetStatsForCampaignDaily(ctx context.Context, campaignId uuid.UUID) ([]models.StatsDaily, error) {
	op := "StatsRepo.GetStatsForCampaignDaily"

	// $1 - campaign id
	query := `
	WITH
		impressions_stats AS
    	(
    		SELECT
				count(*) AS impressions_count,
				COALESCE(sum(profit), 0) AS spent_impressions,
				date
			FROM impressions
			WHERE campaign_id = $1
     		GROUP BY date
    	),
    	clicks_stats AS
    	(
    	 	SELECT
				count(*) AS clicks_count,
				COALESCE(sum(profit), 0) AS spent_clicks,
				date
			FROM clicks
			WHERE campaign_id = $1
    	    GROUP BY date
    	)
	SELECT
		COALESCE(impressions_stats.impressions_count, 0) AS impressions_count,
		COALESCE(impressions_stats.spent_impressions, 0) AS spent_impressions,
		COALESCE(clicks_stats.clicks_count, 0) AS clicks_count,
		COALESCE(clicks_stats.spent_clicks, 0) AS spent_clicks,
		COALESCE(impressions_stats.date, clicks_stats.date) AS date
	FROM impressions_stats
	FULL JOIN clicks_stats ON clicks_stats.date = impressions_stats.date
	ORDER BY COALESCE(impressions_stats.date, clicks_stats.date) ASC
	`

	dailyStats := []models.StatsDaily{}
	if err := sr.db.SelectContext(ctx, &dailyStats, query, campaignId); err != nil {
		return nil, fmt.Errorf("%s: db.SelectContext: %w", op, err)
	}

	for i, stats := range dailyStats {
		if stats.ImpressionsCount == 0 {
			stats.Conversion = 0
		} else {
			stats.Conversion = float64(stats.ClicksCount) / float64(stats.ImpressionsCount) * 100
		}
		stats.SpentTotal = stats.SpentImpressions + stats.SpentClicks
		dailyStats[i] = stats
	}

	return dailyStats, nil
}

func (sr *StatsRepo) GetStatsForAdvertiser(ctx context.Context, advertiserId uuid.UUID) (models.Stats, error) {
	op := "StatsRepo.GetStatsForAdvertiser"

	// $1 - advertiser id
	query := `
	WITH
		impressions_stats AS
		(
			SELECT
				count(*) AS impressions_count,
				COALESCE(sum(profit), 0) AS spent_impressions
			FROM impressions
            JOIN campaigns ON campaigns.id = impressions.campaign_id
			WHERE campaigns.advertiser_id = $1
		),
		clicks_stats AS
		(
			SELECT
				count(*) AS clicks_count,
				COALESCE(sum(profit), 0) AS spent_clicks
			FROM clicks
            JOIN campaigns ON campaigns.id = clicks.campaign_id
			WHERE campaigns.advertiser_id = $1
		)
	SELECT *
	FROM impressions_stats
	JOIN clicks_stats ON true
	`

	var stats models.Stats
	if err := sr.db.GetContext(ctx, &stats, query, advertiserId); err != nil {
		return models.Stats{}, fmt.Errorf("%s: db.GetContext: %w", op, err)
	}

	if stats.ImpressionsCount == 0 {
		stats.Conversion = 0
	} else {
		stats.Conversion = float64(stats.ClicksCount) / float64(stats.ImpressionsCount) * 100
	}

	stats.SpentTotal = stats.SpentImpressions + stats.SpentClicks

	return stats, nil
}

func (sr *StatsRepo) GetStatsForAdvertiserDaily(ctx context.Context, advertiserId uuid.UUID) ([]models.StatsDaily, error) {
	op := "StatsRepo.GetStatsForAdvertiserDaily"

	// $1 - campaign id
	query := `
	WITH
		impressions_stats AS
    	(
    		SELECT
				count(*) AS impressions_count,
				COALESCE(sum(profit), 0) AS spent_impressions,
				date
			FROM impressions
			JOIN campaigns ON campaigns.id = impressions.campaign_id
			WHERE campaigns.advertiser_id = $1
     		GROUP BY date
    	),
    	clicks_stats AS
    	(
    	 	SELECT
				count(*) AS clicks_count,
				COALESCE(sum(profit), 0) AS spent_clicks,
				date
			FROM clicks
			JOIN campaigns ON campaigns.id = clicks.campaign_id
			WHERE campaigns.advertiser_id = $1
    	    GROUP BY date
    	)
	SELECT
		COALESCE(impressions_stats.impressions_count, 0) AS impressions_count,
		COALESCE(impressions_stats.spent_impressions, 0) AS spent_impressions,
		COALESCE(clicks_stats.clicks_count, 0) AS clicks_count,
		COALESCE(clicks_stats.spent_clicks, 0) AS spent_clicks,
		COALESCE(impressions_stats.date, clicks_stats.date) AS date
	FROM impressions_stats
	FULL JOIN clicks_stats ON clicks_stats.date = impressions_stats.date
	ORDER BY COALESCE(impressions_stats.date, clicks_stats.date) ASC
	`

	dailyStats := []models.StatsDaily{}
	if err := sr.db.SelectContext(ctx, &dailyStats, query, advertiserId); err != nil {
		return nil, fmt.Errorf("%s: db.SelectContext: %w", op, err)
	}

	for i, stats := range dailyStats {
		if stats.ImpressionsCount == 0 {
			stats.Conversion = 0
		} else {
			stats.Conversion = float64(stats.ClicksCount) / float64(stats.ImpressionsCount) * 100
		}
		stats.SpentTotal = stats.SpentImpressions + stats.SpentClicks
		dailyStats[i] = stats
	}

	return dailyStats, nil
}
