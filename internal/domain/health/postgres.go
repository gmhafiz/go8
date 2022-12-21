package postgres

import (
	"github.com/jmoiron/sqlx"
)
type Repository interface {
	Readiness() error
}

type repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Readiness() error {
	return r.db.Ping()
}
