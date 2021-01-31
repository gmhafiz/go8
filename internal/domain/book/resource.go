package book

import (
	"reflect"
	"time"

	"github.com/jinzhu/copier"
	"github.com/jinzhu/now"
	"github.com/volatiletech/null/v8"

	"github.com/gmhafiz/go8/internal/models"
)

type Request struct {
	BookID        string `json:"-"`
	Title         string `json:"title" validate:"required"`
	PublishedDate string `json:"published_date" validate:"required"`
	ImageURL      string `json:"image_url" validate:"url"`
	Description   string `json:"description" validate:"required"`
}

type Resource struct {
	BookID        int64       `json:"book_id" deepcopier:"field:book_id" db:"id"`
	Title         string      `json:"title" deepcopier:"field:title" db:"title"`
	PublishedDate time.Time   `json:"published_date" deepcopier:"field:force" db:"published_date"`
	ImageURL      null.String `json:"image_url" deepcopier:"field:image_url" db:"image_url"`
	Description   null.String `json:"description" deepcopier:"field:description"`
}

func ToBook(req *Request) *models.Book {
	return &models.Book{
		Title:         req.Title,
		PublishedDate: now.MustParse(req.PublishedDate),
		ImageURL: null.String{
			String: req.ImageURL,
			Valid:  true,
		},
		Description: req.Description,
	}
}

func Book(book *models.Book) (Resource, error) {
	var resource Resource

	err := copier.Copy(&resource, &book)
	if err != nil {
		return resource, err
	}

	return resource, nil
}

func Books(books []*models.Book) (interface{}, error) {
	var resource Resource

	if len(books) == 0 {
		return make([]string, 0), nil
	}

	rt := reflect.TypeOf(books)
	if rt.Kind() == reflect.Slice {
		var resources []Resource
		for _, book := range books {
			res, _ := Book(book)
			resources = append(resources, res)
		}
		return resources, nil
	}

	err := copier.Copy(&resource, books)
	if err != nil {
		return resource, err
	}

	return resource, nil
}
