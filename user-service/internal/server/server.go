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
	"user-service/migrations"
	"user-service/proto/pb"
)

func RunServer() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	db, err := repository.NewPostgresRepository(cfg)
	if err != nil {
		log.Fatalf("failed to loading database %s", err)
	}
	defer db.Close()

	var groupClient *clients.GroupClient
	go func() {
		for {
			groupClient = clients.NewGroupClient(cfg.Grpc.EducationService.Address)
			if groupClient != nil {
				log.Println("Connected to education service.")
				break
			}
			log.Println("Waiting for education service to be available...")
			time.Sleep(3 * time.Second)
		}
	}()

	migrations.SetUpMigrating(cfg.Database.Action, db)
	userRepo := repository.NewUserRepository(db, groupClient)
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userRepo)

	listen, err := net.Listen("tcp", ":"+strconv.Itoa(cfg.Server.Port))
	if err != nil {
		log.Fatalf("error while listening %s", err)
	}
	server := grpc.NewServer()
	pb.RegisterUserServiceServer(server, userService)
	pb.RegisterAuthServiceServer(server, authService)
	log.Printf("Server listening on port %v", cfg.Server.Port)
	if err := server.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
