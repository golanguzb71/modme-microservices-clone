package repository

import (
	"database/sql"
	"education-service/proto/pb"
)

type CompanyFinanceRepository struct {
	db *sql.DB
}

func NewCompanyFinanceRepository(db *sql.DB) *CompanyFinanceRepository {
	return &CompanyFinanceRepository{db: db}
}

func (r CompanyFinanceRepository) Create(req *pb.CompanyFinance) (*pb.CompanyFinance, error) {
	r.db.Query(`INSERT INTO company_payments(company_id, tariff_id, comment, sum, edited_valid_date) values ($1 ,$2,$3,$4,$5)`, req.CompanyId, req.TariffId)
	return nil, nil
}

func (r CompanyFinanceRepository) Delete(req *pb.DeleteAbsRequest) (*pb.AbsResponse, error) {
	return nil, nil
}

func (r CompanyFinanceRepository) GetAll(req *pb.PageRequest) (*pb.CompanyFinanceList, error) {
	return nil, nil
}

func (r CompanyFinanceRepository) GetByCompany(req *pb.PageRequest) (*pb.CompanyFinanceSelf, error) {
	return nil, nil
}
