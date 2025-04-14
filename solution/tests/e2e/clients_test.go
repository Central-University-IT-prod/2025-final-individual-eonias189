package e2e

import (
	"advertising/tests/helpers"
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gavv/httpexpect/v2"
	"github.com/google/uuid"
)

func TestClients(t *testing.T) {
	ctx := context.Background()
	// advertisingServerUrl := helpers.SetUpInfrastructure(ctx, t, "../../advertising-service/migrations")
	advertisingServerUrl := "http://localhost:8080"

	t.Run("insert client success", func(t *testing.T) {
		e := helpers.ConfigureExpect(t, ctx, advertisingServerUrl)

		n := 20
		clients := make([]helpers.JSON, 0, n)
		for range n {
			clients = append(clients, generateClient())
		}

		upsertClientsSuccess(e, clients...).
			JSON().IsArray().Array().
			IsEqualUnordered(clients)

		for _, client := range clients {
			getClientSuccess(e, client["client_id"].(uuid.UUID)).
				JSON().
				IsEqual(client)
		}
	})

	t.Run("update clients", func(t *testing.T) {
		e := helpers.ConfigureExpect(t, ctx, advertisingServerUrl)

		n := 20

		clientsWas := make([]helpers.JSON, 0, n)
		for range n {
			clientsWas = append(clientsWas, generateClient())
		}

		upsertClientsSuccess(e, clientsWas...)

		clientsBecome := make([]helpers.JSON, 0, n)
		for _, client := range clientsWas {
			updatedClient := generateClient()
			updatedClient["client_id"] = client["client_id"]
			clientsBecome = append(clientsBecome, updatedClient)
		}

		upsertClientsSuccess(e, clientsBecome...).
			JSON().IsArray().Array().
			IsEqualUnordered(clientsBecome)

		for _, client := range clientsBecome {
			getClientSuccess(e, client["client_id"].(uuid.UUID)).
				JSON().
				IsEqual(client)
		}
	})

	t.Run("present client several times", func(t *testing.T) {
		e := helpers.ConfigureExpect(t, ctx, advertisingServerUrl)

		client := generateClient()

		// present same client
		upsertClientsSuccess(e, client, client, client).
			JSON().
			Array().
			ContainsOnly(client)

		// present same client with another fields
		client2 := generateClient()
		client2["client_id"] = client["client_id"]

		upsertClientsSuccess(e, client, client2, client, client2).
			JSON().
			Array().
			ContainsOnly(client2)

		getClientSuccess(e, client["client_id"].(uuid.UUID)).
			JSON().
			IsEqual(client2)
	})

	t.Run("insert invalid client", func(t *testing.T) {
		e := helpers.ConfigureExpect(t, ctx, advertisingServerUrl)

		// insert client with invalid gender
		client := generateClient()
		client["gender"] = "invalid value"

		upsertClients(e, client).
			Expect().
			Status(http.StatusBadRequest)

		// insert client with invalid age
		client = generateClient()
		client["age"] = -1

		upsertClients(e, client).
			Expect().
			Status(http.StatusBadRequest)
	})

	t.Run("get non-existent client", func(t *testing.T) {
		e := helpers.ConfigureExpect(t, ctx, advertisingServerUrl)

		getClient(e, uuid.New()).
			Expect().
			Status(http.StatusNotFound)
	})

}

func generateClient() helpers.JSON {
	return helpers.JSON{
		"client_id": uuid.New(),
		"login":     gofakeit.Username(),
		"age":       gofakeit.IntRange(5, 70),
		"location":  gofakeit.City(),
		"gender":    strings.ToUpper(gofakeit.Gender()),
	}
}

func upsertClients(e *httpexpect.Expect, clients ...helpers.JSON) *httpexpect.Request {
	return e.
		POST("/clients/bulk").
		WithJSON(clients)
}

func upsertClientsSuccess(e *httpexpect.Expect, clients ...helpers.JSON) *httpexpect.Response {
	return upsertClients(e, clients...).
		Expect().
		Status(http.StatusCreated)
}

func getClient(e *httpexpect.Expect, id uuid.UUID) *httpexpect.Request {
	return e.GET("/clients/{client_id}", id)
}

func getClientSuccess(e *httpexpect.Expect, id uuid.UUID) *httpexpect.Response {
	return getClient(e, id).
		Expect().
		Status(http.StatusOK)
}
