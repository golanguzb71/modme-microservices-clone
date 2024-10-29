package server

import (
	"finance-service/config"
	"finance-service/internal/clients"
	"finance-service/internal/repository"
	"finance-service/internal/service"
	"finance-service/internal/utils"
	"finance-service/proto/pb"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"log"
	"net"
	"strconv"
)

func RunServer() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf(err.Error())
	}
	db, err := repository.NewFinanceDB(&cfg.Database)
	if err != nil {
		log.Fatalf(err.Error())
	}
	educationClient, err := clients.NewEducationClient(cfg.Grpc.EducationService.Address)
	if err != nil {
		log.Fatalf(err.Error())
	}
	discountRepo := repository.NewDiscountRepository(db, educationClient)
	discountService := service.NewDiscountService(discountRepo)
	categoryRepo := repository.NewCategoryRepository(db)
	categoryService := service.NewCategoryService(categoryRepo)
	expenseRepo := repository.NewExpenseRepository(db)
	expenseService := service.NewExpenseService(expenseRepo)
	list, err := net.Listen("tcp", ":"+strconv.Itoa(cfg.Server.Port))
	if err != nil {
		log.Fatalf(err.Error())
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(utils.RecoveryInterceptor),
	)

	pb.RegisterDiscountServiceServer(grpcServer, discountService)
	pb.RegisterCategoryServiceServer(grpcServer, categoryService)
	pb.RegisterExpenseServiceServer(grpcServer, expenseService)
	log.Printf("Server listening on port %v", cfg.Server.Port)
	if err := grpcServer.Serve(list); err != nil {
		log.Fatalf("Failed to serve  %v", err)
	}
}
