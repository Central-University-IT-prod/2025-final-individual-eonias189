package handlers

import (
	"advertising/advertising-service/internal/dto"
	"advertising/advertising-service/internal/models"
	"advertising/pkg/logger"
	api "advertising/pkg/ogen/advertising-service"
	"bytes"
	"context"
	"errors"
	"io"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type CampaignsUsecase interface {
	CreateCampaign(ctx context.Context, advertiserId uuid.UUID, data dto.CampaignData) (models.Campaign, error)
	ListCampaignsForAdvertiser(ctx context.Context, advertiserId uuid.UUID, params dto.PaginationParams) ([]models.Campaign, error)
	GetCampaignById(ctx context.Context, advertiserId uuid.UUID, campaignId uuid.UUID) (models.Campaign, error)
	UpdateCampaign(ctx context.Context, advertiserId, campaignId uuid.UUID, data dto.CampaignData) (models.Campaign, error)
	DeleteCampaign(ctx context.Context, advertiserId uuid.UUID, campaignId uuid.UUID) error
	UploadCampaignImage(ctx context.Context, advertiserId, campaignId uuid.UUID, image models.Static) (*string, error)
}

type CampaignsHandler struct {
	cu CampaignsUsecase
}

func NewCampaignsHandler(cu CampaignsUsecase) *CampaignsHandler {
	return &CampaignsHandler{
		cu: cu,
	}
}

// GetCampaign implements getCampaign operation.
//
// Получение кампании по ID.
//
// GET /advertisers/{advertiserId}/campaigns/{campaignId}
func (ch *CampaignsHandler) GetCampaign(ctx context.Context, params api.GetCampaignParams) (api.GetCampaignRes, error) {
	campaign, err := ch.cu.GetCampaignById(ctx, params.AdvertiserId, params.CampaignId)
	if err != nil {
		if errors.Is(err, models.ErrAdvertiserNotFound) {
			return &api.Response404{
				Resource: api.ResourceEnumAdvertiser,
			}, nil
		}
		if errors.Is(err, models.ErrCampaignNotFound) {
			return &api.Response404{
				Resource: api.ResourceEnumCampaign,
			}, nil
		}

		logger.FromCtx(ctx).Error("get campaign", zap.Error(err))
		return nil, err
	}

	res := modelsCampaignToApiCampaign(campaign)
	return &res, nil
}

// CreateCampaign implements createCampaign operation.
//
// Создаёт новую рекламную кампанию для указанного
// рекламодателя.
//
// POST /advertisers/{advertiserId}/campaigns
func (ch *CampaignsHandler) CreateCampaign(ctx context.Context, req *api.CampaignCreate, params api.CreateCampaignParams) (api.CreateCampaignRes, error) {
	data := dto.CampaignData{
		ImpressionsLimit:  req.GetImpressionsLimit(),
		ClicksLimit:       req.GetClicksLimit(),
		CostPerImpression: float64(req.GetCostPerImpression()),
		CostPerClick:      float64(req.GetCostPerClick()),
		AdTitle:           req.GetAdTitle(),
		AdText:            req.GetAdText(),
		StartDate:         int(req.GetStartDate()),
		EndDate:           int(req.GetEndDate()),
	}
	if req.GetTargeting().IsSet() {
		targeting := req.GetTargeting().Value

		if targeting.GetAgeFrom().IsSet() &&
			!targeting.GetAgeFrom().IsNull() &&
			targeting.GetAgeTo().IsSet() &&
			!targeting.GetAgeTo().IsNull() {
			if targeting.GetAgeTo().Value < targeting.GetAgeFrom().Value {
				return &api.Response400{
					Message: api.NewOptString("age_to must be not less than age_from"),
				}, nil
			}
		}

		if targeting.GetGender().IsSet() && !targeting.GetGender().IsNull() {
			data.Gender = pointer(models.Gender((targeting.GetGender().Value)))
		}
		if targeting.GetAgeFrom().IsSet() && !targeting.GetAgeFrom().IsNull() {
			data.AgeFrom = pointer(int(targeting.GetAgeFrom().Value))
		}
		if targeting.GetAgeTo().IsSet() && !targeting.GetAgeTo().IsNull() {
			data.AgeTo = pointer(int(targeting.GetAgeTo().Value))
		}
		if targeting.GetLocation().IsSet() && !targeting.GetLocation().IsNull() {
			data.Location = pointer(targeting.GetLocation().Value)
		}
	}

	if req.GetClicksLimit() > req.GetImpressionsLimit() {
		return &api.Response400{
			Message: api.NewOptString("clicks limit must be not greater than impressions_limit"),
		}, nil
	}

	if req.GetEndDate() < req.GetStartDate() {
		return &api.Response400{
			Message: api.NewOptString("end_date must be not less than start_date"),
		}, nil
	}

	created, err := ch.cu.CreateCampaign(ctx, params.AdvertiserId, data)
	if err != nil {
		if errors.Is(err, models.ErrAdvertiserNotFound) {
			return &api.Response404{
				Resource: api.ResourceEnumAdvertiser,
			}, nil
		}
		if errors.Is(err, models.ErrInvalidStartDate) {
			return &api.Response400{
				Message: api.NewOptString("start_date must be not in past"),
			}, nil
		}

		logger.FromCtx(ctx).Error("create campaign", zap.Error(err))
		return nil, err
	}

	res := modelsCampaignToApiCampaign(created)
	return &res, nil
}

// DeleteCampaign implements deleteCampaign operation.
//
// Удаляет рекламную кампанию рекламодателя по
// заданному campaignId.
//
// DELETE /advertisers/{advertiserId}/campaigns/{campaignId}
func (ch *CampaignsHandler) DeleteCampaign(ctx context.Context, params api.DeleteCampaignParams) (api.DeleteCampaignRes, error) {
	err := ch.cu.DeleteCampaign(ctx, params.AdvertiserId, params.CampaignId)
	if err != nil {
		if errors.Is(err, models.ErrAdvertiserNotFound) {
			return &api.Response404{
				Resource: api.ResourceEnumAdvertiser,
			}, nil
		}
		if errors.Is(err, models.ErrCampaignNotFound) {
			return &api.Response404{
				Resource: api.ResourceEnumCampaign,
			}, nil
		}

		logger.FromCtx(ctx).Error("delete campaign", zap.Error(err))
		return nil, err
	}

	return &api.DeleteCampaignNoContent{}, nil
}

// ListCampaigns implements listCampaigns operation.
//
// Возвращает список рекламных кампаний для указанного
// рекламодателя с пагинацией.
//
// GET /advertisers/{advertiserId}/campaigns
func (ch *CampaignsHandler) ListCampaigns(ctx context.Context, params api.ListCampaignsParams) (api.ListCampaignsRes, error) {
	paginateParams := dto.PaginationParams{
		Size: params.Size.Or(50),
		Page: params.Page.Or(1),
	}

	campaigns, err := ch.cu.ListCampaignsForAdvertiser(ctx, params.AdvertiserId, paginateParams)
	if err != nil {
		if errors.Is(err, models.ErrAdvertiserNotFound) {
			return &api.Response404{
				Resource: api.ResourceEnumAdvertiser,
			}, nil
		}

		logger.FromCtx(ctx).Error("list campaigns", zap.Error(err))
		return nil, err
	}

	res := api.ListCampaignsOKApplicationJSON(make(api.ListCampaignsOKApplicationJSON, 0, len(campaigns)))
	for _, campaign := range campaigns {
		res = append(res, modelsCampaignToApiCampaign(campaign))
	}

	return &res, nil
}

// UpdateCampaign implements updateCampaign operation.
//
// Обновляет разрешённые параметры рекламной кампании
// до её старта.
//
// PUT /advertisers/{advertiserId}/campaigns/{campaignId}
func (ch *CampaignsHandler) UpdateCampaign(ctx context.Context, req *api.CampaignUpdate, params api.UpdateCampaignParams) (api.UpdateCampaignRes, error) {
	data := dto.CampaignData{
		ImpressionsLimit:  req.GetImpressionsLimit(),
		ClicksLimit:       req.GetClicksLimit(),
		CostPerImpression: float64(req.GetCostPerImpression()),
		CostPerClick:      float64(req.GetCostPerClick()),
		AdTitle:           req.GetAdTitle(),
		AdText:            req.GetAdText(),
		StartDate:         int(req.GetStartDate()),
		EndDate:           int(req.GetEndDate()),
	}
	if req.GetTargeting().IsSet() {
		targeting := req.GetTargeting().Value

		if targeting.GetAgeFrom().IsSet() &&
			!targeting.GetAgeFrom().IsNull() &&
			targeting.GetAgeTo().IsSet() &&
			!targeting.GetAgeTo().IsNull() {
			if targeting.GetAgeTo().Value < targeting.GetAgeFrom().Value {
				return &api.Response400{
					Message: api.NewOptString("age_to must be not less than age_from"),
				}, nil
			}
		}

		if targeting.GetGender().IsSet() && !targeting.GetGender().IsNull() {
			data.Gender = pointer(models.Gender((targeting.GetGender().Value)))
		}
		if targeting.GetAgeFrom().IsSet() && !targeting.GetAgeFrom().IsNull() {
			data.AgeFrom = pointer(int(targeting.GetAgeFrom().Value))
		}
		if targeting.GetAgeTo().IsSet() && !targeting.GetAgeTo().IsNull() {
			data.AgeTo = pointer(int(targeting.GetAgeTo().Value))
		}
		if targeting.GetLocation().IsSet() && !targeting.GetLocation().IsNull() {
			data.Location = pointer(targeting.GetLocation().Value)
		}
	}

	if req.GetClicksLimit() > req.GetImpressionsLimit() {
		return &api.Response400{
			Message: api.NewOptString("clicks limit must be not greater than impressions_limit"),
		}, nil
	}

	if req.GetEndDate() < req.GetStartDate() {
		return &api.Response400{
			Message: api.NewOptString("end_date must be not less than start_date"),
		}, nil
	}

	campaign, err := ch.cu.UpdateCampaign(ctx, params.AdvertiserId, params.CampaignId, data)
	if err != nil {
		if errors.Is(err, models.ErrAdvertiserNotFound) {
			return &api.Response404{
				Resource: api.ResourceEnumAdvertiser,
			}, nil
		}
		if errors.Is(err, models.ErrCampaignNotFound) {
			return &api.Response404{
				Resource: api.ResourceEnumCampaign,
			}, nil
		}
		if errors.Is(err, models.ErrInvalidStartDate) {
			return &api.Response400{
				Message: api.NewOptString("start_date must be not in past"),
			}, nil
		}
		if errors.Is(err, models.ErrCantUpdateCampaign) {
			return &api.UpdateCampaignForbidden{}, nil
		}

		logger.FromCtx(ctx).Error("update campaign", zap.Error(err))
		return nil, err
	}

	res := modelsCampaignToApiCampaign(campaign)
	return &res, nil
}

// UploadCampaignImage implements uploadCampaignImage operation.
//
// Загружает изображение рекламного объявления. Если не
// передать изображение, то прежнее удалится.
//
// PUT /advertisers/{advertiserId}/campaigns/{campaignId}/image
func (ch *CampaignsHandler) UploadCampaignImage(ctx context.Context, req api.UploadCampaignImageReq, params api.UploadCampaignImageParams) (api.UploadCampaignImageRes, error) {
	var image models.Static
	switch v := req.(type) {
	case *api.UploadCampaignImageReqImagePNG:
		image.Data = v.Data
		image.ContentType = "image/png"
	case *api.UploadCampaignImageReqImageJpeg:
		image.Data = v.Data
		image.ContentType = "image/jpeg"
	}

	var err error
	image.Data, image.Size, err = getReaderSize(image.Data)
	if err != nil {
		logger.FromCtx(ctx).Error("get image size", zap.Error(err))
		return nil, err
	}

	logger.FromCtx(ctx).Debug("image size", zap.Int64("size", image.Size))

	imageUrl, err := ch.cu.UploadCampaignImage(ctx, params.AdvertiserId, params.CampaignId, image)
	if err != nil {
		if errors.Is(err, models.ErrAdvertiserNotFound) {
			return &api.Response404{
				Resource: api.ResourceEnumAdvertiser,
			}, nil
		}
		if errors.Is(err, models.ErrCampaignNotFound) {
			return &api.Response404{
				Resource: api.ResourceEnumCampaign,
			}, nil
		}

		logger.FromCtx(ctx).Error("load campaign image", zap.Error(err))
		return nil, err
	}

	var res api.NilString
	if imageUrl != nil {
		res.SetTo(*imageUrl)
	} else {
		res.SetToNull()
	}

	return &api.UploadCampaignImageOK{
		AdImageURL: res,
	}, nil
}

func modelsCampaignToApiCampaign(campaign models.Campaign) api.Campaign {
	targetting := api.Targeting{}
	if campaign.Gender != nil {
		targetting.Gender = api.NewOptNilTargetingGender(api.TargetingGender(*campaign.Gender))
	} else {
		targetting.Gender.SetToNull()
	}
	if campaign.AgeFrom != nil {
		targetting.AgeFrom = api.NewOptNilInt(*campaign.AgeFrom)
	} else {
		targetting.AgeFrom.SetToNull()
	}
	if campaign.AgeTo != nil {
		targetting.AgeTo = api.NewOptNilInt(*campaign.AgeTo)
	} else {
		targetting.AgeTo.SetToNull()
	}
	if campaign.Location != nil {
		targetting.Location = api.NewOptNilString(*campaign.Location)
	} else {
		targetting.Location.SetToNull()
	}

	res := api.Campaign{
		CampaignID:        campaign.Id,
		AdvertiserID:      campaign.AdvertiserId,
		ImpressionsLimit:  campaign.ImpressionsLimit,
		ClicksLimit:       campaign.ClicksLimit,
		CostPerImpression: float32(campaign.CostPerImpression),
		CostPerClick:      float32(campaign.CostPerClick),
		AdTitle:           campaign.AdTitle,
		AdText:            campaign.AdText,
		StartDate:         api.Date(campaign.StartDate),
		EndDate:           api.Date(campaign.EndDate),
		Targeting:         targetting,
	}

	if campaign.AdImageUrl != nil {
		res.AdImageURL = api.NewOptString(*campaign.AdImageUrl)
	}

	return res
}

func pointer[T any](v T) *T {
	return &v
}

func getReaderSize(r io.Reader) (io.Reader, int64, error) {
	buf := &bytes.Buffer{}
	size, err := io.Copy(buf, r)
	if err != nil {
		return nil, 0, err
	}

	return buf, size, nil
}
