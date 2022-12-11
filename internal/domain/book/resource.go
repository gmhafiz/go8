package book

import (
	"time"
)

type Res struct {
	ID            int       `json:"id"`
	Title         string    `json:"title"`
	PublishedDate time.Time `json:"published_date"`
	ImageURL      string    `json:"image_url" swaggertype:"string"`
	Description   string    `json:"description" swaggertype:"string"`
}

func Resource(book *Schema) *Res {
	if book == nil {
		return &Res{}
	}
	resource := &Res{
		ID:            book.ID,
		Title:         book.Title,
		PublishedDate: book.PublishedDate,
		ImageURL:      book.ImageURL,
		Description:   book.Description,
	}

	return resource
}

func Resources(books []*Schema) ([]*Res, error) {
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
