package repository

import "database/sql"

type CompanyFinanceRepository struct {
	db *sql.DB
}

func NewCompanyFinanceRepository(db *sql.DB) *CompanyFinanceRepository {
	return &CompanyFinanceRepository{db: db}
}
