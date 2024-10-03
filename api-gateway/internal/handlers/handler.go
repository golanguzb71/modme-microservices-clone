package handlers

import (
	"api-gateway/grpc"
	client "api-gateway/internal/clients"
)

var (
	auditingClient *client.AuditClient
	userClient     *client.UserClient
	businessClient *client.BusinessClient
	lidClient      *client.LidClient
)

func InitClients(client *grpc.Clients) {
	auditingClient = client.AuditClient
	userClient = client.UserClient
	businessClient = client.BusinessClient
	lidClient = client.LidClient
}
