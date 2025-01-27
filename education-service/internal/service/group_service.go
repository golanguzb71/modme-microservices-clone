package service

import (
	"context"
	"education-service/internal/repository"
	"education-service/internal/utils"
	"education-service/proto/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
)

type GroupService struct {
	pb.UnimplementedGroupServiceServer
	repo *repository.GroupRepository
}

func NewGroupService(repo *repository.GroupRepository) *GroupService {
	return &GroupService{repo: repo}
}

func (s *GroupService) CreateGroup(ctx context.Context, req *pb.CreateGroupRequest) (*pb.AbsResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	id, err := s.repo.CreateGroup(companyId, req.Name, req.CourseId, req.TeacherId, req.Type, req.Days, req.RoomId, req.LessonStartTime, req.GroupStartDate, req.GroupEndDate)
	if err != nil {
		return nil, err
	}
	return &pb.AbsResponse{Status: 200, Message: id}, nil
}

func (s *GroupService) UpdateGroup(ctx context.Context, req *pb.GetUpdateGroupAbs) (*pb.AbsResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	err := s.repo.UpdateGroup(companyId, req.Id, req.Name, req.CourseId, req.TeacherId, req.Type, req.Days, req.RoomId, req.LessonStartTime, req.GroupStartDate, req.GroupEndDate)
	if err != nil {
		return nil, err
	}
	return &pb.AbsResponse{Status: 200, Message: "Group updated successfully"}, nil
}

func (s *GroupService) DeleteGroup(ctx context.Context, req *pb.DeleteAbsRequest) (*pb.AbsResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	err := s.repo.DeleteGroup(companyId, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.AbsResponse{Status: 200, Message: "Group toggled successfully"}, nil
}

func (s *GroupService) GetGroups(ctx context.Context, req *pb.GetGroupsRequest) (*pb.GetGroupsResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	group, err := s.repo.GetGroup(ctx, companyId, req.Page.Page, req.Page.Size, req.IsArchived)
	if err != nil {
		log.Printf("Error in GetGroups: %v", err)
		return nil, err
	}
	return group, nil
}

func (s *GroupService) GetGroupById(ctx context.Context, req *pb.GetGroupByIdRequest) (*pb.GetGroupAbsResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	group, err := s.repo.GetGroupById(ctx, companyId, req.Id, req.ActionRole, req.ActionId)
	if err != nil {
		log.Printf("Error in GetGroupById: %v", err)
		return nil, err
	}
	return group, nil
}

func (s *GroupService) GetGroupsByCourseId(ctx context.Context, req *pb.GetGroupByIdRequest) (*pb.GetGroupsByCourseResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	groups, err := s.repo.GetGroupByCourseId(ctx, companyId, req.Id)
	if err != nil {
		log.Printf("Error in GetGroupById: %v", err)
		return nil, err
	}
	return groups, nil
}

func (s *GroupService) GetGroupsByTeacherId(ctx context.Context, req *pb.GetGroupsByTeacherIdRequest) (*pb.GetGroupsByTeacherResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.NotFound, "error while getting company from context")
	}
	return s.repo.GetGroupByTeacherId(companyId, req.TeacherId, req.IsArchived)
}

func (s *GroupService) GetCommonInformationEducation(ctx context.Context, req *emptypb.Empty) (*pb.GetCommonInformationEducationResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	return s.repo.GetCommonInformationEducation(companyId)
}

func (s *GroupService) GetGroupsByStudentId(ctx context.Context, req *pb.StudentIdRequest) (*pb.GetGroupsByStudentResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	return s.repo.GetGroupsByStudentId(companyId, req.StudentId)
}

func (s *GroupService) GetLeftAfterTrialPeriod(ctx context.Context, req *pb.GetLeftAfterTrialPeriodRequest) (*pb.GetLeftAfterTrialPeriodResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	return s.repo.GetLeftAfterTrial(companyId, req.From, req.To, req.Page, req.Size)
}
