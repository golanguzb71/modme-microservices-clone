package service

import (
	"lid-service/internal/repository"
	"lid-service/proto/pb"
)

type SetService struct {
	pb.UnimplementedSetServiceServer
}

func NewSetService(repo *repository.SetRepository) *SetService {
	return &SetService{}
}
