package e2e

import (
	"advertising/tests/helpers"
	"context"
	"math"
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCampaignStats(t *testing.T) {
	ctx := context.Background()
	advertisingServerUrl := "http://localhost:8080"
	// advertisingServerUrl := helpers.SetUpInfrastructure(ctx, t, "../../advertising-service/migrations")

	t.Run("get campaign stats success", func(t *testing.T) {
		e := helpers.ConfigureExpect(t, ctx, advertisingServerUrl)

		// set day
		advanceDaySuccess(e, pointer(0))

		// create advertiser
		advertiser := generateAdvertiser()
		advertiserId := advertiser["advertiser_id"].(uuid.UUID)
		upsertAdvertisersSuccess(e, advertiser)

		// create campaign
		campaign := generateCampaign(advertiserId, helpers.JSON{})
		campaign["start_date"] = 0
		campaign["end_date"] = 10
		campaign["impressions_limit"] = 1000
		campaign["clicks_limit"] = 1000
		campaignIdStr := createCampaignSuccess(e, campaign).JSON().Object().Value("campaign_id").String().Raw()
		campaignId := uuid.MustParse(campaignIdStr)
		t.Cleanup(func() {
			deleteCapaignSuccess(e, advertiserId, campaignId)
		})

		costPerImpression := float64(campaign["cost_per_impression"].(float32))
		costPerClick := float64(campaign["cost_per_click"].(float32))

		// create clients
		clientsNumber := 5
		clients := make([]helpers.JSON, 0, clientsNumber)
		clientsIds := make([]uuid.UUID, 0, clientsNumber)

		for range clientsNumber {
			client := generateClient()
			clientId := client["client_id"].(uuid.UUID)

			clients = append(clients, client)
			clientsIds = append(clientsIds, clientId)
		}

		upsertClientsSuccess(e, clients...)

		// record impressions and clicks

		// day 1: 2 impressions and 1 click
		// day 2: 1 impression 1 click
		// day 3: 2 impression 2 clicks
		// total: 5 impressions, 4 clicks

		// day 1
		advanceDaySuccess(e, pointer(1))
		getAdForClientSuccess(e, clientsIds[0])
		getAdForClientSuccess(e, clientsIds[1])
		recordClickSuccess(e, campaignId, clientsIds[0])

		// day 2
		advanceDaySuccess(e, pointer(2))
		getAdForClientSuccess(e, clientsIds[2])
		recordClickSuccess(e, campaignId, clientsIds[1])

		// day 3
		advanceDaySuccess(e, pointer(3))
		getAdForClientSuccess(e, clientsIds[3])
		getAdForClientSuccess(e, clientsIds[4])
		recordClickSuccess(e, campaignId, clientsIds[2])
		recordClickSuccess(e, campaignId, clientsIds[3])

		// check result
		expected := helpers.JSON{
			"impressions_count": 5,
			"clicks_count":      4,
			"conversion":        float64(4) / float64(5) * 100,
			"spent_impressions": 5 * costPerImpression,
			"spent_clicks":      4 * costPerClick,
			"spent_total":       5*costPerImpression + 4*costPerClick,
		}

		var actual helpers.JSON
		getCampaignStatsSuccess(e, campaignId).JSON().
			IsObject().
			Object().Decode(&actual)

		checkStats(t, expected, actual)
	})

	t.Run("get campaign stats zero values", func(t *testing.T) {
		e := helpers.ConfigureExpect(t, ctx, advertisingServerUrl)

		// set day
		advanceDaySuccess(e, pointer(0))

		// create advertiser
		advertiser := generateAdvertiser()
		advertiserId := advertiser["advertiser_id"].(uuid.UUID)
		upsertAdvertisersSuccess(e, advertiser)

		// create campaign
		campaign := generateCampaign(advertiserId, helpers.JSON{})
		campaign["start_date"] = 0
		campaign["end_date"] = 10
		campaign["impressions_limit"] = 1000
		campaign["clicks_limit"] = 1000
		campaignIdStr := createCampaignSuccess(e, campaign).JSON().Object().Value("campaign_id").String().Raw()
		campaignId := uuid.MustParse(campaignIdStr)
		t.Cleanup(func() {
			deleteCapaignSuccess(e, advertiserId, campaignId)
		})

		getCampaignStatsSuccess(e, campaignId).
			JSON().
			IsObject().
			Object().
			HasValue("impressions_count", 0).
			HasValue("clicks_count", 0).
			HasValue("conversion", 0).
			HasValue("spent_impressions", 0).
			HasValue("spent_clicks", 0).
			HasValue("spent_total", 0)
	})

	t.Run("get stats for non-existent campaign", func(t *testing.T) {
		e := helpers.ConfigureExpect(t, ctx, advertisingServerUrl)

		getCampaignStats(e, uuid.New()).
			Expect().
			Status(http.StatusNotFound)
	})
}

func TestGetCampaignStatsDaily(t *testing.T) {
	ctx := context.Background()
	advertisingServerUrl := "http://localhost:8080"
	// advertisingServerUrl := helpers.SetUpInfrastructure(ctx, t, "../../advertising-service/migrations")

	t.Run("get campaign stats daily success", func(t *testing.T) {
		e := helpers.ConfigureExpect(t, ctx, advertisingServerUrl)

		// set day
		advanceDaySuccess(e, pointer(0))

		// create advertiser
		advertiser := generateAdvertiser()
		advertiserId := advertiser["advertiser_id"].(uuid.UUID)
		upsertAdvertisersSuccess(e, advertiser)

		// create campaign
		campaign := generateCampaign(advertiserId, helpers.JSON{})
		campaign["start_date"] = 0
		campaign["end_date"] = 10
		campaign["impressions_limit"] = 1000
		campaign["clicks_limit"] = 1000
		campaignIdStr := createCampaignSuccess(e, campaign).JSON().Object().Value("campaign_id").String().Raw()
		campaignId := uuid.MustParse(campaignIdStr)
		t.Cleanup(func() {
			deleteCapaignSuccess(e, advertiserId, campaignId)
		})

		costPerImpression := float64(campaign["cost_per_impression"].(float32))
		costPerClick := float64(campaign["cost_per_click"].(float32))

		// create clients
		clientsNumber := 5
		clients := make([]helpers.JSON, 0, clientsNumber)
		clientsIds := make([]uuid.UUID, 0, clientsNumber)

		for range clientsNumber {
			client := generateClient()
			clientId := client["client_id"].(uuid.UUID)

			clients = append(clients, client)
			clientsIds = append(clientsIds, clientId)
		}

		upsertClientsSuccess(e, clients...)

		// record impressions and clicks

		// day 1: 2 impressions and 1 click
		// day 2: 1 impression 1 click
		// day 3: 2 impression 2 clicks
		// total: 5 impressions, 4 clicks

		// day 1
		advanceDaySuccess(e, pointer(1))
		getAdForClientSuccess(e, clientsIds[0])
		getAdForClientSuccess(e, clientsIds[1])
		recordClickSuccess(e, campaignId, clientsIds[0])

		// day 2
		advanceDaySuccess(e, pointer(2))
		getAdForClientSuccess(e, clientsIds[2])
		recordClickSuccess(e, campaignId, clientsIds[1])

		// day 3
		advanceDaySuccess(e, pointer(3))
		getAdForClientSuccess(e, clientsIds[3])
		getAdForClientSuccess(e, clientsIds[4])
		recordClickSuccess(e, campaignId, clientsIds[2])
		recordClickSuccess(e, campaignId, clientsIds[3])

		expected := []helpers.JSON{
			{
				"date":              1,
				"impressions_count": 2,
				"clicks_count":      1,
				"conversion":        float64(50),
				"spent_impressions": float64(2) * costPerImpression,
				"spent_clicks":      float64(1) * costPerClick,
				"spent_total":       float64(2)*costPerImpression + float64(1)*costPerClick,
			},
			{
				"date":              2,
				"impressions_count": 1,
				"clicks_count":      1,
				"conversion":        float64(100),
				"spent_impressions": costPerImpression,
				"spent_clicks":      costPerClick,
				"spent_total":       costPerImpression + costPerClick,
			},
			{
				"date":              3,
				"impressions_count": 2,
				"clicks_count":      2,
				"conversion":        float64(100),
				"spent_impressions": float64(2) * costPerImpression,
				"spent_clicks":      float64(2) * costPerClick,
				"spent_total":       float64(2)*costPerImpression + float64(2)*costPerClick,
			},
		}

		// check result
		var actual []helpers.JSON
		getCampaignStatsDailySuccess(e, campaignId).
			JSON().
			IsArray().
			Array().Decode(&actual)
		checkStatsDaily(t, expected, actual)
	})

	t.Run("get campaign stats daily for campaign with no impressions", func(t *testing.T) {
		e := helpers.ConfigureExpect(t, ctx, advertisingServerUrl)

		// set day
		advanceDaySuccess(e, pointer(0))

		// create advertiser
		advertiser := generateAdvertiser()
		advertiserId := advertiser["advertiser_id"].(uuid.UUID)
		upsertAdvertisersSuccess(e, advertiser)

		// create campaign
		campaign := generateCampaign(advertiserId, helpers.JSON{})
		campaign["start_date"] = 0
		campaign["end_date"] = 10
		campaign["impressions_limit"] = 1000
		campaign["clicks_limit"] = 1000
		campaignIdStr := createCampaignSuccess(e, campaign).JSON().Object().Value("campaign_id").String().Raw()
		campaignId := uuid.MustParse(campaignIdStr)
		t.Cleanup(func() {
			deleteCapaignSuccess(e, advertiserId, campaignId)
		})

		getCampaignStatsDailySuccess(e, campaignId).
			JSON().
			IsArray().
			Array().
			Length().
			IsEqual(0)
	})

	t.Run("get daily stats for non-existent campaign", func(t *testing.T) {
		e := helpers.ConfigureExpect(t, ctx, advertisingServerUrl)

		getCampaignStatsDaily(e, uuid.New()).
			Expect().
			Status(http.StatusNotFound)
	})
}

func TestGetAdvertiserStats(t *testing.T) {
	ctx := context.Background()
	advertisingServerUrl := "http://localhost:8080"
	// advertisingServerUrl := helpers.SetUpInfrastructure(ctx, t, "../../advertising-service/migrations")

	t.Run("get advertiser stats success", func(t *testing.T) {
		e := helpers.ConfigureExpect(t, ctx, advertisingServerUrl)

		// set day
		advanceDaySuccess(e, pointer(0))

		// create advertiser
		advertiser := generateAdvertiser()
		advertiserId := advertiser["advertiser_id"].(uuid.UUID)
		upsertAdvertisersSuccess(e, advertiser)

		// create campaigns
		campaignsNumer := 3

		campaignsMap := map[uuid.UUID]helpers.JSON{}

		for range campaignsNumer {
			campaign := generateCampaign(advertiserId, helpers.JSON{})
			campaign["start_date"] = 0
			campaign["end_date"] = 10
			campaign["impressions_limit"] = 1000
			campaign["clicks_limit"] = 1000
			campaignIdStr := createCampaignSuccess(e, campaign).JSON().Object().Value("campaign_id").String().Raw()
			campaignId := uuid.MustParse(campaignIdStr)
			t.Cleanup(func() {
				deleteCapaignSuccess(e, advertiserId, campaignId)
			})

			campaignsMap[campaignId] = campaign
		}

		// create clients
		clientsNumber := 3
		clients := make([]helpers.JSON, 0, clientsNumber)
		clientsIds := make([]uuid.UUID, 0, clientsNumber)

		for range clientsNumber {
			client := generateClient()
			clientId := client["client_id"].(uuid.UUID)

			clients = append(clients, client)
			clientsIds = append(clientsIds, clientId)
		}

		upsertClientsSuccess(e, clients...)

		// record impressions and clicks

		// day 1: 3 impressions 1 click
		// day 2: 2 impressions 2 clicks
		// day 3: 4 impressions 3 clicks

		var (
			impressionsCount = 9
			clicksCount      = 6
			conversion       = float64(6) / float64(9) * 100
			spentImpressions float64
			spentClicks      float64
			spentTotal       float64
		)

		// day 1
		advanceDaySuccess(e, pointer(1))

		campaignIdStr := getAdForClientSuccess(e, clientsIds[0]).JSON().Object().Value("ad_id").String().Raw()
		campaignId := uuid.MustParse(campaignIdStr)
		campaign := campaignsMap[campaignId]
		spentImpressions += float64(campaign["cost_per_impression"].(float32))

		recordClickSuccess(e, campaignId, clientsIds[0])
		spentClicks += float64(campaign["cost_per_click"].(float32))

		campaignIdStr = getAdForClientSuccess(e, clientsIds[0]).JSON().Object().Value("ad_id").String().Raw()
		campaignId = uuid.MustParse(campaignIdStr)
		campaign = campaignsMap[campaignId]
		spentImpressions += float64(campaign["cost_per_impression"].(float32))

		campaignIdStr = getAdForClientSuccess(e, clientsIds[1]).JSON().Object().Value("ad_id").String().Raw()
		campaignId = uuid.MustParse(campaignIdStr)
		campaign = campaignsMap[campaignId]
		spentImpressions += float64(campaign["cost_per_impression"].(float32))

		// day 2
		advanceDaySuccess(e, pointer(2))

		recordClickSuccess(e, campaignId, clientsIds[1])
		spentClicks += float64(campaign["cost_per_click"].(float32))

		campaignIdStr = getAdForClientSuccess(e, clientsIds[0]).JSON().Object().Value("ad_id").String().Raw()
		campaignId = uuid.MustParse(campaignIdStr)
		campaign = campaignsMap[campaignId]
		spentImpressions += float64(campaign["cost_per_impression"].(float32))

		campaignIdStr = getAdForClientSuccess(e, clientsIds[1]).JSON().Object().Value("ad_id").String().Raw()
		campaignId = uuid.MustParse(campaignIdStr)
		campaign = campaignsMap[campaignId]
		spentImpressions += float64(campaign["cost_per_impression"].(float32))

		recordClickSuccess(e, campaignId, clientsIds[1])
		spentClicks += float64(campaign["cost_per_click"].(float32))

		// day 3
		advanceDaySuccess(e, pointer(3))

		campaignIdStr = getAdForClientSuccess(e, clientsIds[1]).JSON().Object().Value("ad_id").String().Raw()
		campaignId = uuid.MustParse(campaignIdStr)
		campaign = campaignsMap[campaignId]
		spentImpressions += float64(campaign["cost_per_impression"].(float32))

		recordClickSuccess(e, campaignId, clientsIds[1])
		spentClicks += float64(campaign["cost_per_click"].(float32))

		campaignIdStr = getAdForClientSuccess(e, clientsIds[2]).JSON().Object().Value("ad_id").String().Raw()
		campaignId = uuid.MustParse(campaignIdStr)
		campaign = campaignsMap[campaignId]
		spentImpressions += float64(campaign["cost_per_impression"].(float32))

		campaignIdStr = getAdForClientSuccess(e, clientsIds[2]).JSON().Object().Value("ad_id").String().Raw()
		campaignId = uuid.MustParse(campaignIdStr)
		campaign = campaignsMap[campaignId]
		spentImpressions += float64(campaign["cost_per_impression"].(float32))

		recordClickSuccess(e, campaignId, clientsIds[2])
		spentClicks += float64(campaign["cost_per_click"].(float32))

		campaignIdStr = getAdForClientSuccess(e, clientsIds[2]).JSON().Object().Value("ad_id").String().Raw()
		campaignId = uuid.MustParse(campaignIdStr)
		campaign = campaignsMap[campaignId]
		spentImpressions += float64(campaign["cost_per_impression"].(float32))

		recordClickSuccess(e, campaignId, clientsIds[2])
		spentClicks += float64(campaign["cost_per_click"].(float32))

		spentTotal = spentImpressions + spentClicks

		// check

		expected := helpers.JSON{
			"impressions_count": impressionsCount,
			"clicks_count":      clicksCount,
			"conversion":        conversion,
			"spent_impressions": spentImpressions,
			"spent_clicks":      spentClicks,
			"spent_total":       spentTotal,
		}

		var actual helpers.JSON
		getAdvertiserStatsSuccess(e, advertiserId).
			JSON().
			IsObject().
			Object().
			Decode(&actual)

		checkStats(t, expected, actual)
	})

	t.Run("get stats for advertiser with no campaigns", func(t *testing.T) {
		e := helpers.ConfigureExpect(t, ctx, advertisingServerUrl)

		// set day
		advanceDaySuccess(e, pointer(0))

		// create advertiser
		advertiser := generateAdvertiser()
		advertiserId := advertiser["advertiser_id"].(uuid.UUID)
		upsertAdvertisersSuccess(e, advertiser)

		getAdvertiserStatsSuccess(e, advertiserId).
			JSON().
			IsObject().
			Object().
			HasValue("impressions_count", 0).
			HasValue("clicks_count", 0).
			HasValue("conversion", 0).
			HasValue("spent_impressions", 0).
			HasValue("spent_clicks", 0).
			HasValue("spent_total", 0)
	})

	t.Run("get stats for non-existent advertiser", func(t *testing.T) {
		e := helpers.ConfigureExpect(t, ctx, advertisingServerUrl)

		getAdvertiserStats(e, uuid.New()).
			Expect().
			Status(http.StatusNotFound)
	})
}

func TestGetAdvertiserStatsDaily(t *testing.T) {
	ctx := context.Background()
	// advertisingServerUrl := helpers.SetUpInfrastructure(ctx, t, "../../advertising-service/migrations")
	advertisingServerUrl := "http://localhost:8080"

	t.Run("get advertiser stats daily success", func(t *testing.T) {
		e := helpers.ConfigureExpect(t, ctx, advertisingServerUrl)

		// set day
		advanceDaySuccess(e, pointer(0))

		// create advertiser
		advertiser := generateAdvertiser()
		advertiserId := advertiser["advertiser_id"].(uuid.UUID)
		upsertAdvertisersSuccess(e, advertiser)

		// create campaigns
		campaignsNumer := 3

		campaignsMap := map[uuid.UUID]helpers.JSON{}

		for range campaignsNumer {
			campaign := generateCampaign(advertiserId, helpers.JSON{})
			campaign["start_date"] = 0
			campaign["end_date"] = 10
			campaign["impressions_limit"] = 1000
			campaign["clicks_limit"] = 1000
			campaignIdStr := createCampaignSuccess(e, campaign).JSON().Object().Value("campaign_id").String().Raw()
			campaignId := uuid.MustParse(campaignIdStr)
			t.Cleanup(func() {
				deleteCapaignSuccess(e, advertiserId, campaignId)
			})

			campaignsMap[campaignId] = campaign
		}

		// create clients
		clientsNumber := 3
		clients := make([]helpers.JSON, 0, clientsNumber)
		clientsIds := make([]uuid.UUID, 0, clientsNumber)

		for range clientsNumber {
			client := generateClient()
			clientId := client["client_id"].(uuid.UUID)

			clients = append(clients, client)
			clientsIds = append(clientsIds, clientId)
		}

		upsertClientsSuccess(e, clients...)

		// record impressions and clicks

		// day 1: 3 impressions 1 click
		// day 2: 2 impressions 2 clicks
		// day 3: 4 impressions 3 clicks

		expected := []helpers.JSON{
			{
				"date":              1,
				"impressions_count": 3,
				"clicks_count":      1,
				"conversion":        float64(1) / float64(3) * 100,
				"spent_impressions": float64(0),
				"spent_clicks":      float64(0),
				"spent_total":       float64(0),
			},
			{
				"date":              2,
				"impressions_count": 2,
				"clicks_count":      2,
				"conversion":        float64(100),
				"spent_impressions": float64(0),
				"spent_clicks":      float64(0),
				"spent_total":       float64(0),
			},
			{
				"date":              3,
				"impressions_count": 4,
				"clicks_count":      3,
				"conversion":        float64(3) / float64(4) * 100,
				"spent_impressions": float64(0),
				"spent_clicks":      float64(0),
				"spent_total":       float64(0),
			},
		}

		var (
			spentImpressions float64
			spentClicks      float64
		)

		// day 1
		advanceDaySuccess(e, pointer(1))
		campaignIdStr := getAdForClientSuccess(e, clientsIds[0]).JSON().Object().Value("ad_id").String().Raw()
		campaignId := uuid.MustParse(campaignIdStr)
		campaign := campaignsMap[campaignId]
		spentImpressions += float64(campaign["cost_per_impression"].(float32))

		recordClickSuccess(e, campaignId, clientsIds[0])
		spentClicks += float64(campaign["cost_per_click"].(float32))

		campaignIdStr = getAdForClientSuccess(e, clientsIds[0]).JSON().Object().Value("ad_id").String().Raw()
		campaignId = uuid.MustParse(campaignIdStr)
		campaign = campaignsMap[campaignId]
		spentImpressions += float64(campaign["cost_per_impression"].(float32))

		campaignIdStr = getAdForClientSuccess(e, clientsIds[1]).JSON().Object().Value("ad_id").String().Raw()
		campaignId = uuid.MustParse(campaignIdStr)
		campaign = campaignsMap[campaignId]
		spentImpressions += float64(campaign["cost_per_impression"].(float32))

		expected[0]["spent_impressions"] = spentImpressions
		expected[0]["spent_clicks"] = spentClicks
		expected[0]["spent_total"] = spentImpressions + spentClicks
		spentImpressions = 0
		spentClicks = 0

		// day 2
		advanceDaySuccess(e, pointer(2))

		recordClickSuccess(e, campaignId, clientsIds[1])
		spentClicks += float64(campaign["cost_per_click"].(float32))

		campaignIdStr = getAdForClientSuccess(e, clientsIds[0]).JSON().Object().Value("ad_id").String().Raw()
		campaignId = uuid.MustParse(campaignIdStr)
		campaign = campaignsMap[campaignId]
		spentImpressions += float64(campaign["cost_per_impression"].(float32))

		campaignIdStr = getAdForClientSuccess(e, clientsIds[1]).JSON().Object().Value("ad_id").String().Raw()
		campaignId = uuid.MustParse(campaignIdStr)
		campaign = campaignsMap[campaignId]
		spentImpressions += float64(campaign["cost_per_impression"].(float32))

		recordClickSuccess(e, campaignId, clientsIds[1])
		spentClicks += float64(campaign["cost_per_click"].(float32))

		expected[1]["spent_impressions"] = spentImpressions
		expected[1]["spent_clicks"] = spentClicks
		expected[1]["spent_total"] = spentImpressions + spentClicks
		spentImpressions = 0
		spentClicks = 0

		// day 3
		advanceDaySuccess(e, pointer(3))

		campaignIdStr = getAdForClientSuccess(e, clientsIds[1]).JSON().Object().Value("ad_id").String().Raw()
		campaignId = uuid.MustParse(campaignIdStr)
		campaign = campaignsMap[campaignId]
		spentImpressions += float64(campaign["cost_per_impression"].(float32))

		recordClickSuccess(e, campaignId, clientsIds[1])
		spentClicks += float64(campaign["cost_per_click"].(float32))

		campaignIdStr = getAdForClientSuccess(e, clientsIds[2]).JSON().Object().Value("ad_id").String().Raw()
		campaignId = uuid.MustParse(campaignIdStr)
		campaign = campaignsMap[campaignId]
		spentImpressions += float64(campaign["cost_per_impression"].(float32))

		campaignIdStr = getAdForClientSuccess(e, clientsIds[2]).JSON().Object().Value("ad_id").String().Raw()
		campaignId = uuid.MustParse(campaignIdStr)
		campaign = campaignsMap[campaignId]
		spentImpressions += float64(campaign["cost_per_impression"].(float32))

		recordClickSuccess(e, campaignId, clientsIds[2])
		spentClicks += float64(campaign["cost_per_click"].(float32))

		campaignIdStr = getAdForClientSuccess(e, clientsIds[2]).JSON().Object().Value("ad_id").String().Raw()
		campaignId = uuid.MustParse(campaignIdStr)
		campaign = campaignsMap[campaignId]
		spentImpressions += float64(campaign["cost_per_impression"].(float32))

		recordClickSuccess(e, campaignId, clientsIds[2])
		spentClicks += float64(campaign["cost_per_click"].(float32))

		expected[2]["spent_impressions"] = spentImpressions
		expected[2]["spent_clicks"] = spentClicks
		expected[2]["spent_total"] = spentImpressions + spentClicks

		// check
		var actual []helpers.JSON
		getAdvertiserStatsDailySuccess(e, advertiserId).
			JSON().
			IsArray().
			Array().
			Decode(&actual)

		checkStatsDaily(t, expected, actual)
	})

	t.Run("get daily stats for advertiser with no campaigns", func(t *testing.T) {
		e := helpers.ConfigureExpect(t, ctx, advertisingServerUrl)

		// set day
		advanceDaySuccess(e, pointer(0))

		// create advertiser
		advertiser := generateAdvertiser()
		advertiserId := advertiser["advertiser_id"].(uuid.UUID)
		upsertAdvertisersSuccess(e, advertiser)

		getAdvertiserStatsDailySuccess(e, advertiserId).
			JSON().
			IsArray().
			Array().
			Length().
			IsEqual(0)
	})

	t.Run("get daily stats for non-existent advertiser", func(t *testing.T) {
		e := helpers.ConfigureExpect(t, ctx, advertisingServerUrl)

		getAdvertiserStatsDaily(e, uuid.New()).
			Expect().
			Status(http.StatusNotFound)
	})
}

func getCampaignStats(e *httpexpect.Expect, campaignId uuid.UUID) *httpexpect.Request {
	return e.GET("/stats/campaigns/{campaignId}", campaignId)
}

func getCampaignStatsSuccess(e *httpexpect.Expect, campaignId uuid.UUID) *httpexpect.Response {
	return getCampaignStats(e, campaignId).
		Expect().
		Status(http.StatusOK)
}

func getCampaignStatsDaily(e *httpexpect.Expect, campaignId uuid.UUID) *httpexpect.Request {
	return e.GET("/stats/campaigns/{campaignId}/daily", campaignId)
}

func getCampaignStatsDailySuccess(e *httpexpect.Expect, campaignId uuid.UUID) *httpexpect.Response {
	return getCampaignStatsDaily(e, campaignId).
		Expect().
		Status(http.StatusOK)
}

func getAdvertiserStats(e *httpexpect.Expect, advertiser uuid.UUID) *httpexpect.Request {
	return e.GET("/stats/advertisers/{advertiserId}/campaigns", advertiser)
}

func getAdvertiserStatsSuccess(e *httpexpect.Expect, advertiser uuid.UUID) *httpexpect.Response {
	return getAdvertiserStats(e, advertiser).
		Expect().
		Status(http.StatusOK)
}

func getAdvertiserStatsDaily(e *httpexpect.Expect, advertiser uuid.UUID) *httpexpect.Request {
	return e.GET("/stats/advertisers/{advertiserId}/campaigns/daily", advertiser)
}

func getAdvertiserStatsDailySuccess(e *httpexpect.Expect, advertiser uuid.UUID) *httpexpect.Response {
	return getAdvertiserStatsDaily(e, advertiser).
		Expect().
		Status(http.StatusOK)
}

func checkStats(t *testing.T, expected, actual helpers.JSON) {
	accurancy := math.Pow10(-6)

	assert.EqualValues(t, expected["impressions_count"], actual["impressions_count"])
	assert.EqualValues(t, expected["clicks_count"], actual["clicks_count"])
	assert.InDelta(t, expected["conversion"], actual["conversion"], accurancy)
	assert.InDelta(t, expected["spent_impressions"], actual["spent_impressions"], accurancy)
	assert.InDelta(t, expected["spent_clicks"], actual["spent_clicks"], accurancy)
	assert.InDelta(t, expected["spent_total"], actual["spent_total"], accurancy)
}

func checkStatsDaily(t *testing.T, expected, actual []helpers.JSON) {
	require.Equal(t, len(expected), len(actual))
	for i := range len(expected) {
		expectedStats := expected[i]
		actualStats := actual[i]
		require.EqualValues(t, expectedStats["date"], actualStats["date"])
		checkStats(t, expectedStats, actualStats)
	}
}
