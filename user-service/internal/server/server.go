package server

import (
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"log"
	"net"
	"strconv"
	"user-service/config"
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
	migrations.SetUpMigrating(cfg.Database.Action, db)
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)

	listen, err := net.Listen("tcp", ":"+strconv.Itoa(cfg.Server.Port))
	if err != nil {
		log.Fatalf("error while listening %s", err)
	}
	server := grpc.NewServer()
	pb.RegisterUserServiceServer(server, userService)
	log.Printf("Server listening on port %v", cfg.Server.Port)
	if err := server.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
