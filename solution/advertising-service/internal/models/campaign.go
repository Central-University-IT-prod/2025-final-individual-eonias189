package models

import "github.com/google/uuid"

type Campaign struct {
	Id                uuid.UUID `db:"id"`
	AdvertiserId      uuid.UUID `db:"advertiser_id"`
	ImpressionsLimit  int       `db:"impressions_limit"`
	ClicksLimit       int       `db:"clicks_limit"`
	CostPerImpression float64   `db:"cost_per_impression"`
	CostPerClick      float64   `db:"cost_per_click"`
	AdTitle           string    `db:"ad_title"`
	AdText            string    `db:"ad_text"`
	AdImageUrl        *string   `db:"ad_image_url"`
	StartDate         int       `db:"start_date"`
	EndDate           int       `db:"end_date"`
	Gender            *Gender   `db:"gender"`
	AgeFrom           *int      `db:"age_from"`
	AgeTo             *int      `db:"age_to"`
	Location          *string   `db:"location"`
}
