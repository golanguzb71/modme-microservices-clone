package repository

import (
	"context"
	"database/sql"
	"education-service/proto/pb"
	"encoding/json"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

type TariffRepository struct {
	db *sql.DB
}

func NewTariffRepository(db *sql.DB) *TariffRepository {
	return &TariffRepository{db: db}
}

func (r TariffRepository) Get(ctx context.Context, req *emptypb.Empty) (*pb.TariffList, error) {
	query := `
        SELECT id, name, student_count, sum, discounts, created_at 
        FROM tariff 
        WHERE is_deleted = false
        ORDER BY id`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tariffs := make([]*pb.Tariff, 0)
	for rows.Next() {
		var tariff pb.Tariff
		var discountsJSON []byte
		var createdAt time.Time

		err := rows.Scan(
			&tariff.Id,
			&tariff.Name,
			&tariff.StudentCount,
			&tariff.Sum,
			&discountsJSON,
			&createdAt,
		)
		if err != nil {
			return nil, err
		}

		if discountsJSON != nil {
			if err := json.Unmarshal(discountsJSON, &tariff.Discounts); err != nil {
				return nil, err
			}
		}

		tariff.CreatedAt = createdAt.Format(time.RFC3339)
		tariffs = append(tariffs, &tariff)
	}

	return &pb.TariffList{Items: tariffs}, rows.Err()
}

func (r TariffRepository) Create(ctx context.Context, req *pb.Tariff) (*pb.Tariff, error) {
	query := `
        INSERT INTO tariff (name, student_count, sum, discounts)
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at`

	discountsJSON, err := json.Marshal(req.Discounts)
	if err != nil {
		return nil, err
	}

	var createdAt time.Time
	err = r.db.QueryRowContext(ctx, query,
		req.Name,
		req.StudentCount,
		req.Sum,
		discountsJSON,
	).Scan(&req.Id, &createdAt)

	if err != nil {
		return nil, err
	}

	req.CreatedAt = createdAt.Format(time.RFC3339)
	return req, nil
}

func (r TariffRepository) Update(ctx context.Context, req *pb.Tariff) (*pb.Tariff, error) {
	query := `
        UPDATE tariff
        SET name = $1,
            student_count = $2,
            sum = $3,
            discounts = $4
        WHERE id = $5 AND is_deleted = false
        RETURNING created_at`

	discountsJSON, err := json.Marshal(req.Discounts)
	if err != nil {
		return nil, err
	}

	var createdAt time.Time
	err = r.db.QueryRowContext(ctx, query,
		req.Name,
		req.StudentCount,
		req.Sum,
		discountsJSON,
		req.Id,
	).Scan(&createdAt)

	if err != nil {
		return nil, err
	}

	req.CreatedAt = createdAt.Format(time.RFC3339)
	return req, nil
}

func (r TariffRepository) Delete(ctx context.Context, req *pb.Tariff) (*pb.Tariff, error) {
	query := `
        UPDATE tariff
        SET is_deleted = true
        WHERE id = $1 AND is_deleted = false
        RETURNING created_at`

	var createdAt time.Time
	err := r.db.QueryRowContext(ctx, query, req.Id).Scan(&createdAt)
	if err != nil {
		return nil, err
	}

	req.CreatedAt = createdAt.Format(time.RFC3339)
	return req, nil
}
