package e2e

import (
	"advertising/tests/helpers"
	"context"
	"maps"
	"net/http"
	"slices"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gavv/httpexpect/v2"
	"github.com/google/uuid"
)

func TestCampaignCreateAndGetById(t *testing.T) {
	ctx := context.Background()
	// advertisingServerUrl := helpers.SetUpInfrastructure(ctx, t, "../../advertising-service/migrations")
	advertisingServerUrl := "http://localhost:8080"

	t.Run("create campaign success", func(t *testing.T) {
		e := helpers.ConfigureExpect(t, ctx, advertisingServerUrl)

		// set day
		advanceDaySuccess(e, pointer(0))

		// create advertiser
		advertiser := generateAdvertiser()
		upsertAdvertisersSuccess(e, advertiser)
		advertiserId := advertiser["advertiser_id"].(uuid.UUID)

		// create campaign with full targeting
		campaign1 := generateCampaign(advertiserId, generateFullTargeting())
		campaign1IdStr := createCampaignSuccess(e, campaign1).
			JSON().
			IsObject().
			Object().
			ContainsSubset(campaign1).
			Value("campaign_id").
			String().Raw()

		campaign1Id := uuid.MustParse(campaign1IdStr)
		t.Cleanup(func() {
			deleteCapaignSuccess(e, advertiserId, campaign1Id)
		})

		campaign1["campaign_id"] = campaign1Id

		getCampaignSuccess(e, advertiserId, campaign1Id).
			JSON().
			IsObject().
			Object().
			IsEqual(campaign1)

		// create campaign with partial targeting
		campaign2 := generateCampaign(advertiserId, generatePartialTargeting())
		campaign2IdStr := createCampaignSuccess(e, campaign2).
			JSON().
			IsObject().
			Object().
			ContainsSubset(campaign2).
			Value("campaign_id").
			String().Raw()

		campaign2Id := uuid.MustParse(campaign2IdStr)
		t.Cleanup(func() {
			deleteCapaignSuccess(e, advertiserId, campaign2Id)
		})

		campaign2["campaign_id"] = campaign2Id

		getCampaignSuccess(e, advertiserId, campaign2Id).
			JSON().
			IsObject().
			Object().
			IsEqual(campaign2)
	})

	t.Run("create invalid campaign", func(t *testing.T) {
		e := helpers.ConfigureExpect(t, ctx, advertisingServerUrl)

		// set day
		advanceDaySuccess(e, pointer(0))

		// create advertiser
		advertiser := generateAdvertiser()
		upsertAdvertisersSuccess(e, advertiser)
		advertiserId := advertiser["advertiser_id"].(uuid.UUID)

		// create campaign with negative impressions_limit
		campaign := generateCampaign(advertiserId, generatePartialTargeting())
		campaign["impressions_limit"] = -1
		createCampaign(e, campaign).
			Expect().
			Status(http.StatusBadRequest)

		// create campaign with negative clicks_limit
		campaign = generateCampaign(advertiserId, generatePartialTargeting())
		campaign["clicks_limit"] = -1
		createCampaign(e, campaign).
			Expect().
			Status(http.StatusBadRequest)

		// create campaign with negative cost_per_impression
		campaign = generateCampaign(advertiserId, generatePartialTargeting())
		campaign["cost_per_impression"] = -1
		createCampaign(e, campaign).
			Expect().
			Status(http.StatusBadRequest)

		// create campaign with negative cost_per_click
		campaign = generateCampaign(advertiserId, generatePartialTargeting())
		campaign["cost_per_click"] = -1
		createCampaign(e, campaign).
			Expect().
			Status(http.StatusBadRequest)

		// create campaign with negative start_date
		campaign = generateCampaign(advertiserId, generatePartialTargeting())
		campaign["start_date"] = -1
		createCampaign(e, campaign).
			Expect().
			Status(http.StatusBadRequest)

		// create campaign with negative end_date
		campaign = generateCampaign(advertiserId, generatePartialTargeting())
		campaign["end_date"] = -1
		createCampaign(e, campaign).
			Expect().
			Status(http.StatusBadRequest)

		// create campaign with negative age_from
		targeting := generateFullTargeting()
		targeting["age_from"] = -1

		campaign = generateCampaign(advertiserId, targeting)
		createCampaign(e, campaign).
			Expect().
			Status(http.StatusBadRequest)

		// create campaign invalid gender
		targeting = generateFullTargeting()
		targeting["gender"] = "invalid gender"

		campaign = generateCampaign(advertiserId, targeting)
		createCampaign(e, campaign).
			Expect().
			Status(http.StatusBadRequest)

		// create campaign clicks_limit > impressions_limit
		campaign = generateCampaign(advertiserId, generatePartialTargeting())
		campaign["impressions_limit"] = 100
		campaign["clicks_limit"] = 200
		createCampaign(e, campaign).
			Expect().
			Status(http.StatusBadRequest)

		// create campaign end_date < start_date
		campaign = generateCampaign(advertiserId, generatePartialTargeting())
		campaign["end_date"] = 100
		campaign["start_date"] = 200
		createCampaign(e, campaign).
			Expect().
			Status(http.StatusBadRequest)

		// create campaign age_to < age_from
		targeting = generateFullTargeting()
		targeting["age_to"] = 100
		targeting["age_from"] = 200

		campaign = generateCampaign(advertiserId, targeting)
		createCampaign(e, campaign).
			Expect().
			Status(http.StatusBadRequest)

		// create campaign with start_date in past
		advanceDaySuccess(e, pointer(10))

		campaign = generateCampaign(advertiserId, generateFullTargeting())
		campaign["start_date"] = 2
		createCampaign(e, campaign).
			Expect().
			Status(http.StatusBadRequest)

	})

	t.Run("create campaign with non-existent advertiser", func(t *testing.T) {
		e := helpers.ConfigureExpect(t, ctx, advertisingServerUrl)

		// set day
		advanceDaySuccess(e, pointer(0))

		campaign := generateCampaign(uuid.New(), generatePartialTargeting())
		createCampaign(e, campaign).
			Expect().
			Status(http.StatusNotFound)
	})
}

