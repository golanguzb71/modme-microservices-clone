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
	err := s.repo.CreateGroup(req.Name, req.CourseId, req.TeacherId, req.Type, req.Days, req.RoomId, req.LessonStartTime, req.GroupStartDate, req.GroupEndDate)
	if err != nil {
		return nil, err
	}
	return &pb.AbsResponse{Status: 200, Message: "Group created successfully"}, nil
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
	return &pb.AbsResponse{Status: 200, Message: "Group deleted successfully"}, nil
}

func (s *GroupService) GetGroups(ctx context.Context, req *emptypb.Empty) (*pb.GetGroupsResponse, error) {
	log.Println("Received GetGroups request")
	group, err := s.repo.GetGroup()
	if err != nil {
		log.Printf("Error in GetGroups: %v", err)
		return nil, err
	}
	log.Println("Returning groups")
	return group, nil
}

func (s *GroupService) GetGroupById(ctx context.Context, req *pb.GetGroupByIdRequest) (*pb.GetGroupAbsResponse, error) {
	log.Printf("Received GetGroupById request for id: %s", req.Id)
	group, err := s.repo.GetGroupById(req.Id)
	if err != nil {
		log.Printf("Error in GetGroupById: %v", err)
		return nil, err
	}
	log.Printf("Returning group with id: %s", group.Id)
	return group, nil
}
