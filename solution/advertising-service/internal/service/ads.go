package service

import (
	"advertising/advertising-service/internal/models"
	"advertising/advertising-service/internal/repo"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type AdsService struct {
	adsRepo           repo.AdsRepo
	clientsRepo       repo.ClientsRepo
	campaignsRepo     repo.CampaignsRepo
	clientActionsRepo repo.ClientActionsRepo
	timeRepo          repo.TimeRepo
}

func NewAdsService(
	adsRepo repo.AdsRepo,
	clientsRepo repo.ClientsRepo,
	campaignsRepo repo.CampaignsRepo,
	clientActionsRepo repo.ClientActionsRepo,
	timeRepo repo.TimeRepo,
) *AdsService {
	return &AdsService{
		adsRepo:           adsRepo,
		clientsRepo:       clientsRepo,
		campaignsRepo:     campaignsRepo,
		clientActionsRepo: clientActionsRepo,
		timeRepo:          timeRepo,
	}
}

func (as *AdsService) GetAdForClient(ctx context.Context, clientId uuid.UUID) (models.Ad, error) {
	op := "AdsService.GetAdForClient"

	currentDay, err := as.timeRepo.GetDay(ctx)
	if err != nil {
		return models.Ad{}, fmt.Errorf("%s: timeRepo.GetDay: %w", op, err)
	}

	client, err := as.clientsRepo.GetClientById(ctx, clientId)
	if err != nil {
		return models.Ad{}, fmt.Errorf("%s: clientsRepo.GetClientById: %w", op, err)
	}

	ad, err := as.adsRepo.GetAdForClient(ctx, client, currentDay)
	if err != nil {
		return models.Ad{}, fmt.Errorf("%s: adsRepo.GetAdForClient: %w", op, err)
	}

	campaign, err := as.campaignsRepo.GetCampaignById(ctx, ad.CampaignId)
	if err != nil {
		return models.Ad{}, fmt.Errorf("%s: campaignsRepo.GetCampaignById: %w", op, err)
	}

	impression := models.Impression{
		ClientId:   clientId,
		CampaignId: ad.CampaignId,
		Date:       currentDay,
		Profit:     campaign.CostPerImpression,
	}

	err = as.clientActionsRepo.RecordImpression(ctx, impression)
	if err != nil {
		return models.Ad{}, fmt.Errorf("%s: clientActionsRepo.RecordImpression: %w", op, err)
	}

	return ad, nil
}

func (as *AdsService) RecordAdClick(ctx context.Context, clientId uuid.UUID, campaignId uuid.UUID) error {
	op := "AdsService.RecordAdClick"

	currentDay, err := as.timeRepo.GetDay(ctx)
	if err != nil {
		return fmt.Errorf("%s: timeRepo.GetDay: %w", op, err)
	}

	campaign, err := as.campaignsRepo.GetCampaignById(ctx, campaignId)
	if err != nil {
		return fmt.Errorf("%s: campaignsRepo.GetCampaignById: %w", op, err)
	}

	impressed, err := as.clientActionsRepo.CheckImpressed(ctx, clientId, campaignId)
	if err != nil {
		return fmt.Errorf("%s: clientsActionsRepo.CheckImpressed: %w", op, err)
	}

	if !impressed {
		return models.ErrNotImpressed
	}

	click := models.Click{
		ClientId:   clientId,
		CampaignId: campaignId,
		Date:       currentDay,
		Profit:     campaign.CostPerClick,
	}

	err = as.clientActionsRepo.RecordClick(ctx, click)
	if err != nil {
		if errors.Is(err, models.ErrAlreadyClicked) {
			return nil
		}
		return fmt.Errorf("%s: clientActionsRepo.RecordClick: %w", op, err)
	}

	return nil
}
