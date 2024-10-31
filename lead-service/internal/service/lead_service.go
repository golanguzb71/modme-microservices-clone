package service

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"lid-service/internal/repository"
	"lid-service/proto/pb"
)

type LeadService struct {
	pb.UnimplementedLeadServiceServer
	repo *repository.LeadRepository
}

func NewLeadService(repo *repository.LeadRepository) *LeadService {
	return &LeadService{repo: repo}
}

func (s *LeadService) CreateLead(ctx context.Context, req *pb.CreateLeadRequest) (*pb.AbsResponse, error) {
	err := s.repo.CreateLead(req.Title)
	if err != nil {
		return &pb.AbsResponse{Status: 500, Message: "Failed to create lead: " + err.Error()}, err
	}
	return &pb.AbsResponse{Status: 200, Message: "Lead created successfully"}, nil
}

func (s *LeadService) GetLeadCommon(ctx context.Context, req *pb.GetLeadCommonRequest) (*pb.GetLeadCommonResponse, error) {
	return s.repo.GetLeadCommon(req)
}

func (s *LeadService) UpdateLead(ctx context.Context, req *pb.UpdateLeadRequest) (*pb.AbsResponse, error) {
	err := s.repo.UpdateLead(req.Id, req.Title)
	if err != nil {
		return &pb.AbsResponse{Status: 500, Message: "Failed to update lead: " + err.Error()}, err
	}
	return &pb.AbsResponse{Status: 200, Message: "Lead updated successfully"}, nil
}

func (s *LeadService) DeleteLead(ctx context.Context, req *pb.DeleteAbsRequest) (*pb.AbsResponse, error) {
	err := s.repo.DeleteLead(req.Id)
	if err != nil {
		return &pb.AbsResponse{Status: 500, Message: "Failed to delete lead: " + err.Error()}, err
	}
	return &pb.AbsResponse{Status: 200, Message: "Lead deleted successfully"}, nil
}

func (s *LeadService) GetListSection(ctx context.Context, req *emptypb.Empty) (*pb.GetLeadListResponse, error) {
	return s.repo.GetAllLeads()
}

func (s *LeadService) GetLeadReports(ctx context.Context, req *pb.GetLeadReportsRequest) (*pb.GetLeadReportsResponse, error) {
	return s.repo.GetLeadReports(req.EndYear, req.StartYear)
}
