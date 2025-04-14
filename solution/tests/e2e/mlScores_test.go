package e2e

import (
	"advertising/tests/helpers"
	"context"
	"math"
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gavv/httpexpect/v2"
	"github.com/google/uuid"
)

func TestMlScores(t *testing.T) {
	ctx := context.Background()
	// advertisingServerUrl := helpers.SetUpInfrastructure(ctx, t, "../../advertising-service/migrations")
	advertisingServerUrl := "http://localhost:8080"

	t.Run("insert ml_score success", func(t *testing.T) {
		e := helpers.ConfigureExpect(t, ctx, advertisingServerUrl)

		client := generateClient()
		upsertClientsSuccess(e, client)

		advertiser := generateAdvertiser()
		upsertAdvertisersSuccess(e, advertiser)

		upsertMLScoreSuccess(e, helpers.JSON{
			"client_id":     client["client_id"],
			"advertiser_id": advertiser["advertiser_id"],
			"score":         gofakeit.IntRange(0, math.MaxInt32),
		})
	})

	t.Run("update ml_score success", func(t *testing.T) {
		e := helpers.ConfigureExpect(t, ctx, advertisingServerUrl)

		client := generateClient()
		upsertClientsSuccess(e, client)

		advertiser := generateAdvertiser()
		upsertAdvertisersSuccess(e, advertiser)

		upsertMLScoreSuccess(e, helpers.JSON{
			"client_id":     client["client_id"],
			"advertiser_id": advertiser["advertiser_id"],
			"score":         gofakeit.Int32(),
		})

		upsertMLScoreSuccess(e, helpers.JSON{
			"client_id":     client["client_id"],
			"advertiser_id": advertiser["advertiser_id"],
			"score":         gofakeit.Int32(),
		})
	})

	t.Run("insert invalid ml_score", func(t *testing.T) {
		e := helpers.ConfigureExpect(t, ctx, advertisingServerUrl)

		upsertMLScore(e, helpers.JSON{
			"score": -1,
		}).
			Expect().
			Status(http.StatusBadRequest)
	})

	t.Run("insert ml_score with non-existent client", func(t *testing.T) {
		e := helpers.ConfigureExpect(t, ctx, advertisingServerUrl)

		advertiser := generateAdvertiser()
		upsertAdvertisersSuccess(e, advertiser)

		upsertMLScore(e, helpers.JSON{
			"client_id":     uuid.New(),
			"advertiser_id": advertiser["advertiser_id"],
			"score":         gofakeit.IntRange(0, math.MaxInt32),
		}).
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("insert ml_score with non-existent advertiser", func(t *testing.T) {
		e := helpers.ConfigureExpect(t, ctx, advertisingServerUrl)

		client := generateClient()
		upsertClientsSuccess(e, client)

		upsertMLScore(e, helpers.JSON{
			"client_id":     client["client_id"],
			"advertiser_id": uuid.New(),
			"score":         gofakeit.IntRange(0, math.MaxInt32),
		}).
			Expect().
			Status(http.StatusNotFound)
	})
}

func upsertMLScore(e *httpexpect.Expect, mlScore helpers.JSON) *httpexpect.Request {
	return e.POST("/ml-scores").WithJSON(mlScore)
}

func upsertMLScoreSuccess(e *httpexpect.Expect, mlScore helpers.JSON) *httpexpect.Response {
	return upsertMLScore(e, mlScore).
		Expect().
		Status(http.StatusOK)
}
