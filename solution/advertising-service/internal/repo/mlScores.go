package repo

import (
	"advertising/advertising-service/internal/models"
	"context"
)

//go:generate go run github.com/vektra/mockery/v2@v2.52.2 --name MlScoresRepo
type MlScoresRepo interface {
	UpsertMLScore(ctx context.Context, mlScore models.MLScore) error
}
