package service

import (
	"context"
	"education-service/internal/repository"
	"education-service/proto/pb"
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
	id, err := s.repo.CreateGroup(req.Name, req.CourseId, req.TeacherId, req.Type, req.Days, req.RoomId, req.LessonStartTime, req.GroupStartDate, req.GroupEndDate)
	if err != nil {
		return nil, err
	}
	return &pb.AbsResponse{Status: 200, Message: id}, nil
}

func (s *GroupService) UpdateGroup(ctx context.Context, req *pb.GetUpdateGroupAbs) (*pb.AbsResponse, error) {
	err := s.repo.UpdateGroup(req.Id, req.Name, req.CourseId, req.TeacherId, req.Type, req.Days, req.RoomId, req.LessonStartTime, req.GroupStartDate, req.GroupEndDate)
	if err != nil {
		return nil, err
	}
	return &pb.AbsResponse{Status: 200, Message: "Group updated successfully"}, nil
}

func (s *GroupService) DeleteGroup(ctx context.Context, req *pb.DeleteAbsRequest) (*pb.AbsResponse, error) {
	err := s.repo.DeleteGroup(req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.AbsResponse{Status: 200, Message: "Group toggled successfully"}, nil
}

func (s *GroupService) GetGroups(ctx context.Context, req *pb.GetGroupsRequest) (*pb.GetGroupsResponse, error) {
	group, err := s.repo.GetGroup(req.Page.Page, req.Page.Size, req.IsArchived)
	if err != nil {
		log.Printf("Error in GetGroups: %v", err)
		return nil, err
	}
	return group, nil
}

func (s *GroupService) GetGroupById(ctx context.Context, req *pb.GetGroupByIdRequest) (*pb.GetGroupAbsResponse, error) {
	group, err := s.repo.GetGroupById(req.Id, req.ActionRole, req.ActionId)
	if err != nil {
		log.Printf("Error in GetGroupById: %v", err)
		return nil, err
	}
	return group, nil
}

func (s *GroupService) GetGroupsByCourseId(ctx context.Context, req *pb.GetGroupByIdRequest) (*pb.GetGroupsByCourseResponse, error) {
	groups, err := s.repo.GetGroupByCourseId(req.Id)
	if err != nil {
		log.Printf("Error in GetGroupById: %v", err)
		return nil, err
	}
	return groups, nil
}

func (s *GroupService) GetGroupsByTeacherId(ctx context.Context, req *pb.GetGroupsByTeacherIdRequest) (*pb.GetGroupsByTeacherResponse, error) {
	return s.repo.GetGroupByTeacherId(req.TeacherId, req.IsArchived)
}

func (s *GroupService) GetCommonInformationEducation(ctx context.Context, req *emptypb.Empty) (*pb.GetCommonInformationEducationResponse, error) {
	return s.repo.GetCommonInformationEducation()
}

func (s *GroupService) GetGroupsByStudentId(ctx context.Context, req *pb.StudentIdRequest) (*pb.GetGroupsByStudentResponse, error) {
	return s.repo.GetGroupsByStudentId(req.StudentId)
}

func (s *GroupService) GetLeftAfterTrialPeriod(ctx context.Context, req *pb.GetLeftAfterTrialPeriodRequest) (*pb.GetLeftAfterTrialPeriodResponse, error) {
	return s.repo.GetLeftAfterTrial(req.From, req.To, req.Page, req.Size)
}
