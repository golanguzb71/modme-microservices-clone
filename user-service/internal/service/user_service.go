package service

import (
	"context"
	"fmt"
	"user-service/internal/repository"
	"user-service/internal/utils"
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
	companyId := utils.GetCompanyDetails(ctx)

	return u.userRepo.CreateUser(companyId, req.Gender, req.PhoneNumber, req.BirthDate, req.FullName, req.Password, req.Role)
}

func (u *UserService) GetTeachers(ctx context.Context, req *pb.GetTeachersRequest) (*pb.GetTeachersResponse, error) {
	companyId := utils.GetCompanyDetails(ctx)
	fmt.Println("here is the company id ", companyId)

	return u.userRepo.GetTeachers(companyId, req.IsDeleted)
}

func (u *UserService) GetUserById(ctx context.Context, req *pb.UserAbsRequest) (*pb.GetUserByIdResponse, error) {
	companyId := utils.GetCompanyDetails(ctx)

	return u.userRepo.GetUserById(companyId, req.UserId)
}

func (u *UserService) UpdateUserById(ctx context.Context, req *pb.UpdateUserRequest) (*pb.AbsResponse, error) {
	companyId := utils.GetCompanyDetails(ctx)

	return u.userRepo.UpdateUser(companyId, req.Id, req.Name, req.Gender, req.Role, req.BirthDate, req.PhoneNumber)
}

func (u *UserService) DeleteUserById(ctx context.Context, req *pb.UserAbsRequest) (*pb.AbsResponse, error) {
	companyId := utils.GetCompanyDetails(ctx)

	return u.userRepo.DeleteUser(companyId, req.UserId)
}

func (u *UserService) GetAllEmployee(ctx context.Context, req *pb.GetAllEmployeeRequest) (*pb.GetAllEmployeeResponse, error) {
	companyId := utils.GetCompanyDetails(ctx)

	return u.userRepo.GetAllEmployee(companyId, req.IsArchived)
}

func (u *UserService) GetAllStuff(ctx context.Context, req *pb.GetAllEmployeeRequest) (*pb.GetAllStuffResponse, error) {
	companyId := utils.GetCompanyDetails(ctx)

	return u.userRepo.GetAllStuff(companyId, req.IsArchived)
}

func (u *UserService) GetHistoryByUserId(ctx context.Context, req *pb.UserAbsRequest) (*pb.GetHistoryByUserIdResponse, error) {
	companyId := utils.GetCompanyDetails(ctx)

	return u.userRepo.GetHistoryByUserId(companyId, req.UserId)
}

func (u *UserService) UpdateUserPassword(ctx context.Context, req *pb.UpdateUserPasswordRequest) (*pb.AbsResponse, error) {
	companyId := utils.GetCompanyDetails(ctx)

	return u.userRepo.UpdateUserPassword(companyId, req.UserId, req.NewPassword)
}
