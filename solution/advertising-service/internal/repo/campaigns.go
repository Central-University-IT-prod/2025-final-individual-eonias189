package repo

import (
	"advertising/advertising-service/internal/dto"
	"advertising/advertising-service/internal/models"
	"context"

	"github.com/google/uuid"
)

//go:generate go run github.com/vektra/mockery/v2@v2.52.2 --name CampaignsRepo
type CampaignsRepo interface {
	CreateCampaign(ctx context.Context, advertiserId uuid.UUID, data dto.CampaignData) (uuid.UUID, error)
	ListCampaignsForAdvertiser(ctx context.Context, advertiserId uuid.UUID, params dto.PaginationParams) ([]models.Campaign, error)
	GetCampaignById(ctx context.Context, campaignId uuid.UUID) (models.Campaign, error)
	UpdateCampaign(ctx context.Context, campaignId uuid.UUID, data dto.CampaignData) error
	SetCampaignAdImageUrl(ctx context.Context, campaignId uuid.UUID, adImageUrl *string) error
	DeleteCampaign(ctx context.Context, campaignId uuid.UUID) error
}
