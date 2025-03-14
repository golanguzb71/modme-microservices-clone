package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type LeadDataRepository struct {
	db *sql.DB
}

// NewLeadDataRepository initializes a new LeadDataRepository
func NewLeadDataRepository(db *sql.DB) *LeadDataRepository {
	return &LeadDataRepository{db: db}
}

func (r *LeadDataRepository) CreateLeadData(companyId string, phoneNumber, leadID, expectID, setID, comment, name *string) error {
	query := `
		INSERT INTO lead_user (phone_number, lead_id, expect_id, set_id, comment , full_name , company_id) 
		VALUES ($1, $2, $3, $4, $5, $6 , $7)`
	_, err := r.db.Exec(query, phoneNumber, leadID, expectID, setID, comment, name, companyId)
	if err != nil {
		return fmt.Errorf("failed to create lead data: %w", err)
	}
	var title string
	err = r.db.QueryRow(`SELECT title from lead_section where id=$1 and company_id=$2`, leadID, companyId).Scan(&title)
	if err != nil {
		return status.Errorf(codes.Aborted, err.Error())
	}
	var checker bool
	_ = r.db.QueryRow(`SELECT exists(SELECT 1 FROM lead_source_reports where source=$1 and company_id=$2)`, title, companyId).Scan(&checker)
	if checker {
		_, _ = r.db.Exec(`UPDATE lead_source_reports SET lead_count=lead_count+1 where source=$1 and company_id=$2`, title, companyId)
	} else {
		_, _ = r.db.Exec(`INSERT INTO lead_source_reports(id, lead_count, source , created_at , company_id) values ($1 , $2 , $3 , $4 , $5)`, uuid.New(), 1, title, time.Now(), companyId)
	}
	format := time.Now().Format("2006-01")

	_ = r.db.QueryRow(`SELECT exists(SELECT 1 FROM lead_conversion_reports where conversion_date=$1 and company_id=$2)`, format, companyId).Scan(&checker)
	if checker {
		_, _ = r.db.Exec(`UPDATE lead_conversion_reports SET lead_count=lead_count+1 where conversion_date=$1 and company_id=$2`, format, companyId)
	} else {
		_, _ = r.db.Exec(`INSERT INTO lead_conversion_reports(id, lead_count, conversion_date, created_at , company_id) values ($1 , $2 , $3 , $4,  $5)`, uuid.New(), 1, format, time.Now(), companyId)
	}
	return nil
}

func (r *LeadDataRepository) UpdateLeadData(companyId string, id, phoneNumber, comment, name string) error {
	query := `
		UPDATE lead_user 
		SET phone_number = $1, comment = $2 , full_name= $3 WHERE id = $4 and company_id=$5`

	_, err := r.db.Exec(query, phoneNumber, comment, name, id, companyId)
	if err != nil {
		return fmt.Errorf("failed to update lead data: %w", err)
	}
	return nil
}

func (r *LeadDataRepository) DeleteLeadData(companyId string, id string) error {
	query := "DELETE FROM lead_user WHERE id = $1 and company_id=$2"
	_, err := r.db.Exec(query, id, companyId)
	if err != nil {
		return fmt.Errorf("failed to delete lead data: %w", err)
	}
	return nil
}

func (r *LeadDataRepository) ChangeLeadPlace(companyId string, sectionID, sectionType, itemId *string) error {
	query := `UPDATE lead_user SET `
	switch *sectionType {
	case "set":
		query += `set_id=$1 , expect_id=null,lead_id=null`
	case "expectation":
		query += `set_id=null , expect_id=$1 , lead_id=null`
	case "lead":
		query += `set_id=null, expect_id=null,lead_id=$1`
	default:
		return errors.New("section type should include : set , expectation , lead")
	}

	query += ` where id=$2 and company_id=$3`
	_, err := r.db.Exec(query, *sectionID, *itemId, companyId)
	if err != nil {
		return err
	}
	return nil
}
