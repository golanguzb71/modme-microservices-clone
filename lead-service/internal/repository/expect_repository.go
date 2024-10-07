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

func (r *ExpectRepository) CreateExpectation(title string) error {
	query := "INSERT INTO expect_section (title) VALUES ($1)"
	_, err := r.db.Exec(query, title)
	if err != nil {
		return fmt.Errorf("failed to create expectation: %w", err)
	}
	return nil
}

func (r *ExpectRepository) UpdateExpectation(id, title string) error {
	query := "UPDATE expect_section SET title = $1 WHERE id = $2"
	_, err := r.db.Exec(query, title, id)
	if err != nil {
		return fmt.Errorf("failed to update expectation: %w", err)
	}
	return nil
}

func (r *ExpectRepository) DeleteExpectation(id string) error {
	query := "DELETE FROM expect_section WHERE id = $1"
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete expectation: %w", err)
	}
	return nil
}
