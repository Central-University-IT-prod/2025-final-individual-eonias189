package repo

import (
	"advertising/advertising-service/internal/models"
	"context"

	"github.com/google/uuid"
)

//go:generate go run github.com/vektra/mockery/v2@v2.52.2 --name StatsRepo
type StatsRepo interface {
	GetStatsForCampaign(ctx context.Context, campaignId uuid.UUID) (models.Stats, error)
	GetStatsForCampaignDaily(ctx context.Context, campaignId uuid.UUID) ([]models.StatsDaily, error)
	GetStatsForAdvertiser(ctx context.Context, advertiserId uuid.UUID) (models.Stats, error)
	GetStatsForAdvertiserDaily(ctx context.Context, advertiserId uuid.UUID) ([]models.StatsDaily, error)
}
