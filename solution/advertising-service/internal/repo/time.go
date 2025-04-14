package repo

import "context"

//go:generate go run github.com/vektra/mockery/v2@v2.52.2 --name TimeRepo
type TimeRepo interface {
	SetDay(ctx context.Context, date int) error
	GetDay(ctx context.Context) (int, error)
}
