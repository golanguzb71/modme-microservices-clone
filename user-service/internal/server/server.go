package server

import (
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"log"
	"net"
	"strconv"
	"time"
	"user-service/config"
	"user-service/internal/clients"
	"user-service/internal/repository"
	"user-service/internal/service"
	"user-service/internal/utils"
	"user-service/proto/pb"
)

func RunServer() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := repository.NewPostgresRepository(cfg)
	if err != nil {
		log.Fatalf("Failed to load database: %v", err)
	}
	defer db.Close()

	groupClientChan := make(chan *clients.GroupClient)
	go func() {
		time.Sleep(2 * time.Second)
		var client *clients.GroupClient
		for {
			client, err = clients.NewGroupClient(cfg.Grpc.EducationService.Address)
			if err == nil {
				log.Println("Connected to Education Service successfully.")
				groupClientChan <- client
				close(groupClientChan)
				break
			}
			log.Printf("Waiting for Education Service...")
			time.Sleep(2 * time.Second)
		}
	}()

	userRepo := repository.NewUserRepository(db, groupClientChan)
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userRepo)

	listen, err := net.Listen("tcp", ":"+strconv.Itoa(cfg.Server.Port))
	if err != nil {
		log.Fatalf("Error while listening: %v", err)
	}

	server := grpc.NewServer(
		grpc.UnaryInterceptor(utils.RecoveryInterceptor),
	)
	pb.RegisterUserServiceServer(server, userService)
	pb.RegisterAuthServiceServer(server, authService)

	log.Printf("Server listening on port %v", cfg.Server.Port)

	if err := server.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
