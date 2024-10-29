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

func (as *AuthService) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.GetUserByIdResponse, error) {
	claims, err := security.ValidateToken(req.Token)
	if err != nil {
		return nil, err
	}
	user, _, err := as.userRepo.GetUserByPhoneNumber(claims.Username)
	if err != nil {
		return nil, err
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
