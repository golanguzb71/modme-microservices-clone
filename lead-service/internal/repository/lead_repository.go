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

func (r *LeadRepository) CreateLead(companyId string, title string) error {
	query := "INSERT INTO lead_section (title , company_id) VALUES ($1 , $2)"
	_, err := r.db.Exec(query, title, companyId)
	if err != nil {
		return fmt.Errorf("failed to create lead: %w", err)
	}
	return nil
}

func (r *LeadRepository) GetLeadCommon(companyId string, req *pb.GetLeadCommonRequest) (*pb.GetLeadCommonResponse, error) {
	resp := &pb.GetLeadCommonResponse{}
	requestedSections := make(map[string][]string)
	for _, request := range req.Requests {
		requestedSections[request.Type] = append(requestedSections[request.Type], request.Id)
	}
	fetchAllSections(companyId, resp, r.db, requestedSections)

	return resp, nil
}

func fetchAllSections(companyId string, p *pb.GetLeadCommonResponse, db *sql.DB, requestedSections map[string][]string) {
	calculateSet(companyId, p, db, requestedSections["set"])
	calculateExpectations(companyId, p, db, requestedSections["expectation"])
	calculateLeadsWithDetails(companyId, p, db, requestedSections["lead"])
}

func calculateSet(companyId string, p *pb.GetLeadCommonResponse, db *sql.DB, requestedIds []string) {
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
			section.Leads = fetchLeadsForSection(companyId, db, section.Id, "set")
		}
		_ = db.QueryRow(`SELECT count(*) FROM lead_user where set_id=$1 and company_id=$2`, section.Id, companyId).Scan(&section.LeadsCount)
		sections = append(sections, section)
	}
	p.Sets = sections
}

func calculateExpectations(companyId string, p *pb.GetLeadCommonResponse, db *sql.DB, requestedIds []string) {
	query := `
        SELECT es.id, es.title
        FROM expect_section es where company_id=$1
    `
	rows, err := db.Query(query, companyId)
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
			section.Leads = fetchLeadsForSection(companyId, db, section.Id, "expectation")
		}
		_ = db.QueryRow(`SELECT count(*) FROM lead_user where expect_id=$1 and company_id=$2`, section.Id, companyId).Scan(&section.LeadsCount)
		sections = append(sections, section)
	}
	p.Expectations = sections
}

func calculateLeadsWithDetails(companyId string, p *pb.GetLeadCommonResponse, db *sql.DB, requestedIds []string) {
	query := `
        SELECT ls.id, ls.title
        FROM lead_section ls where company_id=$1
    `
	rows, err := db.Query(query, companyId)
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
		if containsString(requestedIds, section.Id) {
			section.Leads = fetchLeadsForSection(companyId, db, section.Id, "lead")
		}
		_ = db.QueryRow(`SELECT count(*) FROM lead_user where lead_id=$1 and company_id=$2`, section.Id, companyId).Scan(&section.LeadsCount)
		sections = append(sections, section)
	}
	p.Leads = sections
}

