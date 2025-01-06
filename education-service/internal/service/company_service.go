package service

import (
	"context"
	"education-service/internal/repository"
	"education-service/proto/pb"
)

type CompanyService struct {
	pb.UnimplementedCompanyServiceServer
	companyRepo *repository.CompanyRepository
}

func NewCompanyService(repo *repository.CompanyRepository) *CompanyService {
	return &CompanyService{
		companyRepo: repo,
	}
}

func (cs *CompanyService) GetCompanyBySubdomain(ctx context.Context, req *pb.GetCompanyRequest) (*pb.GetCompanyResponse, error) {
	return cs.companyRepo.GetCompanyByDomain(req.Domain)
}
