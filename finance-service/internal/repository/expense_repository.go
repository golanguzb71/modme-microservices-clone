package repository

import (
	"database/sql"
	"finance-service/proto/pb"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type ExpenseRepository struct {
	db *sql.DB
}

func (r *ExpenseRepository) CreateExpense(title, givenDate, expenseType, categoryId, userId, sum, createdBy, paymentMethod string) error {
	query := `
		INSERT INTO expense (
			id, title, user_id, category_id, expense_type, sum, given_date, created_at, created_by, payment_method
		) VALUES (
		 $1, $2, $3, $4, $5, $6, NOW(), $7, $8 , $9
		)
	`
	_, err := r.db.Exec(query, uuid.New(), title, userId, categoryId, expenseType, sum, givenDate, createdBy, paymentMethod)
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
func (r *ExpenseRepository) GetAllExpense(page, size int32, from, to, category string) (*pb.GetAllExpenseResponse, error) {
	offset := (page - 1) * size
	query := `
        SELECT 
            e.id,
            e.given_date,
            c.name as category_name,
            e.user_id,
            e.expense_type,
            e.sum,
            e.created_by,
            COUNT(*) OVER() as total_count
        FROM expense e
        LEFT JOIN category c ON e.category_id = c.id
        WHERE 1=1
    `
	args := make([]interface{}, 0)
	argPosition := 1

	if from != "" {
		query += fmt.Sprintf(" AND e.given_date >= $%d", argPosition)
		args = append(args, from)
		argPosition++
	}
	if to != "" {
		query += fmt.Sprintf(" AND e.given_date <= $%d", argPosition)
		args = append(args, to)
		argPosition++
	}

	if category != "" {
		query += fmt.Sprintf(" AND c.name = $%d", argPosition)
		args = append(args, category)
		argPosition++
	}

	query += fmt.Sprintf(" ORDER BY e.given_date DESC LIMIT $%d OFFSET $%d",
		argPosition, argPosition+1)
	args = append(args, size, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	expenses := make([]*pb.GetAllExpenseAbs, 0)
	var totalCount int

	for rows.Next() {
		var expense pb.GetAllExpenseAbs
		var givenDate time.Time
		var sum float64

		err := rows.Scan(
			&expense.Id,
			&givenDate,
			&expense.CategoryName,
			&expense.UserId,
			&expense.ExpenseType,
			&sum,
			&expense.CreatedById,
			&totalCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}

		expense.GivenDate = givenDate.Format("2006-01-02")
		expense.Sum = fmt.Sprintf("%.2f", sum)

		expenses = append(expenses, &expense)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	totalPages := (int32(totalCount) + size - 1) / size

	return &pb.GetAllExpenseResponse{
		TotalPageCount: totalPages,
		Expenses:       expenses,
	}, nil
}
func (r *ExpenseRepository) GetExpenseDiagram(from, to string) (*pb.GetAllExpenseDiagramResponse, error) {
	return nil, nil
}

func NewExpenseRepository(db *sql.DB) *ExpenseRepository {
	return &ExpenseRepository{db: db}
}
