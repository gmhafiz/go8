package handler

import (
	"encoding/json"
	"github.com/gmhafiz/go8/internal/domain/author"
	"io"
)

type request struct {
	author.Filter
	FirstName  string `json:"first_name" validate:"required"`
	MiddleName string `json:"middle_name"`
	LastName   string `json:"last_name" validate:"required"`
	Books      []Book `json:"books"`
}

type CreateRequest struct {
	FirstName  string `json:"first_name" validate:"required"`
	MiddleName string `json:"middle_name"`
	LastName   string `json:"last_name" validate:"required"`
	Books      []Book `json:"books"`
}

type Book struct {
	BookID        int    `json:"id"`
	Title         string `json:"title" validate:"required"`
	PublishedDate string `json:"published_date" validate:"required"`
	Description   string `json:"description" validate:"required"`
}

type Update struct {
	ID         int    `json:"id"`
	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name,omitempty"`
	LastName   string `json:"last_name"`
}

func (r *Update) Bind(body io.ReadCloser) error {
	return json.NewDecoder(body).Decode(r)
}

func (r *CreateRequest) Bind(body io.ReadCloser) error {
	return json.NewDecoder(body).Decode(r)
}

func (r *request) Bind(body io.ReadCloser) error {
	return json.NewDecoder(body).Decode(r)
}
