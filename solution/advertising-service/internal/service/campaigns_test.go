package service

import (
	"advertising/advertising-service/internal/dto"
	"advertising/advertising-service/internal/models"
	"advertising/advertising-service/internal/repo/mocks"
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCampaignsService(t *testing.T) {

	campaignDataSample := dto.CampaignData{
		ImpressionsLimit:  1000,
		ClicksLimit:       100,
		CostPerImpression: 100,
		CostPerClick:      100,
		AdTitle:           "ad title",
		AdText:            "ad text",
		StartDate:         5,
		EndDate:           10,
	}

	t.Run("create campaign success", func(t *testing.T) {
		ctx := context.Background()

		campaignsRepoMock := mocks.NewCampaignsRepo(t)
		advertisersRepoMock := mocks.NewAdvertisersRepo(t)
		timeRepoMock := mocks.NewTimeRepo(t)
		staticRepoMock := mocks.NewStaticRepo(t)

		service := NewCampaignsService(campaignsRepoMock, advertisersRepoMock, timeRepoMock, staticRepoMock, "http://localhost:8080/static")

		// setup mocks
		timeRepoMock.On("GetDay", ctx).Return(0, nil).Once()

		advertiserId := uuid.New()
		campaignData := campaignDataSample

		campaignId := uuid.New()
		expectedCampaign := campaignData.ToCampaign()
		expectedCampaign.Id = campaignId
		expectedCampaign.AdvertiserId = advertiserId

		campaignsRepoMock.On("CreateCampaign", ctx, advertiserId, campaignData).Return(campaignId, nil).Once()

		// check
		actualCampaign, err := service.CreateCampaign(ctx, advertiserId, campaignData)
		require.NoError(t, err)
		require.Equal(t, expectedCampaign, actualCampaign)
	})

	t.Run("create campaign time repo error", func(t *testing.T) {
		ctx := context.Background()

		campaignsRepoMock := mocks.NewCampaignsRepo(t)
		advertisersRepoMock := mocks.NewAdvertisersRepo(t)
		timeRepoMock := mocks.NewTimeRepo(t)
		staticRepoMock := mocks.NewStaticRepo(t)

		service := NewCampaignsService(campaignsRepoMock, advertisersRepoMock, timeRepoMock, staticRepoMock, "http://localhost:8080/static")

		// setup mocks
		expectedError := errors.New("failed to get time")
		timeRepoMock.On("GetDay", ctx).Return(0, expectedError).Once()

		// check
		actualCampaign, err := service.CreateCampaign(ctx, uuid.New(), campaignDataSample)
		require.ErrorIs(t, err, expectedError)
		require.Equal(t, models.Campaign{}, actualCampaign)
	})

	t.Run("create campaign start_date < current day", func(t *testing.T) {
		ctx := context.Background()

		campaignsRepoMock := mocks.NewCampaignsRepo(t)
		advertisersRepoMock := mocks.NewAdvertisersRepo(t)
		timeRepoMock := mocks.NewTimeRepo(t)
		staticRepoMock := mocks.NewStaticRepo(t)

		service := NewCampaignsService(campaignsRepoMock, advertisersRepoMock, timeRepoMock, staticRepoMock, "http://localhost:8080/static")

		// setup mocks
		timeRepoMock.On("GetDay", ctx).Return(5, nil).Once()

		advertiserId := uuid.New()
		campaignData := campaignDataSample
		campaignData.StartDate = 0

		// check
		actualCampaign, err := service.CreateCampaign(ctx, advertiserId, campaignData)
		require.ErrorIs(t, err, models.ErrInvalidStartDate)
		require.Equal(t, models.Campaign{}, actualCampaign)

	})

	t.Run("create campaign campaign repo error", func(t *testing.T) {
		ctx := context.Background()

		campaignsRepoMock := mocks.NewCampaignsRepo(t)
		advertisersRepoMock := mocks.NewAdvertisersRepo(t)
		timeRepoMock := mocks.NewTimeRepo(t)
		staticRepoMock := mocks.NewStaticRepo(t)

		service := NewCampaignsService(campaignsRepoMock, advertisersRepoMock, timeRepoMock, staticRepoMock, "http://localhost:8080/static")

		// setup mocks
		timeRepoMock.On("GetDay", ctx).Return(0, nil)

		advertiserId := uuid.New()
		campaignData := campaignDataSample
		expectedError := errors.New("failed to create campaign")
		expectedCampaign := models.Campaign{}

		campaignsRepoMock.On("CreateCampaign", ctx, advertiserId, campaignData).Return(uuid.New(), expectedError).Once()

		// check
		actualCampaign, err := service.CreateCampaign(ctx, advertiserId, campaignData)
		require.ErrorIs(t, err, expectedError)
		require.Equal(t, expectedCampaign, actualCampaign)
	})

	t.Run("get campaign by id success", func(t *testing.T) {
		ctx := context.Background()

		campaignsRepoMock := mocks.NewCampaignsRepo(t)
		advertisersRepoMock := mocks.NewAdvertisersRepo(t)
		timeRepoMock := mocks.NewTimeRepo(t)
		staticRepoMock := mocks.NewStaticRepo(t)

		service := NewCampaignsService(campaignsRepoMock, advertisersRepoMock, timeRepoMock, staticRepoMock, "http://localhost:8080/static")

		// setup mocks
		advertiserId := uuid.New()
		advertisersRepoMock.On("GetAdvertiserById", ctx, advertiserId).Return(models.Advertiser{
			Id:   advertiserId,
			Name: "name",
		}, nil).Once()

		campaignId := uuid.New()
		expectedCampaign := campaignDataSample.ToCampaign()
		expectedCampaign.Id = campaignId
		expectedCampaign.AdvertiserId = advertiserId

		campaignsRepoMock.On("GetCampaignById", ctx, campaignId).Return(expectedCampaign, nil).Once()

		// check
		actualCampaign, err := service.GetCampaignById(ctx, advertiserId, campaignId)
		require.NoError(t, err)
		require.Equal(t, expectedCampaign, actualCampaign)
	})

	t.Run("get campaign by id advertiser repo error", func(t *testing.T) {
		ctx := context.Background()

		campaignsRepoMock := mocks.NewCampaignsRepo(t)
		advertisersRepoMock := mocks.NewAdvertisersRepo(t)
		timeRepoMock := mocks.NewTimeRepo(t)
		staticRepoMock := mocks.NewStaticRepo(t)

		service := NewCampaignsService(campaignsRepoMock, advertisersRepoMock, timeRepoMock, staticRepoMock, "http://localhost:8080/static")

		// setup mocks
		advertiserId := uuid.New()
		campaignId := uuid.New()
		expectedError := errors.New("failed to get advertiser")

		advertisersRepoMock.On("GetAdvertiserById", ctx, advertiserId).Return(models.Advertiser{}, expectedError).Once()

		// check
		actualCampaign, err := service.GetCampaignById(ctx, advertiserId, campaignId)
		require.ErrorIs(t, err, expectedError)
		require.Equal(t, models.Campaign{}, actualCampaign)
	})

	t.Run("get campaign by id campaign repo error", func(t *testing.T) {
		ctx := context.Background()

		campaignsRepoMock := mocks.NewCampaignsRepo(t)
		advertisersRepoMock := mocks.NewAdvertisersRepo(t)
		timeRepoMock := mocks.NewTimeRepo(t)
		staticRepoMock := mocks.NewStaticRepo(t)

		service := NewCampaignsService(campaignsRepoMock, advertisersRepoMock, timeRepoMock, staticRepoMock, "http://localhost:8080/static")

		// setup mocks
		advertiserId := uuid.New()
		advertisersRepoMock.On("GetAdvertiserById", ctx, advertiserId).Return(models.Advertiser{
			Id:   advertiserId,
			Name: "name",
		}, nil).Once()

		campaignId := uuid.New()
		expectedError := errors.New("failed to get campaign")
		campaignsRepoMock.On("GetCampaignById", ctx, campaignId).Return(models.Campaign{}, expectedError).Once()

		// check
		actualCampaign, err := service.GetCampaignById(ctx, advertiserId, campaignId)
		require.ErrorIs(t, err, expectedError)
		require.Equal(t, models.Campaign{}, actualCampaign)
	})

	t.Run("get campaign by id with campaign.AdvertiserId != advertiserId", func(t *testing.T) {
		ctx := context.Background()

		campaignsRepoMock := mocks.NewCampaignsRepo(t)
		advertisersRepoMock := mocks.NewAdvertisersRepo(t)
		timeRepoMock := mocks.NewTimeRepo(t)
		staticRepoMock := mocks.NewStaticRepo(t)

		service := NewCampaignsService(campaignsRepoMock, advertisersRepoMock, timeRepoMock, staticRepoMock, "http://localhost:8080/static")

		// setup mocks
		advertiserId := uuid.New()
		advertisersRepoMock.On("GetAdvertiserById", ctx, advertiserId).Return(models.Advertiser{
			Id:   advertiserId,
			Name: "name",
		}, nil).Once()

		campaignId := uuid.New()
		campaign := campaignDataSample.ToCampaign()
		campaign.Id = campaignId
		campaign.AdvertiserId = uuid.New()

		campaignsRepoMock.On("GetCampaignById", ctx, campaignId).Return(campaign, nil).Once()

		// check
		actualCampaign, err := service.GetCampaignById(ctx, advertiserId, campaignId)
		require.ErrorIs(t, err, models.ErrCampaignNotFound)
		require.Equal(t, models.Campaign{}, actualCampaign)
	})

	t.Run("list campaigns success", func(t *testing.T) {
		ctx := context.Background()

		campaignsRepoMock := mocks.NewCampaignsRepo(t)
		advertisersRepoMock := mocks.NewAdvertisersRepo(t)
		timeRepoMock := mocks.NewTimeRepo(t)
		staticRepoMock := mocks.NewStaticRepo(t)

		service := NewCampaignsService(campaignsRepoMock, advertisersRepoMock, timeRepoMock, staticRepoMock, "http://localhost:8080/static")

		// setup mocks
		advertiserId := uuid.New()
		advertisersRepoMock.On("GetAdvertiserById", ctx, advertiserId).Return(models.Advertiser{
			Id:   advertiserId,
			Name: "name",
		}, nil).Once()

		campaignsNumber := 5
		expectedcampaigns := make([]models.Campaign, 0, 5)
		for range campaignsNumber {
			campaignId := uuid.New()
			campaign := campaignDataSample.ToCampaign()
			campaign.Id = campaignId
			campaign.AdvertiserId = advertiserId
			expectedcampaigns = append(expectedcampaigns, campaign)
		}
		paginationParams := dto.PaginationParams{
			Size: 5,
			Page: 2,
		}

		campaignsRepoMock.On("ListCampaignsForAdvertiser", ctx, advertiserId, paginationParams).Return(expectedcampaigns, nil).Once()

		// check
		actualCampaigns, err := service.ListCampaignsForAdvertiser(ctx, advertiserId, paginationParams)
		require.NoError(t, err)
		require.Equal(t, expectedcampaigns, actualCampaigns)
	})

	t.Run("list campaigns advertisers repo error", func(t *testing.T) {
		ctx := context.Background()

		campaignsRepoMock := mocks.NewCampaignsRepo(t)
		advertisersRepoMock := mocks.NewAdvertisersRepo(t)
		timeRepoMock := mocks.NewTimeRepo(t)
		staticRepoMock := mocks.NewStaticRepo(t)

		service := NewCampaignsService(campaignsRepoMock, advertisersRepoMock, timeRepoMock, staticRepoMock, "http://localhost:8080/static")

		// setup mocks
		advertiserId := uuid.New()
		expectedError := errors.New("failed to get advertiser")
		advertisersRepoMock.On("GetAdvertiserById", ctx, advertiserId).Return(models.Advertiser{}, expectedError).Once()

		// check
		actualCampagins, err := service.ListCampaignsForAdvertiser(ctx, advertiserId, dto.PaginationParams{})
		require.ErrorIs(t, err, expectedError)
		require.Nil(t, actualCampagins)
	})

	t.Run("list advertisers campaigns repo error", func(t *testing.T) {
		ctx := context.Background()

		campaignsRepoMock := mocks.NewCampaignsRepo(t)
		advertisersRepoMock := mocks.NewAdvertisersRepo(t)
		timeRepoMock := mocks.NewTimeRepo(t)
		staticRepoMock := mocks.NewStaticRepo(t)

		service := NewCampaignsService(campaignsRepoMock, advertisersRepoMock, timeRepoMock, staticRepoMock, "http://localhost:8080/static")

		// setup mocks
		advertiserId := uuid.New()
		advertisersRepoMock.On("GetAdvertiserById", ctx, advertiserId).Return(models.Advertiser{
			Id:   advertiserId,
			Name: "name",
		}, nil).Once()

		params := dto.PaginationParams{
			Size: 4,
			Page: 2,
		}
		expectedError := errors.New("failed to list campaigns")
		campaignsRepoMock.On("ListCampaignsForAdvertiser", ctx, advertiserId, params).Return(nil, expectedError).Once()

		// check
		actualCampaigns, err := service.ListCampaignsForAdvertiser(ctx, advertiserId, params)
		require.ErrorIs(t, err, expectedError)
		require.Nil(t, actualCampaigns)
	})

	t.Run("update campaign success", func(t *testing.T) {
		ctx := context.Background()

		campaignsRepoMock := mocks.NewCampaignsRepo(t)
		advertisersRepoMock := mocks.NewAdvertisersRepo(t)
		timeRepoMock := mocks.NewTimeRepo(t)
		staticRepoMock := mocks.NewStaticRepo(t)

		service := NewCampaignsService(campaignsRepoMock, advertisersRepoMock, timeRepoMock, staticRepoMock, "http://localhost:8080/static")

		// setup mocks
		timeRepoMock.On("GetDay", ctx).Return(0, nil).Once()

		advertiserId := uuid.New()
		advertisersRepoMock.On("GetAdvertiserById", ctx, advertiserId).Return(models.Advertiser{
			Id:   advertiserId,
			Name: "name",
		}, nil).Once()

		campaignId := uuid.New()
		campaignWas := campaignDataSample.ToCampaign()
		campaignWas.StartDate = 100
		campaignWas.EndDate = 200
		campaignWas.Id = campaignId
		campaignWas.AdvertiserId = advertiserId

		campaignsRepoMock.On("GetCampaignById", ctx, campaignId).Return(campaignWas, nil).Once()

		updatedData := dto.CampaignData{
			ImpressionsLimit:  500,
			ClicksLimit:       400,
			CostPerImpression: 200,
			CostPerClick:      5395,
			AdTitle:           "new title",
			AdText:            "new text",
			StartDate:         20,
			EndDate:           50,
		}
		expectedCampaign := updatedData.ToCampaign()
		expectedCampaign.Id = campaignId
		expectedCampaign.AdvertiserId = advertiserId

		campaignsRepoMock.On("UpdateCampaign", ctx, campaignId, updatedData).Return(nil).Once()

		// check
		actualCampaign, err := service.UpdateCampaign(ctx, advertiserId, campaignId, updatedData)
		require.NoError(t, err)
		require.Equal(t, expectedCampaign, actualCampaign)
	})

	t.Run("update campaign time repo error", func(t *testing.T) {
		ctx := context.Background()

		campaignsRepoMock := mocks.NewCampaignsRepo(t)
		advertisersRepoMock := mocks.NewAdvertisersRepo(t)
		timeRepoMock := mocks.NewTimeRepo(t)
		staticRepoMock := mocks.NewStaticRepo(t)

		service := NewCampaignsService(campaignsRepoMock, advertisersRepoMock, timeRepoMock, staticRepoMock, "http://localhost:8080/static")

		// setup mocks
		expectedError := errors.New("falied to get time")
		timeRepoMock.On("GetDay", ctx).Return(0, expectedError).Once()

		// check
		actualCampaign, err := service.UpdateCampaign(ctx, uuid.New(), uuid.New(), dto.CampaignData{})
		require.ErrorIs(t, err, expectedError)
		require.Equal(t, models.Campaign{}, actualCampaign)
	})

	t.Run("update campaign advertisers repo error", func(t *testing.T) {
		ctx := context.Background()

		campaignsRepoMock := mocks.NewCampaignsRepo(t)
		advertisersRepoMock := mocks.NewAdvertisersRepo(t)
		timeRepoMock := mocks.NewTimeRepo(t)
		staticRepoMock := mocks.NewStaticRepo(t)

		service := NewCampaignsService(campaignsRepoMock, advertisersRepoMock, timeRepoMock, staticRepoMock, "http://localhost:8080/static")

		// setup mocks
		timeRepoMock.On("GetDay", ctx).Return(0, nil).Once()

		advertiserId := uuid.New()
		expectedError := errors.New("failed to get advertiser")
		advertisersRepoMock.On("GetAdvertiserById", ctx, advertiserId).Return(models.Advertiser{}, expectedError).Once()

		// check
		actualCampaign, err := service.UpdateCampaign(ctx, advertiserId, uuid.New(), dto.CampaignData{})
		require.ErrorIs(t, err, expectedError)
		require.Equal(t, models.Campaign{}, actualCampaign)
	})

	t.Run("update campaign campaigns repo get campaign error", func(t *testing.T) {
		ctx := context.Background()

		campaignsRepoMock := mocks.NewCampaignsRepo(t)
		advertisersRepoMock := mocks.NewAdvertisersRepo(t)
		timeRepoMock := mocks.NewTimeRepo(t)
		staticRepoMock := mocks.NewStaticRepo(t)

		service := NewCampaignsService(campaignsRepoMock, advertisersRepoMock, timeRepoMock, staticRepoMock, "http://localhost:8080/static")

		// setup mocks
		timeRepoMock.On("GetDay", ctx).Return(0, nil).Once()

		advertiserId := uuid.New()
		advertisersRepoMock.On("GetAdvertiserById", ctx, advertiserId).Return(models.Advertiser{
			Id:   advertiserId,
			Name: "name",
		}, nil).Once()

		campaignId := uuid.New()
		expectedError := errors.New("failed to get campaign")

		campaignsRepoMock.On("GetCampaignById", ctx, campaignId).Return(models.Campaign{}, expectedError).Once()

		// check
		actualCampaign, err := service.UpdateCampaign(ctx, advertiserId, campaignId, dto.CampaignData{})
		require.ErrorIs(t, err, expectedError)
		require.Equal(t, models.Campaign{}, actualCampaign)
	})

	t.Run("update campaign campaigns repo update campaign error", func(t *testing.T) {
		ctx := context.Background()

		campaignsRepoMock := mocks.NewCampaignsRepo(t)
		advertisersRepoMock := mocks.NewAdvertisersRepo(t)
		timeRepoMock := mocks.NewTimeRepo(t)
		staticRepoMock := mocks.NewStaticRepo(t)

		service := NewCampaignsService(campaignsRepoMock, advertisersRepoMock, timeRepoMock, staticRepoMock, "http://localhost:8080/static")

		// setup mocks
		timeRepoMock.On("GetDay", ctx).Return(0, nil).Once()

		advertiserId := uuid.New()
		advertisersRepoMock.On("GetAdvertiserById", ctx, advertiserId).Return(models.Advertiser{
			Id:   advertiserId,
			Name: "name",
		}, nil).Once()

		campaignId := uuid.New()
		campaignWas := campaignDataSample.ToCampaign()
		campaignWas.StartDate = 100
		campaignWas.EndDate = 200
		campaignWas.Id = campaignId
		campaignWas.AdvertiserId = advertiserId

		campaignsRepoMock.On("GetCampaignById", ctx, campaignId).Return(campaignWas, nil).Once()

		updatedData := dto.CampaignData{
			ImpressionsLimit:  500,
			ClicksLimit:       400,
			CostPerImpression: 200,
			CostPerClick:      5395,
			AdTitle:           "new title",
			AdText:            "new text",
			StartDate:         20,
			EndDate:           50,
		}
		expectedError := errors.New("failed to update campaign")

		campaignsRepoMock.On("UpdateCampaign", ctx, campaignId, updatedData).Return(expectedError).Once()

		// check
		actualCampaign, err := service.UpdateCampaign(ctx, advertiserId, campaignId, updatedData)
		require.ErrorIs(t, err, expectedError)
		require.Equal(t, models.Campaign{}, actualCampaign)
	})

	t.Run("update campaign campaignWas.AdvertiserId != advertiserId", func(t *testing.T) {
		ctx := context.Background()

		campaignsRepoMock := mocks.NewCampaignsRepo(t)
		advertisersRepoMock := mocks.NewAdvertisersRepo(t)
		timeRepoMock := mocks.NewTimeRepo(t)
		staticRepoMock := mocks.NewStaticRepo(t)

		service := NewCampaignsService(campaignsRepoMock, advertisersRepoMock, timeRepoMock, staticRepoMock, "http://localhost:8080/static")

		// setup mocks
		timeRepoMock.On("GetDay", ctx).Return(0, nil).Once()

		advertiserId := uuid.New()
		advertisersRepoMock.On("GetAdvertiserById", ctx, advertiserId).Return(models.Advertiser{
			Id:   advertiserId,
			Name: "name",
		}, nil).Once()

		campaignId := uuid.New()
		campaignWas := campaignDataSample.ToCampaign()
		campaignWas.StartDate = 100
		campaignWas.EndDate = 200
		campaignWas.Id = campaignId
		campaignWas.AdvertiserId = uuid.New()

		campaignsRepoMock.On("GetCampaignById", ctx, campaignId).Return(campaignWas, nil).Once()

		// check
		actualCampaign, err := service.UpdateCampaign(ctx, advertiserId, campaignId, dto.CampaignData{})
		require.ErrorIs(t, err, models.ErrCampaignNotFound)
		require.Equal(t, models.Campaign{}, actualCampaign)
	})

	t.Run("update started campaign allowed fields", func(t *testing.T) {
		ctx := context.Background()

		campaignsRepoMock := mocks.NewCampaignsRepo(t)
		advertisersRepoMock := mocks.NewAdvertisersRepo(t)
		timeRepoMock := mocks.NewTimeRepo(t)
		staticRepoMock := mocks.NewStaticRepo(t)

		service := NewCampaignsService(campaignsRepoMock, advertisersRepoMock, timeRepoMock, staticRepoMock, "http://localhost:8080/static")

		// setup mocks
		timeRepoMock.On("GetDay", ctx).Return(5, nil).Once()

		advertiserId := uuid.New()
		advertisersRepoMock.On("GetAdvertiserById", ctx, advertiserId).Return(models.Advertiser{
			Id:   advertiserId,
			Name: "name",
		}, nil).Once()

		campaignId := uuid.New()
		campaignWas := campaignDataSample.ToCampaign()
		campaignWas.StartDate = 0
		campaignWas.EndDate = 10
		campaignWas.Id = campaignId
		campaignWas.AdvertiserId = advertiserId

		campaignsRepoMock.On("GetCampaignById", ctx, campaignId).Return(campaignWas, nil).Once()

		updatedData := dto.CampaignData{
			ImpressionsLimit:  campaignWas.ImpressionsLimit,
			ClicksLimit:       campaignWas.ClicksLimit,
			CostPerImpression: 200,
			CostPerClick:      5395,
			AdTitle:           "new title",
			AdText:            "new text",
			StartDate:         campaignWas.StartDate,
			EndDate:           campaignWas.EndDate,
		}
		expectedCampaign := updatedData.ToCampaign()
		expectedCampaign.Id = campaignId
		expectedCampaign.AdvertiserId = advertiserId

		campaignsRepoMock.On("UpdateCampaign", ctx, campaignId, updatedData).Return(nil).Once()

		// check
		actualCampaign, err := service.UpdateCampaign(ctx, advertiserId, campaignId, updatedData)
		require.NoError(t, err)
		require.Equal(t, expectedCampaign, actualCampaign)
	})

	t.Run("update started campaign forbidden fields", func(t *testing.T) {
		ctx := context.Background()

		campaignsRepoMock := mocks.NewCampaignsRepo(t)
		advertisersRepoMock := mocks.NewAdvertisersRepo(t)
		timeRepoMock := mocks.NewTimeRepo(t)
		staticRepoMock := mocks.NewStaticRepo(t)

		service := NewCampaignsService(campaignsRepoMock, advertisersRepoMock, timeRepoMock, staticRepoMock, "http://localhost:8080/static")

		// setup mocks
		timeRepoMock.On("GetDay", ctx).Return(5, nil)

		advertiserId := uuid.New()
		advertisersRepoMock.On("GetAdvertiserById", ctx, advertiserId).Return(models.Advertiser{
			Id:   advertiserId,
			Name: "name",
		}, nil)

		campaignId := uuid.New()
		campaignWas := campaignDataSample.ToCampaign()
		campaignWas.StartDate = 0
		campaignWas.EndDate = 10
		campaignWas.Id = campaignId
		campaignWas.AdvertiserId = advertiserId

		campaignsRepoMock.On("GetCampaignById", ctx, campaignId).Return(campaignWas, nil)

		cases := []dto.CampaignData{
			{
				ImpressionsLimit:  campaignWas.ImpressionsLimit + 1,
				ClicksLimit:       campaignWas.ClicksLimit,
				CostPerImpression: 200,
				CostPerClick:      5395,
				AdTitle:           "new title",
				AdText:            "new text",
				StartDate:         campaignWas.StartDate,
				EndDate:           campaignWas.EndDate,
			},
			{
				ImpressionsLimit:  campaignWas.ImpressionsLimit,
				ClicksLimit:       campaignWas.ClicksLimit + 1,
				CostPerImpression: 200,
				CostPerClick:      5395,
				AdTitle:           "new title",
				AdText:            "new text",
				StartDate:         campaignWas.StartDate,
				EndDate:           campaignWas.EndDate,
			},
			{
				ImpressionsLimit:  campaignWas.ImpressionsLimit,
				ClicksLimit:       campaignWas.ClicksLimit,
				CostPerImpression: 200,
				CostPerClick:      5395,
				AdTitle:           "new title",
				AdText:            "new text",
				StartDate:         campaignWas.StartDate + 1,
				EndDate:           campaignWas.EndDate,
			},
			{
				ImpressionsLimit:  campaignWas.ImpressionsLimit,
				ClicksLimit:       campaignWas.ClicksLimit,
				CostPerImpression: 200,
				CostPerClick:      5395,
				AdTitle:           "new title",
				AdText:            "new text",
				StartDate:         campaignWas.StartDate,
				EndDate:           campaignWas.EndDate + 1,
			},
		}

		// check
		for _, c := range cases {
			actualCampaign, err := service.UpdateCampaign(ctx, advertiserId, campaignId, c)
			require.ErrorIs(t, err, models.ErrCantUpdateCampaign)
			require.Equal(t, models.Campaign{}, actualCampaign)
		}
	})

	t.Run("update not started campaign with new start_date < current date", func(t *testing.T) {
		ctx := context.Background()

		campaignsRepoMock := mocks.NewCampaignsRepo(t)
		advertisersRepoMock := mocks.NewAdvertisersRepo(t)
		timeRepoMock := mocks.NewTimeRepo(t)
		staticRepoMock := mocks.NewStaticRepo(t)

		service := NewCampaignsService(campaignsRepoMock, advertisersRepoMock, timeRepoMock, staticRepoMock, "http://localhost:8080/static")

		// setup mocks
		timeRepoMock.On("GetDay", ctx).Return(5, nil).Once()

		advertiserId := uuid.New()
		advertisersRepoMock.On("GetAdvertiserById", ctx, advertiserId).Return(models.Advertiser{
			Id:   advertiserId,
			Name: "name",
		}, nil).Once()

		campaignId := uuid.New()
		campaignWas := campaignDataSample.ToCampaign()
		campaignWas.StartDate = 100
		campaignWas.EndDate = 200
		campaignWas.Id = campaignId
		campaignWas.AdvertiserId = advertiserId

		campaignsRepoMock.On("GetCampaignById", ctx, campaignId).Return(campaignWas, nil).Once()

		updatedData := dto.CampaignData{
			ImpressionsLimit:  500,
			ClicksLimit:       400,
			CostPerImpression: 200,
			CostPerClick:      5395,
			AdTitle:           "new title",
			AdText:            "new text",
			StartDate:         0,
			EndDate:           10,
		}

		// check
		actualCampaign, err := service.UpdateCampaign(ctx, advertiserId, campaignId, updatedData)
		require.ErrorIs(t, err, models.ErrInvalidStartDate)
		require.Equal(t, models.Campaign{}, actualCampaign)
	})

	t.Run("delete campaign success", func(t *testing.T) {
		ctx := context.Background()

		campaignsRepoMock := mocks.NewCampaignsRepo(t)
		advertisersRepoMock := mocks.NewAdvertisersRepo(t)
		timeRepoMock := mocks.NewTimeRepo(t)
		staticRepoMock := mocks.NewStaticRepo(t)

		service := NewCampaignsService(campaignsRepoMock, advertisersRepoMock, timeRepoMock, staticRepoMock, "http://localhost:8080/static")

		// setup mocks
		advertiserId := uuid.New()
		advertisersRepoMock.On("GetAdvertiserById", ctx, advertiserId).Return(models.Advertiser{
			Id:   advertiserId,
			Name: "name",
		}, nil).Once()

		campaignId := uuid.New()
		campaignWas := campaignDataSample.ToCampaign()
		campaignWas.Id = campaignId
		campaignWas.AdvertiserId = advertiserId

		campaignsRepoMock.On("GetCampaignById", ctx, campaignId).Return(campaignWas, nil)
		campaignsRepoMock.On("DeleteCampaign", ctx, campaignId).Return(nil)

		// check
		err := service.DeleteCampaign(ctx, advertiserId, campaignId)
		require.NoError(t, err)
	})

	t.Run("delete campaign advertisers repo error", func(t *testing.T) {
		ctx := context.Background()

		campaignsRepoMock := mocks.NewCampaignsRepo(t)
		advertisersRepoMock := mocks.NewAdvertisersRepo(t)
		timeRepoMock := mocks.NewTimeRepo(t)
		staticRepoMock := mocks.NewStaticRepo(t)

		service := NewCampaignsService(campaignsRepoMock, advertisersRepoMock, timeRepoMock, staticRepoMock, "http://localhost:8080/static")

		// setup mocks
		advertiserId := uuid.New()
		expectedError := errors.New("failed to get advertiser")
		advertisersRepoMock.On("GetAdvertiserById", ctx, advertiserId).Return(models.Advertiser{}, expectedError).Once()

		// check
		err := service.DeleteCampaign(ctx, advertiserId, uuid.New())
		require.ErrorIs(t, err, expectedError)
	})

	t.Run("delete campaign camaigns repo get campaign error", func(t *testing.T) {
		ctx := context.Background()

		campaignsRepoMock := mocks.NewCampaignsRepo(t)
		advertisersRepoMock := mocks.NewAdvertisersRepo(t)
		timeRepoMock := mocks.NewTimeRepo(t)
		staticRepoMock := mocks.NewStaticRepo(t)

		service := NewCampaignsService(campaignsRepoMock, advertisersRepoMock, timeRepoMock, staticRepoMock, "http://localhost:8080/static")

		// setup mocks
		advertiserId := uuid.New()
		advertisersRepoMock.On("GetAdvertiserById", ctx, advertiserId).Return(models.Advertiser{
			Id:   advertiserId,
			Name: "name",
		}, nil).Once()

		campaignId := uuid.New()
		expectedError := errors.New("failed to get campaign")

		campaignsRepoMock.On("GetCampaignById", ctx, campaignId).Return(models.Campaign{}, expectedError)

		// check
		err := service.DeleteCampaign(ctx, advertiserId, campaignId)
		require.ErrorIs(t, err, expectedError)
	})

	t.Run("delete campaign campaignWas.AdvertiserId != advertiserId", func(t *testing.T) {
		ctx := context.Background()

		campaignsRepoMock := mocks.NewCampaignsRepo(t)
		advertisersRepoMock := mocks.NewAdvertisersRepo(t)
		timeRepoMock := mocks.NewTimeRepo(t)
		staticRepoMock := mocks.NewStaticRepo(t)

		service := NewCampaignsService(campaignsRepoMock, advertisersRepoMock, timeRepoMock, staticRepoMock, "http://localhost:8080/static")

		// setup mocks
		advertiserId := uuid.New()
		advertisersRepoMock.On("GetAdvertiserById", ctx, advertiserId).Return(models.Advertiser{
			Id:   advertiserId,
			Name: "name",
		}, nil).Once()

		campaignId := uuid.New()
		campaignWas := campaignDataSample.ToCampaign()
		campaignWas.Id = campaignId
		campaignWas.AdvertiserId = uuid.New()

		campaignsRepoMock.On("GetCampaignById", ctx, campaignId).Return(campaignWas, nil)

		// check
		err := service.DeleteCampaign(ctx, advertiserId, campaignId)
		require.ErrorIs(t, err, models.ErrCampaignNotFound)
	})
}
