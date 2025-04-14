package repo

import (
	"advertising/advertising-service/internal/models"
	"context"

	"github.com/google/uuid"
)

//go:generate go run github.com/vektra/mockery/v2@v2.52.2 --name AdvertisersRepo
type AdvertisersRepo interface {
	GetAdvertiserById(ctx context.Context, id uuid.UUID) (models.Advertiser, error)
	UpsertAdvertisers(ctx context.Context, advertisers []models.Advertiser) ([]models.Advertiser, error)
}
