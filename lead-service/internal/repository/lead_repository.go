package repository

import (
	"database/sql"
	"fmt"
	"lid-service/proto/pb"
	"log"
)

type LeadRepository struct {
	db *sql.DB
}

func NewLeadRepository(db *sql.DB) *LeadRepository {
	return &LeadRepository{db: db}
}

func (r *LeadRepository) CreateLead(title string) error {
	query := "INSERT INTO lead_section (title) VALUES ($1)"
	_, err := r.db.Exec(query, title)
	if err != nil {
		return fmt.Errorf("failed to create lead: %w", err)
	}
	return nil
}

func (r *LeadRepository) GetLeadCommon(req *pb.GetLeadCommonRequest) (*pb.GetLeadCommonResponse, error) {
	resp := &pb.GetLeadCommonResponse{}
	requestedSections := make(map[string][]string)
	for _, request := range req.Requests {
		requestedSections[request.Type] = append(requestedSections[request.Type], request.Id)
	}
	fetchAllSections(resp, r.db, requestedSections)

	return resp, nil
}

func fetchAllSections(p *pb.GetLeadCommonResponse, db *sql.DB, requestedSections map[string][]string) {
	calculateSet(p, db, requestedSections["set"])
	calculateExpectations(p, db, requestedSections["expectation"])
	calculateLeadsWithDetails(p, db, requestedSections["lead"])
}

func calculateSet(p *pb.GetLeadCommonResponse, db *sql.DB, requestedIds []string) {
	query := `
        SELECT ss.id, ss.title
        FROM set_section ss
    `
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Error fetching sets: %v", err)
		return
	}
	defer rows.Close()

	var sections []*pb.Section
	for rows.Next() {
		section := &pb.Section{}
		if err := rows.Scan(&section.Id, &section.Name); err != nil {
			log.Printf("Error scanning set section row: %v", err)
			return
		}
		section.Type = "set"

		if containsString(requestedIds, section.Id) {
			section.Leads = fetchLeadsForSection(db, section.Id, "set")
			section.LeadsCount = int32(len(section.Leads))
		}

		sections = append(sections, section)
	}
	p.Sets = sections
}

func calculateExpectations(p *pb.GetLeadCommonResponse, db *sql.DB, requestedIds []string) {
	query := `
        SELECT es.id, es.title
        FROM expect_section es
    `
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Error fetching expectations: %v", err)
		return
	}
	defer rows.Close()

	var sections []*pb.Section
	for rows.Next() {
		section := &pb.Section{}
		if err := rows.Scan(&section.Id, &section.Name); err != nil {
			log.Printf("Error scanning expectation section row: %v", err)
			return
		}
		section.Type = "expectation"

		if containsString(requestedIds, section.Id) {
			section.Leads = fetchLeadsForSection(db, section.Id, "expectation")
			section.LeadsCount = int32(len(section.Leads))
		}

		sections = append(sections, section)
	}
	p.Expectations = sections
}

func calculateLeadsWithDetails(p *pb.GetLeadCommonResponse, db *sql.DB, requestedIds []string) {
	query := `
        SELECT ls.id, ls.title
        FROM lead_section ls
    `
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Error fetching lead sections: %v", err)
		return
	}
	defer rows.Close()

	var sections []*pb.Section
	for rows.Next() {
		section := &pb.Section{}
		if err := rows.Scan(&section.Id, &section.Name); err != nil {
			log.Printf("Error scanning lead section row: %v", err)
			return
		}
		section.Type = "lead"

		// Check if this section's leads should be fetched
		if containsString(requestedIds, section.Id) {
			section.Leads = fetchLeadsForSection(db, section.Id, "lead")
			section.LeadsCount = int32(len(section.Leads))
		}

		sections = append(sections, section)
	}
	p.Leads = sections
}

func fetchLeadsForSection(db *sql.DB, sectionId, sectionType string) []*pb.Lead {
	var query string
	switch sectionType {
	case "set":
		query = `SELECT id, comment, created_at, phone_number FROM lead_user WHERE set_id=$1`
	case "expectation":
		query = `SELECT id, comment, created_at, phone_number FROM lead_user WHERE expect_id=$1`
	case "lead":
		query = `SELECT id, comment, created_at, phone_number FROM lead_user WHERE lead_id=$1`
	default:
		return nil
	}

	rows, err := db.Query(query, sectionId)
	if err != nil {
		log.Printf("Error fetching leads for section: %v", err)
		return nil
	}
	defer rows.Close()

	var leads []*pb.Lead
	for rows.Next() {
		lead := &pb.Lead{}
		if err := rows.Scan(&lead.Id, &lead.Comment, &lead.CreatedAt, &lead.PhoneNumber); err != nil {
			log.Printf("Error scanning lead row: %v", err)
			return nil
		}
		leads = append(leads, lead)
	}
	return leads
}

func containsString(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

func (r *LeadRepository) UpdateLead(id, title string) error {
	query := "UPDATE lead_section SET title = $1 WHERE id = $2"
	_, err := r.db.Exec(query, title, id)
	if err != nil {
		return fmt.Errorf("failed to update lead: %w", err)
	}
	return nil
}

func (r *LeadRepository) DeleteLead(id string) error {
	query := "DELETE FROM lead_section WHERE id = $1"
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete lead: %w", err)
	}
	return nil
}
