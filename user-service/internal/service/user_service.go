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

func (u *UserService) GetUserById(ctx context.Context, req *pb.UserAbsRequest) (*pb.GetUserByIdResponse, error) {
	return u.userRepo.GetUserById(req.UserId)
}

func (u *UserService) UpdateUserById(ctx context.Context, req *pb.UpdateUserRequest) (*pb.AbsResponse, error) {
	return u.userRepo.UpdateUser(req.Id, req.Name, req.Gender, req.Role, req.BirthDate, req.PhoneNumber)
}

func (u *UserService) DeleteUserById(ctx context.Context, req *pb.UserAbsRequest) (*pb.AbsResponse, error) {
	return u.userRepo.DeleteUser(req.UserId)
}

func (u *UserService) GetAllEmployee(ctx context.Context, req *pb.GetAllEmployeeRequest) (*pb.GetAllEmployeeResponse, error) {
	return u.userRepo.GetAllEmployee(req.IsArchived)
}

func (u *UserService) GetAllStuff(ctx context.Context, req *pb.GetAllEmployeeRequest) (*pb.GetAllStuffResponse, error) {
	return u.userRepo.GetAllStuff(req.IsArchived)
}

func (u *UserService) GetHistoryByUserId(ctx context.Context, req *pb.UserAbsRequest) (*pb.GetHistoryByUserIdResponse, error) {
	return u.userRepo.GetHistoryByUserId(req.UserId)
}

func (u *UserService) UpdateUserPassword(ctx context.Context, req *pb.UpdateUserPasswordRequest) (*pb.AbsResponse, error) {
	return u.userRepo.UpdateUserPassword(req.UserId, req.NewPassword)
}
