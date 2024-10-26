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
		log.Fatalf("Failed to load database: %v", err)
	}
	defer db.Close()

	var groupClient *clients.GroupClient
	for {
		groupClient, err = clients.NewGroupClient(cfg.Grpc.EducationService.Address)
		if err != nil {
			log.Println("Failed to connect to education service, retrying in 3 seconds...")
			time.Sleep(3 * time.Second)
			continue
		}
		log.Println("Connected to education service.")
		break
	}

	migrations.SetUpMigrating(cfg.Database.Action, db)

	userRepo := repository.NewUserRepository(db, groupClient)
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
