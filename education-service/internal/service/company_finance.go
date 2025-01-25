package service

import (
	"context"
	"education-service/internal/repository"
	"education-service/proto/pb"
)

type CompanyFinanceService struct {
	pb.UnimplementedCompanyFinanceServiceServer
	companyFinanceRepo *repository.CompanyFinanceRepository
}

func NewCompanyFinanceService(repo *repository.CompanyFinanceRepository) *CompanyFinanceService {
	return &CompanyFinanceService{
		companyFinanceRepo: repo,
	}
}

func (cf *CompanyFinanceService) Create(ctx context.Context, req *pb.CompanyFinance) (*pb.CompanyFinance, error) {
	return cf.companyFinanceRepo.Create(req)
}
func (cf *CompanyFinanceService) Delete(ctx context.Context, req *pb.DeleteAbsRequest) (*pb.AbsResponse, error) {
	return cf.companyFinanceRepo.Delete(req)
}
func (cf *CompanyFinanceService) GetAll(ctx context.Context, req *pb.PageRequest) (*pb.CompanyFinanceList, error) {
	return cf.companyFinanceRepo.GetAll(req)
}
func (cf *CompanyFinanceService) GetByCompany(ctx context.Context, req *pb.PageRequest) (*pb.CompanyFinanceSelfList, error) {
	return cf.companyFinanceRepo.GetByCompany(req)
}

func (cf *CompanyFinanceService) UpdateByCompany(ctx context.Context, req *pb.CompanyFinance) (*pb.CompanyFinance, error) {
	return cf.companyFinanceRepo.UpdateByCompany(req)
}
