package service

import (
	"advertising/advertising-service/internal/models"
	"advertising/advertising-service/internal/repo"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type StatsService struct {
	sr repo.StatsRepo
	cr repo.CampaignsRepo
	ar repo.AdvertisersRepo
}

func NewStatsService(
	sr repo.StatsRepo,
	cr repo.CampaignsRepo,
	ar repo.AdvertisersRepo,
) *StatsService {
	return &StatsService{
		sr: sr,
		cr: cr,
		ar: ar,
	}
}

func (ss *StatsService) GetStatsForCampaign(ctx context.Context, campaignId uuid.UUID) (models.Stats, error) {
	op := "StatsService.GetStatsForCampaign"

	// check campaign existence
	_, err := ss.cr.GetCampaignById(ctx, campaignId)
	if err != nil {
		return models.Stats{}, fmt.Errorf("%s: cr.GetCampaignById: %w", op, err)
	}

	stats, err := ss.sr.GetStatsForCampaign(ctx, campaignId)
	if err != nil {
		return models.Stats{}, fmt.Errorf("%s: sr.GetStatsForCampaign: %w", op, err)
	}

	return stats, nil
}

func (ss *StatsService) GetStatsForCampaignDaily(ctx context.Context, campaignId uuid.UUID) ([]models.StatsDaily, error) {
	op := "StatsService.GetStatsForCampaignDaily"

	// check campaign existence
	_, err := ss.cr.GetCampaignById(ctx, campaignId)
	if err != nil {
		return nil, fmt.Errorf("%s: cr.GetCampaignById: %w", op, err)
	}

	stats, err := ss.sr.GetStatsForCampaignDaily(ctx, campaignId)
	if err != nil {
		return nil, fmt.Errorf("%s: sr.GetStatsForCampaignDaily: %w", op, err)
	}

	return stats, nil
}

func (ss *StatsService) GetStatsForAdvertiser(ctx context.Context, advertiser uuid.UUID) (models.Stats, error) {
	op := "StatsService.GetStatsForAdvertiser"

	// check advertiser existence
	_, err := ss.ar.GetAdvertiserById(ctx, advertiser)
	if err != nil {
		return models.Stats{}, fmt.Errorf("%s: ar.GetAdvertiserById: %w", op, err)
	}

	stats, err := ss.sr.GetStatsForAdvertiser(ctx, advertiser)
	if err != nil {
		return models.Stats{}, fmt.Errorf("%s: sr.GetStatsForAdvertiser: %w", op, err)
	}

	return stats, nil
}

func (ss *StatsService) GetStatsForAdvertiserDaily(ctx context.Context, advertiser uuid.UUID) ([]models.StatsDaily, error) {
	op := "StatsService.GetStatsForAdvertiserDaily"

	// check advertiser existence
	_, err := ss.ar.GetAdvertiserById(ctx, advertiser)
	if err != nil {
		return nil, fmt.Errorf("%s: ar.GetAdvertiserById: %w", op, err)
	}

	stats, err := ss.sr.GetStatsForAdvertiserDaily(ctx, advertiser)
	if err != nil {
		return nil, fmt.Errorf("%s: sr.GetStatsForAdvertiserDaily: %w", op, err)
	}

	return stats, nil
}
