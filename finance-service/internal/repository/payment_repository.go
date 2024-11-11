package repository

import (
	"context"
	"database/sql"
	"errors"
	"finance-service/internal/clients"
	"finance-service/proto/pb"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"strconv"
	"time"
)

type PaymentRepository struct {
	db              *sql.DB
	educationClient *clients.EducationClient
}

func (r *PaymentRepository) AddPayment(givenDate, sum, method, comment, studentId, actionByName, actionById, groupId string, isRefund bool) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
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
		VALUES ($1, $2, $3, $4, $5, $6, $11, $7, $8 , $9 , $10)`
	paymentType := "ADD"
	if isRefund {
		paymentType = "REFUND"
	}

	if groupId == "" {
		_, err = tx.Exec(query, paymentID, studentId, method, amount, parsedDate, comment, actionById, actionByName, time.Now(), nil, paymentType)
	} else {
		_, err = tx.Exec(query, paymentID, studentId, method, amount, parsedDate, comment, actionById, actionByName, time.Now(), groupId, paymentType)
	}
	if err != nil {
		return fmt.Errorf("failed to add payment: %v", err)
	}
	err = r.educationClient.ChangeUserBalanceHistory(studentId, sum, givenDate, comment, "ADD", actionById, actionByName, groupId)
	if err != nil {
		return fmt.Errorf("failed to update user balance history: %v", err)
	}
	return nil
}

func (r *PaymentRepository) TakeOffPayment(date, sum, method, comment, studentId, actionByName, actionById, groupId string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

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
		(id, student_id, method, amount, given_date, comment, payment_type, created_by_id, created_by_name , created_at, group_id)
		VALUES ($1, $2, $3, $4, $5, $6, 'TAKE_OFF', $7, $8 , $9 , $10)`

	if groupId == "" {
		_, err = tx.Exec(query, paymentID, studentId, method, amount, parsedDate, comment, actionById, actionByName, time.Now(), nil)
	} else {
		_, err = tx.Exec(query, paymentID, studentId, method, amount, parsedDate, comment, actionById, actionByName, time.Now(), groupId)
	}
	if err != nil {
		return fmt.Errorf("failed to take off payment: %v", err)
	}

	err = r.educationClient.ChangeUserBalanceHistory(studentId, sum, date, comment, "TAKE_OFF", actionById, actionByName, groupId)
	if err != nil {
		return fmt.Errorf("failed to update user balance history: %v", err)
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
		GroupId       string
	}

	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	query := `SELECT id, student_id, method, amount, given_date, comment, payment_type, created_by_id, created_by_name, created_at , coalesce(group_id , 0)
			  FROM student_payments WHERE id = $1`
	err = tx.QueryRow(query, paymentId).Scan(
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
		&payment.GroupId,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("payment not found")
	} else if err != nil {
		return nil, fmt.Errorf("failed to retrieve payment: %v", err)
	}

	deleteQuery := `DELETE FROM student_payments WHERE id = $1`
	_, err = tx.Exec(deleteQuery, paymentId)
	if err != nil {
		return nil, fmt.Errorf("failed to delete payment: %v", err)
	}
	if payment.PaymentType == "TAKE_OFF" {
		payment.PaymentType = "ADD"
	} else {
		payment.PaymentType = "TAKE_OFF"
	}

	err = r.educationClient.ChangeUserBalanceHistory(payment.StudentID, fmt.Sprintf("%.2f", payment.Amount), payment.GivenDate, payment.Comment, payment.PaymentType, actionById, actionByName, payment.GroupId)
	if err != nil {
		return nil, fmt.Errorf("failed to update user balance history: %v", err)
	}

	return &pb.AbsResponse{
		Status:  http.StatusOK,
		Message: "payment returned successfully",
	}, nil
}

