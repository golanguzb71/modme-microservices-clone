package service

import (
	"context"
	"finance-service/internal/repository"
	"finance-service/proto/pb"
	"google.golang.org/protobuf/types/known/emptypb"
	"net/http"
)

type CategoryService struct {
	pb.UnimplementedCategoryServiceServer
	repo *repository.CategoryRepository
}

func NewCategoryService(repo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (c *CategoryService) CreateCategory(ctx context.Context, req *pb.CreateCategoryRequest) (*pb.AbsResponse, error) {
	if err := c.repo.CreateCategory(req.Name, req.Desc); err != nil {
		return nil, err
	}
	return &pb.AbsResponse{
		Status:  http.StatusCreated,
		Message: "category created",
	}, nil
}
func (c *CategoryService) DeleteCategory(ctx context.Context, req *pb.DeleteAbsRequest) (*pb.AbsResponse, error) {
	if err := c.repo.DeleteCategory(req.Id); err != nil {
		return nil, err
	}
	return &pb.AbsResponse{
		Status:  http.StatusOK,
		Message: "category deleted",
	}, nil
}
func (c *CategoryService) GetAllCategory(ctx context.Context, req *emptypb.Empty) (*pb.GetAllCategoryRequest, error) {
	return c.repo.GetAllCategory()
}
