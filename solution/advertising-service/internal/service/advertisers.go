package service

import (
	"advertising/advertising-service/internal/models"
	"advertising/advertising-service/internal/repo"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type AdvertiserService struct {
	ar  repo.AdvertisersRepo
	msr repo.MlScoresRepo
}

func NewAdvertisersService(ar repo.AdvertisersRepo, msr repo.MlScoresRepo) *AdvertiserService {
	return &AdvertiserService{
		ar:  ar,
		msr: msr,
	}
}

func (as *AdvertiserService) GetAdvertiserById(ctx context.Context, id uuid.UUID) (models.Advertiser, error) {
	op := "AdvertiserService.GetAdvertiserById"

	advertiser, err := as.ar.GetAdvertiserById(ctx, id)
	if err != nil {
		return models.Advertiser{}, fmt.Errorf("%s: ar.GetAdvertiserById: %w", op, err)
	}

	return advertiser, nil
}

func (as *AdvertiserService) UpsertAdvertisers(ctx context.Context, advertisers []models.Advertiser) ([]models.Advertiser, error) {
	op := "AdvertiserService.UpsertAdvertisers"

	advertisersGot, err := as.ar.UpsertAdvertisers(ctx, advertisers)
	if err != nil {
		return nil, fmt.Errorf("%s: ar.UpsertAdvertisers: %w", op, err)
	}

	return advertisersGot, nil
}

func (as *AdvertiserService) UpsertMLScore(ctx context.Context, mlScore models.MLScore) error {
	op := "AdvertiserService.UpsertMLScore"

	err := as.msr.UpsertMLScore(ctx, mlScore)
	if err != nil {
		return fmt.Errorf("%s: msr.UpsertMLScore: %w", op, err)
	}

	return nil
}
