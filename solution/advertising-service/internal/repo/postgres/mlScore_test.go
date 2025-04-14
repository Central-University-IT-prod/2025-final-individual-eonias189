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

func TestMLScoreRepo(t *testing.T) {
	ctx := context.Background()
	db := helpers.SetUpPostgres(ctx, t, "../../../migrations")
	clientsRepo := NewClientRepo(db)
	advertiserRepo := NewAdvertiserRepo(db)
	mlScoreRepo := NewMlScoresRepo(db)

	// check insert scores
	clients := make([]models.Client, 0, 20)
	advertisers := make([]models.Advertiser, 0, 20)
	mlScores := make([]models.MLScore, 0, 400)

	for i := 0; i < 20; i++ {
		clients = append(clients, generateClient())
		advertisers = append(advertisers, generateAdvertiser())
	}

	for _, client := range clients {
		for _, advertiser := range advertisers {
			mlScores = append(mlScores, models.MLScore{
				ClientId:     client.Id,
				AdvertiserId: advertiser.Id,
				Score:        gofakeit.IntRange(0, 999),
			})
		}
	}

	_, err := clientsRepo.UpsertClients(ctx, clients)
	require.NoError(t, err)
	_, err = advertiserRepo.UpsertAdvertisers(ctx, advertisers)
	require.NoError(t, err)

	for _, mlScore := range mlScores {
		err := mlScoreRepo.UpsertMLScore(ctx, mlScore)
		require.NoError(t, err)
	}

	// check update scores
	for i, mlScore := range mlScores {
		mlScore.Score = gofakeit.IntRange(0, 999)
		mlScores[i] = mlScore
	}

	for _, mlScore := range mlScores {
		err := mlScoreRepo.UpsertMLScore(ctx, mlScore)
		require.NoError(t, err)
	}

	// check insert score with non-existent client
	err = mlScoreRepo.UpsertMLScore(ctx, models.MLScore{
		ClientId:     uuid.New(),
		AdvertiserId: advertisers[0].Id,
		Score:        gofakeit.IntRange(0, 999),
	})
	require.ErrorIs(t, err, models.ErrClientNotFound)

	// check insert score with non-existent advertiser
	err = mlScoreRepo.UpsertMLScore(ctx, models.MLScore{
		ClientId:     clients[0].Id,
		AdvertiserId: uuid.New(),
		Score:        gofakeit.IntRange(0, 999),
	})
	require.ErrorIs(t, err, models.ErrAdvertiserNotFound)
}
