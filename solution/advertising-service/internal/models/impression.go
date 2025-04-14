package models

import "github.com/google/uuid"

type Impression struct {
	ClientId   uuid.UUID
	CampaignId uuid.UUID
	Date       int
	Profit     float64
}
