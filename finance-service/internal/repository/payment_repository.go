package repository

import (
	"database/sql"
	"errors"
	"finance-service/proto/pb"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"strconv"
	"time"
)

type PaymentRepository struct {
	db *sql.DB
}

func (r *PaymentRepository) AddPayment(givenDate, sum, method, comment, studentId, actionByName, actionById, groupId string) error {
	validMethods := map[string]bool{"CLICK": true, "CASH": true, "PAYME": true}
	if !validMethods[method] {
		return errors.New("invalid payment method")
	}

	amount, err := strconv.ParseFloat(sum, 64)
	if err != nil {
		return fmt.Errorf("invalid sum amount: %v", err)
	}

	parsedDate, err := time.Parse("2006-01-02", givenDate)
	if err != nil {
		return fmt.Errorf("invalid date format: %v", err)
	}
	paymentID := uuid.New()

	query := `INSERT INTO student_payments 
		(id, student_id, method, amount, given_date, comment, payment_type, created_by_id, created_by_name , created_at , group_id)
		VALUES ($1, $2, $3, $4, $5, $6, 'ADD', $7, $8 , $9 , $9)`

	_, err = r.db.Exec(query, paymentID, studentId, method, amount, parsedDate, comment, actionById, actionByName, time.Now(), groupId)
	if err != nil {
		return fmt.Errorf("failed to add payment: %v", err)
	}

	return nil
}

func (r *PaymentRepository) TakeOffPayment(date, sum, method, comment, studentId, actionByName, actionById string) error {
	validMethods := map[string]bool{"CLICK": true, "CASH": true, "PAYME": true}
	if !validMethods[method] {
		return errors.New("invalid payment method")
	}

	amount, err := strconv.ParseFloat(sum, 64)
	if err != nil {
		return fmt.Errorf("invalid sum amount: %v", err)
	}

	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return fmt.Errorf("invalid date format: %v", err)
	}

	paymentID := uuid.New()

	query := `INSERT INTO student_payments 
		(id, student_id, method, amount, given_date, comment, payment_type, created_by_id, created_by_name)
		VALUES ($1, $2, $3, $4, $5, $6, 'TAKE_OFF', $7, $8)`

	_, err = r.db.Exec(query, paymentID, studentId, method, amount, parsedDate, comment, actionById, actionByName)
	if err != nil {
		return fmt.Errorf("failed to take off payment: %v", err)
	}

	return nil
}

func (r *PaymentRepository) PaymentReturn(paymentId, actionByName, actionById string) (*pb.AbsResponse, error) {
	var payment struct {
		ID            string
		StudentID     string
		Method        string
		Amount        float64
		GivenDate     string
		Comment       string
		PaymentType   string
		CreatedByID   string
		CreatedByName string
		CreatedAt     string
	}

	query := `SELECT id, student_id, method, amount, given_date, comment, payment_type, created_by_id, created_by_name, created_at 
			  FROM student_payments WHERE id = $1`
	err := r.db.QueryRow(query, paymentId).Scan(
		&payment.ID,
		&payment.StudentID,
		&payment.Method,
		&payment.Amount,
		&payment.GivenDate,
		&payment.Comment,
		&payment.PaymentType,
		&payment.CreatedByID,
		&payment.CreatedByName,
		&payment.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("payment not found")
	} else if err != nil {
		return nil, fmt.Errorf("failed to retrieve payment: %v", err)
	}

	deleteQuery := `DELETE FROM student_payments WHERE id = $1`
	_, err = r.db.Exec(deleteQuery, paymentId)
	if err != nil {
		return nil, fmt.Errorf("failed to delete payment: %v", err)
	}

	return &pb.AbsResponse{
		Status:  http.StatusOK,
		Message: fmt.Sprintf("payment returned successfully"),
	}, nil
}

func (r *PaymentRepository) PaymentUpdate(paymentId string, date string, method string, userId string, comment string, debit, actionByName, actionById, groupId string) (*pb.AbsResponse, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM student_payments WHERE id = $1)`
	err := r.db.QueryRow(query, paymentId).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("error checking payment existence: %v", err)
	}

	if !exists {
		return nil, errors.New("payment not found")
	}

	updateQuery := `UPDATE student_payments 
		SET given_date = $1, method = $2, comment = $3, amount = $4, created_by_id = $5, created_by_name = $6 , group_id=$8
		WHERE id = $7`

	_, err = r.db.Exec(updateQuery, date, method, comment, debit, actionById, actionByName, paymentId, groupId)
	if err != nil {
		return nil, fmt.Errorf("failed to update payment: %v", err)
	}

	return &pb.AbsResponse{
		Status:  http.StatusOK,
		Message: "payment updated successfully",
	}, nil
}

func (r *PaymentRepository) GetMonthlyStatus(studentId string) (*pb.GetMonthlyStatusResponse, error) {
	query := `
		SELECT 
			TO_CHAR(given_date, 'YYYY-MM') AS month, 
			SUM(CASE WHEN payment_type = 'ADD' THEN amount ELSE 0 END) AS total_add,
			SUM(CASE WHEN payment_type = 'TAKE_OFF' THEN amount ELSE 0 END) AS total_take_off
		FROM 
			student_payments 
		WHERE 
			student_id = $1
		GROUP BY 
			month 
		ORDER BY 
			month`

	rows, err := r.db.Query(query, studentId)
	if err != nil {
		return nil, fmt.Errorf("error querying monthly status: %v", err)
	}
	defer rows.Close()

	var monthlyStatus []*pb.AbsGetMonthlyStatusResponse

	for rows.Next() {
		var month string
		var totalAdd, totalTakeOff float64

		if err := rows.Scan(&month, &totalAdd, &totalTakeOff); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		balance := totalAdd - totalTakeOff

		monthlyStatus = append(monthlyStatus, &pb.AbsGetMonthlyStatusResponse{
			Month:   month,
			Balance: fmt.Sprintf("%.2f", balance),
		})
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %v", err)
	}

	return &pb.GetMonthlyStatusResponse{
		MonthStatus: monthlyStatus,
	}, nil
}

func (r *PaymentRepository) GetAllPaymentsByMonth(month string, studentId string) (*pb.GetAllPaymentsByMonthResponse, error) {
	query := `
		SELECT 
		    id,
			given_date, 
			payment_type, 
			amount, 
			comment, 
			created_by_id, 
			created_by_name, 
			created_at ,
			group_id
		FROM 
			student_payments 
		WHERE 
			student_id = $1 AND 
			TO_CHAR(given_date, 'YYYY-MM') = $2 
		ORDER BY 
			given_date`

	rows, err := r.db.Query(query, studentId, month)
	if err != nil {
		return nil, fmt.Errorf("error querying payments: %v", err)
	}
	defer rows.Close()

	var payments []*pb.AbsGetAllPaymentsByMonthResponse

	for rows.Next() {
		var payment pb.AbsGetAllPaymentsByMonthResponse

		var givenDate sql.NullTime

		if err := rows.Scan(
			&payment.PaymentId,
			&givenDate,
			&payment.PaymentType,
			&payment.Amount,
			&payment.Comment,
			&payment.CreatedById,
			&payment.CreatedByName,
			&payment.CreatedAt,
			&payment.GroupId,
		); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		if givenDate.Valid {
			payment.GivenDate = givenDate.Time.Format("2006-01-02")
		} else {
			payment.GivenDate = ""
		}

		payments = append(payments, &payment)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %v", err)
	}

	return &pb.GetAllPaymentsByMonthResponse{
		Payments: payments,
	}, nil
}

func NewPaymentRepository(db *sql.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}