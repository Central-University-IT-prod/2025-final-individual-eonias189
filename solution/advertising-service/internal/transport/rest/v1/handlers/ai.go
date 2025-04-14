package handlers

import (
	"advertising/pkg/logger"
	api "advertising/pkg/ogen/advertising-service"
	"context"

	"go.uber.org/zap"
)

type AIUsecase interface {
	GenerateAdText(ctx context.Context, adTitle string) (string, error)
	ModerateAdText(ctx context.Context, adText string) (bool, []string, error)
}

type AIHandler struct {
	aiu AIUsecase
}

func NewAIHandler(aiu AIUsecase) *AIHandler {
	return &AIHandler{
		aiu: aiu,
	}
}

// GenerateAdText implements generateAdText operation.
//
// Генерация текста для рекламной кампании.
//
// POST /ai/generate-ad-text
func (aih *AIHandler) GenerateAdText(ctx context.Context, req *api.GenerateAdTextReq) (api.GenerateAdTextRes, error) {
	adText, err := aih.aiu.GenerateAdText(ctx, req.GetAdTitle())
	if err != nil {
		logger.FromCtx(ctx).Error("generate ad text", zap.Error(err))
		return nil, err
	}

	return &api.GenerateAdTextOK{
		AdText: adText,
	}, nil
}

// ModerateAdText implements moderateAdText operation.
//
// Модерирует текст рекламного объявления.
//
// POST /ai/moderate-ad-text
func (aih *AIHandler) ModerateAdText(ctx context.Context, req *api.ModerateAdTextReq) (api.ModerateAdTextRes, error) {
	ok, illegalPhrases, err := aih.aiu.ModerateAdText(ctx, req.GetAdText())
	if err != nil {
		logger.FromCtx(ctx).Error("moderate ad text", zap.Error(err))
		return nil, err
	}

	return &api.ModerateAdTextOK{
		Ok:             ok,
		IllegalPhrases: illegalPhrases,
	}, nil
}
