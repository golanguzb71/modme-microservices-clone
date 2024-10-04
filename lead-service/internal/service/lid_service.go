package service

import (
	"lid-service/internal/repository"
	"lid-service/proto/pb"
)

type LidService struct {
	pb.UnimplementedLidServiceServer
}

func NewLidService(repo *repository.LidRepository) *LidService {
	return &LidService{}
}