func TestCampaignsList(t *testing.T) {

	var (
		advertiserId uuid.UUID
		campaignsN   int
		campaigns    []helpers.JSON
	)

	ctx := context.Background()
	// advertisingServerUrl := helpers.SetUpInfrastructure(ctx, t, "../../advertising-service/migrations")
	advertisingServerUrl := "http://localhost:8080"

	// seed data
	e := helpers.ConfigureExpect(t, ctx, advertisingServerUrl)

	// set day
	advanceDaySuccess(e, pointer(0))

	// create advertiser
	advertiser := generateAdvertiser()
	advertiserId = advertiser["advertiser_id"].(uuid.UUID)
	upsertAdvertisersSuccess(e, advertiser)

	// create campaigns
	campaignsN = gofakeit.IntRange(15, 40)
	campaigns = make([]helpers.JSON, 0, campaignsN)
	for range campaignsN {
		campaign := generateCampaign(advertiserId, generatePartialTargeting())
		campaignIdStr := createCampaignSuccess(e, campaign).
			JSON().
			Object().
			Value("campaign_id").
			String().
			Raw()
		campaignId := uuid.MustParse(campaignIdStr)
		t.Cleanup(func() {
			deleteCapaignSuccess(e, advertiserId, campaignId)
		})
		campaign["campaign_id"] = campaignId
		campaigns = append(campaigns, campaign)
	}
	slices.Reverse(campaigns)

	t.Run("without pagination", func(t *testing.T) {
		e := helpers.ConfigureExpect(t, ctx, advertisingServerUrl)

		listCampaignsSuccess(e, advertiserId, nil, nil).
			JSON().
			IsArray().
			Array().
			IsEqual(campaigns)
	})

	t.Run("with size", func(t *testing.T) {
		e := helpers.ConfigureExpect(t, ctx, advertisingServerUrl)

		size := gofakeit.IntRange(1, campaignsN-1)
		listCampaignsSuccess(e, advertiserId, &size, nil).
			JSON().
			IsArray().
			Array().
			IsEqual(campaigns[:size])
	})

	t.Run("with size and page", func(t *testing.T) {
		e := helpers.ConfigureExpect(t, ctx, advertisingServerUrl)

		page := gofakeit.IntRange(3, 7)
		size := gofakeit.IntRange(1, campaignsN/(page+2))
		listCampaignsSuccess(e, advertiserId, &size, &page).
			JSON().
			IsArray().
			Array().
			IsEqual(campaigns[(page-1)*size : page*size])
	})

	t.Run("with non-existent advertiser", func(t *testing.T) {
		e := helpers.ConfigureExpect(t, ctx, advertisingServerUrl)

		listCampaigns(e, uuid.New(), nil, nil).
			Expect().
			Status(http.StatusNotFound)
	})
}

