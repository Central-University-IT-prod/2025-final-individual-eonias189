package postgres

import (
	"advertising/advertising-service/internal/dto"
	"advertising/advertising-service/internal/models"
	"advertising/tests/helpers"
	"context"
	"slices"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCampaignCreateAndGetById(t *testing.T) {
	ctx := context.Background()
	db := helpers.SetUpPostgres(ctx, t, "../../../migrations")
	advertisersRepo := NewAdvertiserRepo(db)
	campaignsRepo := NewCampaignsRepo(db)

	advertiserId := uuid.New()

	_, err := advertisersRepo.UpsertAdvertisers(ctx, []models.Advertiser{
		{
			Id:   advertiserId,
			Name: gofakeit.Company(),
		},
	})
	require.NoError(t, err)

	// check create campaign with all fields
	campaign := generateCampaign()
	campaign.AdvertiserId = advertiserId

	campaign.Id, err = campaignsRepo.CreateCampaign(ctx, advertiserId, dto.CampaignDataFromCampaign(campaign))
	require.NoError(t, err)

	campaignGot, err := campaignsRepo.GetCampaignById(ctx, campaign.Id)
	require.NoError(t, err)
	require.Equal(t, campaign, campaignGot)

	// check create campaign with not all fields
	campaign = generateCampaign()
	campaign.AdvertiserId = advertiserId
	campaign.AgeTo = nil
	campaign.Gender = nil

	campaign.Id, err = campaignsRepo.CreateCampaign(ctx, advertiserId, dto.CampaignDataFromCampaign(campaign))
	require.NoError(t, err)

	campaignGot, err = campaignsRepo.GetCampaignById(ctx, campaign.Id)
	require.NoError(t, err)
	require.Equal(t, campaign, campaignGot)

	// check get non-existent campaign
	_, err = campaignsRepo.GetCampaignById(ctx, uuid.New())
	require.ErrorIs(t, err, models.ErrCampaignNotFound)

	// check create campaign with non-existent advertiser
	campaign = generateCampaign()

	_, err = campaignsRepo.CreateCampaign(ctx, uuid.New(), dto.CampaignDataFromCampaign(campaign))
	require.ErrorIs(t, err, models.ErrAdvertiserNotFound)

}

func TestListCampaignsByAdvertiserId(t *testing.T) {
	ctx := context.Background()
	db := helpers.SetUpPostgres(ctx, t, "../../../migrations")
	advertisersRepo := NewAdvertiserRepo(db)
	campaignsRepo := NewCampaignsRepo(db)

	advertiserId := uuid.New()

	_, err := advertisersRepo.UpsertAdvertisers(ctx, []models.Advertiser{
		{
			Id:   advertiserId,
			Name: gofakeit.Company(),
		},
	})
	require.NoError(t, err)

	campaigns := make([]models.Campaign, 0, 20)
	for range 20 {
		campaign := generateCampaign()
		campaign.AdvertiserId = advertiserId
		if gofakeit.IntRange(0, 10) > 6 {
			campaign.Gender = nil
		}
		if gofakeit.IntRange(0, 10) > 6 {
			campaign.AgeFrom = nil
		}
		if gofakeit.IntRange(0, 10) > 6 {
			campaign.AgeTo = nil
		}
		if gofakeit.IntRange(0, 10) > 6 {
			campaign.Location = nil
		}

		campaign.Id, err = campaignsRepo.CreateCampaign(ctx, advertiserId, dto.CampaignDataFromCampaign(campaign))
		require.NoError(t, err)

		campaigns = append(campaigns, campaign)
	}

	slices.Reverse(campaigns)

	// check list with big size and page 1
	campaignsGot, err := campaignsRepo.ListCampaignsForAdvertiser(ctx, advertiserId, dto.PaginationParams{
		Size: 9999,
		Page: 1,
	})
	require.NoError(t, err)
	require.Equal(t, campaigns, campaignsGot)

	// check list with size 10 and page 1
	campaignsGot, err = campaignsRepo.ListCampaignsForAdvertiser(ctx, advertiserId, dto.PaginationParams{
		Size: 10,
		Page: 1,
	})
	require.NoError(t, err)
	require.Equal(t, campaigns[:10], campaignsGot)

	// check list with size 5 and page 2
	campaignsGot, err = campaignsRepo.ListCampaignsForAdvertiser(ctx, advertiserId, dto.PaginationParams{
		Size: 5,
		Page: 2,
	})
	require.NoError(t, err)
	require.ElementsMatch(t, campaigns[5:10], campaignsGot)

	// check list with empty result
	campaignsGot, err = campaignsRepo.ListCampaignsForAdvertiser(ctx, uuid.New(), dto.PaginationParams{})
	require.NoError(t, err)
	require.Empty(t, campaignsGot)

}

func TestUpdateCampaign(t *testing.T) {
	ctx := context.Background()
	db := helpers.SetUpPostgres(ctx, t, "../../../migrations")
	advertisersRepo := NewAdvertiserRepo(db)
	campaignsRepo := NewCampaignsRepo(db)

	advertiserId := uuid.New()

	_, err := advertisersRepo.UpsertAdvertisers(ctx, []models.Advertiser{
		{
			Id:   advertiserId,
			Name: gofakeit.Company(),
		},
	})
	require.NoError(t, err)

	// check update all fields
	campaign := generateCampaign()
	campaign.AdvertiserId = advertiserId

	campaign.Id, err = campaignsRepo.CreateCampaign(ctx, advertiserId, dto.CampaignDataFromCampaign(campaign))
	require.NoError(t, err)

	campaignNew := generateCampaign()
	campaignNew.Id = campaign.Id
	campaignNew.AdvertiserId = advertiserId

	err = campaignsRepo.UpdateCampaign(ctx, campaign.Id, dto.CampaignDataFromCampaign(campaignNew))
	require.NoError(t, err)

	campaignGot, err := campaignsRepo.GetCampaignById(ctx, campaign.Id)
	require.NoError(t, err)
	require.Equal(t, campaignNew, campaignGot)

	// check update not all fields
	campaign = generateCampaign()
	campaign.AdvertiserId = advertiserId

	campaign.Id, err = campaignsRepo.CreateCampaign(ctx, advertiserId, dto.CampaignDataFromCampaign(campaign))
	require.NoError(t, err)

	campaignNew = campaign
	campaignNew.ClicksLimit = gofakeit.IntRange(0, 9999)
	campaignNew.CostPerImpression = gofakeit.Float64Range(0, 999)
	campaignNew.AdText = gofakeit.Sentence(30)

	err = campaignsRepo.UpdateCampaign(ctx, campaign.Id, dto.CampaignDataFromCampaign(campaignNew))
	require.NoError(t, err)

	campaignGot, err = campaignsRepo.GetCampaignById(ctx, campaign.Id)
	require.NoError(t, err)
	require.Equal(t, campaignNew, campaignGot)

	// check update non-existent campaign
	err = campaignsRepo.UpdateCampaign(ctx, uuid.New(), dto.CampaignData{})
	require.ErrorIs(t, err, models.ErrCampaignNotFound)

}

func TestDeleteCampaign(t *testing.T) {
	ctx := context.Background()
	db := helpers.SetUpPostgres(ctx, t, "../../../migrations")
	advertisersRepo := NewAdvertiserRepo(db)
	campaignsRepo := NewCampaignsRepo(db)

	advertiserId := uuid.New()

	_, err := advertisersRepo.UpsertAdvertisers(ctx, []models.Advertiser{
		{
			Id:   advertiserId,
			Name: gofakeit.Company(),
		},
	})
	require.NoError(t, err)

	// check delete existing campaign
	campaign := generateCampaign()
	campaign.AdvertiserId = advertiserId
	campaign.Id, err = campaignsRepo.CreateCampaign(ctx, advertiserId, dto.CampaignDataFromCampaign(campaign))
	require.NoError(t, err)

	err = campaignsRepo.DeleteCampaign(ctx, campaign.Id)
	require.NoError(t, err)

	_, err = campaignsRepo.GetCampaignById(ctx, campaign.Id)
	require.ErrorIs(t, err, models.ErrCampaignNotFound)

	// check delete non-existent campaign
	campaign = generateCampaign()
	campaign.AdvertiserId = advertiserId
	campaign.Id, err = campaignsRepo.CreateCampaign(ctx, advertiserId, dto.CampaignDataFromCampaign(campaign))
	require.NoError(t, err)

	err = campaignsRepo.DeleteCampaign(ctx, uuid.New())
	require.ErrorIs(t, err, models.ErrCampaignNotFound)

}

func pointer[T any](value T) *T {
	return &value
}

func generateGender() models.Gender {
	return models.Gender(gofakeit.RandomString([]string{"MALE", "FEMALE", "ALL"}))
}

func generateCampaign() models.Campaign {
	return models.Campaign{
		ImpressionsLimit:  gofakeit.IntRange(100, 9999),
		ClicksLimit:       gofakeit.IntRange(10, 9999),
		CostPerImpression: gofakeit.Float64Range(0, 9999),
		CostPerClick:      gofakeit.Float64Range(0, 9999),
		AdTitle:           gofakeit.Sentence(8),
		AdText:            gofakeit.Sentence(20),
		StartDate:         gofakeit.IntRange(0, 10),
		EndDate:           gofakeit.IntRange(11, 20),
		Gender:            pointer(generateGender()),
		AgeFrom:           pointer(gofakeit.IntRange(0, 10)),
		AgeTo:             pointer(gofakeit.IntRange(15, 90)),
		Location:          pointer(gofakeit.City()),
	}
}
