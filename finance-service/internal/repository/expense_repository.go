package repository

import (
	"context"
	"database/sql"
	"finance-service/proto/pb"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ExpenseRepository struct {
	db *sql.DB
}

func (r *ExpenseRepository) CreateExpense(title, givenDate, expenseType, categoryId, userId, sum, createdBy, paymentMethod string) error {
	query := `
		INSERT INTO expense (
			id, title, user_id, category_id, expense_type, sum, created_at , given_date, created_by, payment_method
		) VALUES (
		 $1, $2, $3, $4, $5, $6, NOW(), $7, $8 , $9
		)
	`
	var err error
	if expenseType == "USER" {
		_, err = r.db.Exec(query, uuid.New(), title, userId, nil, expenseType, sum, givenDate, createdBy, paymentMethod)
	} else {
		_, err = r.db.Exec(query, uuid.New(), title, nil, categoryId, expenseType, sum, givenDate, createdBy, paymentMethod)
	}

	if err != nil {
		return status.Errorf(codes.Aborted, "error while inserting expense %v", err)
	}
	return nil
}
func (r *ExpenseRepository) DeleteExpense(id string) error {
	_, err := r.db.Exec(`DELETE FROM expense where id=$1`, id)
	if err != nil {
		return status.Errorf(codes.Aborted, "error while deleting expense %v", err)
	}
	return nil
}
func (r *ExpenseRepository) GetExpenseDiagram(from, to string) (*pb.GetAllExpenseDiagramResponse, error) {
	return nil, nil
}
func (r *ExpenseRepository) GetAllExpense(page int32, size int32, from string, to string, idType string, id string) (*pb.GetAllExpenseResponse, error) {
	offset := (page - 1) * size
	baseQuery := `SELECT expense.id, given_date, category.name as category_name, user_id, expense_type, sum, created_by 
                  FROM expense
                  LEFT JOIN category ON expense.category_id = category.id
                  WHERE given_date BETWEEN $1 AND $2`
	params := []interface{}{from, to}
	if idType == "USER" {
		baseQuery += ` AND user_id = $3`
		params = append(params, id)
	} else if idType == "CATEGORY" {
		baseQuery += ` AND category_id = $3`
		params = append(params, id)
	}
	baseQuery += ` ORDER BY given_date DESC LIMIT $4 OFFSET $5`
	params = append(params, size, offset)

	rows, err := r.db.QueryContext(context.Background(), baseQuery, params...)
	if err != nil {
		return nil, fmt.Errorf("error querying expenses: %v", err)
	}
	defer rows.Close()

	var expenses []*pb.GetAllExpenseAbs
	for rows.Next() {
		var expense pb.GetAllExpenseAbs
		var givenDate string
		if err := rows.Scan(&expense.Id, &givenDate, &expense.CategoryName, &expense.UserId, &expense.ExpenseType, &expense.Sum, &expense.CreatedById); err != nil {
			return nil, fmt.Errorf("error scanning expense row: %v", err)
		}
		expense.GivenDate = givenDate
		expenses = append(expenses, &expense)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating expense rows: %v", err)
	}
	var totalCount int32
	countQuery := `SELECT COUNT(*) FROM expense WHERE given_date BETWEEN $1 AND $2`
	if idType == "USER" {
		countQuery += ` AND user_id = $3`
	} else if idType == "CATEGORY" {
		countQuery += ` AND category_id = $3`
	}
	err = r.db.QueryRowContext(context.Background(), countQuery, params[:3]...).Scan(&totalCount)
	if err != nil {
		return nil, fmt.Errorf("error counting expenses: %v", err)
	}
	totalPageCount := (totalCount + size - 1) / size

	response := &pb.GetAllExpenseResponse{
		TotalPageCount: totalPageCount,
		Expenses:       expenses,
	}
	return response, nil
}
func NewExpenseRepository(db *sql.DB) *ExpenseRepository {
	return &ExpenseRepository{db: db}
}
