package repository

import "database/sql"

type SetRepository struct {
	db *sql.DB
}

func NewSetRepository(db *sql.DB) *SetRepository {
	return &SetRepository{db: db}
}
