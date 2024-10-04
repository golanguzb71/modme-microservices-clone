package repository

import "database/sql"

type LeadDataRepository struct {
	db *sql.DB
}

func NewLeadDataRepository(db *sql.DB) *LeadDataRepository {
	return &LeadDataRepository{db: db}
}
