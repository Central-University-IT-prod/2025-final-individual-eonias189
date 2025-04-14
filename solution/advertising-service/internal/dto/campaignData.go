package dto

import "advertising/advertising-service/internal/models"

type CampaignData struct {
	ImpressionsLimit  int            `db:"impressions_limit"`
	ClicksLimit       int            `db:"clicks_limit"`
	CostPerImpression float64        `db:"cost_per_impression"`
	CostPerClick      float64        `db:"cost_per_click"`
	AdTitle           string         `db:"ad_title"`
	AdText            string         `db:"ad_text"`
	StartDate         int            `db:"start_date"`
	EndDate           int            `db:"end_date"`
	Gender            *models.Gender `db:"gender"`
	AgeFrom           *int           `db:"age_from"`
	AgeTo             *int           `db:"age_to"`
	Location          *string        `db:"location"`
}

func CampaignDataFromCampaign(campaign models.Campaign) CampaignData {
	return CampaignData{
		ImpressionsLimit:  campaign.ImpressionsLimit,
		ClicksLimit:       campaign.ClicksLimit,
		CostPerImpression: campaign.CostPerImpression,
		CostPerClick:      campaign.CostPerClick,
		AdTitle:           campaign.AdTitle,
		AdText:            campaign.AdText,
		StartDate:         campaign.StartDate,
		EndDate:           campaign.EndDate,
		Gender:            campaign.Gender,
		AgeFrom:           campaign.AgeFrom,
		AgeTo:             campaign.AgeTo,
		Location:          campaign.Location,
	}
}

func (cd CampaignData) ToCampaign() models.Campaign {
	return models.Campaign{
		ImpressionsLimit:  cd.ImpressionsLimit,
		ClicksLimit:       cd.ClicksLimit,
		CostPerImpression: cd.CostPerImpression,
		CostPerClick:      cd.CostPerClick,
		AdTitle:           cd.AdTitle,
		AdText:            cd.AdText,
		StartDate:         cd.StartDate,
		EndDate:           cd.EndDate,
		Gender:            cd.Gender,
		AgeFrom:           cd.AgeFrom,
		AgeTo:             cd.AgeTo,
		Location:          cd.Location,
	}
}
