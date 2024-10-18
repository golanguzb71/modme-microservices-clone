package service

import (
	"context"
	"education-service/internal/repository"
	"education-service/proto/pb"
)

type StudentService struct {
	pb.UnimplementedStudentServiceServer
	repo *repository.StudentRepository
}

func NewStudentService(repo *repository.StudentRepository) *StudentService {
	return &StudentService{repo: repo}
}

func (s *StudentService) GetAllStudent(ctx context.Context, req *pb.GetAllStudentRequest) (*pb.GetAllStudentResponse, error) {
	return s.repo.GetAllStudent(req.Condition, req.Page, req.Size)
}

func (s *StudentService) CreateStudent(ctx context.Context, req *pb.CreateStudentRequest) (*pb.AbsResponse, error) {
	if err := s.repo.CreateStudent(req.CreatedBy, req.PhoneNumber, req.Name, req.GroupId, req.Address, req.AdditionalContact, req.DateFrom, req.DateOfBirth, req.Gender, req.PassportId, req.TelegramUsername); err != nil {
		return nil, err
	}
	return &pb.AbsResponse{
		Status:  200,
		Message: "student created successfully",
	}, nil
}

func (s *StudentService) UpdateStudent(ctx context.Context, req *pb.UpdateStudentRequest) (*pb.AbsResponse, error) {
	if err := s.repo.UpdateStudent(req.StudentId, req.PhoneNumber, req.Name, req.Address, req.AdditionalContact, req.DateOfBirth, req.Gender, req.PassportId); err != nil {
		return nil, err
	}
	return &pb.AbsResponse{
		Status:  200,
		Message: "student updated successfully",
	}, nil
}

func (s *StudentService) DeleteStudent(ctx context.Context, req *pb.DeleteAbsRequest) (*pb.AbsResponse, error) {
	if err := s.repo.DeleteStudent(req.Id); err != nil {
		return nil, err
	}
	return &pb.AbsResponse{
		Status:  200,
		Message: "accomplished",
	}, nil
}

func (s *StudentService) AddToGroup(ctx context.Context, req *pb.AddToGroupRequest) (*pb.AbsResponse, error) {
	if err := s.repo.AddToGroup(req.GroupId, req.StudentIds, req.CreatedDate); err != nil {
		return nil, err
	}
	return &pb.AbsResponse{
		Status:  200,
		Message: "students added to group",
	}, nil
}