func fetchLeadsForSection(companyId string, db *sql.DB, sectionId, sectionType string) []*pb.Lead {
	var query string
	switch sectionType {
	case "set":
		query = `SELECT id, full_name,  comment, created_at, phone_number FROM lead_user WHERE set_id=$1 and company_id=$2`
	case "expectation":
		query = `SELECT id,full_name, comment, created_at, phone_number FROM lead_user WHERE expect_id=$1 and company_id=$2`
	case "lead":
		query = `SELECT id, full_name,  comment, created_at, phone_number FROM lead_user WHERE lead_id=$1 and company_id=$2`
	default:
		return nil
	}

	rows, err := db.Query(query, sectionId, companyId)
	if err != nil {
		log.Printf("Error fetching leads for section: %v", err)
		return nil
	}
	defer rows.Close()

	var leads []*pb.Lead
	for rows.Next() {
		lead := &pb.Lead{}
		if err := rows.Scan(&lead.Id, &lead.Name, &lead.Comment, &lead.CreatedAt, &lead.PhoneNumber); err != nil {
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

func (r *LeadRepository) UpdateLead(companyId string, id, title string) error {
	query := "UPDATE lead_section SET title = $1 WHERE id = $2 and company_id=$3"
	_, err := r.db.Exec(query, title, id, companyId)
	if err != nil {
		return fmt.Errorf("failed to update lead: %w", err)
	}
	return nil
}

func (r *LeadRepository) DeleteLead(companyId string, id string) error {
	query := "DELETE FROM lead_section WHERE id = $1 and company_id=$2"
	_, err := r.db.Exec(query, id, companyId)
	if err != nil {
		return fmt.Errorf("failed to delete lead: %w", err)
	}
	return nil
}

func (r *LeadRepository) GetAllLeads(companyId string) (*pb.GetLeadListResponse, error) {
	query := `SELECT id, title FROM lead_section where company_id=$1`
	rows, err := r.db.Query(query, companyId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := &pb.GetLeadListResponse{}

	for rows.Next() {
		var id, title string
		if err := rows.Scan(&id, &title); err != nil {
			return nil, err
		}

		section := &pb.DynamicSection{
			Id:   id,
			Name: title,
		}
		result.Sections = append(result.Sections, section)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *LeadRepository) GetLeadReports(companyId string, endYear string, startYear string) (*pb.GetLeadReportsResponse, error) {
	if startYear > endYear {
		return nil, fmt.Errorf("start year must be less than or equal to end year")
	}

	response := &pb.GetLeadReportsResponse{
		LeadConversion:          []*pb.LeadConversion{},
		LeadConversionForSource: []*pb.LeadConversionForSource{},
	}

	conversionQuery := `
        SELECT conversion_date, SUM(lead_count) as total_leads
        FROM lead_conversion_reports
        WHERE conversion_date >= $1 AND conversion_date <= $2 and company_id=$3
        GROUP BY conversion_date
        ORDER BY conversion_date
    `
	conversionStartDate := startYear
	conversionEndDate := endYear
	conversionRows, err := r.db.Query(conversionQuery, conversionStartDate, conversionEndDate, companyId)
	if err != nil {
		return nil, fmt.Errorf("error querying lead conversions: %v", err)
	}
	defer conversionRows.Close()

	for conversionRows.Next() {
		var conversionDate string
		var leadCount int32
		if err := conversionRows.Scan(&conversionDate, &leadCount); err != nil {
			return nil, fmt.Errorf("error scanning lead conversion row: %v", err)
		}

		response.LeadConversion = append(response.LeadConversion, &pb.LeadConversion{
			ConversionDate: conversionDate,
			LeadCount:      leadCount,
		})
	}

	sourceQuery := `
        SELECT source, SUM(lead_count) as total_leads
        FROM lead_source_reports
        WHERE created_at >= $1 AND created_at <= $2 and company_id=$3
        GROUP BY source
        ORDER BY total_leads DESC
    `
	sourceRows, err := r.db.Query(sourceQuery, conversionStartDate, conversionEndDate, companyId)
	if err != nil {
		return nil, fmt.Errorf("error querying lead sources: %v", err)
	}
	defer sourceRows.Close()

	for sourceRows.Next() {
		var source string
		var leadsCount int32
		if err := sourceRows.Scan(&source, &leadsCount); err != nil {
			return nil, fmt.Errorf("error scanning lead source row: %v", err)
		}

		response.LeadConversionForSource = append(response.LeadConversionForSource, &pb.LeadConversionForSource{
			Source:     source,
			LeadsCount: leadsCount,
		})
	}

	return response, nil
}

func (r *LeadRepository) GetActiveLeadCount(companyId string) (*pb.GetActiveLeadCountResponse, error) {
	activeLeadCount := 0
	r.db.QueryRow(`SELECT COUNT(*) FROM lead_user where company_id=$1`).Scan(&activeLeadCount, companyId)
	return &pb.GetActiveLeadCountResponse{ActiveLeadCount: int32(activeLeadCount)}, nil
}
