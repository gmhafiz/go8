package author

import (
	"encoding/json"
	"io"

	"github.com/gmhafiz/go8/internal/models"
)

type Request struct {
	Filter
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
}

func (r *Request) Decode(body io.ReadCloser) error {
	return json.NewDecoder(body).Decode(r)
}

func ToAuthor(req *Request) *models.Author {
	return &models.Author{
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}
}
