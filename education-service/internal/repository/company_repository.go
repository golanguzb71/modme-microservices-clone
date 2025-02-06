package repository

import (
	"context"
	"database/sql"
	"education-service/internal/clients"
	"education-service/proto/pb"
	"errors"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

type CompanyRepository struct {
	db         *sql.DB
	userClient *clients.UserClient
}

func NewCompanyRepository(db *sql.DB, uc *clients.UserClient) *CompanyRepository {
	return &CompanyRepository{db: db, userClient: uc}
}

func (r *CompanyRepository) GetCompanyByDomain(domain, companyId string) (*pb.GetCompanyResponse, error) {
	var query string
	var row *sql.Row
	if domain != "" {
		query = `
		SELECT 
			c.id, c.title, c.avatar, c.start_time, c.end_time, 
			c.company_phone, c.subdomain, c.valid_date, 
			t.id AS tariff_id, t.name AS tariff_name, t.sum AS tariff_price, t.discounts,
			coalesce(c.discount_id , '0'), c.is_demo, c.created_at , (SELECT count(*) FROM students where condition = 'ACTIVE' and company_id=c.id) as studentcount
		FROM 
			company c
		LEFT JOIN 
			tariff t ON c.tariff_id = t.id
		WHERE 
			c.subdomain = $1
	`
		row = r.db.QueryRow(query, domain)
	} else {
		query = `
		SELECT 
			c.id, c.title, c.avatar, c.start_time, c.end_time, 
			c.company_phone, c.subdomain, c.valid_date, 
			t.id AS tariff_id, t.name AS tariff_name, t.sum AS tariff_price, t.discounts,
			coalesce(c.discount_id , '0'), c.is_demo, c.created_at , (SELECT count(*) FROM students where condition = 'ACTIVE' and company_id=c.id) as studentcount
		FROM 
			company c
		LEFT JOIN 
			tariff t ON c.tariff_id = t.id
		WHERE 
			c.id = $1
	`
		row = r.db.QueryRow(query, companyId)
	}

	var company pb.GetCompanyResponse
	var tariff pb.Tariff
	err := row.Scan(
		&company.Id, &company.Title, &company.AvatarUrl,
		&company.StartTime, &company.EndTime, &company.CompanyPhone,
		&company.Subdomain, &company.ValidDate,
		&tariff.Id, &tariff.Name, &tariff.Sum, &tariff.Discounts,
		&company.DiscountId, &company.IsDemo, &company.CreatedAt, &company.ActiveStudentCount,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("company with domain %s not found", domain)
		}
	}
	id, err := r.userClient.GetUserByCompanyId(context.Background(), company.Id, "CEO")
	if err != nil {
		return nil, status.Error(codes.Aborted, err.Error())
	} else {
		company.CeoId = id.UserId
	}
	company.Tariff = &tariff
	return &company, nil
}

func (r *CompanyRepository) CreateCompany(req *pb.CreateCompanyRequest) (*pb.AbsResponse, error) {
	var exists bool
	if err := r.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM company where subdomain=$1)`, req.Subdomain).Scan(&exists); err != nil || exists {
		return nil, status.Error(codes.Aborted, "this subdomain already have got in database")
	}
	_, err := r.db.Exec(`INSERT INTO company(title, avatar, start_time, end_time, company_phone, subdomain, valid_date, tariff_id, discount_id, is_demo) VALUES ($1,$2, $3, $4, $5, $6 , $7,  $8 , $9 , $10)`,
		req.Title,
		req.AvatarUrl,
		req.StartTime,
		req.EndTime,
		req.CompanyPhone,
		req.Subdomain,
		req.ValidDate,
		req.TariffId,
		req.DiscountId,
		req.IsDemo,
	)
	if err != nil {
		return nil, err
	}
	return &pb.AbsResponse{
		Status:  http.StatusOK,
		Message: "company create",
	}, nil
}

func (r *CompanyRepository) GetAll(page int32, size int32, filter string) (*pb.GetAllResponse, error) {
	offset := (page - 1) * size
	var filterCondition string
	switch filter {
	case "demo":
		filterCondition = "WHERE c.is_demo = true"
	case "active":
		filterCondition = "WHERE c.valid_date > NOW() and c.is_demo = false"
	case "no_active":
		filterCondition = "WHERE c.valid_date <= NOW() "
	default:
		filterCondition = ""
	}

	query := fmt.Sprintf(`
		SELECT 
			c.id, c.title, c.avatar, c.start_time, c.end_time, 
			c.company_phone, c.subdomain, c.valid_date, 
			t.id AS tariff_id, t.name AS tariff_name, t.sum AS tariff_price,  t.discounts,
			COALESCE(c.discount_id, '0'), c.is_demo, c.created_at, 
			(SELECT COUNT(*) FROM students WHERE condition = 'ACTIVE' AND company_id = c.id) AS student_count
		FROM 
			company c
		LEFT JOIN 
			tariff t ON c.tariff_id = t.id
		%s
		ORDER BY 
			c.id
		LIMIT $1 OFFSET $2
	`, filterCondition)

	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM company c %s`, filterCondition)

	var totalCount int32
	err := r.db.QueryRow(countQuery).Scan(&totalCount)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %v", err)
	}

	rows, err := r.db.Query(query, size, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch companies: %v", err)
	}
	defer rows.Close()

	var items []*pb.GetCompanyResponse
	for rows.Next() {
		var company pb.GetCompanyResponse
		var tariff pb.Tariff

		err := rows.Scan(
			&company.Id, &company.Title, &company.AvatarUrl,
			&company.StartTime, &company.EndTime, &company.CompanyPhone,
			&company.Subdomain, &company.ValidDate,
			&tariff.Id, &tariff.Name, &tariff.Sum, &tariff.Discounts,
			&company.DiscountId, &company.IsDemo, &company.CreatedAt, &company.ActiveStudentCount,
		)
		if err != nil {
			continue
		}

		company.Tariff = &tariff
		items = append(items, &company)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return &pb.GetAllResponse{
		Items:      items,
		TotalCount: totalCount,
	}, nil
}

func (r *CompanyRepository) UpdateCompany(req *pb.UpdateCompanyRequest) (*pb.AbsResponse, error) {
	var exists bool
	if err := r.db.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM company WHERE subdomain=$1 AND id<>$2)`,
		req.Subdomain, req.Id,
	).Scan(&exists); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check subdomain existence: %v", err)
	}
	if exists {
		return nil, status.Error(codes.Aborted, "this subdomain already exists in the database")
	}

	_, err := r.db.Exec(
		`UPDATE company SET title=$1, avatar=$2, start_time=$3, end_time=$4, company_phone=$5, subdomain=$6, valid_date=$7, tariff_id=$8, discount_id=$9, is_demo=$10 WHERE id=$11`,
		req.Title, req.AvatarUrl, req.StartTime, req.EndTime, req.CompanyPhone, req.Subdomain,
		req.ValidDate, req.TariffId, req.DiscountId, req.IsDemo, req.Id,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update company: %v", err)
	}

	return &pb.AbsResponse{
		Status:  http.StatusOK,
		Message: "Company updated successfully",
	}, nil
}
