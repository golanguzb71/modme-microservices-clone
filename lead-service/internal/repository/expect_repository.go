package repository

import (
	"database/sql"
	"fmt"
)

type ExpectRepository struct {
	db *sql.DB
}

func NewExpectRepository(db *sql.DB) *ExpectRepository {
	return &ExpectRepository{db: db}
}

func (r *ExpectRepository) CreateExpectation(companyId, title string) error {
	query := "INSERT INTO expect_section (title , company_id) VALUES ($1 , $2)"
	_, err := r.db.Exec(query, title, companyId)
	if err != nil {
		return fmt.Errorf("failed to create expectation: %w", err)
	}
	return nil
}

func (r *ExpectRepository) UpdateExpectation(companyId, id, title string) error {
	query := "UPDATE expect_section SET title = $1 WHERE id = $2 and company_id=$3"
	_, err := r.db.Exec(query, title, id, companyId)
	if err != nil {
		return fmt.Errorf("failed to update expectation: %w", err)
	}
	return nil
}

func (r *ExpectRepository) DeleteExpectation(companyId, id string) error {
	query := "DELETE FROM expect_section WHERE id = $1 and company_id=$2"
	_, err := r.db.Exec(query, id, companyId)
	if err != nil {
		return fmt.Errorf("failed to delete expectation: %w", err)
	}
	return nil
}
