package service

import (
	"lid-service/internal/repository"
	"lid-service/proto/pb"
)

type ExpectService struct {
	pb.UnimplementedExpectServiceServer
}

func NewExpectService(repo *repository.ExpectRepository) *ExpectService {
	return &ExpectService{}
}
