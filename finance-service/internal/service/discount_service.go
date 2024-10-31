package service

import (
	"context"
	"finance-service/internal/repository"
	"finance-service/proto/pb"
	"net/http"
)

type DiscountService struct {
	pb.UnimplementedDiscountServiceServer
	repo *repository.DiscountRepository
}

func NewDiscountService(repo *repository.DiscountRepository) *DiscountService {
	return &DiscountService{repo: repo}
}

func (ds *DiscountService) CreateDiscount(ctx context.Context, req *pb.AbsDiscountRequest) (*pb.AbsResponse, error) {
	if err := ds.repo.CreateDiscount(req.GroupId, req.StudentId, req.DiscountPrice, req.Comment, req.StartDate, req.EndDate, req.WithTeacher); err != nil {
		return nil, err
	}
	return &pb.AbsResponse{
		Status:  http.StatusOK,
		Message: "discount created",
	}, nil
}

func (ds *DiscountService) DeleteDiscount(ctx context.Context, req *pb.AbsDiscountRequest) (*pb.AbsResponse, error) {
	if err := ds.repo.DeleteDiscount(req.GroupId, req.StudentId); err != nil {
		return nil, err
	}
	return &pb.AbsResponse{
		Status:  http.StatusOK,
		Message: "Discount deleted",
	}, nil
}

func (ds *DiscountService) GetAllInformationDiscount(ctx context.Context, req *pb.GetInformationDiscountRequest) (*pb.GetInformationDiscountResponse, error) {
	return ds.repo.GetAllDiscountByGroup(req.GroupId)
}
