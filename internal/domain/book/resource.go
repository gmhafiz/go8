package book

import (
	"time"

	"github.com/volatiletech/null/v8"

	"github.com/gmhafiz/go8/internal/models"
)

type Res struct {
	ID            int64       `json:"id"`
	Title         string      `json:"title"`
	PublishedDate time.Time   `json:"published_date"`
	ImageURL      null.String `json:"image_url" swaggertype:"string"`
	Description   null.String `json:"description" swaggertype:"string"`
}

func Resource(book *models.Book) *Res {
	resource := &Res{
		ID:            book.ID,
		Title:         book.Title,
		PublishedDate: book.PublishedDate,
		ImageURL:      book.ImageURL,
		Description: null.String{
			String: book.Description,
			Valid:  true,
		},
	}

	return resource
}

func Resources(books []*models.Book) ([]*Res, error) {
	if len(books) == 0 {
		return make([]*Res, 0), nil
	}

	var resources []*Res
	for _, book := range books {
		res := Resource(book)
		resources = append(resources, res)
	}
	return resources, nil
}
