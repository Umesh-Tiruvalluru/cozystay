package repository

import (
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

func NewRepositoryUser(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}