func TestCampaignsUpdate(t *testing.T) {
	ctx := context.Background()
	// advertisingServerUrl := helpers.SetUpInfrastructure(ctx, t, "../../advertising-service/migrations")
	advertisingServerUrl := "http://localhost:8080"

	t.Run("update success", func(t *testing.T) {
		e := helpers.ConfigureExpect(t, ctx, advertisingServerUrl)

		// set day
		advanceDaySuccess(e, pointer(0))

		// create advertiser
		advertiser := generateAdvertiser()
		advertiserId := advertiser["advertiser_id"].(uuid.UUID)
		upsertAdvertisersSuccess(e, advertiser)

		// create not started campaign
		campaign := generateCampaign(advertiserId, generateFullTargeting())
		campaign["start_date"] = 5
		campaign["end_date"] = 10
		campaignIdStr := createCampaignSuccess(e, campaign).JSON().Object().Value("campaign_id").String().Raw()
		campaignId := uuid.MustParse(campaignIdStr)
		t.Cleanup(func() {
			deleteCapaignSuccess(e, advertiserId, campaignId)
		})

		// update all fields
		campaignNew := generateCampaign(advertiserId, generateFullTargeting())
		campaignNew["start_date"] = 6
		campaignNew["end_date"] = 11
		campaignNew["campaign_id"] = campaignId
		updateCampaignSuccess(e, advertiserId, campaignId, campaignNew).
			JSON().
			IsEqual(campaignNew)

		getCampaignSuccess(e, advertiserId, campaignId).
			JSON().
			IsEqual(campaignNew)

		// set targetting to null
		campaignNew = generateCampaign(advertiserId, helpers.JSON{
			"gender":   nil,
			"age_from": nil,
			"age_to":   nil,
			"location": nil,
		})
		campaignNew["start_date"] = 6
		campaignNew["end_date"] = 11
		campaignNew["campaign_id"] = campaignId
		updateCampaignSuccess(e, advertiserId, campaignId, campaignNew).
			JSON().
			IsEqual(campaignNew)

		getCampaignSuccess(e, advertiserId, campaignId).
			JSON().
			IsEqual(campaignNew)
	})

	t.Run("update active campaign", func(t *testing.T) {
		e := helpers.ConfigureExpect(t, ctx, advertisingServerUrl)

		// set day
		advanceDaySuccess(e, pointer(0))

		// create advertiser
		advertiser := generateAdvertiser()
		upsertAdvertisersSuccess(e, advertiser)
		advertiserId := advertiser["advertiser_id"].(uuid.UUID)

		// create campaign
		campaign := generateCampaign(advertiserId, generatePartialTargeting())
		campaign["start_date"] = 0
		campaign["end_date"] = 10
		campaignIdStr := createCampaignSuccess(e, campaign).
			JSON().
			Object().
			Value("campaign_id").
			String().Raw()
		campaignId := uuid.MustParse(campaignIdStr)
		t.Cleanup(func() {
			deleteCapaignSuccess(e, advertiserId, campaignId)
		})
		campaign["campaign_id"] = campaignId

		// update allowed fields
		campaignNew := generateCampaign(advertiserId, generateFullTargeting())
		campaignNew["impressions_limit"] = campaign["impressions_limit"]
		campaignNew["clicks_limit"] = campaign["clicks_limit"]
		campaignNew["start_date"] = campaign["start_date"]
		campaignNew["end_date"] = campaign["end_date"]
		campaignNew["campaign_id"] = campaignId
		updateCampaignSuccess(e, advertiserId, campaignId, campaignNew).
			JSON().
			IsEqual(campaignNew)

		getCampaignSuccess(e, advertiserId, campaignId).
			JSON().
			IsEqual(campaignNew)

		// update impressions_limit
		campaignNew = campaign
		newImpressionsLimit := gofakeit.IntRange(campaign["clicks_limit"].(int), 9999)
		for newImpressionsLimit == campaign["impressions_limit"].(int) {
			newImpressionsLimit = gofakeit.IntRange(campaign["clicks_limit"].(int), 9999)
		}
		campaignNew["impressions_limit"] = newImpressionsLimit
		updateCampaign(e, advertiserId, campaignId, campaignNew).
			Expect().
			Status(http.StatusForbidden)

		// update clicks_limit
		campaignNew = campaign
		newClicksLimit := gofakeit.IntRange(0, campaign["impressions_limit"].(int))
		for newClicksLimit == campaign["clicks_limit"].(int) {
			newClicksLimit = gofakeit.IntRange(0, campaign["impressions_limit"].(int))
		}
		campaignNew["clicks_limit"] = newClicksLimit
		updateCampaign(e, advertiserId, campaignId, campaignNew).
			Expect().
			Status(http.StatusForbidden)

		// update start_date
		campaignNew = campaign
		newStartDate := gofakeit.IntRange(0, campaign["end_date"].(int))
		for newStartDate == campaign["start_date"].(int) {
			newStartDate = gofakeit.IntRange(0, campaign["end_date"].(int))
		}
		campaignNew["start_date"] = newStartDate
		updateCampaign(e, advertiserId, campaignId, campaignNew).
			Expect().
			Status(http.StatusForbidden)

		// update end_date
		campaignNew = campaign
		newEndDate := gofakeit.IntRange(campaign["start_date"].(int), 999)
		for newEndDate == campaign["end_date"].(int) {
			newEndDate = gofakeit.IntRange(campaign["start_date"].(int), 999)
		}
		campaignNew["end_date"] = newEndDate
		updateCampaign(e, advertiserId, campaignId, campaignNew).
			Expect().
			Status(http.StatusForbidden)
	})

	t.Run("update campaign with invalid body", func(t *testing.T) {
		e := helpers.ConfigureExpect(t, ctx, advertisingServerUrl)

		// set dat
		advanceDaySuccess(e, pointer(0))

		// create advertiser and campaign
		advertiserId, campaignId, campaign := setupCampaignHelper(t, e)

		// update campaign with negative impressions_limit
		campaignNew := helpers.JSON{}
		maps.Copy(campaignNew, campaign)
		campaignNew["impressions_limit"] = -1
		updateCampaign(e, advertiserId, campaignId, campaignNew).
			Expect().
			Status(http.StatusBadRequest)

		// update campaign with negative clicks_limit
		campaignNew = helpers.JSON{}
		maps.Copy(campaignNew, campaign)
		campaignNew["clicks_limit"] = -1
		updateCampaign(e, advertiserId, campaignId, campaignNew).
			Expect().
			Status(http.StatusBadRequest)

		// update campaign with negative cost_per_impression
		campaignNew = helpers.JSON{}
		maps.Copy(campaignNew, campaign)
		campaignNew["cost_per_impression"] = -1
		updateCampaign(e, advertiserId, campaignId, campaignNew).
			Expect().
			Status(http.StatusBadRequest)

		// update campaign with negative cost_per_click
		campaignNew = helpers.JSON{}
		maps.Copy(campaignNew, campaign)
		campaignNew["cost_per_click"] = -1
		updateCampaign(e, advertiserId, campaignId, campaignNew).
			Expect().
			Status(http.StatusBadRequest)

		// update campaign with negative start_date
		campaignNew = helpers.JSON{}
		maps.Copy(campaignNew, campaign)
		campaignNew["start_date"] = -1
		updateCampaign(e, advertiserId, campaignId, campaignNew).
			Expect().
			Status(http.StatusBadRequest)

		// update campaign with negative end_date
		campaignNew = helpers.JSON{}
		maps.Copy(campaignNew, campaign)
		campaignNew["end_date"] = -1
		updateCampaign(e, advertiserId, campaignId, campaignNew).
			Expect().
			Status(http.StatusBadRequest)

		// update campaign with negative age_from
		targeting := generateFullTargeting()
		targeting["age_from"] = -1

		campaignNew = helpers.JSON{}
		maps.Copy(campaignNew, campaign)
		campaignNew["targeting"] = targeting
		updateCampaign(e, advertiserId, campaignId, campaignNew).
			Expect().
			Status(http.StatusBadRequest)

		// update campaign invalid gender
		targeting = generateFullTargeting()
		targeting["gender"] = "invalid gender"

		campaignNew = helpers.JSON{}
		maps.Copy(campaignNew, campaign)
		campaignNew["targeting"] = targeting
		updateCampaign(e, advertiserId, campaignId, campaignNew).
			Expect().
			Status(http.StatusBadRequest)

		// update campaign clicks_limit > impressions_limit
		campaignNew = helpers.JSON{}
		maps.Copy(campaignNew, campaign)
		campaignNew["impressions_limit"] = 100
		campaignNew["clicks_limit"] = 200
		updateCampaign(e, advertiserId, campaignId, campaignNew).
			Expect().
			Status(http.StatusBadRequest)

		// update campaign end_date < start_date
		campaignNew = helpers.JSON{}
		maps.Copy(campaignNew, campaign)
		campaignNew["end_date"] = 100
		campaignNew["start_date"] = 200
		updateCampaign(e, advertiserId, campaignId, campaignNew).
			Expect().
			Status(http.StatusBadRequest)

		// update campaign age_to < age_from
		targeting = generateFullTargeting()
		targeting["age_to"] = 100
		targeting["age_from"] = 200

		campaignNew = helpers.JSON{}
		maps.Copy(campaignNew, campaign)
		campaignNew["targeting"] = targeting
		updateCampaign(e, advertiserId, campaignId, campaignNew).
			Expect().
			Status(http.StatusBadRequest)
	})

	t.Run("update not active campaign, set start_date < curent day", func(t *testing.T) {
		e := helpers.ConfigureExpect(t, ctx, advertisingServerUrl)

		// set day
		advanceDaySuccess(e, pointer(5))

		// create advertiser
		advertiser := generateAdvertiser()
		advertiserId := advertiser["advertiser_id"].(uuid.UUID)
		upsertAdvertisersSuccess(e, advertiser)

		// create campaign
		campaign := generateCampaign(advertiserId, generateFullTargeting())
		campaign["start_date"] = 10
		campaign["end_date"] = 20
		campaignIdStr := createCampaignSuccess(e, campaign).JSON().Object().Value("campaign_id").String().Raw()
		campaignId := uuid.MustParse(campaignIdStr)
		t.Cleanup(func() {
			deleteCapaignSuccess(e, advertiserId, campaignId)
		})

		// update campaign
		campaignNew := generateCampaign(advertiserId, generateFullTargeting())
		campaignNew["start_date"] = 3
		campaignNew["end_date"] = 7

		updateCampaign(e, advertiserId, campaignId, campaignNew).
			Expect().
			Status(http.StatusBadRequest)
	})

	t.Run("update non-existent campaign", func(t *testing.T) {
		e := helpers.ConfigureExpect(t, ctx, advertisingServerUrl)

		// set day
		advanceDaySuccess(e, pointer(0))

		advertiserId, campaignId, campaign := setupCampaignHelper(t, e)

		// campaign not found
		updateCampaign(e, advertiserId, uuid.New(), campaign).
			Expect().
			Status(http.StatusNotFound)

		advertiserAnother := generateAdvertiser()
		upsertAdvertisersSuccess(e, advertiserAnother)
		advertiserAnotherId := advertiserAnother["advertiser_id"].(uuid.UUID)

		updateCampaign(e, advertiserAnotherId, campaignId, campaign).
			Expect().
			Status(http.StatusNotFound)

		// advertiser not found
		updateCampaign(e, uuid.New(), campaignId, campaign).
			Expect().
			Status(http.StatusNotFound)
	})
}

