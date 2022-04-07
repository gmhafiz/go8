package book

import (
	"encoding/json"
	"io"
)

type CreateRequest struct {
	Title         string `json:"title" validate:"required"`
	PublishedDate string `json:"published_date" validate:"required"`
	ImageURL      string `json:"image_url" validate:"url"`
	Description   string `json:"description" validate:"required"`
}

type UpdateRequest struct {
	ID            int    `json:"-"`
	Title         string `json:"title" validate:"required"`
	PublishedDate string `json:"published_date" validate:"required"`
	ImageURL      string `json:"image_url" validate:"url"`
	Description   string `json:"description" validate:"required"`
}

func Bind(body io.ReadCloser, b *CreateRequest) error {
	return json.NewDecoder(body).Decode(b)
}
