package service

import (
	"context"
	"finance-service/internal/repository"
	"finance-service/proto/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

type PaymentService struct {
	pb.UnimplementedPaymentServiceServer
	repo *repository.PaymentRepository
}

func NewPaymentService(repo *repository.PaymentRepository) *PaymentService {
	return &PaymentService{repo: repo}
}

func (ps *PaymentService) PayStudent(ctx context.Context, req *pb.PayStudentRequest) (*pb.AbsResponse, error) {
	if req.Type == "PAID" {
		if err := ps.repo.PaidStudent(req.UserId, req.Comment, req.Sum, req.Date, req.Method, req.CreatedBy); err != nil {
			return nil, status.Errorf(codes.Canceled, err.Error())
		}
		return &pb.AbsResponse{
			Status:  http.StatusCreated,
			Message: "student paid added",
		}, nil
	} else if req.Type == "UNPAID" {
		if err := ps.repo.UnPaidStudent(req.UserId, req.Comment, req.Sum, req.Date, req.Method, req.CreatedBy); err != nil {
			return nil, status.Errorf(codes.Canceled, err.Error())
		}
		return &pb.AbsResponse{
			Status:  http.StatusCreated,
			Message: "student unpaid",
		}, nil
	}
	return nil, status.Error(codes.Aborted, "request type invalid")
}
