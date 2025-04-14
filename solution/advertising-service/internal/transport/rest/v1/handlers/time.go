package handlers

import (
	"advertising/pkg/logger"
	api "advertising/pkg/ogen/advertising-service"
	"context"

	"go.uber.org/zap"
)

type TimeUsecase interface {
	AdvanceDay(ctx context.Context, currentdaty *int) (int, error)
}

type TimeHandler struct {
	tu TimeUsecase
}

func NewTimeHandler(tu TimeUsecase) *TimeHandler {
	return &TimeHandler{
		tu: tu,
	}
}

// AdvanceDay implements advanceDay operation.
//
// Устанавливает текущий день.
//
// POST /time/advance
func (th *TimeHandler) AdvanceDay(ctx context.Context, req api.OptAdvanceDayReq) (api.AdvanceDayRes, error) {
	var currentDay *int

	if req.IsSet() && req.Value.CurrentDate.IsSet() {
		curDayInt := int(req.Value.CurrentDate.Value)
		currentDay = &curDayInt
	}

	res, err := th.tu.AdvanceDay(ctx, currentDay)
	if err != nil {
		logger.FromCtx(ctx).Error("advance day", zap.Error(err))
		return nil, err
	}

	return &api.AdvanceDayOK{
		CurrentDate: api.NewOptDate(api.Date(res)),
	}, nil
}
