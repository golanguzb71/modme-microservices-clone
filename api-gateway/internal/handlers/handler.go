package handlers

import (
	"api-gateway/grpc"
	client "api-gateway/internal/clients"
)

var (
	auditingClient  *client.AuditClient
	userClient      *client.UserClient
	educationClient *client.EducationClient
	leadClient      *client.LidClient
)

func InitClients(client *grpc.Clients) {
	auditingClient = client.AuditClient
	userClient = client.UserClient
	educationClient = client.EducationClient
	leadClient = client.LidClient
}
