package service

import (
	"advertising/advertising-service/internal/models"
	"advertising/advertising-service/internal/repo/mocks"
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestAdsService(t *testing.T) {
	t.Run("get ad for client success", func(t *testing.T) {
		ctx := context.Background()

		adsRepoMock := mocks.NewAdsRepo(t)
		clientsRepoMock := mocks.NewClientsRepo(t)
		campaignsRepoMock := mocks.NewCampaignsRepo(t)
		clientActionsRepoMock := mocks.NewClientActionsRepo(t)
		timeRepoMock := mocks.NewTimeRepo(t)

		service := NewAdsService(adsRepoMock, clientsRepoMock, campaignsRepoMock, clientActionsRepoMock, timeRepoMock)

		// setup mocks
		currentDay := 5
		timeRepoMock.On("GetDay", ctx).Return(currentDay, nil).Once()

		clientId := uuid.New()
		client := models.Client{Id: clientId}
		clientsRepoMock.On("GetClientById", ctx, clientId).Return(client, nil).Once()

		ad := models.Ad{CampaignId: uuid.New()}
		adsRepoMock.On("GetAdForClient", ctx, client, currentDay).Return(ad, nil).Once()

		campaign := models.Campaign{CostPerImpression: 100}
		campaignsRepoMock.On("GetCampaignById", ctx, ad.CampaignId).Return(campaign, nil).Once()

		clientActionsRepoMock.On("RecordImpression", ctx, models.Impression{
			ClientId:   clientId,
			CampaignId: ad.CampaignId,
			Date:       currentDay,
			Profit:     campaign.CostPerImpression,
		}).Return(nil).Once()

		// check
		actualAd, err := service.GetAdForClient(ctx, clientId)
		require.NoError(t, err)
		require.Equal(t, ad, actualAd)
	})

	t.Run("get ad for client time repo error", func(t *testing.T) {
		ctx := context.Background()

		adsRepoMock := mocks.NewAdsRepo(t)
		clientsRepoMock := mocks.NewClientsRepo(t)
		campaignsRepoMock := mocks.NewCampaignsRepo(t)
		clientActionsRepoMock := mocks.NewClientActionsRepo(t)
		timeRepoMock := mocks.NewTimeRepo(t)

		service := NewAdsService(adsRepoMock, clientsRepoMock, campaignsRepoMock, clientActionsRepoMock, timeRepoMock)

		// setup mocks
		expectedError := errors.New("failed to get time")
		timeRepoMock.On("GetDay", ctx).Return(0, expectedError).Once()

		// check
		actualAd, err := service.GetAdForClient(ctx, uuid.New())
		require.ErrorIs(t, err, expectedError)
		require.Equal(t, models.Ad{}, actualAd)
	})

	t.Run("get ad for client clients repo error", func(t *testing.T) {
		ctx := context.Background()

		adsRepoMock := mocks.NewAdsRepo(t)
		clientsRepoMock := mocks.NewClientsRepo(t)
		campaignsRepoMock := mocks.NewCampaignsRepo(t)
		clientActionsRepoMock := mocks.NewClientActionsRepo(t)
		timeRepoMock := mocks.NewTimeRepo(t)

		service := NewAdsService(adsRepoMock, clientsRepoMock, campaignsRepoMock, clientActionsRepoMock, timeRepoMock)

		// setup mocks
		currentDay := 5
		timeRepoMock.On("GetDay", ctx).Return(currentDay, nil).Once()

		clientId := uuid.New()
		expectedError := errors.New("failed to get client")
		clientsRepoMock.On("GetClientById", ctx, clientId).Return(models.Client{}, expectedError).Once()

		// check
		actualAd, err := service.GetAdForClient(ctx, clientId)
		require.ErrorIs(t, err, expectedError)
		require.Equal(t, models.Ad{}, actualAd)
	})

	t.Run("get ad for client ads repo error", func(t *testing.T) {
		ctx := context.Background()

		adsRepoMock := mocks.NewAdsRepo(t)
		clientsRepoMock := mocks.NewClientsRepo(t)
		campaignsRepoMock := mocks.NewCampaignsRepo(t)
		clientActionsRepoMock := mocks.NewClientActionsRepo(t)
		timeRepoMock := mocks.NewTimeRepo(t)

		service := NewAdsService(adsRepoMock, clientsRepoMock, campaignsRepoMock, clientActionsRepoMock, timeRepoMock)

		// setup mocks
		currentDay := 5
		timeRepoMock.On("GetDay", ctx).Return(currentDay, nil).Once()

		clientId := uuid.New()
		client := models.Client{Id: clientId}
		clientsRepoMock.On("GetClientById", ctx, clientId).Return(client, nil).Once()

		expectedError := errors.New("failed to get ad")
		adsRepoMock.On("GetAdForClient", ctx, client, currentDay).Return(models.Ad{}, expectedError).Once()

		// check
		actualAd, err := service.GetAdForClient(ctx, clientId)
		require.ErrorIs(t, err, expectedError)
		require.Equal(t, models.Ad{}, actualAd)
	})

	t.Run("get ad for client campaigns repo error", func(t *testing.T) {
		ctx := context.Background()

		adsRepoMock := mocks.NewAdsRepo(t)
		clientsRepoMock := mocks.NewClientsRepo(t)
		campaignsRepoMock := mocks.NewCampaignsRepo(t)
		clientActionsRepoMock := mocks.NewClientActionsRepo(t)
		timeRepoMock := mocks.NewTimeRepo(t)

		service := NewAdsService(adsRepoMock, clientsRepoMock, campaignsRepoMock, clientActionsRepoMock, timeRepoMock)

		// setup mocks
		currentDay := 5
		timeRepoMock.On("GetDay", ctx).Return(currentDay, nil).Once()

		clientId := uuid.New()
		client := models.Client{Id: clientId}
		clientsRepoMock.On("GetClientById", ctx, clientId).Return(client, nil).Once()

		ad := models.Ad{CampaignId: uuid.New()}
		adsRepoMock.On("GetAdForClient", ctx, client, currentDay).Return(ad, nil).Once()

		expectedError := errors.New("failed to get campaign")
		campaignsRepoMock.On("GetCampaignById", ctx, ad.CampaignId).Return(models.Campaign{}, expectedError).Once()

		// check
		actualAd, err := service.GetAdForClient(ctx, clientId)
		require.ErrorIs(t, err, expectedError)
		require.Equal(t, models.Ad{}, actualAd)
	})

	t.Run("get ad for client record impression error", func(t *testing.T) {
		ctx := context.Background()

		adsRepoMock := mocks.NewAdsRepo(t)
		clientsRepoMock := mocks.NewClientsRepo(t)
		campaignsRepoMock := mocks.NewCampaignsRepo(t)
		clientActionsRepoMock := mocks.NewClientActionsRepo(t)
		timeRepoMock := mocks.NewTimeRepo(t)

		service := NewAdsService(adsRepoMock, clientsRepoMock, campaignsRepoMock, clientActionsRepoMock, timeRepoMock)

		// setup mocks
		currentDay := 5
		timeRepoMock.On("GetDay", ctx).Return(currentDay, nil).Once()

		clientId := uuid.New()
		client := models.Client{Id: clientId}
		clientsRepoMock.On("GetClientById", ctx, clientId).Return(client, nil).Once()

		ad := models.Ad{CampaignId: uuid.New()}
		adsRepoMock.On("GetAdForClient", ctx, client, currentDay).Return(ad, nil).Once()

		campaign := models.Campaign{CostPerImpression: 100}
		campaignsRepoMock.On("GetCampaignById", ctx, ad.CampaignId).Return(campaign, nil).Once()

		expectedError := errors.New("failed to record impression")
		clientActionsRepoMock.On("RecordImpression", ctx, models.Impression{
			ClientId:   clientId,
			CampaignId: ad.CampaignId,
			Date:       currentDay,
			Profit:     campaign.CostPerImpression,
		}).Return(expectedError).Once()

		// check
		actualAd, err := service.GetAdForClient(ctx, clientId)
		require.ErrorIs(t, err, expectedError)
		require.Equal(t, models.Ad{}, actualAd)
	})

	t.Run("record ad click success", func(t *testing.T) {
		ctx := context.Background()

		adsRepoMock := mocks.NewAdsRepo(t)
		clientsRepoMock := mocks.NewClientsRepo(t)
		campaignsRepoMock := mocks.NewCampaignsRepo(t)
		clientActionsRepoMock := mocks.NewClientActionsRepo(t)
		timeRepoMock := mocks.NewTimeRepo(t)

		service := NewAdsService(adsRepoMock, clientsRepoMock, campaignsRepoMock, clientActionsRepoMock, timeRepoMock)

		// setup mocks
		currentDay := 5
		timeRepoMock.On("GetDay", ctx).Return(currentDay, nil).Once()

		clientId := uuid.New()
		campaignId := uuid.New()

		campaign := models.Campaign{CostPerClick: 100}
		campaignsRepoMock.On("GetCampaignById", ctx, campaignId).Return(campaign, nil).Once()

		clientActionsRepoMock.On("CheckImpressed", ctx, clientId, campaignId).Return(true, nil).Once()

		clientActionsRepoMock.On("RecordClick", ctx, models.Click{
			ClientId:   clientId,
			CampaignId: campaignId,
			Date:       currentDay,
			Profit:     campaign.CostPerClick,
		}).Return(nil).Once()

		// check
		err := service.RecordAdClick(ctx, clientId, campaignId)
		require.NoError(t, err)
	})

	t.Run("record ad click several times", func(t *testing.T) {
		ctx := context.Background()

		adsRepoMock := mocks.NewAdsRepo(t)
		clientsRepoMock := mocks.NewClientsRepo(t)
		campaignsRepoMock := mocks.NewCampaignsRepo(t)
		clientActionsRepoMock := mocks.NewClientActionsRepo(t)
		timeRepoMock := mocks.NewTimeRepo(t)

		service := NewAdsService(adsRepoMock, clientsRepoMock, campaignsRepoMock, clientActionsRepoMock, timeRepoMock)

		// setup mocks
		currentDay := 5
		timeRepoMock.On("GetDay", ctx).Return(currentDay, nil).Once()

		clientId := uuid.New()
		campaignId := uuid.New()

		campaign := models.Campaign{CostPerClick: 100}
		campaignsRepoMock.On("GetCampaignById", ctx, campaignId).Return(campaign, nil).Once()

		clientActionsRepoMock.On("CheckImpressed", ctx, clientId, campaignId).Return(true, nil).Once()

		clientActionsRepoMock.On("RecordClick", ctx, models.Click{
			ClientId:   clientId,
			CampaignId: campaignId,
			Date:       currentDay,
			Profit:     campaign.CostPerClick,
		}).Return(models.ErrAlreadyClicked).Once()

		// check
		err := service.RecordAdClick(ctx, clientId, campaignId)
		require.NoError(t, err)
	})

	t.Run("record ad click time repo error", func(t *testing.T) {
		ctx := context.Background()

		adsRepoMock := mocks.NewAdsRepo(t)
		clientsRepoMock := mocks.NewClientsRepo(t)
		campaignsRepoMock := mocks.NewCampaignsRepo(t)
		clientActionsRepoMock := mocks.NewClientActionsRepo(t)
		timeRepoMock := mocks.NewTimeRepo(t)

		service := NewAdsService(adsRepoMock, clientsRepoMock, campaignsRepoMock, clientActionsRepoMock, timeRepoMock)

		// setup mocks
		expectedError := errors.New("failed to get time")
		timeRepoMock.On("GetDay", ctx).Return(0, expectedError).Once()

		// check
		err := service.RecordAdClick(ctx, uuid.New(), uuid.New())
		require.ErrorIs(t, err, expectedError)
	})

	t.Run("record ad click campaigns repo error", func(t *testing.T) {
		ctx := context.Background()

		adsRepoMock := mocks.NewAdsRepo(t)
		clientsRepoMock := mocks.NewClientsRepo(t)
		campaignsRepoMock := mocks.NewCampaignsRepo(t)
		clientActionsRepoMock := mocks.NewClientActionsRepo(t)
		timeRepoMock := mocks.NewTimeRepo(t)

		service := NewAdsService(adsRepoMock, clientsRepoMock, campaignsRepoMock, clientActionsRepoMock, timeRepoMock)

		// setup mocks
		currentDay := 5
		timeRepoMock.On("GetDay", ctx).Return(currentDay, nil).Once()

		clientId := uuid.New()
		campaignId := uuid.New()

		expectedError := errors.New("failed to get campaign")
		campaignsRepoMock.On("GetCampaignById", ctx, campaignId).Return(models.Campaign{}, expectedError).Once()

		// check
		err := service.RecordAdClick(ctx, clientId, campaignId)
		require.ErrorIs(t, err, expectedError)
	})

	t.Run("record ad click check impressed error", func(t *testing.T) {
		ctx := context.Background()

		adsRepoMock := mocks.NewAdsRepo(t)
		clientsRepoMock := mocks.NewClientsRepo(t)
		campaignsRepoMock := mocks.NewCampaignsRepo(t)
		clientActionsRepoMock := mocks.NewClientActionsRepo(t)
		timeRepoMock := mocks.NewTimeRepo(t)

		service := NewAdsService(adsRepoMock, clientsRepoMock, campaignsRepoMock, clientActionsRepoMock, timeRepoMock)

		// setup mocks
		currentDay := 5
		timeRepoMock.On("GetDay", ctx).Return(currentDay, nil).Once()

		clientId := uuid.New()
		campaignId := uuid.New()

		campaign := models.Campaign{CostPerClick: 100}
		campaignsRepoMock.On("GetCampaignById", ctx, campaignId).Return(campaign, nil).Once()

		expectedError := errors.New("failed to check impression")
		clientActionsRepoMock.On("CheckImpressed", ctx, clientId, campaignId).Return(false, expectedError).Once()

		// check
		err := service.RecordAdClick(ctx, clientId, campaignId)
		require.ErrorIs(t, err, expectedError)
	})

	t.Run("record ad click not impressed", func(t *testing.T) {
		ctx := context.Background()

		adsRepoMock := mocks.NewAdsRepo(t)
		clientsRepoMock := mocks.NewClientsRepo(t)
		campaignsRepoMock := mocks.NewCampaignsRepo(t)
		clientActionsRepoMock := mocks.NewClientActionsRepo(t)
		timeRepoMock := mocks.NewTimeRepo(t)

		service := NewAdsService(adsRepoMock, clientsRepoMock, campaignsRepoMock, clientActionsRepoMock, timeRepoMock)

		// setup mocks
		currentDay := 5
		timeRepoMock.On("GetDay", ctx).Return(currentDay, nil).Once()

		clientId := uuid.New()
		campaignId := uuid.New()

		campaign := models.Campaign{CostPerClick: 100}
		campaignsRepoMock.On("GetCampaignById", ctx, campaignId).Return(campaign, nil).Once()

		clientActionsRepoMock.On("CheckImpressed", ctx, clientId, campaignId).Return(false, nil).Once()

		// check
		err := service.RecordAdClick(ctx, clientId, campaignId)
		require.ErrorIs(t, err, models.ErrNotImpressed)
	})

	t.Run("record ad click record click error", func(t *testing.T) {
		ctx := context.Background()

		adsRepoMock := mocks.NewAdsRepo(t)
		clientsRepoMock := mocks.NewClientsRepo(t)
		campaignsRepoMock := mocks.NewCampaignsRepo(t)
		clientActionsRepoMock := mocks.NewClientActionsRepo(t)
		timeRepoMock := mocks.NewTimeRepo(t)

		service := NewAdsService(adsRepoMock, clientsRepoMock, campaignsRepoMock, clientActionsRepoMock, timeRepoMock)

		// setup mocks
		currentDay := 5
		timeRepoMock.On("GetDay", ctx).Return(currentDay, nil).Once()

		clientId := uuid.New()
		campaignId := uuid.New()

		campaign := models.Campaign{CostPerClick: 100}
		campaignsRepoMock.On("GetCampaignById", ctx, campaignId).Return(campaign, nil).Once()

		clientActionsRepoMock.On("CheckImpressed", ctx, clientId, campaignId).Return(true, nil).Once()

		expectedError := errors.New("failed to record click")
		clientActionsRepoMock.On("RecordClick", ctx, models.Click{
			ClientId:   clientId,
			CampaignId: campaignId,
			Date:       currentDay,
			Profit:     campaign.CostPerClick,
		}).Return(expectedError).Once()

		// check
		err := service.RecordAdClick(ctx, clientId, campaignId)
		require.ErrorIs(t, err, expectedError)
	})
}
