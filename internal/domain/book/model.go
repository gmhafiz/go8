package book

import (
	"database/sql"
	"time"
)

type DB struct {
	ID            int          `db:"id"  json:"id"`
	Title         string       `db:"title" json:"title"`
	PublishedDate time.Time    `db:"published_date" json:"published_date"`
	ImageURL      string       `db:"image_url" json:"image_url"`
	Description   string       `db:"description" json:"description"`
	CreatedAt     time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt     sql.NullTime `db:"updated_at" json:"updated_at"`
	DeletedAt     sql.NullTime `db:"deleted_at" json:"deleted_at"`
}
