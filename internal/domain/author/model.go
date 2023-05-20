package author

import (
	"time"

	"github.com/gmhafiz/go8/internal/domain/book"
)

type Schema struct {
	ID         uint64
	FirstName  string
	MiddleName string
	LastName   string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
	Books      []*book.Schema
}
