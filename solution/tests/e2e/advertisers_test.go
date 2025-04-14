package e2e

import (
	"advertising/tests/helpers"
	"context"
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gavv/httpexpect/v2"
	"github.com/google/uuid"
)

func TestAdvertisers(t *testing.T) {
	ctx := context.Background()
	// advertisingServerUrl := helpers.SetUpInfrastructure(ctx, t, "../../advertising-service/migrations")
	advertisingServerUrl := "http://localhost:8080"

	t.Run("insert advertiser success", func(t *testing.T) {
		e := helpers.ConfigureExpect(t, ctx, advertisingServerUrl)

		n := 20

		advertisers := make([]helpers.JSON, 0, n)
		for range n {
			advertisers = append(advertisers, generateAdvertiser())
		}

		upsertAdvertisersSuccess(e, advertisers...).
			JSON().IsArray().Array().
			IsEqualUnordered(advertisers)

		for _, advertiser := range advertisers {
			getAdvertiser(e, advertiser["advertiser_id"].(uuid.UUID)).
				Expect().
				Status(http.StatusOK).
				JSON().
				IsEqual(advertiser)
		}
	})

	t.Run("update advertisers", func(t *testing.T) {
		e := helpers.ConfigureExpect(t, ctx, advertisingServerUrl)

		n := 20

		advertisersWas := make([]helpers.JSON, 0, n)
		for range n {
			advertisersWas = append(advertisersWas, generateAdvertiser())
		}

		upsertAdvertisersSuccess(e, advertisersWas...)

		advertisersBecome := make([]helpers.JSON, 0, n)
		for _, advertiser := range advertisersWas {
			updatedadvertiser := generateAdvertiser()
			updatedadvertiser["advertiser_id"] = advertiser["advertiser_id"]
			advertisersBecome = append(advertisersBecome, updatedadvertiser)
		}

		upsertAdvertisersSuccess(e, advertisersBecome...).
			JSON().IsArray().Array().
			IsEqualUnordered(advertisersBecome)

		for _, advertiser := range advertisersBecome {
			getAdvertiser(e, advertiser["advertiser_id"].(uuid.UUID)).
				Expect().
				Status(http.StatusOK).
				JSON().
				IsEqual(advertiser)
		}
	})

	t.Run("present advertiser several times", func(t *testing.T) {
		e := helpers.ConfigureExpect(t, ctx, advertisingServerUrl)

		advertiser := generateAdvertiser()

		// present same advertiser
		upsertAdvertisersSuccess(e, advertiser, advertiser, advertiser).
			JSON().
			Array().
			ContainsOnly(advertiser)

		// present same advertiser with another fields
		advertiser2 := generateAdvertiser()
		advertiser2["advertiser_id"] = advertiser["advertiser_id"]

		upsertAdvertisersSuccess(e, advertiser, advertiser2, advertiser, advertiser2).
			JSON().
			Array().
			ContainsOnly(advertiser2)

		getAdvertiserSuccess(e, advertiser["advertiser_id"].(uuid.UUID)).
			JSON().
			IsEqual(advertiser2)
	})

	t.Run("get non-existent advertiser", func(t *testing.T) {
		e := helpers.ConfigureExpect(t, ctx, advertisingServerUrl)

		getAdvertiser(e, uuid.New()).
			Expect().
			Status(http.StatusNotFound)
	})

}

func generateAdvertiser() helpers.JSON {
	return helpers.JSON{
		"advertiser_id": uuid.New(),
		"name":          gofakeit.Company(),
	}
}

func upsertAdvertisersSuccess(e *httpexpect.Expect, advertisers ...helpers.JSON) *httpexpect.Response {
	return e.POST("/advertisers/bulk").
		WithJSON(advertisers).
		Expect().
		Status(http.StatusCreated)
}

func getAdvertiser(e *httpexpect.Expect, id uuid.UUID) *httpexpect.Request {
	return e.GET("/advertisers/{advertiser_id}", id)
}

func getAdvertiserSuccess(e *httpexpect.Expect, id uuid.UUID) *httpexpect.Response {
	return getAdvertiser(e, id).
		Expect().
		Status(http.StatusOK)
}