func TestCampaignsDelete(t *testing.T) {
	ctx := context.Background()
	// advertisingServerUrl := helpers.SetUpInfrastructure(ctx, t, "../../advertising-service/migrations")
	advertisingServerUrl := "http://localhost:8080"

	// set day
	e := helpers.ConfigureExpect(t, ctx, advertisingServerUrl)

	advanceDaySuccess(e, pointer(0))

	// create advertiser
	advertiser := generateAdvertiser()
	advertiserId := advertiser["advertiser_id"].(uuid.UUID)
	upsertAdvertisersSuccess(e, advertiser)

	t.Run("delete success", func(t *testing.T) {
		e := helpers.ConfigureExpect(t, ctx, advertisingServerUrl)

		// create campaign
		campaign := generateCampaign(advertiserId, generateFullTargeting())
		campaignIdStr := createCampaignSuccess(e, campaign).JSON().Object().Value("campaign_id").String().Raw()
		campaignId := uuid.MustParse(campaignIdStr)

		deleteCapaignSuccess(e, advertiserId, campaignId)

		getCampaign(e, advertiserId, campaignId).
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("delete with non-existent campaign", func(t *testing.T) {
		e := helpers.ConfigureExpect(t, ctx, advertisingServerUrl)

		advertiserId, campaignId, _ := setupCampaignHelper(t, e)

		// campaign not found
		deleteCapaign(e, advertiserId, uuid.New()).
			Expect().
			Status(http.StatusNotFound)

		advertiserAnother := generateAdvertiser()
		upsertAdvertisersSuccess(e, advertiserAnother)
		advertiserAnotherId := advertiserAnother["advertiser_id"].(uuid.UUID)

		deleteCapaign(e, advertiserAnotherId, campaignId).
			Expect().
			Status(http.StatusNotFound)

		// advertiser not found
		deleteCapaign(e, uuid.New(), campaignId).
			Expect().
			Status(http.StatusNotFound)
	})
}

func createCampaign(e *httpexpect.Expect, campaign helpers.JSON) *httpexpect.Request {
	return e.POST("/advertisers/{advertiser_id}/campaigns", campaign["advertiser_id"]).
		WithJSON(campaign)
}

func createCampaignSuccess(e *httpexpect.Expect, campaign helpers.JSON) *httpexpect.Response {
	return createCampaign(e, campaign).
		Expect().
		Status(http.StatusCreated)
}

func getCampaign(e *httpexpect.Expect, advertiserId uuid.UUID, campaignId uuid.UUID) *httpexpect.Request {
	return e.GET("/advertisers/{advertiser_id}/campaigns/{campaign_id}", advertiserId, campaignId)
}

func getCampaignSuccess(e *httpexpect.Expect, advertiserId uuid.UUID, campaignId uuid.UUID) *httpexpect.Response {
	return getCampaign(e, advertiserId, campaignId).
		Expect().
		Status(http.StatusOK)
}

func deleteCapaign(e *httpexpect.Expect, advertiserId uuid.UUID, campaignId uuid.UUID) *httpexpect.Request {
	return e.DELETE("/advertisers/{advertiser_id}/campaigns/{campaign_id}", advertiserId, campaignId)
}

func deleteCapaignSuccess(e *httpexpect.Expect, advertiserId uuid.UUID, campaignId uuid.UUID) *httpexpect.Response {
	return deleteCapaign(e, advertiserId, campaignId).
		Expect().
		Status(http.StatusNoContent)
}

func listCampaigns(e *httpexpect.Expect, advertiserId uuid.UUID, size *int, page *int) *httpexpect.Request {
	req := e.GET("/advertisers/{advertiser_id}/campaigns", advertiserId)

	if size != nil {
		req = req.WithQuery("size", *size)
	}
	if page != nil {
		req = req.WithQuery("page", *page)
	}

	return req
}

func listCampaignsSuccess(e *httpexpect.Expect, advertiserId uuid.UUID, size *int, page *int) *httpexpect.Response {
	return listCampaigns(e, advertiserId, size, page).
		Expect().
		Status(http.StatusOK)
}

func updateCampaign(e *httpexpect.Expect, advertiserId uuid.UUID, campaignId uuid.UUID, updated helpers.JSON) *httpexpect.Request {
	return e.PUT("/advertisers/{advertiser_id}/campaigns/{campaign_id}", advertiserId, campaignId).
		WithJSON(updated)
}

func updateCampaignSuccess(e *httpexpect.Expect, advertiserId uuid.UUID, campaignId uuid.UUID, updated helpers.JSON) *httpexpect.Response {
	return updateCampaign(e, advertiserId, campaignId, updated).
		Expect().
		Status(http.StatusOK)
}

func generateCampaign(advertiserId uuid.UUID, targeting helpers.JSON) helpers.JSON {
	impressionsLimit := gofakeit.IntRange(100, 2000)
	clicksLimit := gofakeit.IntRange(20, impressionsLimit)
	startDate := gofakeit.IntRange(0, 300)
	endDate := gofakeit.IntRange(startDate, 999)
	return helpers.JSON{
		"advertiser_id":       advertiserId,
		"impressions_limit":   impressionsLimit,
		"clicks_limit":        clicksLimit,
		"cost_per_impression": gofakeit.Float32Range(0, 10000),
		"cost_per_click":      gofakeit.Float32Range(0, 100000),
		"ad_title":            gofakeit.Phrase(),
		"ad_text":             gofakeit.Sentence(gofakeit.IntRange(8, 20)),
		"start_date":          startDate,
		"end_date":            endDate,
		"targeting":           targeting,
	}
}

func generateGender() string {
	return gofakeit.RandomString([]string{"MALE", "FEMALE", "ALL"})
}

func generateFullTargeting() helpers.JSON {
	ageFrom := gofakeit.IntRange(0, 30)
	ageTo := gofakeit.IntRange(ageFrom, 70)
	return helpers.JSON{
		"gender":   generateGender(),
		"age_from": ageFrom,
		"age_to":   ageTo,
		"location": gofakeit.City(),
	}
}

func generatePartialTargeting() helpers.JSON {
	fullTargeting := generateFullTargeting()
	fields := []string{"gender", "age_from", "age_to", "location"}
	for _, field := range fields {
		if gofakeit.Float32Range(0, 1) >= 0.5 {
			fullTargeting[field] = nil
		}
	}
	return fullTargeting
}

func setupCampaignHelper(t *testing.T, e *httpexpect.Expect) (uuid.UUID, uuid.UUID, helpers.JSON) {
	advertiser := generateAdvertiser()
	upsertAdvertisersSuccess(e, advertiser)
	advertiserId := advertiser["advertiser_id"].(uuid.UUID)

	campaign := generateCampaign(advertiserId, generatePartialTargeting())
	campaign["start_date"] = 5
	campaign["end_date"] = 10
	campaignIdStr := createCampaignSuccess(e, campaign).
		JSON().
		Object().
		Value("campaign_id").
		String().Raw()
	campaignId := uuid.MustParse(campaignIdStr)
	t.Cleanup(func() {
		deleteCapaignSuccess(e, advertiserId, campaignId)
	})
	campaign["campaign_id"] = campaignId
	return advertiserId, campaignId, campaign
}
