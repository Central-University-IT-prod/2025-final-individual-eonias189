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

func TestGetAdForClient(t *testing.T) {
	ctx := context.Background()
	// advertisingServerUrl := helpers.SetUpInfrastructure(ctx, t, "../../advertising-service/migrations")
	advertisingServerUrl := "http://localhost:8080"

	// create advertiser
	e := helpers.ConfigureExpect(t, ctx, advertisingServerUrl)

	advertiser := generateAdvertiser()
	upsertAdvertisersSuccess(e, advertiser)
	advertiserId := advertiser["advertiser_id"].(uuid.UUID)

	t.Run("get ad success", func(t *testing.T) {
		e := helpers.ConfigureExpect(t, ctx, advertisingServerUrl)

		// set day
		advanceDaySuccess(e, pointer(0))

		// create campaign
		campaign := generateCampaign(advertiserId, helpers.JSON{
			"gender":   "ALL",
			"age_from": 22,
		})
		campaign["start_date"] = 0
		campaign["impressions_limit"] = 1000
		campaign["clicks_limit"] = 900
		campaignIdStr := createCampaignSuccess(e, campaign).JSON().Object().Value("campaign_id").String().Raw()
		campaignId := uuid.MustParse(campaignIdStr)
		t.Cleanup(func() {
			deleteCapaignSuccess(e, advertiserId, campaignId)
		})

		// create client
		client := generateClient()
		client["age"] = 27
		upsertClientsSuccess(e, client)
		clientId := client["client_id"].(uuid.UUID)

		// shoud get created campaign
		getAdForClientSuccess(e, clientId).
			JSON().
			IsObject().
			Object().
			HasValue("ad_id", campaignId).
			HasValue("advertiser_id", advertiserId).
			HasValue("ad_title", campaign["ad_title"]).
			HasValue("ad_text", campaign["ad_text"])
	})

	t.Run("check doesn`t return not active campaign", func(t *testing.T) {
		e := helpers.ConfigureExpect(t, ctx, advertisingServerUrl)

		// set day
		advanceDaySuccess(e, pointer(0))

		// create campaigns

		// campaign has not started
		campaign1 := generateCampaign(advertiserId, helpers.JSON{})
		campaign1["start_date"] = 15
		campaign1["end_date"] = 20
		campaign1["impressions_limit"] = 1000
		campaign1["clicks_limit"] = 900
		campaign1IdStr := createCampaignSuccess(e, campaign1).JSON().Object().Value("campaign_id").String().Raw()
		campaign1Id := uuid.MustParse(campaign1IdStr)
		t.Cleanup(func() {
			deleteCapaignSuccess(e, advertiserId, campaign1Id)
		})

		// campaign already finished
		campaign2 := generateCampaign(advertiserId, helpers.JSON{})
		campaign2["start_date"] = 2
		campaign2["end_date"] = 7
		campaign2["impressions_limit"] = 1000
		campaign2["clicks_limit"] = 900
		campaign2IdStr := createCampaignSuccess(e, campaign2).JSON().Object().Value("campaign_id").String().Raw()
		campaign2Id := uuid.MustParse(campaign2IdStr)
		t.Cleanup(func() {
			deleteCapaignSuccess(e, advertiserId, campaign2Id)
		})

		// update day
		advanceDaySuccess(e, pointer(10))

		// create client
		client := generateClient()
		upsertClientsSuccess(e, client)
		clientId := client["client_id"].(uuid.UUID)

		// shoud return 404
		getAdForClient(e, clientId).
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("check doesn`t return campaign that doesn`t match client", func(t *testing.T) {
		e := helpers.ConfigureExpect(t, ctx, advertisingServerUrl)

		// set day
		advanceDaySuccess(e, pointer(0))

		// create client
		clientId := uuid.New()
		client := helpers.JSON{
			"client_id": clientId,
			"login":     gofakeit.Username(),
			"gender":    "MALE",
			"age":       42,
			"location":  "Rostov",
		}
		upsertClientsSuccess(e, client)

		// create campaigns

		// campaign with another gender targeting
		campaign1 := generateCampaign(advertiserId, helpers.JSON{
			"gender": "FEMALE",
		})
		campaign1["start_date"] = 0
		campaign1["end_date"] = 10
		campaign1["impressions_limit"] = 1000
		campaign1["clicks_limit"] = 900
		campaign1IdStr := createCampaignSuccess(e, campaign1).JSON().Object().Value("campaign_id").String().Raw()
		campaign1Id := uuid.MustParse(campaign1IdStr)
		t.Cleanup(func() {
			deleteCapaignSuccess(e, advertiserId, campaign1Id)
		})

		// campaign with another location targeting
		campaign2 := generateCampaign(advertiserId, helpers.JSON{
			"location": "Kemerovo",
		})
		campaign2["start_date"] = 0
		campaign2["end_date"] = 10
		campaign2["impressions_limit"] = 1000
		campaign2["clicks_limit"] = 900
		campaign2IdStr := createCampaignSuccess(e, campaign2).JSON().Object().Value("campaign_id").String().Raw()
		campaign2Id := uuid.MustParse(campaign2IdStr)
		t.Cleanup(func() {
			deleteCapaignSuccess(e, advertiserId, campaign2Id)
		})

		// campaigns with another age targeting
		campaign3 := generateCampaign(advertiserId, helpers.JSON{
			"age_from": 43,
		})
		campaign3["start_date"] = 0
		campaign3["end_date"] = 10
		campaign3["impressions_limit"] = 1000
		campaign3["clicks_limit"] = 900
		campaign3IdStr := createCampaignSuccess(e, campaign3).JSON().Object().Value("campaign_id").String().Raw()
		campaign3Id := uuid.MustParse(campaign3IdStr)
		t.Cleanup(func() {
			deleteCapaignSuccess(e, advertiserId, campaign3Id)
		})

		campaign4 := generateCampaign(advertiserId, helpers.JSON{
			"age_to": 41,
		})
		campaign4["start_date"] = 0
		campaign4["end_date"] = 10
		campaign4["impressions_limit"] = 1000
		campaign4["clicks_limit"] = 900
		campaign4IdStr := createCampaignSuccess(e, campaign4).JSON().Object().Value("campaign_id").String().Raw()
		campaign4Id := uuid.MustParse(campaign4IdStr)
		t.Cleanup(func() {
			deleteCapaignSuccess(e, advertiserId, campaign4Id)
		})

		// shoud return 404
		getAdForClient(e, clientId).
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("check if doesn`t return already impressed campaign", func(t *testing.T) {
		e := helpers.ConfigureExpect(t, ctx, advertisingServerUrl)

		// set day
		advanceDaySuccess(e, pointer(0))

		// create client
		client := generateClient()
		clientId := client["client_id"].(uuid.UUID)
		upsertClientsSuccess(e, client)

		// create campaign
		campaign := generateCampaign(advertiserId, helpers.JSON{})
		campaign["start_date"] = 0
		campaign["end_date"] = 10
		campaign["impressions_limit"] = 1000
		campaign["clicks_limit"] = 900
		campaignIdStr := createCampaignSuccess(e, campaign).JSON().Object().Value("campaign_id").String().Raw()
		campaignId := uuid.MustParse(campaignIdStr)
		t.Cleanup(func() {
			deleteCapaignSuccess(e, advertiserId, campaignId)
		})

		// shoud get created campaign
		getAdForClientSuccess(e, clientId).JSON().
			Object().
			HasValue("ad_id", campaignId).
			HasValue("advertiser_id", advertiserId)

		// shoud get 404
		getAdForClient(e, clientId).
			Expect().
			Status(http.StatusNotFound)

	})

	t.Run("get ad for non-existent client", func(t *testing.T) {
		e := helpers.ConfigureExpect(t, ctx, advertisingServerUrl)

		getAdForClient(e, uuid.New()).
			Expect().
			Status(http.StatusNotFound)
	})
}

func TestRecordClick(t *testing.T) {
	ctx := context.Background()
	// advertisingServerUrl := helpers.SetUpInfrastructure(ctx, t, "../../advertising-service/migrations")
	advertisingServerUrl := "http://localhost:8080"

	// create advertiser
	e := helpers.ConfigureExpect(t, ctx, advertisingServerUrl)

	advertiser := generateAdvertiser()
	advertiserId := advertiser["advertiser_id"].(uuid.UUID)
	upsertAdvertisersSuccess(e, advertiser)

	t.Run("record click success", func(t *testing.T) {
		e := helpers.ConfigureExpect(t, ctx, advertisingServerUrl)

		// set day
		advanceDaySuccess(e, pointer(0))

		// create campaign
		campaign := generateCampaign(advertiserId, helpers.JSON{})
		campaign["start_date"] = 0
		campaign["end_date"] = 10
		campaign["impressions_limit"] = 1000
		campaign["clicks_limit"] = 900
		campaignIdStr := createCampaignSuccess(e, campaign).JSON().Object().Value("campaign_id").String().Raw()
		campaignId := uuid.MustParse(campaignIdStr)
		t.Cleanup(func() {
			deleteCapaignSuccess(e, advertiserId, campaignId)
		})

		// create client
		client := generateClient()
		clientId := client["client_id"].(uuid.UUID)
		upsertClientsSuccess(e, client)

		// impress
		getAdForClientSuccess(e, clientId).
			JSON().
			Object().
			HasValue("ad_id", campaignId).
			HasValue("advertiser_id", advertiserId)

		// click
		recordClickSuccess(e, campaignId, clientId)
	})

	t.Run("click several times", func(t *testing.T) {
		e := helpers.ConfigureExpect(t, ctx, advertisingServerUrl)

		// set day
		advanceDaySuccess(e, pointer(0))

		// create campaign
		campaign := generateCampaign(advertiserId, helpers.JSON{})
		campaign["start_date"] = 0
		campaign["end_date"] = 10
		campaign["impressions_limit"] = 1000
		campaign["clicks_limit"] = 900
		campaignIdStr := createCampaignSuccess(e, campaign).JSON().Object().Value("campaign_id").String().Raw()
		campaignId := uuid.MustParse(campaignIdStr)
		t.Cleanup(func() {
			deleteCapaignSuccess(e, advertiserId, campaignId)
		})

		// create client
		client := generateClient()
		clientId := client["client_id"].(uuid.UUID)
		upsertClientsSuccess(e, client)

		// impress
		getAdForClientSuccess(e, clientId).
			JSON().
			Object().
			HasValue("ad_id", campaignId).
			HasValue("advertiser_id", advertiserId)

		// click
		recordClickSuccess(e, campaignId, clientId)

		// click again
		recordClickSuccess(e, campaignId, clientId)
	})

	t.Run("click without impression", func(t *testing.T) {
		e := helpers.ConfigureExpect(t, ctx, advertisingServerUrl)

		// set day
		advanceDaySuccess(e, pointer(0))

		// create campaign
		campaign := generateCampaign(advertiserId, helpers.JSON{})
		campaign["start_date"] = 0
		campaign["end_date"] = 10
		campaign["impressions_limit"] = 1000
		campaign["clicks_limit"] = 900
		campaignIdStr := createCampaignSuccess(e, campaign).JSON().Object().Value("campaign_id").String().Raw()
		campaignId := uuid.MustParse(campaignIdStr)
		t.Cleanup(func() {
			deleteCapaignSuccess(e, advertiserId, campaignId)
		})

		// create client
		client := generateClient()
		clientId := client["client_id"].(uuid.UUID)
		upsertClientsSuccess(e, client)

		// click shoud return 404
		recordClick(e, campaignId, clientId).
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("click with non-existent client", func(t *testing.T) {
		e := helpers.ConfigureExpect(t, ctx, advertisingServerUrl)

		// set day
		advanceDaySuccess(e, pointer(0))

		// create campaign
		campaign := generateCampaign(advertiserId, helpers.JSON{})
		campaign["start_date"] = 0
		campaign["end_date"] = 10
		campaign["impressions_limit"] = 1000
		campaign["clicks_limit"] = 900
		campaignIdStr := createCampaignSuccess(e, campaign).JSON().Object().Value("campaign_id").String().Raw()
		campaignId := uuid.MustParse(campaignIdStr)
		t.Cleanup(func() {
			deleteCapaignSuccess(e, advertiserId, campaignId)
		})

		// click shoud return 404
		recordClick(e, campaignId, uuid.New()).
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("click with non-existent ad", func(t *testing.T) {
		e := helpers.ConfigureExpect(t, ctx, advertisingServerUrl)

		// create client
		client := generateClient()
		clientId := client["client_id"].(uuid.UUID)
		upsertClientsSuccess(e, client)

		// click shoud return 404
		recordClick(e, uuid.New(), clientId).
			Expect().
			Status(http.StatusNotFound)
	})
}

func getAdForClient(e *httpexpect.Expect, clientId uuid.UUID) *httpexpect.Request {
	return e.GET("/ads").WithQuery("client_id", clientId)
}

func getAdForClientSuccess(e *httpexpect.Expect, clientId uuid.UUID) *httpexpect.Response {
	return getAdForClient(e, clientId).
		Expect().
		Status(http.StatusOK)
}

func recordClick(e *httpexpect.Expect, adId, clientId uuid.UUID) *httpexpect.Request {
	return e.POST("/ads/{adId}/click", adId).
		WithJSON(helpers.JSON{
			"client_id": clientId,
		})
}

func recordClickSuccess(e *httpexpect.Expect, adId, clientId uuid.UUID) *httpexpect.Response {
	return recordClick(e, adId, clientId).
		Expect().
		Status(http.StatusNoContent)
}
