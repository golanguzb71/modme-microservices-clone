package server

import (
	"education-service/config"
	"education-service/internal/clients"
	"education-service/internal/repository"
	"education-service/internal/service"
	"education-service/internal/utils"
	"education-service/migrations"
	"education-service/proto/pb"
	"fmt"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"log"
	"net"
	"strconv"
	"time"
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

	userClient, err := clients.NewUserClient(cfg.Grpc.UserService.Address)
	if err != nil {
		log.Fatalf("error %v", err)
	}

	financeClientChan := make(chan *clients.FinanceClient)
	go func() {
		time.Sleep(2 * time.Second)
		var client *clients.FinanceClient
		for {
			fmt.Println(cfg.Grpc.FinanceService.Address)
			client, err = clients.NewFinanceClient(cfg.Grpc.FinanceService.Address)
			if err == nil {
				log.Println("Connected to Finance Service successfully.")
				financeClientChan <- client
				close(financeClientChan)
				break
			}
			log.Printf("Waiting for Finance Service...")
			time.Sleep(2 * time.Second)
		}
	}()

	roomRepo := repository.NewRoomRepository(db)
	roomService := service.NewRoomService(roomRepo)
	courseRepo := repository.NewCourseRepository(db)
	courseService := service.NewCourseService(courseRepo)
	groupRepo := repository.NewGroupRepository(db, userClient)
	groupService := service.NewGroupService(groupRepo)
	attendanceRepo := repository.NewAttendanceRepository(db, financeClientChan)
	attendanceService := service.NewAttendanceService(attendanceRepo)
	studentRepo := repository.NewStudentRepository(db, userClient, financeClientChan)
	studentService := service.NewStudentService(studentRepo)
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(cfg.Server.Port))
	if err != nil {
		log.Fatalf("Failed to listen on port %v: %v", cfg.Server.Port, err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(utils.RecoveryInterceptor),
	)
	pb.RegisterRoomServiceServer(grpcServer, roomService)
	pb.RegisterCourseServiceServer(grpcServer, courseService)
	pb.RegisterGroupServiceServer(grpcServer, groupService)
	pb.RegisterAttendanceServiceServer(grpcServer, attendanceService)
	pb.RegisterStudentServiceServer(grpcServer, studentService)
	log.Printf("Server listening on port %v", cfg.Server.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
