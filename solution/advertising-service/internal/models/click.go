package models

import "github.com/google/uuid"

type Click struct {
	ClientId   uuid.UUID
	CampaignId uuid.UUID
	Date       int
	Profit     float64
}
