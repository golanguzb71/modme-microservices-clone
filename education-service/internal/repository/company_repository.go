package repository

import (
	"database/sql"
	"education-service/proto/pb"
	"errors"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
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

func (r *CompanyRepository) CreateCompany(req *pb.CreateCompanyRequest) (*pb.AbsResponse, error) {
	var exists bool
	if err := r.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM company where subdomain=$1)`, req.Subdomain).Scan(&exists); err != nil {
		return nil, status.Error(codes.Aborted, "this subdomain already have got in database")
	}
	_, err := r.db.Exec(`INSERT INTO company(title, avatar, start_time, end_time, company_phone, subdomain) VALUES ($1,$2, $3, $4, $5, $6)`, req.Title, req.AvatarUrl, req.StartTime, req.EndTime, req.CompanyPhone, req.Subdomain)
	if err != nil {
		return nil, err
	}
	return &pb.AbsResponse{
		Status:  http.StatusOK,
		Message: "company create",
	}, nil
}

func (r *CompanyRepository) GetAll(page int32, size int32) (*pb.GetAllResponse, error) {
	return nil, nil
}

func (r *CompanyRepository) UpdateCompany(req *pb.UpdateCompanyRequest) (*pb.AbsResponse, error) {
	return nil, nil
}
