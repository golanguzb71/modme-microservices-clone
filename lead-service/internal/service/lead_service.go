package service

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"lid-service/internal/repository"
	"lid-service/internal/utils"
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
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	err := s.repo.CreateLead(companyId, req.Title)
	if err != nil {
		return &pb.AbsResponse{Status: 500, Message: "Failed to create lead: " + err.Error()}, err
	}
	return &pb.AbsResponse{Status: 200, Message: "Lead created successfully"}, nil
}

func (s *LeadService) GetLeadCommon(ctx context.Context, req *pb.GetLeadCommonRequest) (*pb.GetLeadCommonResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	return s.repo.GetLeadCommon(companyId, req)
}

func (s *LeadService) UpdateLead(ctx context.Context, req *pb.UpdateLeadRequest) (*pb.AbsResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	err := s.repo.UpdateLead(companyId, req.Id, req.Title)
	if err != nil {
		return &pb.AbsResponse{Status: 500, Message: "Failed to update lead: " + err.Error()}, err
	}
	return &pb.AbsResponse{Status: 200, Message: "Lead updated successfully"}, nil
}

func (s *LeadService) DeleteLead(ctx context.Context, req *pb.DeleteAbsRequest) (*pb.AbsResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	err := s.repo.DeleteLead(companyId, req.Id)
	if err != nil {
		return &pb.AbsResponse{Status: 500, Message: "Failed to delete lead: " + err.Error()}, err
	}
	return &pb.AbsResponse{Status: 200, Message: "Lead deleted successfully"}, nil
}

func (s *LeadService) GetListSection(ctx context.Context, req *emptypb.Empty) (*pb.GetLeadListResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	return s.repo.GetAllLeads(companyId)
}

func (s *LeadService) GetLeadReports(ctx context.Context, req *pb.GetLeadReportsRequest) (*pb.GetLeadReportsResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	return s.repo.GetLeadReports(companyId, req.EndYear, req.StartYear)
}

func (s *LeadService) GetActiveLeadCount(ctx context.Context, req *emptypb.Empty) (*pb.GetActiveLeadCountResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	return s.repo.GetActiveLeadCount(companyId)
}
