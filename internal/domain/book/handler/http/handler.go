package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/volatiletech/null/v8"

	"github.com/gmhafiz/go8/internal/domain/book"
	"github.com/gmhafiz/go8/internal/resource"
)

type Handler struct {
	useCase  book.UseCase
	validate *validator.Validate
}

func NewHandler(useCase book.UseCase) *Handler {
	return &Handler{
		useCase:  useCase,
		validate: validator.New(),
	}
}

type BookRequest struct {
	Title         string `json:"title" validate:"required"`
	PublishedDate string `json:"published_date" validate:"required"`
	ImageURL      string `json:"image_url" validate:"url"`
	Description   string `json:"description" validate:"required"`
}

type BookResource struct {
	BookID        int64       `json:"book_id" deepcopier:"field:book_id" db:"id"`
	Title         string      `json:"title" deepcopier:"field:title" db:"title"`
	PublishedDate time.Time   `json:"published_date" deepcopier:"field:force" db:"published_date"`
	ImageURL      null.String `json:"image_url" deepcopier:"field:image_url" db:"image_url"`
	Description   null.String `json:"description" deepcopier:"field:description"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var bookRequest resource.BookRequest
	err := json.NewDecoder(r.Body).Decode(&bookRequest)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, nil)
		return
	}

	err = h.validate.Struct(bookRequest)
	if err != nil {
		// this check is only needed when your code could produce
		// an invalid value for validation such as interface with nil
		// value most including myself do not usually have code like this.
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
			return
		}

		var errs []string
		for _, err := range err.(validator.ValidationErrors) {
			errs = append(errs, fmt.Sprintf("%s is %s", err.StructNamespace(), err.Tag()))
		}
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string][]string{"errors": errs})
		return
	}

	bk, err := h.useCase.Create(context.Background(), bookRequest.Title, bookRequest.Description, bookRequest.ImageURL, bookRequest.PublishedDate)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, nil)
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, bk)
}

func (h *Handler) All(w http.ResponseWriter, r *http.Request) {
	resp, err := h.useCase.All(r.Context())
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, err.Error())
		return
	}

	list, err := resource.Books(resp)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, err.Error())
		return
	}

	render.JSON(w, r, list)
}
