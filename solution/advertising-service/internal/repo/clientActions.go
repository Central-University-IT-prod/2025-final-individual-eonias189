package repo

import (
	"advertising/advertising-service/internal/models"
	"context"

	"github.com/google/uuid"
)

//go:generate go run github.com/vektra/mockery/v2@v2.52.2 --name ClientActionsRepo
type ClientActionsRepo interface {
	RecordImpression(ctx context.Context, impression models.Impression) error
	RecordClick(ctx context.Context, click models.Click) error
	CheckImpressed(ctx context.Context, clientId, campaignId uuid.UUID) (bool, error)
}
