package service

import (
	"context"
	"lid-service/internal/repository"
	"lid-service/proto/pb"
)

type LeadDataService struct {
	pb.UnimplementedLeadDataServiceServer
	repo *repository.LeadDataRepository
}

func NewLeadDataService(repo *repository.LeadDataRepository) *LeadDataService {
	return &LeadDataService{repo: repo}
}

// CreateLeadData handles lead data creation
func (s *LeadDataService) CreateLeadData(ctx context.Context, req *pb.CreateLeadDataRequest) (*pb.AbsResponse, error) {
	err := s.repo.CreateLeadData(&req.PhoneNumber, &req.LeadId, nil, nil, &req.Comment, &req.Name)
	if err != nil {
		return &pb.AbsResponse{Status: 500, Message: "Failed to create lead data: " + err.Error()}, err
	}
	return &pb.AbsResponse{Status: 200, Message: "Lead data created successfully"}, nil
}

// UpdateLeadData handles lead data updates
func (s *LeadDataService) UpdateLeadData(ctx context.Context, req *pb.UpdateLeadDataRequest) (*pb.AbsResponse, error) {
	err := s.repo.UpdateLeadData(req.Id, req.PhoneNumber, req.Comment, req.Type, req.SectionId)
	if err != nil {
		return &pb.AbsResponse{Status: 500, Message: "Failed to update lead data: " + err.Error()}, err
	}
	return &pb.AbsResponse{Status: 200, Message: "Lead data updated successfully"}, nil
}

// DeleteLeadData handles lead data deletion
func (s *LeadDataService) DeleteLeadData(ctx context.Context, req *pb.DeleteAbsRequest) (*pb.AbsResponse, error) {
	err := s.repo.DeleteLeadData(req.Id)
	if err != nil {
		return &pb.AbsResponse{Status: 500, Message: "Failed to delete lead data: " + err.Error()}, err
	}
	return &pb.AbsResponse{Status: 200, Message: "Lead data deleted successfully"}, nil
}

func (s *LeadDataService) ChangeLeadPlace(ctx context.Context, req *pb.ChangeLeadPlaceRequest) (*pb.AbsResponse, error) {
	err := s.repo.ChangeLeadPlace(&req.ChangedSet.Id, &req.ChangedSet.SectionType, &req.LeadDataId)
	if err != nil {
		return &pb.AbsResponse{Status: 500, Message: "Failed to change lead data: " + err.Error()}, err
	}
	return &pb.AbsResponse{Status: 200, Message: "Lead data changed successfully"}, nil
}
