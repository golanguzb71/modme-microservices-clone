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

	// Iterate over each request
	for _, request := range req.Requests {
		id := request.Id
		sectionType := request.Type

		// Call appropriate calculation based on the section type
		switch sectionType {
		case "set":
			calculateSet(resp, r.db, &id, &sectionType)
		case "expectation":
			calculateExpectations(resp, r.db, &id, &sectionType)
		case "lead":
			calculateLeads(resp, r.db, &id, &sectionType)
		default:
			log.Printf("Unknown section type: %s", sectionType)
		}
	}

	return resp, nil
}

func calculateSet(p *pb.GetLeadCommonResponse, db *sql.DB, id, sectionType *string) {
	query := `
        SELECT ss.id, ss.title, COUNT(lu.id) as leads_count
        FROM set_section ss
        LEFT JOIN lead_user lu ON lu.set_id = ss.id
        GROUP BY ss.id
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
		if err := rows.Scan(&section.Id, &section.Name, &section.LeadsCount); err != nil {
			log.Printf("Error scanning set section row: %v", err)
			return
		}

		if id != nil && *sectionType == "set" && section.Id == *id {
			section.Leads = fetchLeadsForSection(db, section.Id, "set")
		} else {
			section.Leads = []*pb.Lead{}
		}

		section.Type = "set"
		sections = append(sections, section)
	}
	p.Sets = sections
}

func calculateExpectations(p *pb.GetLeadCommonResponse, db *sql.DB, id, sectionType *string) {
	query := `
        SELECT es.id, es.title, COUNT(lu.id) as leads_count
        FROM expect_section es
        LEFT JOIN lead_user lu ON lu.expect_id = es.id
        GROUP BY es.id, es.title
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
		if err := rows.Scan(&section.Id, &section.Name, &section.LeadsCount); err != nil {
			log.Printf("Error scanning expectation section row: %v", err)
			return
		}

		if id != nil && *sectionType == "expectation" && section.Id == *id {
			section.Leads = fetchLeadsForSection(db, section.Id, "expectation")
		} else {
			section.Leads = []*pb.Lead{}
		}

		section.Type = "expectation"
		sections = append(sections, section)
	}
	p.Expectations = sections
}

func calculateLeads(p *pb.GetLeadCommonResponse, db *sql.DB, id, sectionType *string) {
	query := `
        SELECT ls.id, ls.title, COUNT(lu.id) as leads_count
        FROM lead_section ls
        LEFT JOIN lead_user lu ON lu.lead_id = ls.id
        GROUP BY ls.id
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
		if err := rows.Scan(&section.Id, &section.Name, &section.LeadsCount); err != nil {
			log.Printf("Error scanning lead section row: %v", err)
			return
		}

		if id != nil && *sectionType == "lead" && section.Id == *id {
			section.Leads = fetchLeadsForSection(db, section.Id, "lead")
		} else {
			section.Leads = []*pb.Lead{}
		}

		section.Type = "lead"
		sections = append(sections, section)
	}
	p.Leads = sections
}

func fetchLeadsForSection(db *sql.DB, sectionId, sectionType string) []*pb.Lead {
	query := `
        SELECT id, comment, created_at, phone_number
        FROM lead_user WHERE 
    `
	switch sectionType {
	case "set":
		query += ` set_id=$1`
	case "expectation":
		query += ` expect_id=$1`
	case "lead":
		query += ` lead_id=$1`
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