func (r *PaymentRepository) PaymentUpdate(paymentId string, date string, method string, userId string, comment string, debit, actionByName, actionById, groupId string) (*pb.AbsResponse, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	var paymentType string
	var oldDebit string
	query := `SELECT payment_type , amount FROM student_payments where id=$1`
	err = tx.QueryRow(query, paymentId).Scan(&paymentType, &oldDebit)
	if err != nil {
		return nil, fmt.Errorf("error checking payment existence: %v", err)
	}
	updateQuery := `UPDATE student_payments 
					SET given_date = $1, method = $2, comment = $3, amount = $4, created_by_id = $5, created_by_name = $6, group_id = $8 
					WHERE id = $7`
	if groupId == "" {
		_, err = tx.Exec(updateQuery, date, method, comment, debit, actionById, actionByName, paymentId, nil)
	} else {
		_, err = tx.Exec(updateQuery, date, method, comment, debit, actionById, actionByName, paymentId, groupId)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to update payment: %v", err)
	}

	_, err = r.educationClient.ChangeUserBalanceHistoryByDebit(context.TODO(), userId, oldDebit, date, comment, paymentType, actionById, actionByName, groupId, debit)
	if err != nil {
		return nil, fmt.Errorf("failed to update user balance history: %v", err)
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
			SUM(CASE WHEN payment_type = 'ADD' OR payment_type='REFUND' THEN amount ELSE 0 END) AS total_add,
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
			coalesce(group_id , 0),
			method
		FROM 
			student_payments 
		WHERE 
			student_id = $1 AND 
			TO_CHAR(given_date, 'YYYY-MM') = $2 
		ORDER BY 
			created_at desc`

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
			&payment.Method,
		); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		if givenDate.Valid {
			payment.GivenDate = givenDate.Time.Format("2006-01-02")
		} else {
			payment.GivenDate = ""
		}

		payment.GroupName = r.educationClient.GetGroupNameById(payment.GroupId)
		payments = append(payments, &payment)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %v", err)
	}

	return &pb.GetAllPaymentsByMonthResponse{
		Payments: payments,
	}, nil
}

func (r *PaymentRepository) GetAllPaymentTakeOff(from string, to string) (*pb.GetAllPaymentTakeOffResponse, error) {
	query := `
        SELECT 
            id, 
            given_date, 
            student_id, 
            comment, 
            created_by_id, 
            created_by_name
        FROM 
            student_payments
        WHERE 
            payment_type = 'TAKE_OFF' 
            AND created_by_id != '00000000-0000-0000-0000-000000000000'
            AND given_date BETWEEN $1 AND $2;
    `

	rows, err := r.db.Query(query, from, to)
	if err != nil {
		return nil, fmt.Errorf("error querying payments: %w", err)
	}
	defer rows.Close()

	response := &pb.GetAllPaymentTakeOffResponse{}

	for rows.Next() {
		var payment pb.AbsPaymentTakeOff

		err := rows.Scan(
			&payment.PaymentId,
			&payment.GivenDate,
			&payment.StudentId,
			&payment.Comment,
			&payment.CreatorId,
			&payment.CreatorName,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		name, _, _, err := r.educationClient.GetStudentById(payment.StudentId)
		if err != nil {
			payment.StudentName = "error while getting this student name"
		}
		payment.StudentName = name
		response.Pennies = append(response.Pennies, &payment)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return response, nil
}

func (r *PaymentRepository) GetAllPaymentTakeOffChart(from string, to string) (*pb.GetAllPaymentTakeOffChartResponse, error) {
	query := `
        SELECT 
            given_date, 
            SUM(amount)
        FROM 
            student_payments
        WHERE 
            payment_type = 'TAKE_OFF'
            AND created_by_id != '00000000-0000-0000-0000-000000000000'
            AND given_date BETWEEN $1 AND $2
        GROUP BY 
            given_date
        ORDER BY 
            given_date;
    `

	rows, err := r.db.Query(query, from, to)
	if err != nil {
		return nil, fmt.Errorf("error querying payment chart data: %w", err)
	}
	defer rows.Close()

	response := &pb.GetAllPaymentTakeOffChartResponse{}

	for rows.Next() {
		var chartEntry pb.AbsTakeOfChartResponse

		err := rows.Scan(
			&chartEntry.YearMonth,
			&chartEntry.Amount,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning chart data row: %w", err)
		}

		response.ChartResponse = append(response.ChartResponse, &chartEntry)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return response, nil
}

func (r *PaymentRepository) GetAllStudentPayments(from string, to string) (*pb.GetAllStudentPaymentsResponse, error) {
	query := `
SELECT 
       student_id,
       method,
       amount,
       given_date,
       comment,
       created_by_id,
       created_by_name
FROM student_payments
where given_date between $1 and $2
  and payment_type = 'ADD' order by created_at desc 
`
	rows, err := r.db.Query(query, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	resp := pb.GetAllStudentPaymentsResponse{}
	for rows.Next() {
		el := pb.AbsStudentPayments{}
		err := rows.Scan(&el.StudentId, &el.Method, &el.Amount, &el.GivenDate, &el.Comment, &el.CreatorId, &el.CreatorName)
		if err != nil {
			return nil, err
		}
		name, _, _, _ := r.educationClient.GetStudentById(el.StudentId)
		el.StudentName = name
		resp.Payments = append(resp.Payments, &el)
	}
	return &resp, nil
}

func (r *PaymentRepository) GetAllStudentPaymentsChart(from string, to string) (*pb.GetAllStudentPaymentsChartResponse, error) {
	var (
		cash  float64
		payme float64
		click float64
	)
	resp := pb.GetAllStudentPaymentsChartResponse{}
	query := `
SELECT coalesce((SELECT sum(amount)
                 FROM student_payments
                 where method = 'CASH' and payment_type != 'TAKE_OFF'
                   and given_date between $1 and $2), 0),
       coalesce((SELECT sum(amount)
                 FROM student_payments
                 where method = 'PAYME' and payment_type != 'TAKE_OFF'
                   and given_date between $1 and $2), 0),
       coalesce((SELECT sum(amount)
                 FROM student_payments
                 where method = 'CLICK' and payment_type != 'TAKE_OFF'
                   and given_date between $1 and $2), 0)
`

	err := r.db.QueryRow(query, from, to).Scan(&cash, &payme, &click)
	if err != nil {
		return nil, err
	}

	resp.TotalRevenue = strconv.FormatFloat(cash+payme+click, 'f', 2, 64)
	resp.Cash = strconv.FormatFloat(cash, 'f', 2, 64)
	resp.Click = strconv.FormatFloat(click, 'f', 2, 64)
	resp.Payme = strconv.FormatFloat(payme, 'f', 2, 64)

	query = `
        SELECT 
            given_date, 
            SUM(amount)
        FROM 
            student_payments
        WHERE 
            payment_type = 'ADD'
            AND given_date BETWEEN $1 AND $2
        GROUP BY 
            given_date
        ORDER BY 
            given_date;
    `

	rows, err := r.db.Query(query, from, to)
	if err != nil {
		return nil, fmt.Errorf("error querying payment chart data: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var chartEntry pb.AbsTakeOfChartResponse

		err := rows.Scan(
			&chartEntry.YearMonth,
			&chartEntry.Amount,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning chart data row: %w", err)
		}

		resp.PaymentsChart = append(resp.PaymentsChart, &chartEntry)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return &resp, nil
}

func (r *PaymentRepository) GetAllDebtsInformation(from, to string, page, size int32) (*pb.GetAllDebtsInformationResponse, error) {
	if size <= 0 {
		return nil, fmt.Errorf("invalid page size: must be greater than zero")
	}

	query := `
		SELECT 
			student_id AS debtor_id,
			COALESCE(SUM(CASE WHEN payment_type = 'ADD' OR payment_type = 'REFUND' THEN amount ELSE 0 END), 0) -
			COALESCE(SUM(CASE WHEN payment_type = 'TAKE_OFF' THEN amount ELSE 0 END), 0) AS total_on_period
		FROM 
			student_payments
		WHERE 
			given_date BETWEEN $1 AND $2
		GROUP BY 
			student_id
		LIMIT $3 OFFSET $4;
	`

	offset := (page - 1) * size

	rows, err := r.db.Query(query, from, to, size, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var debts []*pb.AbsDebtsInformation

	for rows.Next() {
		var debt pb.AbsDebtsInformation
		if err := rows.Scan(&debt.DebtorId, &debt.TotalOnPeriod); err != nil {
			return nil, err
		}
		name, phoneNumber, balance, err := r.educationClient.GetStudentById(debt.DebtorId)
		if err != nil {
			name = "unknown name"
			phoneNumber = "unknown phoneNumber"
			balance = 0
		}
		if balance >= 0 {
			continue
		}
		debt.DebtorName = name
		debt.PhoneNumber = phoneNumber
		debt.Balance = strconv.FormatFloat(balance, 'f', 2, 64)
		debts = append(debts, &debt)
	}

	var totalRecords int32
	if err := r.db.QueryRow(`SELECT COUNT(DISTINCT student_id) FROM student_payments WHERE given_date BETWEEN $1 AND $2`, from, to).Scan(&totalRecords); err != nil {
		return nil, err
	}

	if totalRecords == 0 {
		totalRecords = 1
	}
	totalPageCount := (totalRecords + size - 1) / size

	return &pb.GetAllDebtsInformationResponse{
		TotalPageCount: totalPageCount,
		Debts:          debts,
	}, nil
}

func (r *PaymentRepository) GetCommonFinanceInformation() (*pb.GetCommonInformationResponse, error) {
	var response *pb.GetCommonInformationResponse
	var payInCurrentMonth int32

	err := r.db.QueryRow(`SELECT COUNT(id) 
FROM student_payments 
WHERE payment_type = 'ADD' 
  AND EXTRACT(MONTH FROM given_date) = EXTRACT(MONTH FROM CURRENT_DATE) 
  AND EXTRACT(YEAR FROM given_date) = EXTRACT(YEAR FROM CURRENT_DATE);
`).Scan(&payInCurrentMonth)
	if err != nil {
		payInCurrentMonth = 0
	}
	response.DebtorsCount = 0
	response.PayInCurrentMonth = payInCurrentMonth
	return response, nil
}

func NewPaymentRepository(db *sql.DB, client *clients.EducationClient) *PaymentRepository {
	return &PaymentRepository{db: db, educationClient: client}
}
