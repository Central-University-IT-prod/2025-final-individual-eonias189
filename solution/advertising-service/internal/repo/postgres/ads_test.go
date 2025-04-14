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

func TestGetAdForClient(t *testing.T) {
	ctx := context.Background()
	db := helpers.SetUpPostgres(ctx, t, "../../../migrations")

	clientsRepo := NewClientRepo(db)
	client := generateClient()
	_, err := clientsRepo.UpsertClients(ctx, []models.Client{client})
	require.NoError(t, err)

	advertisersRepo := NewAdvertiserRepo(db)
	advertiser := generateAdvertiser()
	advertiserId := advertiser.Id
	_, err = advertisersRepo.UpsertAdvertisers(ctx, []models.Advertiser{advertiser})
	require.NoError(t, err)

	// check if works without ml_score
	campaignsRepo := NewCampaignsRepo(db)
	campaign := generateCampaign()
	campaign.AdvertiserId = advertiser.Id
	campaign.StartDate = 0
	campaign.Gender = nil
	campaign.Location = nil
	campaign.AgeFrom = nil
	campaign.AgeTo = nil
	campaign.Id, err = campaignsRepo.CreateCampaign(ctx, advertiserId, dto.CampaignDataFromCampaign(campaign))
	require.NoError(t, err)

	adsRepo := NewAdsRepo(db)
	ad, err := adsRepo.GetAdForClient(ctx, client, 0)
	require.NoError(t, err)

	require.Equal(t, campaign.Id, ad.CampaignId)
	require.Equal(t, campaign.AdvertiserId, ad.AdvertiserId)
	require.Equal(t, campaign.AdTitle, ad.AdTitle)
	require.Equal(t, campaign.AdText, ad.AdText)

	err = campaignsRepo.DeleteCampaign(ctx, campaign.Id)
	require.NoError(t, err)

	// check if works with ml_score
	mlScoreRepo := NewMlScoresRepo(db)
	mlScore := models.MLScore{
		ClientId:     client.Id,
		AdvertiserId: advertiser.Id,
		Score:        int(gofakeit.Int32()),
	}
	err = mlScoreRepo.UpsertMLScore(ctx, mlScore)
	require.NoError(t, err)

	campaign = generateCampaign()
	campaign.AdvertiserId = advertiser.Id
	campaign.StartDate = 0
	campaign.Gender = nil
	campaign.Location = nil
	campaign.AgeFrom = nil
	campaign.AgeTo = nil
	campaign.Id, err = campaignsRepo.CreateCampaign(ctx, advertiserId, dto.CampaignDataFromCampaign(campaign))
	require.NoError(t, err)

	ad, err = adsRepo.GetAdForClient(ctx, client, 0)
	require.NoError(t, err)

	require.Equal(t, campaign.Id, ad.CampaignId)
	require.Equal(t, campaign.AdvertiserId, ad.AdvertiserId)
	require.Equal(t, campaign.AdTitle, ad.AdTitle)
	require.Equal(t, campaign.AdText, ad.AdText)

	err = campaignsRepo.DeleteCampaign(ctx, campaign.Id)
	require.NoError(t, err)

	// check if returns campaign with gender='ALL'
	campaign = generateCampaign()
	campaign.AdvertiserId = advertiser.Id
	campaign.StartDate = 0
	campaign.Gender = &models.GenderAll
	campaign.Location = nil
	campaign.AgeFrom = nil
	campaign.AgeTo = nil
	campaign.Id, err = campaignsRepo.CreateCampaign(ctx, advertiserId, dto.CampaignDataFromCampaign(campaign))
	require.NoError(t, err)

	ad, err = adsRepo.GetAdForClient(ctx, client, 0)
	require.NoError(t, err)

	require.Equal(t, campaign.Id, ad.CampaignId)
	require.Equal(t, campaign.AdvertiserId, ad.AdvertiserId)
	require.Equal(t, campaign.AdTitle, ad.AdTitle)
	require.Equal(t, campaign.AdText, ad.AdText)

	err = campaignsRepo.DeleteCampaign(ctx, campaign.Id)
	require.NoError(t, err)

	// check if returns error models.ErrNoAdsForClient
	_, err = adsRepo.GetAdForClient(ctx, client, 0)
	require.ErrorIs(t, err, models.ErrNoAdsForClient)

	// check if don`t return not active campaign

	// campaign hasn`t started
	campaign1 := generateCampaign()
	campaign1.AdvertiserId = advertiser.Id
	campaign1.StartDate = 100
	campaign1.Gender = nil
	campaign1.Location = nil
	campaign1.AgeFrom = nil
	campaign1.AgeTo = nil
	campaign1.Id, err = campaignsRepo.CreateCampaign(ctx, advertiserId, dto.CampaignDataFromCampaign(campaign1))
	require.NoError(t, err)

	// campaign has already finished
	campaign2 := generateCampaign()
	campaign2.AdvertiserId = advertiser.Id
	campaign2.StartDate = 0
	campaign2.EndDate = 5
	campaign2.Gender = nil
	campaign2.Location = nil
	campaign2.AgeFrom = nil
	campaign2.AgeTo = nil
	campaign2.Id, err = campaignsRepo.CreateCampaign(ctx, advertiserId, dto.CampaignDataFromCampaign(campaign2))
	require.NoError(t, err)

	_, err = adsRepo.GetAdForClient(ctx, client, 6)
	require.ErrorIs(t, err, models.ErrNoAdsForClient)

	err = campaignsRepo.DeleteCampaign(ctx, campaign1.Id)
	require.NoError(t, err)
	err = campaignsRepo.DeleteCampaign(ctx, campaign2.Id)
	require.NoError(t, err)

	// check if don`t return campaign that doesn`t match client
	client = models.Client{
		Id:       uuid.New(),
		Login:    "login123",
		Age:      42,
		Location: "Kemerovo",
		Gender:   models.GenderMale,
	}
	_, err = clientsRepo.UpsertClients(ctx, []models.Client{client})
	require.NoError(t, err)

	// campaign with another gender targeting
	campaign1 = generateCampaign()
	campaign1.AdvertiserId = advertiser.Id
	campaign1.StartDate = 0
	campaign1.Gender = &models.GenderFemale
	campaign1.Location = nil
	campaign1.AgeFrom = nil
	campaign1.AgeTo = nil
	campaign1.Id, err = campaignsRepo.CreateCampaign(ctx, advertiserId, dto.CampaignDataFromCampaign(campaign1))
	require.NoError(t, err)

	// campaign with another location targeting
	campaign2 = generateCampaign()
	campaign2.AdvertiserId = advertiser.Id
	campaign2.StartDate = 0
	campaign2.Gender = nil
	campaign2.Location = pointer("Rostov")
	campaign2.AgeFrom = nil
	campaign2.AgeTo = nil
	campaign2.Id, err = campaignsRepo.CreateCampaign(ctx, advertiserId, dto.CampaignDataFromCampaign(campaign2))
	require.NoError(t, err)

	// campaigns with another age targeting
	campaign3 := generateCampaign()
	campaign3.AdvertiserId = advertiser.Id
	campaign3.StartDate = 0
	campaign3.Gender = nil
	campaign3.Location = nil
	campaign3.AgeFrom = pointer(43)
	campaign3.AgeTo = nil
	campaign3.Id, err = campaignsRepo.CreateCampaign(ctx, advertiserId, dto.CampaignDataFromCampaign(campaign3))
	require.NoError(t, err)

	campaign4 := generateCampaign()
	campaign4.AdvertiserId = advertiser.Id
	campaign4.StartDate = 0
	campaign4.Gender = nil
	campaign4.Location = nil
	campaign4.AgeFrom = nil
	campaign4.AgeTo = pointer(41)
	campaign4.Id, err = campaignsRepo.CreateCampaign(ctx, advertiserId, dto.CampaignDataFromCampaign(campaign4))
	require.NoError(t, err)

	_, err = adsRepo.GetAdForClient(ctx, client, 0)
	require.ErrorIs(t, err, models.ErrNoAdsForClient)

	err = campaignsRepo.DeleteCampaign(ctx, campaign1.Id)
	require.NoError(t, err)
	err = campaignsRepo.DeleteCampaign(ctx, campaign2.Id)
	require.NoError(t, err)
	err = campaignsRepo.DeleteCampaign(ctx, campaign3.Id)
	require.NoError(t, err)
	err = campaignsRepo.DeleteCampaign(ctx, campaign4.Id)
	require.NoError(t, err)

	// check if don`t return already impressed campaign

	campaign = generateCampaign()
	campaign.AdvertiserId = advertiser.Id
	campaign.StartDate = 0
	campaign.Gender = nil
	campaign.Location = nil
	campaign.AgeFrom = nil
	campaign.AgeTo = nil
	campaign.Id, err = campaignsRepo.CreateCampaign(ctx, advertiserId, dto.CampaignDataFromCampaign(campaign))
	require.NoError(t, err)

	client = generateClient()
	_, err = clientsRepo.UpsertClients(ctx, []models.Client{client})
	require.NoError(t, err)

	clientActionsRepo := NewClientActionsRepo(db)
	err = clientActionsRepo.RecordImpression(ctx, models.Impression{
		ClientId:   client.Id,
		CampaignId: campaign.Id,
		Date:       0,
		Profit:     100,
	})
	require.NoError(t, err)

	_, err = adsRepo.GetAdForClient(ctx, client, 0)
	require.ErrorIs(t, err, models.ErrNoAdsForClient)

	err = campaignsRepo.DeleteCampaign(ctx, campaign.Id)
	require.NoError(t, err)
}
