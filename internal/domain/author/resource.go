package author

import (
	"github.com/gmhafiz/go8/ent/gen"
)

type CreateResponse struct {
	ID         uint        `json:"id"`
	FirstName  string      `json:"first_name"`
	MiddleName string      `json:"middle_name"`
	LastName   string      `json:"last_name"`
	Books      []*gen.Book `json:"books,omitempty"`
}
