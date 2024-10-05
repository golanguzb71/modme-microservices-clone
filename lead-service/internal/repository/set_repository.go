package repository

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
)

type SetRepository struct {
	db *sql.DB
}

// NewSetRepository initializes a new SetRepository
func NewSetRepository(db *sql.DB) *SetRepository {
	return &SetRepository{db: db}
}

func (r *SetRepository) CreateSet(title, courseId, teacherId, dateType string, date []string, lessonStartTime string) error {
	query := `
		INSERT INTO set_section (title, course_id, teacher_id, date_type, days, start_time) 
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.Exec(query, title, courseId, teacherId, dateType, pq.Array(date), lessonStartTime)
	if err != nil {
		return fmt.Errorf("failed to create set: %w", err)
	}
	return nil
}

func (r *SetRepository) UpdateSet(id, title, courseId, teacherId, dateType string, date []string, lessonStartTime string) error {
	query := `
		UPDATE set_section 
		SET title = $1, course_id = $2, teacher_id = $3, date_type = $4, days = $5, start_time = $6 
		WHERE id = $7`
	_, err := r.db.Exec(query, title, courseId, teacherId, dateType, pq.Array(date), lessonStartTime, id)
	if err != nil {
		return fmt.Errorf("failed to update set: %w", err)
	}
	return nil
}

func (r *SetRepository) DeleteSet(id string) error {
	query := "DELETE FROM set_section WHERE id = $1"
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete set: %w", err)
	}
	return nil
}
