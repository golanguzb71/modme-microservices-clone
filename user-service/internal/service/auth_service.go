package service

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"user-service/internal/repository"
	"user-service/internal/security"
	"user-service/internal/utils"
	"user-service/proto/pb"
)

type AuthService struct {
	pb.UnimplementedAuthServiceServer
	userRepo *repository.UserRepository
}

func NewAuthService(repo *repository.UserRepository) *AuthService {
	return &AuthService{
		userRepo: repo,
	}
}

func (as *AuthService) Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, password, err := as.userRepo.GetUserByPhoneNumber(request.CompanyId, request.PhoneNumber)
	if err != nil {
		return nil, errors.New("notog'ri login yoki parol")
	}
	if user.IsDeleted {
		return nil, status.Error(codes.Unauthenticated, "forbidden operation. deleted user request detect")
	}
	err = utils.ComparePasswords(password, request.Password)
	if err != nil {
		return nil, errors.New("notog'ri login yoki parol")
	}
	token, err := security.GenerateToken(user)
	if err != nil {
		return nil, err
	}
	return &pb.LoginResponse{
		User:  user,
		Token: token,
		IsOk:  true,
	}, nil
}

func (as *AuthService) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.GetUserByIdResponse, error) {
	claims, err := security.ValidateToken(req.Token)
	if err != nil {
		return nil, err
	}
	user, _, err := as.userRepo.GetUserByIdFilter(claims.Username)
	if err != nil {
		return nil, err
	}
	if user.IsDeleted {
		return nil, status.Error(codes.Aborted, "forbidden operation. deleted user request detect")
	}
	var checker = false
	for _, role := range req.RequiredRoles {
		if user.Role == role {
			checker = true
			break
		}
	}
	if !checker {
		return nil, errors.New("this user not valid for this endpoint")
	}
	return user, nil
}
