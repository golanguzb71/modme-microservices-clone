package repository

import (
	"database/sql"
)

type ExpectRepository struct {
	db *sql.DB
}

func NewExpectRepository(db *sql.DB) *ExpectRepository {
	return &ExpectRepository{db: db}
}
