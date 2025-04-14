package service

import (
	"advertising/advertising-service/internal/repo"
	"context"
	"fmt"
)

type TimeService struct {
	tr repo.TimeRepo
}

func NewTimeService(tr repo.TimeRepo) *TimeService {
	return &TimeService{
		tr: tr,
	}
}

func (ts *TimeService) AdvanceDay(ctx context.Context, currentDay *int) (int, error) {
	op := "TimeService.AdvanceDay"

	if currentDay != nil {
		err := ts.tr.SetDay(ctx, *currentDay)
		if err != nil {
			return 0, fmt.Errorf("%s: set day: %w", op, err)
		}

		return *currentDay, nil
	}

	curDay, err := ts.tr.GetDay(ctx)
	if err != nil {
		return 0, fmt.Errorf("%s: get day: %w", op, err)
	}

	err = ts.tr.SetDay(ctx, curDay+1)
	if err != nil {
		return 0, fmt.Errorf("%s: increment day: %w", op, err)
	}

	return curDay + 1, nil
}
