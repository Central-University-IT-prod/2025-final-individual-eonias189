package helpers

import (
	"context"
	"testing"
	"time"

	pg_helper "advertising/pkg/postgres"

	"github.com/docker/go-connections/nat"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func SetUpPostgres(ctx context.Context, t *testing.T, migrationsDir string) *sqlx.DB {
	pgCont, err := startPostgresContainer(ctx, nil)
	require.NoError(t, err, "start postgres container")

	t.Cleanup(func() {
		pgCont.Terminate(ctx)
	})

	pgPort, err := pgCont.MappedPort(ctx, nat.Port("5432"))
	require.NoError(t, err, "get mapped port")

	pgCfg := pg_helper.Config{
		Host:            "localhost",
		Port:            pgPort.Int(),
		DB:              "test-db",
		User:            "test-user",
		Password:        "test-password",
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxLifetime: time.Minute * 5,
		ConnMaxIdleTime: time.Minute * 3,
	}

	err = pg_helper.Migrate(migrationsDir, pgCfg.GetConnString())
	require.NoError(t, err, "migrate postgres")

	db, err := pg_helper.Connect(ctx, pgCfg)
	require.NoError(t, err, "connect to postgres")

	return db
}

// starts postrges container with database test-db,
// user test-user and password test-password
func startPostgresContainer(ctx context.Context, network *string) (testcontainers.Container, error) {
	networks := []string{}
	if network != nil {
		networks = append(networks, *network)
	}

	return postgres.Run(ctx,
		"postgres:16.6-alpine",
		postgres.WithDatabase("test-db"),
		postgres.WithUsername("test-user"),
		postgres.WithPassword("test-password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(15*time.Second),
		),
		testcontainers.CustomizeRequest(testcontainers.GenericContainerRequest{
			ContainerRequest: testcontainers.ContainerRequest{
				ExposedPorts: []string{"5432/tcp"},
				Hostname:     "postgres",
				Networks:     networks,
			},
		}),
	)
}
