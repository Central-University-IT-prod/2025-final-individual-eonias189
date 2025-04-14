package models

import "github.com/google/uuid"

type MLScore struct {
	ClientId     uuid.UUID `db:"client_id"`
	AdvertiserId uuid.UUID `db:"advertiser_id"`
	Score        int       `db:"score"`
}
