package book

import (
	"reflect"
	"strconv"
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

type DB struct {
	BookID        int64       `db:"book_id"`
	Title         string      `db:"title"`
	PublishedDate time.Time   `db:"published_date"`
	ImageURL      null.String `db:"image_url"`
	Description   string      `db:"description"`
	CreatedAt     null.Time   `db:"created_at"`
	UpdatedAt     null.Time   `db:"updated_at"`
	DeletedAt     null.Time   `db:"deleted_at"`
}

func ToBook(req *Request) *models.Book {
	id, err := strconv.ParseInt(req.BookID, 10, 64)
	if err != nil {
		return nil
	}
	return &models.Book{
		BookID:        id,
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
