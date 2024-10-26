package service

import (
	"context"
	"errors"
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
	user, password, err := as.userRepo.GetUserByPhoneNumber(request.PhoneNumber)
	if err != nil || user.IsDeleted {
		return nil, errors.New("notog'ri login yoki parol")
	}
	err = utils.ComparePasswords(password, request.Password)
	if err != nil {
		return nil, errors.New("notog'ri login yoki parol")
	}
	token, err := security.GenerateToken(user.PhoneNumber)
	if err != nil {
		return nil, err
	}
	return &pb.LoginResponse{
		User:  user,
		Token: token,
		IsOk:  true,
	}, nil
}
