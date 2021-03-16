package book

import (
	"encoding/json"
	"io"
	"time"

	"github.com/jinzhu/copier"
	"github.com/jinzhu/now"
	"github.com/volatiletech/null/v8"

	"github.com/gmhafiz/go8/internal/models"
	"github.com/gmhafiz/go8/internal/utility/filter"
)

type Request struct {
	BookID        int64  `json:"-"`
	Title         string `json:"title" validate:"required"`
	PublishedDate string `json:"published_date" validate:"required"`
	ImageURL      string `json:"image_url" validate:"url"`
	Description   string `json:"description" validate:"required"`
	filter.Filter
}

type Res struct {
	BookID        int64       `json:"book_id" deepcopier:"field:book_id" db:"id"`
	Title         string      `json:"title" deepcopier:"field:title" db:"title"`
	PublishedDate time.Time   `json:"published_date" deepcopier:"field:force" db:"published_date"`
	ImageURL      null.String `json:"image_url" deepcopier:"field:image_url" db:"image_url"`
	Description   null.String `json:"description" deepcopier:"field:description"`
}

func Decode(body io.ReadCloser, b *Request) error {
	return json.NewDecoder(body).Decode(b)
}

func ToBook(req *Request) *models.Book {
	return &models.Book{
		BookID:        req.BookID,
		Title:         req.Title,
		PublishedDate: now.MustParse(req.PublishedDate),
		ImageURL: null.String{
			String: req.ImageURL,
			Valid:  true,
		},
		Description: req.Description,
	}
}

func Resource(book *models.Book) (Res, error) {
	var resource Res

	err := copier.Copy(&resource, &book)
	if err != nil {
		return resource, err
	}

	return resource, nil
}

func Resources(books []*models.Book) (interface{}, error) {
	if len(books) == 0 {
		return make([]string, 0), nil
	}

	var resources []Res
	for _, book := range books {
		res, _ := Resource(book)
		resources = append(resources, res)
	}
	return resources, nil
}
