package helpers

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func SetUpRedis(ctx context.Context, t *testing.T) *redis.Client {
	redisContainer, err := startRedisContainer(ctx, nil)
	require.NoError(t, err, "start redis container")

	t.Cleanup(func() {
		redisContainer.Terminate(ctx)
	})

	port, err := redisContainer.MappedPort(ctx, nat.Port("6379"))
	require.NoError(t, err, "get redis container port")

	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("localhost:%d", port.Int()),
		DB:   0,
	})

	err = rdb.Ping(ctx).Err()
	require.NoError(t, err, "ping redis")

	return rdb
}

// starts redis container with default settings
func startRedisContainer(ctx context.Context, network *string) (testcontainers.Container, error) {
	networks := []string{}
	if network != nil {
		networks = append(networks, *network)
	}

	return testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "redis:7.4-alpine",
			Hostname:     "redis",
			ExposedPorts: []string{"6379/tcp"},
			WaitingFor:   wait.ForLog("Ready to accept connections tcp").WithStartupTimeout(10 * time.Second),
			Networks:     networks,
		},
		Started: true,
	})
}
