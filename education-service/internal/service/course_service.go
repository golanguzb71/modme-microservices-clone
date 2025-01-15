package service

import (
	"context"
	"education-service/internal/repository"
	"education-service/internal/utils"
	"education-service/proto/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CourseService struct {
	pb.UnimplementedCourseServiceServer
	repo *repository.CourseRepository
}

func NewCourseService(repo *repository.CourseRepository) *CourseService {
	return &CourseService{repo: repo}
}

func (s *CourseService) CreateCourse(ctx context.Context, req *pb.CreateCourseRequest) (*pb.AbsResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	err := s.repo.CreateCourse(companyId, req.Name, req.Description, req.LessonDuration, req.CourseDuration, req.Price)
	if err != nil {
		return nil, err
	}
	return &pb.AbsResponse{Status: 200, Message: "Course created successfully"}, nil
}

func (s *CourseService) UpdateCourse(ctx context.Context, req *pb.AbsCourse) (*pb.AbsResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	err := s.repo.UpdateCourse(companyId, req.Name, req.Description, req.Id, req.LessonDuration, req.CourseDuration, req.Price)
	if err != nil {
		return nil, err
	}
	return &pb.AbsResponse{Status: 200, Message: "Course updated successfully"}, nil
}

func (s *CourseService) DeleteCourse(ctx context.Context, req *pb.DeleteAbsRequest) (*pb.AbsResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	err := s.repo.DeleteCourse(companyId, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.AbsResponse{Status: 200, Message: "Expectation deleted successfully"}, nil
}

func (s *CourseService) GetCourses(ctx context.Context, req *emptypb.Empty) (*pb.GetUpdateCourseAbs, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	return s.repo.GetCourse(companyId)
}

func (s *CourseService) GetCourseById(ctx context.Context, req *pb.GetCourseByIdRequest) (*pb.GetCourseByIdResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	return s.repo.GetCourseById(companyId, req.Id)
}
