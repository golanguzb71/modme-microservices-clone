package service

import (
	"context"
	"finance-service/internal/repository"
	"finance-service/internal/utils"
	"finance-service/proto/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	if err := ds.repo.CreateDiscount(ctx, companyId, req.GroupId, req.StudentId, req.DiscountPrice, req.Comment, req.StartDate, req.EndDate, req.WithTeacher); err != nil {
		return nil, err
	}
	return &pb.AbsResponse{
		Status:  http.StatusOK,
		Message: "discount created",
	}, nil
}

func (ds *DiscountService) DeleteDiscount(ctx context.Context, req *pb.AbsDiscountRequest) (*pb.AbsResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	if err := ds.repo.DeleteDiscount(companyId, req.GroupId, req.StudentId); err != nil {
		return nil, err
	}
	return &pb.AbsResponse{
		Status:  http.StatusOK,
		Message: "Discount deleted",
	}, nil
}

func (ds *DiscountService) GetHistoryDiscount(ctx context.Context, req *pb.GetHistoryDiscountRequest) (*pb.GetHistoryDiscountResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	return ds.repo.GetHistoryDiscount(ctx, companyId, req.StudentId)
}

func (ds *DiscountService) GetAllInformationDiscount(ctx context.Context, req *pb.GetInformationDiscountRequest) (*pb.GetInformationDiscountResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	return ds.repo.GetAllDiscountByGroup(ctx, companyId, req.GroupId)
}

func (ds *DiscountService) GetDiscountByStudentId(ctx context.Context, req *pb.GetDiscountByStudentIdRequest) (*pb.GetDiscountByStudentIdResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	return ds.repo.GetDiscountByStudentId(companyId, req.StudentId, req.GroupId)
}
