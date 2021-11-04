package author

import (
	"github.com/jinzhu/copier"

	"github.com/gmhafiz/go8/internal/domain/book"
	"github.com/gmhafiz/go8/internal/models"
)

type Res struct {
	ID         uint64 `json:"id" deepcopier:"field:id"`
	FirstName  string `json:"first_name" deepcopier:"field:first_name"`
	MiddleName string `json:"middle_name" deepcopier:"field:middle_name"`
	LastName   string `json:"last_name" deepcopier:"field:last_name"`
}

type ResWithBooks struct {
	ID         uint64     `json:"id" `
	FirstName  string     `json:"first_name"`
	MiddleName string     `json:"middle_name" swaggertype:"string"`
	LastName   string     `json:"last_name"`
	Books      []book.Res `json:"books"`
}

func ResourceUpdate(author *models.Author) (Res, error) {
	return Res{
		ID:         uint64(author.ID),
		FirstName:  author.FirstName,
		MiddleName: author.MiddleName.String,
		LastName:   author.LastName,
	}, nil
}

func Resource(author *models.Author) (Res, error) {
	var resource Res

	err := copier.Copy(&resource, &author)
	if err != nil {
		return resource, err
	}

	return resource, nil
}

func ResourceWithBooks(authorWithBooks *WithBooks) ResWithBooks {
	return ResWithBooks{
		ID:         uint64(authorWithBooks.ID),
		FirstName:  authorWithBooks.FirstName,
		MiddleName: authorWithBooks.MiddleName.String,
		LastName:   authorWithBooks.LastName,
		Books:      authorWithBooks.Books,
	}
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

type CreateResponse struct {
	ID         int64  `json:"id"`
	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name"`
	LastName   string `json:"last_name"`
	Books      []Book `json:"books,omitempty"`
}
