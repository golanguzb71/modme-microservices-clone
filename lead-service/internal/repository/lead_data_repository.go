package repository

import (
	"database/sql"
	"errors"
	"fmt"
)

type LeadDataRepository struct {
	db *sql.DB
}

// NewLeadDataRepository initializes a new LeadDataRepository
func NewLeadDataRepository(db *sql.DB) *LeadDataRepository {
	return &LeadDataRepository{db: db}
}

func (r *LeadDataRepository) CreateLeadData(phoneNumber, leadID, expectID, setID, comment, name *string) error {
	query := `
		INSERT INTO lead_user (phone_number, lead_id, expect_id, set_id, comment , full_name) 
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.Exec(query, phoneNumber, leadID, expectID, setID, comment, name)
	if err != nil {
		return fmt.Errorf("failed to create lead data: %w", err)
	}
	return nil
}

func (r *LeadDataRepository) UpdateLeadData(id, phoneNumber, comment, sectionType, sectionId string) error {
	query := `
		UPDATE lead_user 
		SET phone_number = $1, comment = $2 `

	switch sectionType {
	case "SET":
		query += `, set_id=$3 , expect_id=null , lead_id=null `
	case "EXPECTATION":
		query += `, expect_id=$3 , set_id=null , lead_id=null `
	case "LEAD":
		query += `, lead_id=$3, expect_id=null , set_id=null `
	default:
		return errors.New("section type should include : SET , EXPECTATION , LEAD")
	}

	query += ` WHERE id = $4`
	_, err := r.db.Exec(query, phoneNumber, comment, sectionId, id)
	if err != nil {
		return fmt.Errorf("failed to update lead data: %w", err)
	}
	return nil
}

func (r *LeadDataRepository) DeleteLeadData(id string) error {
	query := "DELETE FROM lead_user WHERE id = $1"
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete lead data: %w", err)
	}
	return nil
}
