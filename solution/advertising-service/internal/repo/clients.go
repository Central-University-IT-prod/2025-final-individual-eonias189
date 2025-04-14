package repo

import (
	"advertising/advertising-service/internal/models"
	"context"

	"github.com/google/uuid"
)

//go:generate go run github.com/vektra/mockery/v2@v2.52.2 --name ClientsRepo
type ClientsRepo interface {
	GetClientById(ctx context.Context, id uuid.UUID) (models.Client, error)
	UpsertClients(ctx context.Context, clients []models.Client) ([]models.Client, error)
}
