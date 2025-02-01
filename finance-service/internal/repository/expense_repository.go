package repository

import (
	"context"
	"database/sql"
	"finance-service/internal/clients"
	"finance-service/internal/utils"
	"finance-service/proto/pb"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ExpenseRepository struct {
	db         *sql.DB
	userClient *clients.UserClient
}

func (r *ExpenseRepository) CreateExpense(ctx context.Context, companyId, title, givenDate, expenseType, categoryId, userId, sum, createdBy, paymentMethod string) error {
	query := `
		INSERT INTO expense (
			id, title, user_id, category_id, expense_type, sum, created_at , given_date, created_by, payment_method,company_id
		) VALUES (
		 $1, $2, $3, $4, $5, $6, NOW(), $7, $8 , $9, $10
		)
	`
	var err error
	if expenseType == "USER" {
		_, err = r.db.Exec(query, uuid.New(), title, userId, nil, expenseType, sum, givenDate, createdBy, paymentMethod, companyId)
	} else {
		_, err = r.db.Exec(query, uuid.New(), title, nil, categoryId, expenseType, sum, givenDate, createdBy, paymentMethod, companyId)
	}

	if err != nil {
		return status.Errorf(codes.Aborted, "error while inserting expense %v", err)
	}
	return nil
}
func (r *ExpenseRepository) DeleteExpense(companyId, id string) error {
	_, err := r.db.Exec(`DELETE FROM expense where id=$1 and company_id=$2`, id, companyId)
	if err != nil {
		return status.Errorf(codes.Aborted, "error while deleting expense %v", err)
	}
	return nil
}
func (r *ExpenseRepository) GetExpenseDiagram(ctx context.Context, companyId, to, from string) (*pb.GetAllExpenseDiagramResponse, error) {
	query := `
		SELECT 
			CASE 
				WHEN expense_type = 'USER' THEN user_id::text
				WHEN expense_type = 'CATEGORY' THEN c.name
			END AS userOrCategories,
			SUM(e.sum) AS userOrCategoriesAmount,
			TO_CHAR(e.given_date, 'YYYY-MM') AS month,
			SUM(SUM(e.sum)) OVER () AS amountCommonExpense,
			expense_type,
			user_id
		FROM expense e 
		LEFT JOIN category c ON e.category_id = c.id
		WHERE e.given_date BETWEEN $1 AND $2 and e.company_id=$3
		GROUP BY userOrCategories, month, expense_type, user_id
		ORDER BY month;
		`

	rows, err := r.db.Query(query, from, to, companyId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var response pb.GetAllExpenseDiagramResponse
	userOrCategoriesMap := make(map[string]float64)
	monthAmountMap := make(map[string]float64)
	var commonExpense float64
	ctx, cancelFunc := utils.NewTimoutContext(ctx, companyId)
	defer cancelFunc()
	for rows.Next() {
		var (
			userOrCategory       string
			userOrCategoryAmount float64
			month                string
			amountCommonExpense  float64
			expenseType          string
			userID               sql.NullString
		)
		if err := rows.Scan(&userOrCategory, &userOrCategoryAmount, &month, &amountCommonExpense, &expenseType, &userID); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		if expenseType == "USER" && userID.Valid {
			userResp, err := r.userClient.GetUserById(ctx, userID.String)
			if err != nil {
				return nil, fmt.Errorf("error fetching user name: %v", err)
			}
			userOrCategory = userResp.Name
		}

		userOrCategoriesMap[userOrCategory] += userOrCategoryAmount
		monthAmountMap[month] += userOrCategoryAmount
		commonExpense = amountCommonExpense
	}

	for userOrCategory, amount := range userOrCategoriesMap {
		response.UserOrCategories = append(response.UserOrCategories, userOrCategory)
		response.UserOrCategoriesAmount = append(response.UserOrCategoriesAmount, fmt.Sprintf("%.2f", amount))
	}

	for month, amount := range monthAmountMap {
		response.Months = append(response.Months, month)
		response.MonthAmount = append(response.MonthAmount, fmt.Sprintf("%.2f", amount))
	}

	response.AmountCommonExpense = fmt.Sprintf("%.2f", commonExpense)

	return &response, nil
}
func (r *ExpenseRepository) GetAllExpense(ctx context.Context, companyId string, page, size int32, from, to, idType string, id interface{}) (*pb.GetAllExpenseResponse, error) {
	offset := (page - 1) * size

	// Validate input parameters
	if page < 1 || size < 1 {
		return nil, fmt.Errorf("invalid pagination parameters: page=%d, size=%d", page, size)
	}

	baseQuery := `
        SELECT 
            e.id,
            e.title,
            e.given_date,
            COALESCE(c.id, 0) as category_id,
            COALESCE(c.name, '') as category_name,
            COALESCE(c.description, '') as category_description,
            COALESCE(CAST(e.user_id AS TEXT), '') as user_id,
            e.expense_type,
            e.sum,
            e.created_by,
            e.payment_method,
            e.created_at
        FROM expense e
        LEFT JOIN category c ON e.category_id = c.id
        WHERE e.given_date BETWEEN $1 AND $2 AND e.company_id = $3`

	args := []interface{}{from, to, companyId}
	paramCount := 3

	// Add filters based on idType
	switch idType {
	case "USER", "CATEGORY":
		paramCount++
		fieldName := "e.user_id"
		if idType == "CATEGORY" {
			fieldName = "e.category_id"
		}
		baseQuery += fmt.Sprintf(" AND %s = $%d", fieldName, paramCount)
		args = append(args, id)
	case "":
		// No filter needed
	default:
		return nil, fmt.Errorf("invalid idType: %s", idType)
	}

	// Add pagination
	baseQuery += fmt.Sprintf(" ORDER BY e.created_at DESC LIMIT $%d OFFSET $%d", paramCount+1, paramCount+2)
	args = append(args, size, offset)

	// Create a timeout context
	ctx, cancel := utils.NewTimoutContext(ctx, companyId)
	defer cancel()

	// Execute query
	expenses, err := r.fetchExpenses(ctx, baseQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch expenses: %w", err)
	}

	// Get total count
	totalCount, err := r.getTotalCount(ctx, from, to, companyId, idType, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	totalPageCount := (totalCount + size - 1) / size
	return &pb.GetAllExpenseResponse{
		TotalPageCount: totalPageCount,
		Expenses:       expenses,
	}, nil
}

func (r *ExpenseRepository) fetchExpenses(ctx context.Context, query string, args ...interface{}) ([]*pb.GetAllExpenseAbs, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error querying expenses: %w", err)
	}
	defer rows.Close()

	var expenses []*pb.GetAllExpenseAbs
	for rows.Next() {
		expense, err := r.scanExpense(ctx, rows)
		if err != nil {
			return nil, err
		}
		expenses = append(expenses, expense)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating expense rows: %w", err)
	}

	return expenses, nil
}

func (r *ExpenseRepository) scanExpense(ctx context.Context, rows *sql.Rows) (*pb.GetAllExpenseAbs, error) {
	var expense pb.GetAllExpenseAbs
	var (
		categoryID                            int
		categoryName, categoryDesc            string
		userID                                sql.NullString
		sum                                   float64
		paymentMethod, expenseType, createdBy string
	)

	err := rows.Scan(
		&expense.Id,
		&expense.Title,
		&expense.GivenDate,
		&categoryID,
		&categoryName,
		&categoryDesc,
		&userID,
		&expenseType,
		&sum,
		&createdBy,
		&paymentMethod,
		&expense.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("error scanning expense row: %w", err)
	}

	expense.ExpenseType = expenseType
	expense.Sum = fmt.Sprintf("%.2f", sum)
	expense.PaymentType = paymentMethod

	if categoryID > 0 {
		expense.Category = &pb.AbsCategory{
			Id:   fmt.Sprintf("%d", categoryID),
			Name: categoryName,
			Desc: categoryDesc,
		}
	}

	// Fetch user details if needed
	if err := r.enrichUserDetails(ctx, &expense, userID, createdBy); err != nil {
		return nil, err
	}

	return &expense, nil
}

func (r *ExpenseRepository) enrichUserDetails(ctx context.Context, expense *pb.GetAllExpenseAbs, userID sql.NullString, createdBy string) error {
	if userID.Valid && userID.String != "" {
		userResp, err := r.userClient.GetUserById(ctx, userID.String)
		if err != nil {
			return fmt.Errorf("error fetching user details: %w", err)
		}
		expense.User = userResp
	}

	creatorResp, err := r.userClient.GetUserById(ctx, createdBy)
	if err != nil {
		return fmt.Errorf("error fetching creator details: %w", err)
	}
	expense.Creator = creatorResp

	return nil
}

func (r *ExpenseRepository) getTotalCount(ctx context.Context, from, to, companyId, idType string, id interface{}) (int32, error) {
	countQuery := `
        SELECT COUNT(*)
        FROM expense e
        LEFT JOIN category c ON e.category_id = c.id
        WHERE e.given_date BETWEEN $1 AND $2 AND e.company_id = $3`

	args := []interface{}{from, to, companyId}

	if idType == "USER" || idType == "CATEGORY" {
		fieldName := "e.user_id"
		if idType == "CATEGORY" {
			fieldName = "e.category_id"
		}
		countQuery += fmt.Sprintf(" AND %s = $4", fieldName)
		args = append(args, id)
	}

	var totalCount int32
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		return 0, fmt.Errorf("error counting expenses: %w", err)
	}

	return totalCount, nil
}
func NewExpenseRepository(db *sql.DB, userClient *clients.UserClient) *ExpenseRepository {
	return &ExpenseRepository{db: db, userClient: userClient}
}
