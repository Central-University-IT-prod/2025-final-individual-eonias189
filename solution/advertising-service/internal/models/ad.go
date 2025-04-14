package models

import "github.com/google/uuid"

type Ad struct {
	CampaignId   uuid.UUID `db:"campaign_id"`
	AdvertiserId uuid.UUID `db:"advertiser_id"`
	AdTitle      string    `db:"ad_title"`
	AdText       string    `db:"ad_text"`
	AdImageUrl   *string   `db:"ad_image_url"`
}
