package grpc

import (
	"api-gateway/config"
	client "api-gateway/internal/clients"
	"log"
)

type Clients struct {
	UserClient      *client.UserClient
	LidClient       *client.LidClient
	EducationClient *client.EducationClient
	FinanceClient   *client.FinanceClient
}

func InitializeGrpcClients(cfg *config.Config) *Clients {
	educationClient, err := client.NewEducationClient(cfg.Grpc.EducationService.Address)
	if err != nil {
		log.Fatalf("%v", err)
	}
	userClient, err := client.NewUserClient(cfg.Grpc.UserService.Address)
	if err != nil {
		log.Fatalf("%v", err)
	}
	lidClient, err := client.NewLidClient(cfg.Grpc.LidService.Address)
	if err != nil {
		log.Fatalf("%v", err)
	}
	financeClient, err := client.NewFinanceClient(cfg.Grpc.FinanceService.Address)
	if err != nil {
		log.Fatalf("%v", err)
	}

	return &Clients{
		UserClient:      userClient,
		LidClient:       lidClient,
		EducationClient: educationClient,
		FinanceClient:   financeClient,
	}
}
