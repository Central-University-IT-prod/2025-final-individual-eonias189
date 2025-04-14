package helpers

import (
	"context"
	"fmt"
	"os"
	"testing"

	pg_helper "advertising/pkg/postgres"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/network"
	"github.com/testcontainers/testcontainers-go/wait"
)

// returns advertising service base url
func SetUpInfrastructure(ctx context.Context, t *testing.T, advertisingServiceMidrationsDir string) string {
	ntwrk, err := network.New(ctx)
	require.NoError(t, err, "create network")

	t.Cleanup(func() {
		ntwrk.Remove(ctx)
	})

	pgCont, err := startPostgresContainer(ctx, &ntwrk.Name)
	require.NoError(t, err, "start postgres container")

	t.Cleanup(func() {
		pgCont.Terminate(ctx)
	})

	pgPort, err := pgCont.MappedPort(ctx, "5432")
	require.NoError(t, err, "get postgres port")

	err = pg_helper.Migrate(advertisingServiceMidrationsDir, pg_helper.Config{
		Host:     "localhost",
		Port:     pgPort.Int(),
		DB:       "test-db",
		User:     "test-user",
		Password: "test-password",
	}.GetConnString())
	require.NoError(t, err, "run postgres migrations")

	redisCont, err := startRedisContainer(ctx, &ntwrk.Name)
	require.NoError(t, err, "start redis container")

	t.Cleanup(func() {
		redisCont.Terminate(ctx)
	})

	advertisingServiceImage := os.Getenv("ADVERTISING_SERVICE_IMAGE")
	if advertisingServiceImage == "" {
		advertisingServiceImage = "advertising-service:e2e"
	}

	req := testcontainers.ContainerRequest{
		Image:        advertisingServiceImage,
		ExposedPorts: []string{"8080/tcp"},
		WaitingFor:   wait.ForLog("starting server"),
		Env: map[string]string{
			"SERVER_PORT":       "8080",
			"POSTGRES_HOST":     "postgres",
			"POSTGRES_DB":       "test-db",
			"POSTGRES_USER":     "test-user",
			"POSTGRES_PASSWORD": "test-password",
			"REDIS_HOST":        "redis",
			"REDIS_PORT":        "6379",
			"LOG_LEVEL":         "debug",
		},
		Networks: []string{ntwrk.Name},
	}

	advertisingServiceCont, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err, "start advertising service container")

	t.Cleanup(func() {
		advertisingServiceCont.Terminate(ctx)
	})

	port, err := advertisingServiceCont.MappedPort(ctx, "8080")
	require.NoError(t, err, "get advertising service port")

	return fmt.Sprintf("http://localhost:%d", port.Int())
}
