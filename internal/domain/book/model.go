package book

import (
	"database/sql"
	"time"
)

type Schema struct {
	ID            int          `db:"id"`
	Title         string       `db:"title"`
	PublishedDate time.Time    `db:"published_date"`
	ImageURL      string       `db:"image_url"`
	Description   string       `db:"description"`
	CreatedAt     time.Time    `db:"created_at"`
	UpdatedAt     time.Time    `db:"updated_at"`
	DeletedAt     sql.NullTime `db:"deleted_at" swaggertype:"string"`
}
