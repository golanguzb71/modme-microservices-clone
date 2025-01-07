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

func (cs *CompanyService) CreateCompany(ctx context.Context, req *pb.CreateCompanyRequest) (*pb.AbsResponse, error) {
	return cs.companyRepo.CreateCompany(req)
}

func (cs *CompanyService) GetAll(ctx context.Context, req *pb.PageRequest) (*pb.GetAllResponse, error) {
	return cs.companyRepo.GetAll(req.Page, req.Size)
}
func (cs *CompanyService) UpdateCompany(ctx context.Context, req *pb.UpdateCompanyRequest) (*pb.AbsResponse, error) {
	return cs.companyRepo.UpdateCompany(req)
}
