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

func (r *SetRepository) CreateSet(companyId, title, courseId, teacherId, dateType string, date []string, lessonStartTime string) error {
	query := `
		INSERT INTO set_section (title, course_id, teacher_id, date_type, days, start_time , company_id) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.Exec(query, title, courseId, teacherId, dateType, pq.Array(date), lessonStartTime, companyId)
	if err != nil {
		return fmt.Errorf("failed to create set: %w", err)
	}
	return nil
}

func (r *SetRepository) UpdateSet(companyId, id, title, courseId, teacherId, dateType string, date []string, lessonStartTime string) error {
	query := `
		UPDATE set_section 
		SET title = $1, course_id = $2, teacher_id = $3, date_type = $4, days = $5, start_time = $6 
		WHERE id = $7 and company_id=$8`
	_, err := r.db.Exec(query, title, courseId, teacherId, dateType, pq.Array(date), lessonStartTime, id, companyId)
	if err != nil {
		return fmt.Errorf("failed to update set: %w", err)
	}
	return nil
}

func (r *SetRepository) DeleteSet(companyId, id string) error {
	query := "DELETE FROM set_section WHERE id = $1 and company_id=$2"
	_, err := r.db.Exec(query, id, companyId)
	if err != nil {
		return fmt.Errorf("failed to delete set: %w", err)
	}
	return nil
}

func (r *SetRepository) GetLeadDataBySetId(companyId, setId string) ([]string, []string, error) {
	queryLeadData := `SELECT full_name , phone_number FROM lead_user where set_id=$1 and company_id=$2 and company_id=$3`
	rows, err := r.db.Query(queryLeadData, setId, companyId)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	var names, phoneNumbers []string
	for rows.Next() {
		var name, phoneNumber string
		err = rows.Scan(&name, &phoneNumber)
		if err != nil {
			return nil, nil, err
		}

		names = append(names, name)
		phoneNumbers = append(phoneNumbers, phoneNumber)
	}
	return names, phoneNumbers, nil
}
