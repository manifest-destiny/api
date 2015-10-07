package api

import (
	"github.com/jmoiron/sqlx"
)

// DB wrapper for sqlx.DB.
type DB struct {
	*sqlx.DB
}

// NewDB constructor for database connection.
func NewDB(driver, info string) (*DB, error) {
	db, err := sqlx.Open(driver, info)
	if err != nil {
		return &DB{}, err
	}
	return &DB{db}, nil
}
