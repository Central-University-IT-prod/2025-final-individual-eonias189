package handlers

import (
	"advertising/advertising-service/internal/models"
	"advertising/pkg/logger"
	api "advertising/pkg/ogen/advertising-service"
	"context"
	"errors"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ClientsUsecase interface {
	GetClientById(ctx context.Context, id uuid.UUID) (models.Client, error)
	UpsertClients(ctx context.Context, clients []models.Client) ([]models.Client, error)
}

type ClientsHandler struct {
	cu ClientsUsecase
}

func NewClientsHandler(cu ClientsUsecase) *ClientsHandler {
	return &ClientsHandler{
		cu: cu,
	}
}

// GetClientById implements getClientById operation.
//
// Возвращает информацию о клиенте по его ID.
//
// GET /clients/{clientId}
func (ch *ClientsHandler) GetClientById(ctx context.Context, params api.GetClientByIdParams) (api.GetClientByIdRes, error) {
	client, err := ch.cu.GetClientById(ctx, params.ClientId)
	if err != nil {
		if errors.Is(err, models.ErrClientNotFound) {
			return &api.Response404{
				Resource: api.ResourceEnumClient,
			}, nil
		}

		logger.FromCtx(ctx).Error("get client by id", zap.Error(err))
		return nil, err
	}

	res := modelsClientToApiClientModel(client)
	return &res, nil
}

// UpsertClients implements upsertClients operation.
//
// Создаёт новых или обновляет существующих клиентов.
//
// POST /clients/bulk
func (ch *ClientsHandler) UpsertClients(ctx context.Context, req []api.ClientUpsert) (api.UpsertClientsRes, error) {
	clients := make([]models.Client, 0, len(req))
	for _, client := range req {
		clients = append(clients, apiClientUpsertToModelsClient(client))
	}

	clientsGot, err := ch.cu.UpsertClients(ctx, clients)
	if err != nil {
		logger.FromCtx(ctx).Error("upsert clients", zap.Error(err))
		return nil, err
	}
	res := api.UpsertClientsCreatedApplicationJSON(make([]api.ClientModel, 0, len(clientsGot)))
	for _, client := range clientsGot {
		res = append(res, modelsClientToApiClientModel(client))
	}

	return &res, nil
}

func modelsClientToApiClientModel(client models.Client) api.ClientModel {
	return api.ClientModel{
		ClientID: client.Id,
		Login:    client.Login,
		Age:      client.Age,
		Location: client.Location,
		Gender:   api.ClientModelGender(client.Gender),
	}
}

func apiClientUpsertToModelsClient(client api.ClientUpsert) models.Client {
	return models.Client{
		Id:       client.GetClientID(),
		Login:    client.GetLogin(),
		Age:      int(client.GetAge()),
		Location: client.GetLocation(),
		Gender:   models.Gender(client.GetGender()),
	}
}
