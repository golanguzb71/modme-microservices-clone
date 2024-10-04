package repository

import "database/sql"

type LidUserRepository struct {
	db *sql.DB
}

func NewLidUserRepository(db *sql.DB) *LidUserRepository {
	return &LidUserRepository{db: db}
}
