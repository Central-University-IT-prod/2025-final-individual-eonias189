package postgres

import (
	"advertising/advertising-service/internal/dto"
	"advertising/advertising-service/internal/models"
	"advertising/tests/helpers"
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestRecordImpression(t *testing.T) {
	ctx := context.Background()
	db := helpers.SetUpPostgres(ctx, t, "../../../migrations")

	clientsRepo := NewClientRepo(db)
	advertiserRepo := NewAdvertiserRepo(db)
	campaignsRepo := NewCampaignsRepo(db)
	clientActionsRepo := NewClientActionsRepo(db)

	client := generateClient()
	_, err := clientsRepo.UpsertClients(ctx, []models.Client{client})
	require.NoError(t, err)

	advertiser := generateAdvertiser()
	advertiserId := advertiser.Id
	_, err = advertiserRepo.UpsertAdvertisers(ctx, []models.Advertiser{advertiser})
	require.NoError(t, err)

	campaign := generateCampaign()
	campaign.AdvertiserId = advertiser.Id
	campaign.Id, err = campaignsRepo.CreateCampaign(ctx, advertiserId, dto.CampaignDataFromCampaign(campaign))
	require.NoError(t, err)

	// record impression success
	err = clientActionsRepo.RecordImpression(ctx, models.Impression{
		ClientId:   client.Id,
		CampaignId: campaign.Id,
		Date:       gofakeit.IntRange(0, 999),
		Profit:     gofakeit.Float64Range(0, 999),
	})
	require.NoError(t, err)

	// check returns models.ErrAlreadyImpressed
	err = clientActionsRepo.RecordImpression(ctx, models.Impression{
		ClientId:   client.Id,
		CampaignId: campaign.Id,
		Date:       gofakeit.IntRange(0, 999),
		Profit:     gofakeit.Float64Range(0, 999),
	})
	require.ErrorIs(t, err, models.ErrAlreadyImpressed)

	// check returns models.ErrClientNotFound
	err = clientActionsRepo.RecordImpression(ctx, models.Impression{
		ClientId:   uuid.New(),
		CampaignId: campaign.Id,
		Date:       gofakeit.IntRange(0, 999),
		Profit:     gofakeit.Float64Range(0, 999),
	})
	require.ErrorIs(t, err, models.ErrClientNotFound)

	// check returns models.ErrCampaignNotFound
	err = clientActionsRepo.RecordImpression(ctx, models.Impression{
		ClientId:   client.Id,
		CampaignId: uuid.New(),
		Date:       gofakeit.IntRange(0, 999),
		Profit:     gofakeit.Float64Range(0, 999),
	})
	require.ErrorIs(t, err, models.ErrCampaignNotFound)
}

func TestRecordClick(t *testing.T) {
	ctx := context.Background()
	db := helpers.SetUpPostgres(ctx, t, "../../../migrations")

	clientsRepo := NewClientRepo(db)
	advertiserRepo := NewAdvertiserRepo(db)
	campaignsRepo := NewCampaignsRepo(db)
	clientActionsRepo := NewClientActionsRepo(db)

	client := generateClient()
	_, err := clientsRepo.UpsertClients(ctx, []models.Client{client})
	require.NoError(t, err)

	advertiser := generateAdvertiser()
	advertiserId := advertiser.Id
	_, err = advertiserRepo.UpsertAdvertisers(ctx, []models.Advertiser{advertiser})
	require.NoError(t, err)

	campaign := generateCampaign()
	campaign.AdvertiserId = advertiser.Id
	campaign.Id, err = campaignsRepo.CreateCampaign(ctx, advertiserId, dto.CampaignDataFromCampaign(campaign))
	require.NoError(t, err)

	// record click success
	err = clientActionsRepo.RecordClick(ctx, models.Click{
		ClientId:   client.Id,
		CampaignId: campaign.Id,
		Date:       gofakeit.IntRange(0, 999),
		Profit:     gofakeit.Float64Range(0, 999),
	})
	require.NoError(t, err)

	// check returns models.ErrAlreadyClicked
	err = clientActionsRepo.RecordClick(ctx, models.Click{
		ClientId:   client.Id,
		CampaignId: campaign.Id,
		Date:       gofakeit.IntRange(0, 999),
		Profit:     gofakeit.Float64Range(0, 999),
	})
	require.ErrorIs(t, err, models.ErrAlreadyClicked)

	// check returns models.ErrClientNotFound
	err = clientActionsRepo.RecordClick(ctx, models.Click{
		ClientId:   uuid.New(),
		CampaignId: campaign.Id,
		Date:       gofakeit.IntRange(0, 999),
		Profit:     gofakeit.Float64Range(0, 999),
	})
	require.ErrorIs(t, err, models.ErrClientNotFound)

	// check returns models.ErrCampaignNotFound
	err = clientActionsRepo.RecordClick(ctx, models.Click{
		ClientId:   client.Id,
		CampaignId: uuid.New(),
		Date:       gofakeit.IntRange(0, 999),
		Profit:     gofakeit.Float64Range(0, 999),
	})
	require.ErrorIs(t, err, models.ErrCampaignNotFound)
}

func TestCheckImpressed(t *testing.T) {
	ctx := context.Background()
	db := helpers.SetUpPostgres(ctx, t, "../../../migrations")

	clientsRepo := NewClientRepo(db)
	advertiserRepo := NewAdvertiserRepo(db)
	campaignsRepo := NewCampaignsRepo(db)
	clientActionsRepo := NewClientActionsRepo(db)

	client := generateClient()
	_, err := clientsRepo.UpsertClients(ctx, []models.Client{client})
	require.NoError(t, err)

	advertiser := generateAdvertiser()
	advertiserId := advertiser.Id
	_, err = advertiserRepo.UpsertAdvertisers(ctx, []models.Advertiser{advertiser})
	require.NoError(t, err)

	campaign := generateCampaign()
	campaign.AdvertiserId = advertiser.Id
	campaign.Id, err = campaignsRepo.CreateCampaign(ctx, advertiserId, dto.CampaignDataFromCampaign(campaign))
	require.NoError(t, err)

	err = clientActionsRepo.RecordImpression(ctx, models.Impression{
		ClientId:   client.Id,
		CampaignId: campaign.Id,
		Date:       gofakeit.IntRange(0, 999),
		Profit:     gofakeit.Float64Range(0, 999),
	})
	require.NoError(t, err)

	// check impressed true
	impressed, err := clientActionsRepo.CheckImpressed(ctx, client.Id, campaign.Id)
	require.NoError(t, err)
	require.True(t, impressed)

	// check impressed false
	impressed, err = clientActionsRepo.CheckImpressed(ctx, uuid.New(), campaign.Id)
	require.NoError(t, err)
	require.False(t, impressed)

	impressed, err = clientActionsRepo.CheckImpressed(ctx, client.Id, uuid.New())
	require.NoError(t, err)
	require.False(t, impressed)
}
