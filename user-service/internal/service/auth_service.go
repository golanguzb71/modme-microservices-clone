package service

import (
	"context"
	"user-service/internal/repository"
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

}
