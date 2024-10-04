package repository

import (
	"database/sql"
)

type LidRepository struct {
	db *sql.DB
}

func NewLidRepository(db *sql.DB) *LidRepository {
	return &LidRepository{db: db}
}
