package service

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "company id required")
	}
	return u.userRepo.CreateUser(companyId, req.Gender, req.PhoneNumber, req.BirthDate, req.FullName, req.Password, req.Role)
}

func (u *UserService) GetTeachers(ctx context.Context, req *pb.GetTeachersRequest) (*pb.GetTeachersResponse, error) {
	companyId := utils.GetCompanyDetails(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, fmt.Sprintf("company id required  %v and %v", companyId, companyId == ""))
	}
	return u.userRepo.GetTeachers(ctx, companyId, req.IsDeleted)
}

func (u *UserService) GetUserById(ctx context.Context, req *pb.UserAbsRequest) (*pb.GetUserByIdResponse, error) {
	companyId := utils.GetCompanyDetails(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "company id required")
	}
	return u.userRepo.GetUserById(companyId, req.UserId)
}

func (u *UserService) UpdateUserById(ctx context.Context, req *pb.UpdateUserRequest) (*pb.AbsResponse, error) {
	companyId := utils.GetCompanyDetails(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "company id required")
	}
	return u.userRepo.UpdateUser(companyId, req.Id, req.Name, req.Gender, req.Role, req.BirthDate, req.PhoneNumber, req.Password)
}

func (u *UserService) DeleteUserById(ctx context.Context, req *pb.UserAbsRequest) (*pb.AbsResponse, error) {
	companyId := utils.GetCompanyDetails(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "company id required")
	}

	return u.userRepo.DeleteUser(companyId, req.UserId)
}

func (u *UserService) GetAllEmployee(ctx context.Context, req *pb.GetAllEmployeeRequest) (*pb.GetAllEmployeeResponse, error) {
	companyId := utils.GetCompanyDetails(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "company id required")
	}

	return u.userRepo.GetAllEmployee(companyId, req.IsArchived)
}

func (u *UserService) GetAllStuff(ctx context.Context, req *pb.GetAllEmployeeRequest) (*pb.GetAllStuffResponse, error) {
	companyId := utils.GetCompanyDetails(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "company id required")
	}

	return u.userRepo.GetAllStuff(companyId, req.IsArchived)
}

func (u *UserService) GetHistoryByUserId(ctx context.Context, req *pb.UserAbsRequest) (*pb.GetHistoryByUserIdResponse, error) {
	companyId := utils.GetCompanyDetails(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "company id required")
	}

	return u.userRepo.GetHistoryByUserId(companyId, req.UserId)
}

func (u *UserService) UpdateUserPassword(ctx context.Context, req *pb.UpdateUserPasswordRequest) (*pb.AbsResponse, error) {
	companyId := utils.GetCompanyDetails(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "company id required")
	}

	return u.userRepo.UpdateUserPassword(companyId, req.UserId, req.NewPassword)
}

func (u *UserService) GetUserByCompanyId(ctx context.Context, req *pb.GetUserByCompanyIdRequest) (*pb.GetUserByCompanyIdResponse, error) {
	return u.userRepo.GetUserCompanyId(req.CompanyId, req.Role)
}
