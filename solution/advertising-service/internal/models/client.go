package models

import "github.com/google/uuid"

type Client struct {
	Id       uuid.UUID `db:"id"`
	Login    string    `db:"login"`
	Age      int       `db:"age"`
	Location string    `db:"location"`
	Gender   Gender    `db:"gender"`
}
