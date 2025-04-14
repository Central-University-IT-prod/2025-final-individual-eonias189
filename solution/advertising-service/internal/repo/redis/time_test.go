package redis

import (
	"advertising/tests/helpers"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTimeRepo(t *testing.T) {
	ctx := context.Background()

	rdb := helpers.SetUpRedis(ctx, t)

	timeRepo := NewTimeRepo(rdb)

	// check get day if day not set
	curDay, err := timeRepo.GetDay(ctx)
	require.NoError(t, err, "get day")
	require.Equal(t, 0, curDay)

	// check set day
	err = timeRepo.SetDay(ctx, 5)
	require.NoError(t, err)

	// check get day
	curDay, err = timeRepo.GetDay(ctx)
	require.NoError(t, err)
	require.Equal(t, 5, curDay)
}
