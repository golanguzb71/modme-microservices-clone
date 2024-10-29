package repository

import (
	"database/sql"
	"finance-service/proto/pb"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CategoryRepository struct {
	db *sql.DB
}

func (r *CategoryRepository) CreateCategory(name string, desc string) error {
	_, err := r.db.Exec(`INSERT INTO category(name, description) values ($1, $2)`, name, desc)
	if err != nil {
		return status.Errorf(codes.Aborted, "error while inserting category %v", err)
	}
	return nil
}

func (r *CategoryRepository) DeleteCategory(id string) error {
	_, err := r.db.Exec(`DELETE FROM category where id=$1`, id)
	if err != nil {
		return status.Errorf(codes.Aborted, "error while deleting category %v", err)
	}
	return nil
}

func (r *CategoryRepository) GetAllCategory() (*pb.GetAllCategoryRequest, error) {
	query := "SELECT id, name, description FROM category"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query categories: %w", err)
	}
	defer rows.Close()
	var response pb.GetAllCategoryRequest
	for rows.Next() {
		var category pb.AbsCategory
		var id int32
		if err := rows.Scan(&id, &category.Name, &category.Desc); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		category.Id = fmt.Sprintf("%d", id)
		response.Categories = append(response.Categories, &category)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return &response, nil
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}
