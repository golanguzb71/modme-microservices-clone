package service

import (
	"lid-service/internal/repository"
	"lid-service/proto/pb"
)

type LeadDataService struct {
	pb.UnimplementedLeadDataServiceServer
}

func NewLeadDataService(repo *repository.LeadDataRepository) *LeadDataService {
	return &LeadDataService{}
}
