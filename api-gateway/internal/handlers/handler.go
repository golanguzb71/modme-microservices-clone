package handlers

import (
	"api-gateway/grpc"
	client "api-gateway/internal/clients"
)

var (
	userClient      *client.UserClient
	educationClient *client.EducationClient
	leadClient      *client.LidClient
	financeClient   *client.FinanceClient
)

func InitClients(client *grpc.Clients) {
	userClient = client.UserClient
	educationClient = client.EducationClient
	leadClient = client.LidClient
	financeClient = client.FinanceClient
}
