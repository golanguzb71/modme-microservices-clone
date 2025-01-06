package repository

import (
	"database/sql"
	"education-service/proto/pb"
	"errors"
	"fmt"
)

type CompanyRepository struct {
	db *sql.DB
}

func NewCompanyRepository(db *sql.DB) *CompanyRepository {
	return &CompanyRepository{db: db}
}

func (r *CompanyRepository) GetCompanyByDomain(domain string) (*pb.GetCompanyResponse, error) {
	query := `
        SELECT id, title, avatar, start_time, end_time, company_phone, subdomain
        FROM company
        WHERE subdomain = $1
    `
	var company pb.GetCompanyResponse
	row := r.db.QueryRow(query, domain)
	err := row.Scan(
		&company.Id,
		&company.Title,
		&company.AvatarUrl,
		&company.StartTime,
		&company.EndTime,
		&company.CompanyPhone,
		&company.Subdomain,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("company with domain %s not found", domain)
		}
		return nil, err
	}
	return &company, nil
}
