package service

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"lid-service/internal/clients"
	"lid-service/internal/repository"
	"lid-service/internal/utils"
	"lid-service/proto/pb"
	"strconv"
	"time"
)

type SetService struct {
	pb.UnimplementedSetServiceServer
	repo          *repository.SetRepository
	groupClient   *clients.GroupClient
	studentClient *clients.StudentClient
	userClient    *clients.UserClient
}

func NewSetService(repo *repository.SetRepository, client *clients.GroupClient, studentClient *clients.StudentClient, userClient *clients.UserClient) *SetService {
	return &SetService{repo: repo, groupClient: client, studentClient: studentClient, userClient: userClient}
}

func (s *SetService) CreateSet(ctx context.Context, req *pb.CreateSetRequest) (*pb.AbsResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	err := s.repo.CreateSet(companyId, req.Title, req.CourseId, req.TeacherId, req.DateType, req.Date, req.LessonStartTime)
	if err != nil {
		return &pb.AbsResponse{Status: 500, Message: "Failed to create set: " + err.Error()}, err
	}
	return &pb.AbsResponse{Status: 200, Message: "Set created successfully"}, nil
}

func (s *SetService) UpdateSet(ctx context.Context, req *pb.UpdateSetRequest) (*pb.AbsResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	err := s.repo.UpdateSet(companyId, req.Id, req.Title, req.CourseId, req.TeacherId, req.DateType, req.Date, req.LessonStartTime)
	if err != nil {
		return &pb.AbsResponse{Status: 500, Message: "Failed to update set: " + err.Error()}, err
	}
	return &pb.AbsResponse{Status: 200, Message: "Set updated successfully"}, nil
}

func (s *SetService) DeleteSet(ctx context.Context, req *pb.DeleteAbsRequest) (*pb.AbsResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	err := s.repo.DeleteSet(companyId, req.Id)
	if err != nil {
		return &pb.AbsResponse{Status: 500, Message: "Failed to delete set: " + err.Error()}, err
	}
	return &pb.AbsResponse{Status: 200, Message: "Set deleted successfully"}, nil
}

func (s *SetService) ChangeToSet(ctx context.Context, req *pb.ChangeToSetRequest) (*pb.AbsResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	if s.groupClient == nil {
		return nil, fmt.Errorf("uninitialized group detected")
	}

	if s.studentClient == nil {
		return nil, fmt.Errorf("uninitialized studentClient detected")
	}
	if s.repo == nil {
		return nil, fmt.Errorf("uninitialized repo detected")
	}
	courseId, err := strconv.ParseInt(req.CourseId, 10, 32)
	if err != nil {
		return nil, err
	}
	roomId, err := strconv.ParseInt(req.RoomId, 10, 32)
	if err != nil {
		return nil, err
	}
	createGroupReq := pb.CreateGroupRequest{
		Name:            req.Name,
		CourseId:        int32(courseId),
		TeacherId:       req.TeacherId,
		Type:            req.DateType,
		Days:            req.Days,
		RoomId:          int32(roomId),
		LessonStartTime: req.StartTime,
		GroupStartDate:  req.StartDate,
		GroupEndDate:    req.EndDate,
	}
	ctx, cancelFunc := utils.NewTimoutContext(ctx, companyId)
	defer cancelFunc()
	err, groupId := s.groupClient.CreateGroup(ctx, &createGroupReq)
	if err != nil {
		return nil, err
	}
	currentDate := time.Now().Format("2006-01-02")
	names, phoneNumbers, err := s.repo.GetLeadDataBySetId(companyId, req.SetId)
	if err != nil {
		return nil, err
	}
	length := min(len(names), len(phoneNumbers))
	for i := 0; i < length; i++ {
		_, err = s.studentClient.CreateStudent(ctx, phoneNumbers[i], names[i], "2006-12-14", groupId, currentDate, "1b39d121-7840-4411-bcfe-87c1beb9422b", true)
		if err != nil {
			return nil, err
		}
	}
	err = s.repo.DeleteSet(companyId, req.SetId)
	if err != nil {
		return nil, err
	}
	return &pb.AbsResponse{Status: 200, Message: "Set changed to group successfully"}, nil
}
func (s *SetService) GetById(ctx context.Context, req *pb.DeleteAbsRequest) (*pb.SetDataResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	resp, err := s.repo.GetById(companyId, req.Id)
	if err != nil {
		return nil, err
	}
	ctx, cancelFunc := utils.NewTimoutContext(ctx, companyId)
	defer cancelFunc()
	course, err := s.groupClient.GetCourse(ctx, resp.CourseId)
	if err != nil {
		return nil, err
	}
	user, err := s.userClient.GetUserById(ctx, resp.TeacherId)
	if err != nil {
		return nil, err
	}
	resp.CourseName = course.Name
	resp.TeacherName = user.Name
	return resp, nil
}
