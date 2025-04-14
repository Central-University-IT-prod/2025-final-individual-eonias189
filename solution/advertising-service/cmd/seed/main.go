package main

import (
	"advertising/advertising-service/internal/dto"
	"advertising/advertising-service/internal/models"
	"advertising/advertising-service/internal/repo/postgres"
	"advertising/pkg/logger"
	pg_helper "advertising/pkg/postgres"
	"context"
	"log"
	"strings"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	l, err := logger.Get("info")
	if err != nil {
		log.Fatal("get logger", err)
	}

	var cfg pg_helper.Config
	err = cleanenv.ReadEnv(&cfg)
	if err != nil {
		l.Fatal("get config", zap.Error(err))
	}

	db, err := pg_helper.Connect(ctx, cfg)
	if err != nil {
		l.Fatal("connect to postrges", zap.Error(err))
	}

	l.Info("seeding clients")
	clients, err := seedClients(ctx, db, 100)
	if err != nil {
		l.Error("seed clients", zap.Error(err))
	}

	l.Info("seeding advertisers")
	advertisers, err := seedAdvertisers(ctx, db, 100)
	if err != nil {
		l.Error("seed advertisers", zap.Error(err))
	}

	l.Info("seeding ml scores")
	if err := seedMLScores(ctx, db, 800, clients, advertisers); err != nil {
		l.Error("seed ml scores", zap.Error(err))
	}

	l.Info("seeding campaigns")
	campaigns, err := seedCampaigns(ctx, db, 400, advertisers)
	if err != nil {
		l.Error("seed campaigns", zap.Error(err))
	}

	l.Info("seeding impressions")
	impressions, err := seedImpressions(ctx, db, clients, campaigns)
	if err != nil {
		l.Error("seed impressions", zap.Error(err))
	}

	l.Info("seeding clicks")
	_, err = seedClicks(ctx, db, impressions, campaigns)
	if err != nil {
		l.Error("seed campaigns", zap.Error(err))
	}

	l.Info("seeded successfully")

}

func seedClients(ctx context.Context, db *sqlx.DB, n int) ([]models.Client, error) {
	cr := postgres.NewClientRepo(db)

	clients := make([]models.Client, 0, n)
	for range n {
		clients = append(clients, generateClient())
	}

	_, err := cr.UpsertClients(ctx, clients)
	if err != nil {
		return nil, err
	}

	return clients, nil
}

func seedAdvertisers(ctx context.Context, db *sqlx.DB, n int) ([]models.Advertiser, error) {
	ar := postgres.NewAdvertiserRepo(db)

	advertisers := make([]models.Advertiser, 0, n)
	for range n {
		advertisers = append(advertisers, generateAdvertiser())
	}

	_, err := ar.UpsertAdvertisers(ctx, advertisers)
	if err != nil {
		return nil, err
	}

	return advertisers, nil
}

func seedMLScores(ctx context.Context, db *sqlx.DB, n int, clients []models.Client, advertisers []models.Advertiser) error {
	msr := postgres.NewMlScoresRepo(db)

	inserted := map[struct {
		ClientId     uuid.UUID
		AdvertiserId uuid.UUID
	}]bool{}

	for range n {
		var (
			clientId     uuid.UUID
			advertiserId uuid.UUID
		)
		for {
			clientId = clients[gofakeit.IntRange(0, len(clients)-1)].Id
			advertiserId = advertisers[gofakeit.IntRange(0, len(advertisers)-1)].Id
			if _, ok := inserted[struct {
				ClientId     uuid.UUID
				AdvertiserId uuid.UUID
			}{ClientId: clientId, AdvertiserId: advertiserId}]; !ok {
				break
			}
		}
		if err := msr.UpsertMLScore(ctx, models.MLScore{
			ClientId:     clientId,
			AdvertiserId: advertiserId,
			Score:        int(gofakeit.Int32()),
		}); err != nil {
			return err
		}
		inserted[struct {
			ClientId     uuid.UUID
			AdvertiserId uuid.UUID
		}{
			ClientId:     clientId,
			AdvertiserId: advertiserId,
		}] = true
	}

	return nil
}

func seedCampaigns(ctx context.Context, db *sqlx.DB, n int, advertisers []models.Advertiser) ([]models.Campaign, error) {
	cr := postgres.NewCampaignsRepo(db)
	campaigns := make([]models.Campaign, 0, n)

	for range n {
		advertiserId := advertisers[gofakeit.IntRange(0, len(advertisers)-1)].Id
		campaign := generateCampaign(advertiserId)
		campaignId, err := cr.CreateCampaign(ctx, advertiserId, dto.CampaignDataFromCampaign(campaign))
		if err != nil {
			return nil, err
		}
		campaign.Id = campaignId
		campaigns = append(campaigns, campaign)
	}

	return campaigns, nil
}

