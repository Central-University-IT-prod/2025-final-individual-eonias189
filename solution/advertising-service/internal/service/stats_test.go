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

func TestStatsService(t *testing.T) {
	t.Run("get stats for campaign success", func(t *testing.T) {
		ctx := context.Background()

		statsRepoMock := mocks.NewStatsRepo(t)
		campaignsRepoMock := mocks.NewCampaignsRepo(t)
		advertisersRepoMock := mocks.NewAdvertisersRepo(t)

		service := NewStatsService(statsRepoMock, campaignsRepoMock, advertisersRepoMock)

		// setup mocks
		campaignId := uuid.New()
		campaignsRepoMock.On("GetCampaignById", ctx, campaignId).Return(models.Campaign{}, nil).Once()

		expectedStats := models.Stats{}
		statsRepoMock.On("GetStatsForCampaign", ctx, campaignId).Return(expectedStats, nil).Once()

		// check
		actualStats, err := service.GetStatsForCampaign(ctx, campaignId)
		require.NoError(t, err)
		require.Equal(t, expectedStats, actualStats)
	})

	t.Run("get stats for campaign campaigns repo error", func(t *testing.T) {
		ctx := context.Background()

		statsRepoMock := mocks.NewStatsRepo(t)
		campaignsRepoMock := mocks.NewCampaignsRepo(t)
		advertisersRepoMock := mocks.NewAdvertisersRepo(t)

		service := NewStatsService(statsRepoMock, campaignsRepoMock, advertisersRepoMock)

		// setup mocks
		campaignId := uuid.New()
		expectedError := errors.New("failed to get campaign")
		campaignsRepoMock.On("GetCampaignById", ctx, campaignId).Return(models.Campaign{}, expectedError).Once()

		// check
		actualStats, err := service.GetStatsForCampaign(ctx, campaignId)
		require.ErrorIs(t, err, expectedError)
		require.Equal(t, models.Stats{}, actualStats)
	})

	t.Run("get stats for campaign stats repo error", func(t *testing.T) {
		ctx := context.Background()

		statsRepoMock := mocks.NewStatsRepo(t)
		campaignsRepoMock := mocks.NewCampaignsRepo(t)
		advertisersRepoMock := mocks.NewAdvertisersRepo(t)

		service := NewStatsService(statsRepoMock, campaignsRepoMock, advertisersRepoMock)

		// setup mocks
		campaignId := uuid.New()
		campaignsRepoMock.On("GetCampaignById", ctx, campaignId).Return(models.Campaign{}, nil).Once()

		expectedError := errors.New("failed to get stats")
		statsRepoMock.On("GetStatsForCampaign", ctx, campaignId).Return(models.Stats{}, expectedError).Once()

		// check
		actualStats, err := service.GetStatsForCampaign(ctx, campaignId)
		require.ErrorIs(t, err, expectedError)
		require.Equal(t, models.Stats{}, actualStats)
	})

	t.Run("get stats for campaign daily success", func(t *testing.T) {
		ctx := context.Background()

		statsRepoMock := mocks.NewStatsRepo(t)
		campaignsRepoMock := mocks.NewCampaignsRepo(t)
		advertisersRepoMock := mocks.NewAdvertisersRepo(t)

		service := NewStatsService(statsRepoMock, campaignsRepoMock, advertisersRepoMock)

		// setup mocks
		campaignId := uuid.New()
		campaignsRepoMock.On("GetCampaignById", ctx, campaignId).Return(models.Campaign{}, nil).Once()

		expectedStats := []models.StatsDaily{}
		statsRepoMock.On("GetStatsForCampaignDaily", ctx, campaignId).Return(expectedStats, nil).Once()

		// check
		actualStats, err := service.GetStatsForCampaignDaily(ctx, campaignId)
		require.NoError(t, err)
		require.Equal(t, expectedStats, actualStats)
	})

	t.Run("get stats for campaign daily campaigns repo error", func(t *testing.T) {
		ctx := context.Background()

		statsRepoMock := mocks.NewStatsRepo(t)
		campaignsRepoMock := mocks.NewCampaignsRepo(t)
		advertisersRepoMock := mocks.NewAdvertisersRepo(t)

		service := NewStatsService(statsRepoMock, campaignsRepoMock, advertisersRepoMock)

		// setup mocks
		campaignId := uuid.New()
		expectedError := errors.New("failed to get campaign")
		campaignsRepoMock.On("GetCampaignById", ctx, campaignId).Return(models.Campaign{}, expectedError).Once()

		// check
		actualStats, err := service.GetStatsForCampaignDaily(ctx, campaignId)
		require.ErrorIs(t, err, expectedError)
		require.Nil(t, actualStats)
	})

	t.Run("get stats for campaign daily stats repo error", func(t *testing.T) {
		ctx := context.Background()

		statsRepoMock := mocks.NewStatsRepo(t)
		campaignsRepoMock := mocks.NewCampaignsRepo(t)
		advertisersRepoMock := mocks.NewAdvertisersRepo(t)

		service := NewStatsService(statsRepoMock, campaignsRepoMock, advertisersRepoMock)

		// setup mocks
		campaignId := uuid.New()
		campaignsRepoMock.On("GetCampaignById", ctx, campaignId).Return(models.Campaign{}, nil).Once()

		expectedError := errors.New("failed to get stats")
		statsRepoMock.On("GetStatsForCampaignDaily", ctx, campaignId).Return(nil, expectedError).Once()

		// check
		actualStats, err := service.GetStatsForCampaignDaily(ctx, campaignId)
		require.ErrorIs(t, err, expectedError)
		require.Nil(t, actualStats)
	})

	t.Run("get stats for advertiser success", func(t *testing.T) {
		ctx := context.Background()

		statsRepoMock := mocks.NewStatsRepo(t)
		campaignsRepoMock := mocks.NewCampaignsRepo(t)
		advertisersRepoMock := mocks.NewAdvertisersRepo(t)

		service := NewStatsService(statsRepoMock, campaignsRepoMock, advertisersRepoMock)

		// setup mocks
		advertiserId := uuid.New()
		advertisersRepoMock.On("GetAdvertiserById", ctx, advertiserId).Return(models.Advertiser{}, nil).Once()

		expectedStats := models.Stats{}
		statsRepoMock.On("GetStatsForAdvertiser", ctx, advertiserId).Return(expectedStats, nil).Once()

		// check
		actualStats, err := service.GetStatsForAdvertiser(ctx, advertiserId)
		require.NoError(t, err)
		require.Equal(t, expectedStats, actualStats)
	})

	t.Run("get stats for advertiser advertisers repo error", func(t *testing.T) {
		ctx := context.Background()

		statsRepoMock := mocks.NewStatsRepo(t)
		campaignsRepoMock := mocks.NewCampaignsRepo(t)
		advertisersRepoMock := mocks.NewAdvertisersRepo(t)

		service := NewStatsService(statsRepoMock, campaignsRepoMock, advertisersRepoMock)

		// setup mocks
		advertiserId := uuid.New()
		expectedError := errors.New("failed to get advertiser")
		advertisersRepoMock.On("GetAdvertiserById", ctx, advertiserId).Return(models.Advertiser{}, expectedError).Once()

		// check
		actualStats, err := service.GetStatsForAdvertiser(ctx, advertiserId)
		require.ErrorIs(t, err, expectedError)
		require.Equal(t, models.Stats{}, actualStats)
	})

	t.Run("get stats for advertiser stats repo error", func(t *testing.T) {
		ctx := context.Background()

		statsRepoMock := mocks.NewStatsRepo(t)
		campaignsRepoMock := mocks.NewCampaignsRepo(t)
		advertisersRepoMock := mocks.NewAdvertisersRepo(t)

		service := NewStatsService(statsRepoMock, campaignsRepoMock, advertisersRepoMock)

		// setup mocks
		advertiserId := uuid.New()
		advertisersRepoMock.On("GetAdvertiserById", ctx, advertiserId).Return(models.Advertiser{}, nil).Once()

		expectedError := errors.New("failed to get stats")
		statsRepoMock.On("GetStatsForAdvertiser", ctx, advertiserId).Return(models.Stats{}, expectedError).Once()

		// check
		actualStats, err := service.GetStatsForAdvertiser(ctx, advertiserId)
		require.ErrorIs(t, err, expectedError)
		require.Equal(t, models.Stats{}, actualStats)
	})

	t.Run("get stats for advertiser daily success", func(t *testing.T) {
		ctx := context.Background()

		statsRepoMock := mocks.NewStatsRepo(t)
		campaignsRepoMock := mocks.NewCampaignsRepo(t)
		advertisersRepoMock := mocks.NewAdvertisersRepo(t)

		service := NewStatsService(statsRepoMock, campaignsRepoMock, advertisersRepoMock)

		// setup mocks
		advertiserId := uuid.New()
		advertisersRepoMock.On("GetAdvertiserById", ctx, advertiserId).Return(models.Advertiser{}, nil).Once()

		expectedStats := []models.StatsDaily{}
		statsRepoMock.On("GetStatsForAdvertiserDaily", ctx, advertiserId).Return(expectedStats, nil).Once()

		// check
		actualStats, err := service.GetStatsForAdvertiserDaily(ctx, advertiserId)
		require.NoError(t, err)
		require.Equal(t, expectedStats, actualStats)
	})

	t.Run("get stats for advertiser daily advertisers repo error", func(t *testing.T) {
		ctx := context.Background()

		statsRepoMock := mocks.NewStatsRepo(t)
		campaignsRepoMock := mocks.NewCampaignsRepo(t)
		advertisersRepoMock := mocks.NewAdvertisersRepo(t)

		service := NewStatsService(statsRepoMock, campaignsRepoMock, advertisersRepoMock)

		// setup mocks
		advertiserId := uuid.New()
		expectedError := errors.New("failed to get advertiser")
		advertisersRepoMock.On("GetAdvertiserById", ctx, advertiserId).Return(models.Advertiser{}, expectedError).Once()

		// check
		actualStats, err := service.GetStatsForAdvertiserDaily(ctx, advertiserId)
		require.ErrorIs(t, err, expectedError)
		require.Nil(t, actualStats)
	})

	t.Run("get stats for advertiser daily stats repo error", func(t *testing.T) {
		ctx := context.Background()

		statsRepoMock := mocks.NewStatsRepo(t)
		campaignsRepoMock := mocks.NewCampaignsRepo(t)
		advertisersRepoMock := mocks.NewAdvertisersRepo(t)

		service := NewStatsService(statsRepoMock, campaignsRepoMock, advertisersRepoMock)

		// setup mocks
		advertiserId := uuid.New()
		advertisersRepoMock.On("GetAdvertiserById", ctx, advertiserId).Return(models.Advertiser{}, nil).Once()

		expectedError := errors.New("failed to get stats")
		statsRepoMock.On("GetStatsForAdvertiserDaily", ctx, advertiserId).Return(nil, expectedError).Once()

		// check
		actualStats, err := service.GetStatsForAdvertiserDaily(ctx, advertiserId)
		require.ErrorIs(t, err, expectedError)
		require.Nil(t, actualStats)
	})
}
