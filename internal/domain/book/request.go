package book

import (
	"encoding/json"
	"github.com/gmhafiz/go8/internal/models"
	"github.com/jinzhu/now"
	"github.com/volatiletech/null/v8"
	"io"
)

type Request struct {
	BookID        int64  `json:"-"`
	Title         string `json:"title" validate:"required"`
	PublishedDate string `json:"published_date" validate:"required"`
	ImageURL      string `json:"image_url" validate:"url"`
	Description   string `json:"description" validate:"required"`
	Filter
}

func Bind(body io.ReadCloser, b *Request) error {
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
