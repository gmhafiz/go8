package author

import (
	"github.com/gmhafiz/go8/internal/domain/book"
)

type GetResponse struct {
	ID         uint           `json:"id"`
	FirstName  string         `json:"first_name"`
	MiddleName string         `json:"middle_name"`
	LastName   string         `json:"last_name"`
	Books      []*book.Schema `json:"books"`
}

func Resource(a *Schema) *GetResponse {
	if a == nil {
		return &GetResponse{}
	}

	return &GetResponse{
		ID:         a.ID,
		FirstName:  a.FirstName,
		MiddleName: a.MiddleName,
		LastName:   a.LastName,
		Books:      a.Books,
	}
}

func Resources(books []*Schema) []*GetResponse {
	if len(books) == 0 {
		return make([]*GetResponse, 0)
	}

	var resources []*GetResponse
	for _, a := range books {
		res := Resource(a)
		resources = append(resources, res)
	}
	return resources
}
