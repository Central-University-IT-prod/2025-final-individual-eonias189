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

type StatsUsecase interface {
	GetStatsForCampaign(ctx context.Context, campaignId uuid.UUID) (models.Stats, error)
	GetStatsForCampaignDaily(ctx context.Context, campaignId uuid.UUID) ([]models.StatsDaily, error)
	GetStatsForAdvertiser(ctx context.Context, advertiserId uuid.UUID) (models.Stats, error)
	GetStatsForAdvertiserDaily(ctx context.Context, advertiserId uuid.UUID) ([]models.StatsDaily, error)
}

type StatsHandler struct {
	su StatsUsecase
}

func NewStatsHandler(su StatsUsecase) *StatsHandler {
	return &StatsHandler{
		su: su,
	}
}

// GetAdvertiserCampaignsStats implements getAdvertiserCampaignsStats operation.
//
// Возвращает сводную статистику по всем рекламным
// кампаниям, принадлежащим заданному рекламодателю.
//
// GET /stats/advertisers/{advertiserId}/campaigns
func (sh *StatsHandler) GetAdvertiserCampaignsStats(ctx context.Context, params api.GetAdvertiserCampaignsStatsParams) (api.GetAdvertiserCampaignsStatsRes, error) {
	stats, err := sh.su.GetStatsForAdvertiser(ctx, params.AdvertiserId)
	if err != nil {
		if errors.Is(err, models.ErrAdvertiserNotFound) {
			return &api.Response404{
				Resource: api.ResourceEnumAdvertiser,
			}, nil
		}

		logger.FromCtx(ctx).Error("get advertiser campaigns stats", zap.Error(err))
		return nil, err
	}

	res := modelsStatsToApiStats(stats)
	return &res, nil
}

// GetAdvertiserDailyStats implements getAdvertiserDailyStats operation.
//
// Возвращает массив ежедневной сводной статистики по
// всем рекламным кампаниям заданного рекламодателя.
//
// GET /stats/advertisers/{advertiserId}/campaigns/daily
func (sh *StatsHandler) GetAdvertiserDailyStats(ctx context.Context, params api.GetAdvertiserDailyStatsParams) (api.GetAdvertiserDailyStatsRes, error) {
	stats, err := sh.su.GetStatsForAdvertiserDaily(ctx, params.AdvertiserId)
	if err != nil {
		if errors.Is(err, models.ErrAdvertiserNotFound) {
			return &api.Response404{
				Resource: api.ResourceEnumAdvertiser,
			}, nil
		}

		logger.FromCtx(ctx).Error("get advertiser daily stats", zap.Error(err))
		return nil, err
	}

	res := api.GetAdvertiserDailyStatsOKApplicationJSON(modelsStatsDailyToApiDailyStats(stats))
	return &res, nil
}

// GetCampaignDailyStats implements getCampaignDailyStats operation.
//
// Возвращает массив ежедневной статистики для
// указанной рекламной кампании.
//
// GET /stats/campaigns/{campaignId}/daily
func (sh *StatsHandler) GetCampaignDailyStats(ctx context.Context, params api.GetCampaignDailyStatsParams) (api.GetCampaignDailyStatsRes, error) {
	stats, err := sh.su.GetStatsForCampaignDaily(ctx, params.CampaignId)
	if err != nil {
		if errors.Is(err, models.ErrCampaignNotFound) {
			return &api.Response404{
				Resource: api.ResourceEnumCampaign,
			}, nil
		}

		logger.FromCtx(ctx).Error("get campaign daily stats", zap.Error(err))
		return nil, err
	}

	res := api.GetCampaignDailyStatsOKApplicationJSON(modelsStatsDailyToApiDailyStats(stats))
	return &res, nil
}

// GetCampaignStats implements getCampaignStats operation.
//
// Возвращает агрегированную статистику (показы,
// переходы, затраты и конверсию) для заданной рекламной
// кампании.
//
// GET /stats/campaigns/{campaignId}
func (sh *StatsHandler) GetCampaignStats(ctx context.Context, params api.GetCampaignStatsParams) (api.GetCampaignStatsRes, error) {
	stats, err := sh.su.GetStatsForCampaign(ctx, params.CampaignId)
	if err != nil {
		if errors.Is(err, models.ErrCampaignNotFound) {
			return &api.Response404{
				Resource: api.ResourceEnumCampaign,
			}, nil
		}

		logger.FromCtx(ctx).Error("get campaign stats", zap.Error(err))
		return nil, err
	}

	res := modelsStatsToApiStats(stats)
	return &res, nil
}

func modelsStatsToApiStats(stats models.Stats) api.Stats {
	return api.Stats{
		ImpressionsCount: stats.ImpressionsCount,
		ClicksCount:      stats.ClicksCount,
		Conversion:       stats.Conversion,
		SpentImpressions: stats.SpentImpressions,
		SpentClicks:      stats.SpentClicks,
		SpentTotal:       stats.SpentTotal,
	}
}

func modelsStatsDailyToApiDailyStats(statsDaily []models.StatsDaily) []api.DailyStats {
	res := make([]api.DailyStats, 0, len(statsDaily))
	for _, stats := range statsDaily {
		res = append(res, api.DailyStats{
			ImpressionsCount: stats.ImpressionsCount,
			ClicksCount:      stats.ClicksCount,
			Conversion:       stats.Conversion,
			SpentImpressions: stats.SpentImpressions,
			SpentClicks:      stats.SpentClicks,
			SpentTotal:       stats.SpentTotal,
			Date:             api.Date(stats.Date),
		})
	}
	return res
}
