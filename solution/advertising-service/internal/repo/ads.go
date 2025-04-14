package repo

import (
	"advertising/advertising-service/internal/models"
	"context"
)

//go:generate go run github.com/vektra/mockery/v2@v2.52.2 --name AdsRepo
type AdsRepo interface {
	GetAdForClient(ctx context.Context, client models.Client, currentDay int) (models.Ad, error)
}
