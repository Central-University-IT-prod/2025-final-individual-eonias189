package handlers

import (
	"advertising/advertising-service/internal/models"
	"advertising/pkg/logger"
	api "advertising/pkg/ogen/advertising-service"
	"context"
	"errors"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AdsUsecase interface {
	GetAdForClient(ctx context.Context, clientId uuid.UUID) (models.Ad, error)
	RecordAdClick(ctx context.Context, clientId uuid.UUID, campaignId uuid.UUID) error
}

type AdsHandler struct {
	au AdsUsecase
}

func NewAdsHandler(au AdsUsecase) *AdsHandler {
	return &AdsHandler{
		au: au,
	}
}

// GetAdForClient implements getAdForClient operation.
//
// Возвращает рекламное объявление, подходящее для
// показа клиенту с учетом таргетинга и ML скора.
//
// GET /ads
func (ah *AdsHandler) GetAdForClient(ctx context.Context, params api.GetAdForClientParams) (api.GetAdForClientRes, error) {
	ad, err := ah.au.GetAdForClient(ctx, params.ClientID)
	if err != nil {
		if errors.Is(err, models.ErrClientNotFound) {
			return &api.Response404{
				Resource: api.ResourceEnumClient,
			}, nil
		}

		if errors.Is(err, models.ErrNoAdsForClient) {
			return &api.Response404{
				Resource: api.ResourceEnumAd,
			}, nil
		}

		logger.FromCtx(ctx).Error("get ad for client", zap.Error(err))
		return nil, err
	}

	res := api.Ad{
		AdID:         ad.CampaignId,
		AdvertiserID: ad.AdvertiserId,
		AdTitle:      ad.AdTitle,
		AdText:       ad.AdText,
	}

	if ad.AdImageUrl != nil {
		res.AdImageURL = api.NewOptString(*ad.AdImageUrl)
	}

	return &res, nil
}

// RecordAdClick implements recordAdClick operation.
//
// Фиксирует клик (переход) клиента по рекламному
// объявлению.
//
// POST /ads/{adId}/click
func (ah *AdsHandler) RecordAdClick(ctx context.Context, req *api.RecordAdClickReq, params api.RecordAdClickParams) (api.RecordAdClickRes, error) {
	err := ah.au.RecordAdClick(ctx, req.GetClientID(), params.AdId)
	if err != nil {
		if errors.Is(err, models.ErrClientNotFound) {
			return &api.Response404{
				Resource: api.ResourceEnumClient,
			}, nil
		}
		if errors.Is(err, models.ErrCampaignNotFound) || errors.Is(err, models.ErrNotImpressed) {
			return &api.Response404{
				Resource: api.ResourceEnumAd,
			}, nil
		}

		logger.FromCtx(ctx).Error("record ad click", zap.Error(err))
		return nil, err
	}

	return &api.RecordAdClickNoContent{}, nil
}
