package server

import (
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"lid-service/config"
	"lid-service/internal/repository"
	"lid-service/internal/service"
	"lid-service/migrations"
	"lid-service/proto/pb"
	"log"
	"net"
	"strconv"
)

func RunServer() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	db, err := repository.NewPostgresDB(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()
	migrations.SetUpMigrating(cfg.Database.Action, db)

	lidRepo := repository.NewLidRepository(db)
	expectRepo := repository.NewExpectRepository(db)
	setRepo := repository.NewSetRepository(db)
	lidUserRepo := repository.NewLidUserRepository(db)

	lidService := service.NewLidService(lidRepo)
	expectService := service.NewExpectService(expectRepo)
	setService := service.NewSetService(setRepo)
	lidUserService := service.NewLidUserService(lidUserRepo)

	lis, err := net.Listen("tcp", ":"+strconv.Itoa(cfg.Server.Port))
	if err != nil {
		log.Fatalf("Failed to listen on port %v: %v", cfg.Server.Port, err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterLidUserServiceServer(grpcServer, lidUserService)
	pb.RegisterLidServiceServer(grpcServer, lidService)
	pb.RegisterExpectServiceServer(grpcServer, expectService)
	pb.RegisterSetServiceServer(grpcServer, setService)

	log.Printf("Server listening on port %v", cfg.Server.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
