package service

import (
	"lid-service/internal/repository"
	"lid-service/proto/pb"
)

type LidUserService struct {
	pb.UnimplementedLidUserServiceServer
}

func NewLidUserService(repo *repository.LidUserRepository) *LidUserService {
	return &LidUserService{}
}
