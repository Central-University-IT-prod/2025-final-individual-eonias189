package service

import (
	"advertising/advertising-service/internal/dto"
	"advertising/advertising-service/internal/models"
	"advertising/advertising-service/internal/repo"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type CampaignsService struct {
	cr            repo.CampaignsRepo
	ar            repo.AdvertisersRepo
	tr            repo.TimeRepo
	sr            repo.StaticRepo
	staticBaseUrl string
}

func NewCampaignsService(
	cr repo.CampaignsRepo,
	ar repo.AdvertisersRepo,
	tr repo.TimeRepo,
	sr repo.StaticRepo,
	staticBaseUrl string,
) *CampaignsService {
	return &CampaignsService{
		cr:            cr,
		ar:            ar,
		tr:            tr,
		sr:            sr,
		staticBaseUrl: staticBaseUrl,
	}
}

func (cs *CampaignsService) CreateCampaign(ctx context.Context, advertiserId uuid.UUID, data dto.CampaignData) (models.Campaign, error) {
	op := "CampaignsService.CreateCampaign"

	dayNow, err := cs.tr.GetDay(ctx)
	if err != nil {
		return models.Campaign{}, fmt.Errorf("%s: tr.GetDay: %w", op, err)
	}

	if data.StartDate < dayNow {
		return models.Campaign{}, models.ErrInvalidStartDate
	}

	createdId, err := cs.cr.CreateCampaign(ctx, advertiserId, data)
	if err != nil {
		return models.Campaign{}, fmt.Errorf("%s: cr.CreateCampaign: %w", op, err)
	}

	campaign := data.ToCampaign()
	campaign.AdvertiserId = advertiserId
	campaign.Id = createdId
	return campaign, nil
}

func (cs *CampaignsService) ListCampaignsForAdvertiser(
	ctx context.Context,
	advertiserId uuid.UUID,
	params dto.PaginationParams,
) ([]models.Campaign, error) {
	op := "CampaignsService.ListCampaignsForAdvertiser"

	// check advertiser existence
	_, err := cs.ar.GetAdvertiserById(ctx, advertiserId)
	if err != nil {
		return nil, fmt.Errorf("%s: ar.GetAdvertiserById: %w", op, err)
	}

	campaigns, err := cs.cr.ListCampaignsForAdvertiser(ctx, advertiserId, params)
	if err != nil {
		return nil, fmt.Errorf("%s: cr.ListCampaignsForAdvertiser: %w", op, err)
	}

	return campaigns, nil
}

func (cs *CampaignsService) GetCampaignById(
	ctx context.Context,
	advertiserId uuid.UUID,
	campaignId uuid.UUID,
) (models.Campaign, error) {
	op := "CampaignsService.GetCampaignById"

	// check advertiser existence
	_, err := cs.ar.GetAdvertiserById(ctx, advertiserId)
	if err != nil {
		return models.Campaign{}, fmt.Errorf("%s: ar.GetAdvertiserById: %w", op, err)
	}

	campaign, err := cs.cr.GetCampaignById(ctx, campaignId)
	if err != nil {
		return models.Campaign{}, fmt.Errorf("%s: cr.GetCampaignById: %w", op, err)
	}

	if campaign.AdvertiserId != advertiserId {
		return models.Campaign{}, models.ErrCampaignNotFound
	}

	return campaign, nil
}

func (cs *CampaignsService) UpdateCampaign(ctx context.Context, advertiserId, campaignId uuid.UUID, data dto.CampaignData) (models.Campaign, error) {
	op := "CampaignsService.UpdateCampaign"

	dayNow, err := cs.tr.GetDay(ctx)
	if err != nil {
		return models.Campaign{}, fmt.Errorf("%s: tr.GetDay: %w", op, err)
	}

	// check advertiser existence
	_, err = cs.ar.GetAdvertiserById(ctx, advertiserId)
	if err != nil {
		return models.Campaign{}, fmt.Errorf("%s: ar.GetAdvertiserById: %w", op, err)
	}

	campaignWas, err := cs.cr.GetCampaignById(ctx, campaignId)
	if err != nil {
		return models.Campaign{}, fmt.Errorf("%s: cr.GetCampaignById: %w", op, err)
	}

	if campaignWas.AdvertiserId != advertiserId {
		return models.Campaign{}, models.ErrCampaignNotFound
	}

	if campaignWas.StartDate <= dayNow {
		if data.ImpressionsLimit != campaignWas.ImpressionsLimit ||
			data.ClicksLimit != campaignWas.ClicksLimit ||
			data.StartDate != campaignWas.StartDate ||
			data.EndDate != campaignWas.EndDate {
			return models.Campaign{}, models.ErrCantUpdateCampaign
		}
	} else {
		if data.StartDate < dayNow {
			return models.Campaign{}, models.ErrInvalidStartDate
		}
	}

	err = cs.cr.UpdateCampaign(ctx, campaignId, data)
	if err != nil {
		return models.Campaign{}, fmt.Errorf("%s: cr.UpdateCampaign: %w", op, err)
	}

	res := data.ToCampaign()
	res.AdvertiserId = advertiserId
	res.Id = campaignId

	return res, nil
}

func (cs *CampaignsService) DeleteCampaign(ctx context.Context, advertiserId uuid.UUID, campaignId uuid.UUID) error {
	op := "CampaignsService.DeleteCampaign"

	// check advertiser existence
	_, err := cs.ar.GetAdvertiserById(ctx, advertiserId)
	if err != nil {
		return fmt.Errorf("%s: ar.GetAdvertiserById: %w", op, err)
	}

	campaignWas, err := cs.cr.GetCampaignById(ctx, campaignId)
	if err != nil {
		return fmt.Errorf("%s: cr.GetCampaignById: %w", op, err)
	}

	if campaignWas.AdvertiserId != advertiserId {
		return models.ErrCampaignNotFound
	}

	err = cs.cr.DeleteCampaign(ctx, campaignId)
	if err != nil {
		return fmt.Errorf("%s: cr.DeleteCampaign: %w", op, err)
	}

	return nil
}

func (cs *CampaignsService) UploadCampaignImage(ctx context.Context, advertiserId, campaignId uuid.UUID, image models.Static) (*string, error) {
	op := "CampaignsService.UploadCampaignImage"

	// check advertiser existence
	_, err := cs.ar.GetAdvertiserById(ctx, advertiserId)
	if err != nil {
		return nil, fmt.Errorf("%s: ar.GetAdvertiserById: %w", op, err)
	}

	// check campaign exists
	campaignWas, err := cs.cr.GetCampaignById(ctx, campaignId)
	if err != nil {
		return nil, fmt.Errorf("%s: cr.GetCampaignById: %w", op, err)
	}

	if campaignWas.AdvertiserId != advertiserId {
		return nil, models.ErrCampaignNotFound
	}

	name := getCampaignImageName(campaignId)

	// if no image uploaded, then remove image
	if image.Size == 0 {
		if err := cs.sr.DeleteStatic(ctx, name); err != nil {
			return nil, fmt.Errorf("%s: sr.DeleteStatic: %w", op, err)
		}

		if err := cs.cr.SetCampaignAdImageUrl(ctx, campaignId, nil); err != nil {
			return nil, fmt.Errorf("%s: cr.SetCampaignAdImageUrl: %w", op, err)
		}

		return nil, nil
	}

	if err := cs.sr.SaveStatic(ctx, name, image); err != nil {
		return nil, fmt.Errorf("%s: sr.SaveStatic: %w", op, err)
	}

	url := fmt.Sprintf("%s/%s", cs.staticBaseUrl, name)

	if err := cs.cr.SetCampaignAdImageUrl(ctx, campaignId, &url); err != nil {
		return nil, fmt.Errorf("%s: cr.SetCampaignAdImageUrl: %w", op, err)
	}

	return &url, nil

}

func getCampaignImageName(campaignId uuid.UUID) string {
	return fmt.Sprintf("campaign-%s-image", campaignId)
}
