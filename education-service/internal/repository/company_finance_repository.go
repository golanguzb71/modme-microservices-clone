package repository

import (
	"database/sql"
	"education-service/proto/pb"
	"fmt"
)

type CompanyFinanceRepository struct {
	db *sql.DB
}

func NewCompanyFinanceRepository(db *sql.DB) *CompanyFinanceRepository {
	return &CompanyFinanceRepository{db: db}
}

func (r CompanyFinanceRepository) Create(req *pb.CompanyFinance) (*pb.CompanyFinance, error) {
	var validDate string
	err := r.db.QueryRow(`SELECT valid_date FROM company WHERE id = $1`, req.GetCompanyId()).Scan(&validDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("company with id %d not found", req.GetCompanyId())
		}
		return nil, err
	}

	editedValidDate := req.GetEditedValidDate()
	if editedValidDate <= validDate {
		return nil, fmt.Errorf("edited_valid_date (%s) must be greater than valid_date (%s)", editedValidDate, validDate)
	}
	tx, err := r.db.Begin()
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	_, err = tx.Exec(`INSERT INTO company_payments(company_id, tariff_id, comment, sum, edited_valid_date , discount_name , discount_id) values ($1 ,$2,$3,$4,$5 , $6,$7)`, req.CompanyId, req.TariffId, req.Comment, req.Sum, req.EditedValidDate, req.DiscountName, req.DiscountId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	_, err = tx.Exec(`UPDATE company
SET 
    valid_date = $1,
    is_demo = CASE WHEN is_demo = true THEN false ELSE is_demo END
WHERE id = $2;
`, req.EditedValidDate, req.CompanyId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (r CompanyFinanceRepository) Delete(req *pb.DeleteAbsRequest) (*pb.AbsResponse, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}

	var exists bool
	err = tx.QueryRow(`
		SELECT EXISTS(
			SELECT 1 
			FROM company 
			WHERE valid_date = (SELECT edited_valid_date FROM company_payments WHERE id = $1)
		)
	`, req.Id).Scan(&exists)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if exists {
		_, err = tx.Exec(`
			UPDATE company
			SET valid_date = (
				SELECT edited_valid_date
				FROM (
					SELECT edited_valid_date
					FROM company_payments
					WHERE company_id = (
						SELECT company_id
						FROM company_payments
						WHERE id = $1
					)
					AND id != $1
					ORDER BY created_at DESC
					LIMIT 2
				) subquery
				ORDER BY created_at ASC
				LIMIT 1
			)
			WHERE id = (
				SELECT company_id
				FROM company_payments
				WHERE id = $1
			);
		`, req.Id)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	_, err = tx.Exec(`DELETE FROM company_payments WHERE id = $1`, req.Id)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &pb.AbsResponse{
		Status:  200,
		Message: "ok",
	}, nil
}
func (r CompanyFinanceRepository) GetAll(req *pb.PageRequest) (*pb.CompanyFinanceList, error) {
	page := req.GetPage()
	if page <= 0 {
		page = 1
	}

	size := req.GetSize()
	if size <= 0 {
		size = 2
	}

	offset := (page - 1) * size

	query := `
		SELECT
			cp.id,
			c.title AS company_name,
			cp.company_id,
			cp.created_at AS start_from,
			cp.edited_valid_date AS finished_to,
			t.id AS tariff_id,
			t.name AS tariff_name,
			cp.sum,
			coalesce(cp.discount_id , ''),
			coalesce(cp.discount_name , '')
		FROM
			company_payments cp
		LEFT JOIN
			company c ON c.id = cp.company_id
		LEFT JOIN
			tariff t ON t.id = cp.tariff_id
		WHERE
			($1::TIMESTAMP IS NULL OR cp.created_at >= $1)
			AND ($2::TIMESTAMP IS NULL OR cp.created_at <= $2)
		ORDER BY
			cp.created_at DESC
		LIMIT $3 OFFSET $4;
	`

	countQuery := `
		SELECT COUNT(*)
		FROM
			company_payments
	`

	from := req.GetFrom()
	to := req.GetTo()

	// Get the total record count
	var totalCount int32
	err := r.db.QueryRow(countQuery).Scan(&totalCount)
	if err != nil {
		return nil, err
	}

	// Calculate the total page count
	totalPageCount := (totalCount + int32(size) - 1) / int32(size) // Ceiling division

	rows, err := r.db.Query(query, from, to, size, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*pb.CompanyFinanceForList
	for rows.Next() {
		var item pb.CompanyFinanceForList
		err := rows.Scan(
			&item.Id,
			&item.CompanyName,
			&item.CompanyId,
			&item.StartFrom,
			&item.FinishedTo,
			&item.TariffId,
			&item.TariffName,
			&item.Sum,
			&item.DiscountId,
			&item.DiscountName,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &pb.CompanyFinanceList{
		Count: totalPageCount, // Total page count
		Items: items,
	}, nil
}

func (r CompanyFinanceRepository) GetByCompany(req *pb.PageRequest) (*pb.CompanyFinanceSelfList, error) {
	page := req.GetPage()
	if page <= 0 {
		page = 1
	}

	size := req.GetSize()
	if size <= 0 {
		size = 2
	}

	offset := (page - 1) * size

	from := req.GetFrom()
	to := req.GetTo()
	companyId := req.GetCompanyId()

	query := `
		SELECT 
			cp.id, 
			cp.tariff_id, 
			t.sum AS tariff_sum, 
			coalesce(cp.comment, ''), 
			cp.sum, 
			cp.edited_valid_date, 
			cp.created_at, 
			coalesce(cp.discount_id, ''), 
			coalesce(cp.discount_name, '')
		FROM 
			company_payments cp
		LEFT JOIN 
			tariff t ON t.id = cp.tariff_id
		WHERE 
			cp.company_id = $1 
			AND ($2::TIMESTAMP IS NULL OR cp.created_at >= $2::TIMESTAMP)
			AND ($3::TIMESTAMP IS NULL OR cp.created_at <= $3::TIMESTAMP)
		ORDER BY 
			cp.created_at DESC
		LIMIT $4 OFFSET $5;
	`

	countQuery := `
		SELECT COUNT(*)
		FROM 
			company_payments cp
		WHERE 
			cp.company_id = $1 
	`

	sumQuery := `
		SELECT 
			coalesce(SUM(cp.sum), 0) AS sum_amount_period, 
			coalesce(t.name, '') AS tariff_name, 
			coalesce(cp.discount_name, '') AS discount_name,
			t.sum AS required_sum
		FROM 
			company_payments cp
		LEFT JOIN 
			tariff t ON t.id = cp.tariff_id
		WHERE 
			cp.company_id = $1 
			AND ($2::TIMESTAMP IS NULL OR cp.created_at >= $2::TIMESTAMP)
			AND ($3::TIMESTAMP IS NULL OR cp.created_at <= $3::TIMESTAMP)
		GROUP BY 
			t.name, cp.discount_name, t.sum
		LIMIT 1;
	`

	var totalCount int32
	err := r.db.QueryRow(countQuery, companyId).Scan(&totalCount)
	if err != nil {
		return nil, err
	}

	totalPageCount := (totalCount + int32(size) - 1) / int32(size) // Ceiling division

	var sumAmountPeriod float32
	var tariffName string
	var discountName string
	var requiredSum float32
	err = r.db.QueryRow(sumQuery, companyId, from, to).Scan(&sumAmountPeriod, &tariffName, &discountName, &requiredSum)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(query, companyId, from, to, size, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*pb.CompanyFinanceSelf
	for rows.Next() {
		var item pb.CompanyFinanceSelf
		err := rows.Scan(
			&item.Id,
			&item.TariffId,
			&item.TariffSum,
			&item.Comment,
			&item.Sum,
			&item.EditValidDate,
			&item.CreatedAt,
			&item.DiscountId,
			&item.DiscountName,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &pb.CompanyFinanceSelfList{
		Count:           totalPageCount,
		SumAmountPeriod: sumAmountPeriod,
		TariffName:      tariffName,
		DiscountName:    discountName,
		RequiredSum:     requiredSum,
		Items:           items,
	}, nil
}