func seedImpressions(ctx context.Context, db *sqlx.DB, clients []models.Client, campaigns []models.Campaign) ([]models.Impression, error) {
	clientActionsRepo := postgres.NewClientActionsRepo(db)
	res := []models.Impression{}

	impressions := map[uuid.UUID]int{}

	for _, client := range clients {
		for _, campaign := range campaigns {
			if campaign.Gender != nil && *campaign.Gender != models.GenderAll && *campaign.Gender != client.Gender {
				continue
			}
			if campaign.Location != nil && *campaign.Location != client.Location {
				continue
			}
			if campaign.AgeFrom != nil && *campaign.AgeFrom > client.Age {
				continue
			}
			if campaign.AgeTo != nil && *campaign.AgeTo < client.Age {
				continue
			}

			impressionsWas := impressions[campaign.Id]
			if impressionsWas == campaign.ImpressionsLimit {
				continue
			}

			if gofakeit.IntRange(0, 100) > 80 {
				continue
			}

			date := gofakeit.IntRange(campaign.StartDate, campaign.EndDate)

			impression := models.Impression{
				ClientId:   client.Id,
				CampaignId: campaign.Id,
				Date:       date,
				Profit:     campaign.CostPerImpression,
			}

			err := clientActionsRepo.RecordImpression(ctx, impression)
			if err != nil {
				return nil, err
			}
			impressions[campaign.Id] = impressionsWas + 1

			res = append(res, impression)
		}
	}

	return res, nil

}

func seedClicks(ctx context.Context, db *sqlx.DB, impressions []models.Impression, campaigns []models.Campaign) ([]models.Click, error) {
	res := []models.Click{}
	clientActionsRepo := postgres.NewClientActionsRepo(db)

	clicks := map[uuid.UUID]int{}

	campaignsMap := make(map[uuid.UUID]models.Campaign, len(campaigns))
	for _, campaign := range campaigns {
		campaignsMap[campaign.Id] = campaign
	}

	for _, impression := range impressions {
		campaign := campaignsMap[impression.CampaignId]
		clicksWas := clicks[impression.CampaignId]
		if clicksWas == campaign.ClicksLimit {
			continue
		}

		if gofakeit.IntRange(0, 100) > 75 {
			continue
		}

		date := gofakeit.IntRange(impression.Date, campaign.EndDate)

		click := models.Click{
			ClientId:   impression.ClientId,
			CampaignId: impression.CampaignId,
			Date:       date,
			Profit:     campaign.CostPerClick,
		}

		err := clientActionsRepo.RecordClick(ctx, click)
		if err != nil {
			return nil, err
		}

		clicks[impression.CampaignId] = clicksWas + 1

		res = append(res, click)

	}

	return res, nil

}

func generateClient() models.Client {
	return models.Client{
		Id:       uuid.New(),
		Login:    gofakeit.Username(),
		Age:      gofakeit.IntRange(5, 90),
		Location: gofakeit.City(),
		Gender:   models.Gender(strings.ToUpper(gofakeit.Gender())),
	}
}

func generateAdvertiser() models.Advertiser {
	return models.Advertiser{
		Id:   uuid.New(),
		Name: gofakeit.Company(),
	}
}

func generateCampaign(advertiserId uuid.UUID) models.Campaign {
	impressionsLimit := gofakeit.IntRange(1000, 99999)
	startDate := gofakeit.IntRange(0, 50)
	res := models.Campaign{
		AdvertiserId:      advertiserId,
		ImpressionsLimit:  impressionsLimit,
		ClicksLimit:       gofakeit.IntRange(0, impressionsLimit),
		CostPerImpression: gofakeit.Float64Range(0, 99999),
		CostPerClick:      gofakeit.Float64Range(0, 99999),
		AdTitle:           gofakeit.Sentence(gofakeit.IntRange(3, 10)),
		AdText:            gofakeit.Sentence(gofakeit.IntRange(10, 50)),
		StartDate:         startDate,
		EndDate:           gofakeit.IntRange(startDate+1, 100),
	}

	if gofakeit.IntRange(0, 100) >= 50 {
		res.Gender = pointer(generateGender())
	}
	var ageFrom *int
	if gofakeit.IntRange(0, 100) >= 50 {
		ageFrom = pointer(gofakeit.IntRange(0, 50))
		res.AgeFrom = ageFrom
	}
	if gofakeit.IntRange(0, 100) >= 50 {
		var ageTo *int
		if ageFrom != nil {
			ageTo = pointer(gofakeit.IntRange(*ageFrom+1, 100))
		} else {
			ageTo = pointer(gofakeit.IntRange(5, 100))
		}
		res.AgeTo = ageTo
	}
	if gofakeit.IntRange(0, 100) >= 50 {
		res.Location = pointer(gofakeit.City())
	}

	return res
}

func generateGender() models.Gender {
	return models.Gender(gofakeit.RandomString([]string{"MALE", "FEMALE", "ALL"}))
}

func pointer[T any](v T) *T {
	return &v
}
