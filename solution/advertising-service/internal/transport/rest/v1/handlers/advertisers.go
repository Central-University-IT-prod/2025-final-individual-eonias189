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

type AdvertisersUsecase interface {
	GetAdvertiserById(ctx context.Context, id uuid.UUID) (models.Advertiser, error)
	UpsertAdvertisers(ctx context.Context, advertisers []models.Advertiser) ([]models.Advertiser, error)
	UpsertMLScore(ctx context.Context, mlScore models.MLScore) error
}

type AdvertisersHandler struct {
	au AdvertisersUsecase
}

func NewAdvertisersHandler(au AdvertisersUsecase) *AdvertisersHandler {
	return &AdvertisersHandler{
		au: au,
	}
}

// GetAdvertiserById implements getAdvertiserById operation.
//
// Возвращает информацию о рекламодателе по его ID.
//
// GET /advertisers/{advertiserId}
func (ah *AdvertisersHandler) GetAdvertiserById(ctx context.Context, params api.GetAdvertiserByIdParams) (api.GetAdvertiserByIdRes, error) {
	advertiser, err := ah.au.GetAdvertiserById(ctx, params.AdvertiserId)
	if err != nil {
		if errors.Is(err, models.ErrAdvertiserNotFound) {
			return &api.Response404{
				Resource: api.ResourceEnumAdvertiser,
			}, nil
		}

		logger.FromCtx(ctx).Error("get advertiser by id", zap.Error(err))
		return nil, err
	}

	res := modelsAdvertiserToApiAdvertiser(advertiser)
	return &res, nil
}

// UpsertAdvertisers implements upsertAdvertisers operation.
//
// Создаёт новых или обновляет существующих
// рекламодателей.
//
// POST /advertisers/bulk
func (ah *AdvertisersHandler) UpsertAdvertisers(ctx context.Context, req []api.AdvertiserUpsert) (api.UpsertAdvertisersRes, error) {
	advertisers := make([]models.Advertiser, 0, len(req))
	for _, advertiser := range req {
		advertisers = append(advertisers, apiAdvertiserUpsertToModelsAdvertiser(advertiser))
	}

	advertisersGot, err := ah.au.UpsertAdvertisers(ctx, advertisers)
	if err != nil {
		logger.FromCtx(ctx).Error("upsert advertisers", zap.Error(err))
		return nil, err
	}

	res := api.UpsertAdvertisersCreatedApplicationJSON(make([]api.Advertiser, 0, len(advertisersGot)))
	for _, advertiser := range advertisersGot {
		res = append(res, modelsAdvertiserToApiAdvertiser(advertiser))
	}
	return &res, nil
}

// UpsertMLScore implements upsertMLScore operation.
//
// Добавляет или обновляет ML скор для указанной пары
// клиент-рекламодатель.
//
// POST /ml-scores
func (ah *AdvertisersHandler) UpsertMLScore(ctx context.Context, req *api.MLScore) (api.UpsertMLScoreRes, error) {
	err := ah.au.UpsertMLScore(ctx, models.MLScore{
		ClientId:     req.GetClientID(),
		AdvertiserId: req.GetAdvertiserID(),
		Score:        req.GetScore(),
	})
	if err != nil {
		if errors.Is(err, models.ErrClientNotFound) {
			return &api.Response404{
				Resource: api.ResourceEnumClient,
			}, nil
		}
		if errors.Is(err, models.ErrAdvertiserNotFound) {
			return &api.Response404{
				Resource: api.ResourceEnumAdvertiser,
			}, nil
		}

		logger.FromCtx(ctx).Error("upsert ml score", zap.Error(err))
		return nil, err
	}

	return &api.UpsertMLScoreOK{}, nil
}

func modelsAdvertiserToApiAdvertiser(advertiser models.Advertiser) api.Advertiser {
	return api.Advertiser{
		AdvertiserID: advertiser.Id,
		Name:         advertiser.Name,
	}
}

func apiAdvertiserUpsertToModelsAdvertiser(advertiser api.AdvertiserUpsert) models.Advertiser {
	return models.Advertiser{
		Id:   advertiser.GetAdvertiserID(),
		Name: advertiser.GetName(),
	}
}
