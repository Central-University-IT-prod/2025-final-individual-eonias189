package rest

import api "advertising/pkg/ogen/advertising-service"

type Handler struct {
	api.AdsHandler
	api.AdvertisersHandler
	api.CampaignsHandler
	api.ClientsHandler
	api.StatisticsHandler
	api.TimeHandler
	api.AIHandler
}

func NewHandler(
	adsHandler api.AdsHandler,
	advertisersHandler api.AdvertisersHandler,
	campaignsHandler api.CampaignsHandler,
	clientsHandler api.ClientsHandler,
	statisticsHandler api.StatisticsHandler,
	timeHandler api.TimeHandler,
	aiHandler api.AIHandler,
) *Handler {
	return &Handler{
		AdsHandler:         adsHandler,
		AdvertisersHandler: advertisersHandler,
		CampaignsHandler:   campaignsHandler,
		ClientsHandler:     clientsHandler,
		StatisticsHandler:  statisticsHandler,
		TimeHandler:        timeHandler,
		AIHandler:          aiHandler,
	}
}
