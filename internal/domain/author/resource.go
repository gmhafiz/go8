package author

import (
	"github.com/gmhafiz/go8/internal/models"
	"github.com/jinzhu/copier"
)

type Res struct {
	AuthorID  uint64 `json:"author_id" deepcopier:"field:author_id" db:"id"`
	FirstName string `json:"first_name" deepcopier:"field:first_name" db:"first_name"`
	LastName  string `json:"last_name" deepcopier:"field:last_name" db:"last_name"`
}

func Resource(author *models.Author) (Res, error) {
	var resource Res

	err := copier.Copy(&resource, &author)
	if err != nil {
		return resource, err
	}

	return resource, nil
}

func Resources(authors []*models.Author) (interface{}, error) {
	if len(authors) == 0 {
		return make([]string, 0), nil
	}

	var resources []Res
	for _, author := range authors {
		res, _ := Resource(author)
		resources = append(resources, res)
	}
	return resources, nil
}
