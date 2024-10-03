package grpc

import (
	"api-gateway/config"
	client "api-gateway/internal/clients"
	"log"
)

type Clients struct {
	AuditClient    *client.AuditClient
	UserClient     *client.UserClient
	LidClient      *client.LidClient
	BusinessClient *client.BusinessClient
}

func InitializeGrpcClients(cfg *config.Config) *Clients {
	businessClient, err := client.NewBusinessClient(cfg.Grpc.BusinessService.Address)
	if err != nil {
		log.Fatalf("%v", err)
	}
	auditClient, err := client.NewAuditClient(cfg.Grpc.AuditingService.Address)
	if err != nil {
		log.Fatalf("%v", err)
	}
	lidClient, err := client.NewLidClient(cfg.Grpc.LidService.Address)
	if err != nil {
		log.Fatalf("%v", err)
	}
	userClient, err := client.NewUserClient(cfg.Grpc.UserService.Address)
	if err != nil {
		return nil
	}
	return &Clients{
		AuditClient:    auditClient,
		UserClient:     userClient,
		LidClient:      lidClient,
		BusinessClient: businessClient,
	}
}
