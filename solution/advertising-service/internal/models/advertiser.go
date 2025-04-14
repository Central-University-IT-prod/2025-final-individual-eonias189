package models

import "github.com/google/uuid"

type Advertiser struct {
	Id   uuid.UUID `db:"id"`
	Name string    `db:"name"`
}
