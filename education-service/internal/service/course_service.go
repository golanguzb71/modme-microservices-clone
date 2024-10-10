package service

import (
	"context"
	"education-service/internal/repository"
	"education-service/proto/pb"
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
	err := s.repo.CreateCourse(req.Name, req.Description, req.LessonDuration, req.CourseDuration, req.Price)
	if err != nil {
		return nil, err
	}
	return &pb.AbsResponse{Status: 200, Message: "Course created successfully"}, nil
}

func (s *CourseService) UpdateCourse(ctx context.Context, req *pb.AbsCourse) (*pb.AbsResponse, error) {
	err := s.repo.UpdateCourse(req.Name, req.Description, req.Id, req.LessonDuration, req.CourseDuration, req.Price)
	if err != nil {
		return nil, err
	}
	return &pb.AbsResponse{Status: 200, Message: "Course updated successfully"}, nil
}

func (s *CourseService) DeleteCourse(ctx context.Context, req *pb.DeleteAbsRequest) (*pb.AbsResponse, error) {
	err := s.repo.DeleteCourse(req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.AbsResponse{Status: 200, Message: "Expectation deleted successfully"}, nil
}

func (s *CourseService) GetCourses(ctx context.Context, req *emptypb.Empty) (*pb.GetUpdateCourseAbs, error) {
	return s.repo.GetCourse()
}
