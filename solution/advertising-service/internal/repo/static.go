package repo

import (
	"advertising/advertising-service/internal/models"
	"context"
)

//go:generate go run github.com/vektra/mockery/v2@v2.52.2 --name StaticRepo
type StaticRepo interface {
	SaveStatic(ctx context.Context, name string, static models.Static) error
	LoadStatic(ctx context.Context, name string) (models.Static, error)
	DeleteStatic(ctx context.Context, name string) error
}
