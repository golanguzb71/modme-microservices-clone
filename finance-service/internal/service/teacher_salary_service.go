package service

import (
	"context"
	"finance-service/internal/repository"
	"finance-service/proto/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type TeacherSalaryService struct {
	pb.UnimplementedTeacherSalaryServiceServer
	repo *repository.TeacherSalaryRepository
}

func NewTeacherSalaryService(repo *repository.TeacherSalaryRepository) *TeacherSalaryService {
	return &TeacherSalaryService{
		repo: repo,
	}
}

func (ts *TeacherSalaryService) CreateTeacherSalary(ctx context.Context, req *pb.CreateTeacherSalaryRequest) (*pb.AbsResponse, error) {
	return ts.repo.CreateTeacherSalary(req.Amount, req.TeacherId, req.Type)
}
func (ts *TeacherSalaryService) DeleteTeacherSalary(ctx context.Context, req *pb.DeleteTeacherSalaryRequest) (*pb.AbsResponse, error) {
	return ts.repo.DeleteTeacherSalary(req.TeacherId)
}
func (ts *TeacherSalaryService) GetTeacherSalary(ctx context.Context, req *emptypb.Empty) (*pb.GetTeachersSalaryRequest, error) {
	return ts.repo.GetTeacherSalary()
}

func (ts *TeacherSalaryService) GetTeacherSalaryByTeacherID(ctx context.Context, req *pb.DeleteTeacherSalaryRequest) (*pb.AbsGetTeachersSalary, error) {
	return ts.repo.GetTeacherSalaryByTeacherID(req.TeacherId)
}
