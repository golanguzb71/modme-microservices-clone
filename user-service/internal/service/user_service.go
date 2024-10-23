package service

import (
	"context"
	"user-service/internal/repository"
	"user-service/proto/pb"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	userRepo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		userRepo: repo,
	}
}

func (u *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.AbsResponse, error) {
	return u.userRepo.CreateUser(req.Gender, req.PhoneNumber, req.BirthDate, req.FullName, req.Password, req.Role)
}

func (u *UserService) GetTeachers(ctx context.Context, req *pb.GetTeachersRequest) (*pb.GetTeachersResponse, error) {
	return u.userRepo.GetTeachers(req.IsDeleted)
}
