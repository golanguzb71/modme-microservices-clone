package grpc

import (
	"api-gateway/config"
	client "api-gateway/internal/clients"
	"log"
)

type Clients struct {
	AuditClient     *client.AuditClient
	UserClient      *client.UserClient
	LidClient       *client.LidClient
	EducationClient *client.EducationClient
}

func InitializeGrpcClients(cfg *config.Config) *Clients {
	educationClient, err := client.NewEducationClient(cfg.Grpc.EducationService.Address)
	if err != nil {
		log.Fatalf("%v", err)
	}
	userClient, err := client.NewUserClient(cfg.Grpc.UserService.Address)
	if err != nil {
		return nil
	}
	lidClient, err := client.NewLidClient(cfg.Grpc.LidService.Address)
	if err != nil {
		log.Fatalf("%v", err)
	}

	return &Clients{
		AuditClient:     &client.AuditClient{},
		UserClient:      userClient,
		LidClient:       lidClient,
		EducationClient: educationClient,
	}
}
