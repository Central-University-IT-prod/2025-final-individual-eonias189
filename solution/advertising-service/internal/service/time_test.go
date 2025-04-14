package service

import (
	"advertising/advertising-service/internal/repo/mocks"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestTimeService(t *testing.T) {
	ctx := context.Background()

	tr := mocks.NewTimeRepo(t)
	ts := NewTimeService(tr)

	// check increment
	tr.On("GetDay", mock.Anything).Return(0, nil).Once()
	tr.On("SetDay", mock.Anything, 1).Return(nil).Once()

	curDay, err := ts.AdvanceDay(ctx, nil)
	require.NoError(t, err, "increment day")
	require.Equal(t, 1, curDay)

	// check set day
	tr.On("SetDay", mock.Anything, 42).Return(nil).Once()

	setDay := 42
	curDay, err = ts.AdvanceDay(ctx, &setDay)
	require.NoError(t, err, "set day")
	require.Equal(t, setDay, curDay)

	// check returing error
	targetError := errors.New("target error")
	tr.On("GetDay", mock.Anything).Return(setDay, nil).Once()
	tr.On("SetDay", mock.Anything, mock.AnythingOfType("int")).Return(targetError).Once()

	_, err = ts.AdvanceDay(ctx, nil)
	require.ErrorIs(t, err, targetError)

	tr.On("GetDay", mock.Anything).Return(setDay, targetError).Once()

	_, err = ts.AdvanceDay(ctx, nil)
	require.ErrorIs(t, err, targetError)

	tr.On("SetDay", mock.Anything, mock.AnythingOfType("int")).Return(targetError).Once()

	_, err = ts.AdvanceDay(ctx, &setDay)
	require.ErrorIs(t, err, targetError)

}
