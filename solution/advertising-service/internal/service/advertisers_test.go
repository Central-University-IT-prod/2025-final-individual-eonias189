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

func TestAdvertisersService(t *testing.T) {
	t.Run("get advertiser by id success", func(t *testing.T) {
		ctx := context.Background()

		advertisersRepoMock := mocks.NewAdvertisersRepo(t)
		mlscoresRepoMock := mocks.NewMlScoresRepo(t)
		service := NewAdvertisersService(advertisersRepoMock, mlscoresRepoMock)

		// setup mocks
		advertiserId := uuid.New()
		expectedAdvertiser := models.Advertiser{
			Id:   advertiserId,
			Name: "name",
		}

		advertisersRepoMock.On("GetAdvertiserById", ctx, advertiserId).Return(expectedAdvertiser, nil).Once()

		// check
		actualAdvertiser, err := service.GetAdvertiserById(ctx, advertiserId)
		require.NoError(t, err)
		require.Equal(t, expectedAdvertiser, actualAdvertiser)
	})

	t.Run("get advertiser by id error", func(t *testing.T) {
		ctx := context.Background()

		advertisersRepoMock := mocks.NewAdvertisersRepo(t)
		mlscoresRepoMock := mocks.NewMlScoresRepo(t)
		service := NewAdvertisersService(advertisersRepoMock, mlscoresRepoMock)

		// setup mocks
		advertiserId := uuid.New()
		expectedAdvertiser := models.Advertiser{}
		expectedError := errors.New("failed to get advertiser")

		advertisersRepoMock.On("GetAdvertiserById", ctx, advertiserId).Return(expectedAdvertiser, expectedError).Once()

		// check
		actualAdvertiser, actualError := service.GetAdvertiserById(ctx, advertiserId)
		require.ErrorIs(t, actualError, expectedError)
		require.Equal(t, expectedAdvertiser, actualAdvertiser)
	})

	t.Run("upsert advertisers success", func(t *testing.T) {
		ctx := context.Background()

		advertisersRepoMock := mocks.NewAdvertisersRepo(t)
		mlscoresRepoMock := mocks.NewMlScoresRepo(t)
		service := NewAdvertisersService(advertisersRepoMock, mlscoresRepoMock)

		// setup mocks
		advertisers := []models.Advertiser{
			{
				Id:   uuid.New(),
				Name: "name 1",
			},
			{
				Id:   uuid.New(),
				Name: "name 2",
			},
			{
				Id:   uuid.New(),
				Name: "name 3",
			},
		}
		expectedAdvertisers := advertisers

		advertisersRepoMock.On("UpsertAdvertisers", ctx, advertisers).Return(expectedAdvertisers, nil).Once()

		// check
		actualAdvertisers, err := service.UpsertAdvertisers(ctx, advertisers)
		require.NoError(t, err)
		require.Equal(t, expectedAdvertisers, actualAdvertisers)
	})

	t.Run("upsert advertisers error", func(t *testing.T) {
		ctx := context.Background()

		advertisersRepoMock := mocks.NewAdvertisersRepo(t)
		mlscoresRepoMock := mocks.NewMlScoresRepo(t)
		service := NewAdvertisersService(advertisersRepoMock, mlscoresRepoMock)

		// setup mocks
		advertisers := []models.Advertiser{
			{
				Id:   uuid.New(),
				Name: "name 1",
			},
			{
				Id:   uuid.New(),
				Name: "name 2",
			},
			{
				Id:   uuid.New(),
				Name: "name 3",
			},
		}
		var expectedAdvertisers []models.Advertiser = nil
		expectedError := errors.New("failed to upsert advertisers")

		advertisersRepoMock.On("UpsertAdvertisers", ctx, advertisers).Return(expectedAdvertisers, expectedError).Once()

		// check
		actualAdvertisers, actualError := service.UpsertAdvertisers(ctx, advertisers)
		require.ErrorIs(t, actualError, expectedError)
		require.Equal(t, expectedAdvertisers, actualAdvertisers)
	})

	t.Run("upsert ml_score success", func(t *testing.T) {
		ctx := context.Background()

		advertisersRepoMock := mocks.NewAdvertisersRepo(t)
		mlscoresRepoMock := mocks.NewMlScoresRepo(t)
		service := NewAdvertisersService(advertisersRepoMock, mlscoresRepoMock)

		// setup mocks
		mlScore := models.MLScore{
			ClientId:     uuid.New(),
			AdvertiserId: uuid.New(),
			Score:        100,
		}

		mlscoresRepoMock.On("UpsertMLScore", ctx, mlScore).Return(nil).Once()

		// check
		err := service.UpsertMLScore(ctx, mlScore)
		require.NoError(t, err)
	})

	t.Run("upsert ml_score client error", func(t *testing.T) {
		ctx := context.Background()

		advertisersRepoMock := mocks.NewAdvertisersRepo(t)
		mlscoresRepoMock := mocks.NewMlScoresRepo(t)
		service := NewAdvertisersService(advertisersRepoMock, mlscoresRepoMock)

		// setup mocks
		mlScore := models.MLScore{
			ClientId:     uuid.New(),
			AdvertiserId: uuid.New(),
			Score:        100,
		}
		expectedError := errors.New("failed to upsert ml_score")

		mlscoresRepoMock.On("UpsertMLScore", ctx, mlScore).Return(expectedError).Once()

		// check
		err := service.UpsertMLScore(ctx, mlScore)
		require.ErrorIs(t, err, expectedError)
	})
}
