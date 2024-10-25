package server

import (
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"lid-service/config"
	"lid-service/internal/clients"
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

	// lead_service_clients_start
	groupClient := clients.NewGroupClient(cfg.Grpc.EducationService.Address)
	studentClient := clients.NewStudentClient(cfg.Grpc.EducationService.Address)

	// lead_service_services_start
	expectRepo := repository.NewExpectRepository(db)
	setRepo := repository.NewSetRepository(db)
	leadRepo := repository.NewLeadRepository(db)
	leadDataRepo := repository.NewLeadDataRepository(db)

	leadService := service.NewLeadService(leadRepo)
	expectService := service.NewExpectService(expectRepo)
	setService := service.NewSetService(setRepo, groupClient, studentClient)
	leadDataService := service.NewLeadDataService(leadDataRepo)
	// lead_service_services_end

	lis, err := net.Listen("tcp", ":"+strconv.Itoa(cfg.Server.Port))
	if err != nil {
		log.Fatalf("Failed to listen on port %v: %v", cfg.Server.Port, err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterLeadServiceServer(grpcServer, leadService)
	pb.RegisterLeadDataServiceServer(grpcServer, leadDataService)
	pb.RegisterExpectServiceServer(grpcServer, expectService)
	pb.RegisterSetServiceServer(grpcServer, setService)

	log.Printf("Server listening on port %v", cfg.Server.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
