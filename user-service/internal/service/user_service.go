package service

import (
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
