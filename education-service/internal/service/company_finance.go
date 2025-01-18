package service

import (
	"context"
	"education-service/internal/repository"
	"education-service/proto/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (cf *CompanyFinanceService) Delete(ctx context.Context, req *pb.DeleteAbsRequest) (*pb.AbsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (cf *CompanyFinanceService) GetAll(ctx context.Context, req *pb.PageRequest) (*pb.CompanyFinanceList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAll not implemented")
}
func (cf *CompanyFinanceService) GetByCompany(ctx context.Context, req *pb.PageRequest) (*pb.CompanyFinanceSelf, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetByCompany not implemented")
}
