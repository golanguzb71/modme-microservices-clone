package service

import (
	"context"
	"lid-service/internal/repository"
	"lid-service/proto/pb"
)

type SetService struct {
	pb.UnimplementedSetServiceServer
	repo *repository.SetRepository
}

func NewSetService(repo *repository.SetRepository) *SetService {
	return &SetService{repo: repo}
}

func (s *SetService) CreateSet(ctx context.Context, req *pb.CreateSetRequest) (*pb.AbsResponse, error) {
	err := s.repo.CreateSet(req.Title, req.CourseId, req.TeacherId, req.DateType, req.Date, req.LessonStartTime)
	if err != nil {
		return &pb.AbsResponse{Status: 500, Message: "Failed to create set: " + err.Error()}, err
	}
	return &pb.AbsResponse{Status: 200, Message: "Set created successfully"}, nil
}

func (s *SetService) UpdateSet(ctx context.Context, req *pb.UpdateSetRequest) (*pb.AbsResponse, error) {
	err := s.repo.UpdateSet(req.Id, req.Title, req.CourseId, req.TeacherId, req.DateType, req.Date, req.LessonStartTime)
	if err != nil {
		return &pb.AbsResponse{Status: 500, Message: "Failed to update set: " + err.Error()}, err
	}
	return &pb.AbsResponse{Status: 200, Message: "Set updated successfully"}, nil
}

func (s *SetService) DeleteSet(ctx context.Context, req *pb.DeleteAbsRequest) (*pb.AbsResponse, error) {
	err := s.repo.DeleteSet(req.Id)
	if err != nil {
		return &pb.AbsResponse{Status: 500, Message: "Failed to delete set: " + err.Error()}, err
	}
	return &pb.AbsResponse{Status: 200, Message: "Set deleted successfully"}, nil
}
