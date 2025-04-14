package postgres

import (
	"advertising/advertising-service/internal/dto"
	"advertising/advertising-service/internal/models"
	"advertising/tests/helpers"
	"context"
	"math"
	"slices"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

var (
	advertiserId uuid.UUID

	campaign1Id uuid.UUID
	campaign2Id uuid.UUID
	campaign3Id uuid.UUID

	campaign1Stats      models.Stats
	campaign1DailyStats []models.StatsDaily

	campaign2Stats      models.Stats
	campaign2DailyStats []models.StatsDaily

	campaign3Stats      models.Stats
	campaign3DailyStats []models.StatsDaily

	advertiserStats      models.Stats
	advertiserDailyStats []models.StatsDaily
)

func TestStatsRepo(t *testing.T) {
	ctx := context.Background()
	db := helpers.SetUpPostgres(ctx, t, "../../../migrations")
	seed(t, ctx, db)
	statsRepo := NewStatsRepo(db)

	// check stats
	campaign1StatsGot, err := statsRepo.GetStatsForCampaign(ctx, campaign1Id)
	require.NoError(t, err)
	checkStats(t, campaign1Stats, campaign1StatsGot)

	campaign2StatsGot, err := statsRepo.GetStatsForCampaign(ctx, campaign2Id)
	require.NoError(t, err)
	checkStats(t, campaign2Stats, campaign2StatsGot)

	campaign3StatsGot, err := statsRepo.GetStatsForCampaign(ctx, campaign3Id)
	require.NoError(t, err)
	checkStats(t, campaign3Stats, campaign3StatsGot)

	// check daily stats
	campaign1DailyStatsGot, err := statsRepo.GetStatsForCampaignDaily(ctx, campaign1Id)
	require.NoError(t, err)
	checkStatsDaily(t, campaign1DailyStats, campaign1DailyStatsGot)

	campaign2DailyStatsGot, err := statsRepo.GetStatsForCampaignDaily(ctx, campaign2Id)
	require.NoError(t, err)
	checkStatsDaily(t, campaign2DailyStats, campaign2DailyStatsGot)

	campaign3DailyStatsGot, err := statsRepo.GetStatsForCampaignDaily(ctx, campaign3Id)
	require.NoError(t, err)
	checkStatsDaily(t, campaign3DailyStats, campaign3DailyStatsGot)

	// check advertiser stats
	advertiserStatsGot, err := statsRepo.GetStatsForAdvertiser(ctx, advertiserId)
	require.NoError(t, err)
	checkStats(t, advertiserStats, advertiserStatsGot)

	// check advertiser stats daily
	advertiserStatsDailyGot, err := statsRepo.GetStatsForAdvertiserDaily(ctx, advertiserId)
	require.NoError(t, err)
	checkStatsDaily(t, advertiserDailyStats, advertiserStatsDailyGot)
}

func checkStats(t *testing.T, expected, actual models.Stats) {
	require.Equal(t, expected.ImpressionsCount, actual.ImpressionsCount)
	require.Equal(t, expected.ClicksCount, actual.ClicksCount)
	checkFloat64(t, expected.Conversion, actual.Conversion)
	checkFloat64(t, expected.SpentImpressions, actual.SpentImpressions)
	checkFloat64(t, expected.SpentClicks, actual.SpentClicks)
	checkFloat64(t, expected.SpentTotal, actual.SpentTotal)
}

func checkStatsDaily(t *testing.T, expected, actual []models.StatsDaily) {
	require.Equal(t, len(expected), len(actual))
	for i := 0; i < len(expected); i++ {
		expectedStats := expected[i]
		actualStats := actual[i]
		checkStats(t, expectedStats.Stats, actualStats.Stats)
	}
}

func checkFloat64(t *testing.T, expected, actual float64) {
	accurancy := math.Pow10(-6)

	diff := math.Abs(expected - actual)

	require.LessOrEqual(t, diff, accurancy)
}

func seed(t *testing.T, ctx context.Context, db *sqlx.DB) {
	clientsRepo := NewClientRepo(db)
	advertiserRepo := NewAdvertiserRepo(db)
	campaignsRepo := NewCampaignsRepo(db)
	clientActionsRepo := NewClientActionsRepo(db)

	advertiser := generateAdvertiser()
	_, err := advertiserRepo.UpsertAdvertisers(ctx, []models.Advertiser{advertiser})
	require.NoError(t, err)
	advertiserId = advertiser.Id

	campaign1 := generateCampaign()
	campaign1.AdvertiserId = advertiserId
	campaign1.Id, err = campaignsRepo.CreateCampaign(ctx, advertiserId, dto.CampaignDataFromCampaign(campaign1))
	campaign1Id = campaign1.Id
	require.NoError(t, err)

	campaign2 := generateCampaign()
	campaign2.AdvertiserId = advertiserId
	campaign2.Id, err = campaignsRepo.CreateCampaign(ctx, advertiserId, dto.CampaignDataFromCampaign(campaign2))
	campaign2Id = campaign2.Id
	require.NoError(t, err)

	campaign3 := generateCampaign()
	campaign3.AdvertiserId = advertiserId
	campaign3.Id, err = campaignsRepo.CreateCampaign(ctx, advertiserId, dto.CampaignDataFromCampaign(campaign3))
	campaign3Id = campaign3.Id
	require.NoError(t, err)

	campaign1Days := []int{1, 2, 3}
	campaign2Days := []int{2, 3, 4}
	campaign3Days := []int{3, 4, 5}

	campaign1Stats, campaign1DailyStats = createCampaignStats(t, ctx, clientsRepo, clientActionsRepo, campaign1Days, campaign1)
	campaign2Stats, campaign2DailyStats = createCampaignStats(t, ctx, clientsRepo, clientActionsRepo, campaign2Days, campaign2)
	campaign3Stats, campaign3DailyStats = createCampaignStats(t, ctx, clientsRepo, clientActionsRepo, campaign3Days, campaign3)

	advertiserStats, advertiserDailyStats = aggrDailyStats(campaign1DailyStats, campaign2DailyStats, campaign3DailyStats)
}

func createCampaignStats(
	t *testing.T,
	ctx context.Context,
	clientsRepo *ClientsRepo,
	clientActionsRepo *ClientActionsRepo,
	days []int,
	campaign models.Campaign,
) (models.Stats, []models.StatsDaily) {
	res := models.Stats{}
	resDaily := []models.StatsDaily{}
	for _, day := range days {
		impressionsCount := gofakeit.IntRange(0, 50)
		clicksCount := gofakeit.IntRange(0, 50)

		clientsToCreate := max(impressionsCount, clicksCount)
		var impressed int
		var clicked int

		for range clientsToCreate {
			client := generateClient()
			_, err := clientsRepo.UpsertClients(ctx, []models.Client{client})
			require.NoError(t, err)

			if impressed < impressionsCount {
				err = clientActionsRepo.RecordImpression(ctx, models.Impression{
					ClientId:   client.Id,
					CampaignId: campaign.Id,
					Date:       day,
					Profit:     campaign.CostPerImpression,
				})
				require.NoError(t, err)
				impressed++
			}

			if clicked < clicksCount {
				err = clientActionsRepo.RecordClick(ctx, models.Click{
					ClientId:   client.Id,
					CampaignId: campaign.Id,
					Date:       day,
					Profit:     campaign.CostPerClick,
				})
				require.NoError(t, err)
				clicked++
			}
		}

		spentImpressions := float64(impressionsCount) * campaign.CostPerImpression
		spentClicks := float64(clicksCount) * campaign.CostPerClick
		var conversion float64
		if impressionsCount != 0 {
			conversion = float64(clicksCount) / float64(impressionsCount) * 100
		} else {
			conversion = 0
		}
		spentTotal := spentImpressions + spentClicks

		resDaily = append(resDaily, models.StatsDaily{
			Stats: models.Stats{
				ImpressionsCount: impressionsCount,
				ClicksCount:      clicksCount,
				Conversion:       conversion,
				SpentImpressions: spentImpressions,
				SpentClicks:      spentClicks,
				SpentTotal:       spentTotal,
			},
			Date: day,
		})
		res.ImpressionsCount += impressionsCount
		res.ClicksCount += clicksCount
		res.SpentImpressions += spentImpressions
		res.SpentClicks += spentClicks
		res.SpentTotal += spentTotal
	}

	if res.ImpressionsCount != 0 {
		res.Conversion = float64(res.ClicksCount) / float64(res.ImpressionsCount) * 100
	} else {
		res.Conversion = 0
	}

	return res, resDaily
}

func aggrDailyStats(statsRows ...[]models.StatsDaily) (models.Stats, []models.StatsDaily) {
	res := models.Stats{}
	resDailyMap := map[int]models.Stats{}

	for _, stats := range statsRows {
		for _, statsDaily := range stats {
			date := statsDaily.Date
			if _, ok := resDailyMap[date]; !ok {
				resDailyMap[date] = models.Stats{}
			}
			statsWas := resDailyMap[date]
			statsBecome := models.Stats{
				ImpressionsCount: statsWas.ImpressionsCount + statsDaily.ImpressionsCount,
				ClicksCount:      statsWas.ClicksCount + statsDaily.ClicksCount,
				SpentImpressions: statsWas.SpentImpressions + statsDaily.SpentImpressions,
				SpentClicks:      statsWas.SpentClicks + statsDaily.SpentClicks,
				SpentTotal:       statsWas.SpentTotal + statsDaily.SpentTotal,
			}
			if (statsWas.ImpressionsCount + statsDaily.ImpressionsCount) != 0 {
				statsBecome.Conversion = float64(statsWas.ClicksCount+statsDaily.ClicksCount) / float64(statsWas.ImpressionsCount+statsDaily.ImpressionsCount) * 100
			} else {
				statsBecome.Conversion = 0
			}
			resDailyMap[date] = statsBecome

			res.ImpressionsCount += statsDaily.ImpressionsCount
			res.ClicksCount += statsDaily.ClicksCount
			res.SpentImpressions += statsDaily.SpentImpressions
			res.SpentClicks += statsDaily.SpentClicks
			res.SpentTotal += statsDaily.SpentTotal
		}
	}

	if res.ImpressionsCount != 0 {
		res.Conversion = float64(res.ClicksCount) / float64(res.ImpressionsCount) * 100
	} else {
		res.Conversion = 0
	}

	resDaily := make([]models.StatsDaily, 0, len(resDailyMap))
	for date, stats := range resDailyMap {
		resDaily = append(resDaily, models.StatsDaily{
			Stats: stats,
			Date:  date,
		})
	}

	slices.SortFunc(resDaily, func(s1, s2 models.StatsDaily) int {
		return s1.Date - s2.Date
	})

	return res, resDaily
}
