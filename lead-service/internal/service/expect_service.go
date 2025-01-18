package service

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"lid-service/internal/repository"
	"lid-service/internal/utils"
	"lid-service/proto/pb"
)

type ExpectService struct {
	pb.UnimplementedExpectServiceServer
	repo *repository.ExpectRepository
}

func NewExpectService(repo *repository.ExpectRepository) *ExpectService {
	return &ExpectService{repo: repo}
}

func (s *ExpectService) CreateExpect(ctx context.Context, req *pb.CreateExpectRequest) (*pb.AbsResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	err := s.repo.CreateExpectation(companyId, req.Title)
	if err != nil {
		return nil, err
	}
	return &pb.AbsResponse{Status: 200, Message: "Expectation created successfully"}, nil
}

func (s *ExpectService) UpdateExpect(ctx context.Context, req *pb.UpdateExpectRequest) (*pb.AbsResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	err := s.repo.UpdateExpectation(companyId, req.Id, req.Title)
	if err != nil {
		return nil, err
	}
	return &pb.AbsResponse{Status: 200, Message: "Expectation updated successfully"}, nil
}

func (s *ExpectService) DeleteExpect(ctx context.Context, req *pb.DeleteAbsRequest) (*pb.AbsResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	err := s.repo.DeleteExpectation(companyId, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.AbsResponse{Status: 200, Message: "Expectation deleted successfully"}, nil
}
