package postgres

import (
	"advertising/advertising-service/internal/models"
	"advertising/tests/helpers"
	"context"
	"strings"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestClientsRepo(t *testing.T) {
	ctx := context.Background()
	db := helpers.SetUpPostgres(ctx, t, "../../../migrations")
	clientsRepo := NewClientRepo(db)

	// check UpsertClients trivial
	clients := make([]models.Client, 0, 100)
	for i := 0; i < 20; i++ {
		clients = append(clients, generateClient())
	}

	clientsGot, err := clientsRepo.UpsertClients(ctx, clients)
	require.NoError(t, err, "insert clients")
	require.ElementsMatch(t, clients, clientsGot)

	// check GetClientById trivial
	for _, client := range clients {
		clientGot, err := clientsRepo.GetClientById(ctx, client.Id)
		require.NoError(t, err, "get client by id")
		require.Equal(t, client, clientGot)
	}

	// check GetClientById with id which does not exists
	_, err = clientsRepo.GetClientById(ctx, uuid.New())
	require.Error(t, err)
	require.ErrorIs(t, err, models.ErrClientNotFound)

	// check UpsertClients with update

	// update existing clients
	newClients := clients[10:]
	for i, client := range newClients {
		newClient := generateClient()
		newClient.Id = client.Id
		newClients[i] = newClient
	}

	// add new clients
	for i := 0; i < 10; i++ {
		newClients = append(newClients, generateClient())
	}

	newClientsGot, err := clientsRepo.UpsertClients(ctx, newClients)
	require.NoError(t, err, "upsert clients")
	require.ElementsMatch(t, newClients, newClientsGot)
}

func generateClient() models.Client {
	return models.Client{
		Id:       uuid.New(),
		Login:    gofakeit.Username(),
		Age:      gofakeit.IntRange(4, 99),
		Location: gofakeit.City(),
		Gender:   models.Gender(strings.ToUpper(gofakeit.Gender())),
	}
}
