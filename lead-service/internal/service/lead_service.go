package service

import (
	"lid-service/internal/repository"
	"lid-service/proto/pb"
)

type LeadService struct {
	pb.UnimplementedLeadServiceServer
}

func NewLeadService(repo *repository.LeadRepository) *LeadService {
	return &LeadService{}
}
