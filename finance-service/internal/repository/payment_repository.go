package repository

import (
	"database/sql"
	"finance-service/internal/utils"
	"github.com/google/uuid"
	"time"
)

type PaymentRepository struct {
	db *sql.DB
}

func (r *PaymentRepository) PaidStudent(studentID, comment, sum, date, method, createdBy string) error {
	query := `
        INSERT INTO student_pay (id, student_id, payment_type, amount, given_date, comment, created_by , created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7 , $8)
    `
	id := uuid.New()
	amount, err := utils.ParseAmount(sum)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(query, id, studentID, method, amount, date, comment, createdBy, time.Now())
	return err
}

func (r *PaymentRepository) UnPaidStudent(studentID, comment, sum, date, method, createdBy string) error {
	query := `
        INSERT INTO student_unpay (id, student_id, payment_type, amount, given_date, comment, created_by , created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7,  $8)
    `
	id := uuid.New()
	amount, err := utils.ParseAmount(sum)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(query, id, studentID, method, amount, date, comment, createdBy, time.Now())
	return err
}

func NewPaymentRepository(db *sql.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}
