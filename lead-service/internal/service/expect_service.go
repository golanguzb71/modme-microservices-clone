package service

import (
	"context"
	"lid-service/internal/repository"
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
	err := s.repo.CreateExpectation(req.Title)
	if err != nil {
		return nil, err
	}
	return &pb.AbsResponse{Status: 200, Message: "Expectation created successfully"}, nil
}

func (s *ExpectService) UpdateExpect(ctx context.Context, req *pb.UpdateExpectRequest) (*pb.AbsResponse, error) {
	err := s.repo.UpdateExpectation(req.Id, req.Title)
	if err != nil {
		return nil, err
	}
	return &pb.AbsResponse{Status: 200, Message: "Expectation updated successfully"}, nil
}

func (s *ExpectService) DeleteExpect(ctx context.Context, req *pb.DeleteAbsRequest) (*pb.AbsResponse, error) {
	err := s.repo.DeleteExpectation(req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.AbsResponse{Status: 200, Message: "Expectation deleted successfully"}, nil
}
