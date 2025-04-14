package postgres

import (
	"advertising/advertising-service/internal/models"
	"advertising/tests/helpers"
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestAdvertiserRepo(t *testing.T) {
	ctx := context.Background()
	db := helpers.SetUpPostgres(ctx, t, "../../../migrations")
	advertiserRepo := NewAdvertiserRepo(db)

	// check insert
	advertisers := make([]models.Advertiser, 0, 20)
	for i := 0; i < 20; i++ {
		advertisers = append(advertisers, generateAdvertiser())
	}

	advertisersGot, err := advertiserRepo.UpsertAdvertisers(ctx, advertisers)
	require.NoError(t, err)
	require.ElementsMatch(t, advertisers, advertisersGot)

	// check get advertiser by id
	for _, advertiser := range advertisers {
		advertiserGot, err := advertiserRepo.GetAdvertiserById(ctx, advertiser.Id)
		require.NoError(t, err)
		require.Equal(t, advertiser, advertiserGot)
	}

	// check get non-existent advertiser
	_, err = advertiserRepo.GetAdvertiserById(ctx, uuid.New())
	require.ErrorIs(t, err, models.ErrAdvertiserNotFound)

	// check upsert advertiser
	advertisers = advertisers[9:]
	for i, advertiser := range advertisers {
		newAdvertiser := generateAdvertiser()
		newAdvertiser.Id = advertiser.Id
		advertisers[i] = newAdvertiser
	}
	for i := 0; i < 10; i++ {
		advertisers = append(advertisers, generateAdvertiser())
	}

	advertisersGot, err = advertiserRepo.UpsertAdvertisers(ctx, advertisers)
	require.NoError(t, err)
	require.ElementsMatch(t, advertisers, advertisersGot)

	// check insert advertisers with repeated values
	advertisers = make([]models.Advertiser, 20)
	for i := 0; i < 10; i++ {
		newAdvertiser := generateAdvertiser()
		advertisers[i] = newAdvertiser
		advertisers[10+i] = newAdvertiser
	}

	advertisersGot, err = advertiserRepo.UpsertAdvertisers(ctx, advertisers)
	require.NoError(t, err)
	require.ElementsMatch(t, advertisers[:10], advertisersGot)

}

func generateAdvertiser() models.Advertiser {
	return models.Advertiser{
		Id:   uuid.New(),
		Name: gofakeit.Company(),
	}
}
