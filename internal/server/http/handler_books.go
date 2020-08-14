package http

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/friendsofgo/errors"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/volatiletech/null/v8"

	"eight/internal/models"
	"eight/internal/util/converter"
	"eight/pkg/validation"
)

type bookRequest struct {
	Title         string      `json:"title" validate:"required"`
	PublishedDate string      `json:"published_date" validate:"required"`
	ImageURL      null.String `json:"image_url" validate:"url"`
	Description   null.String `json:"description" validate:"required"`
}

// GetAllBooks godoc
// @Summary Show all books
// @Description Get all books. By default it gets first page with 10 items.
// @Accept json
// @Produce json
// @Param page query string false "page number"
// @Param size query string false "size"
// @Success 200 {object} []models.Book
// @Router /books [get]
func (h *Handlers) GetAllBooks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		books, err := h.Api.GetAllBooks(r.Context())

		if err != nil {
			render.Status(r, http.StatusBadRequest)
			return
		}

		render.JSON(w, r, books)
	}
}

// GetBook godoc
// @Summary Create a Book
// @Description Get a book with JSON payload
// @Accept json
// @Produce json
// @Param Book body bookRequest true "Book Request"
// @Success 201 {object} models.Book
// @Router /book [post]
func (h *Handlers) CreateBook() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var bookR bookRequest

		err := json.NewDecoder(r.Body).Decode(&bookR)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			return
		}

		validationErrors := validation.Validate(h.Validation, bookR)
		if validationErrors != nil {
			render.JSON(w, r, map[string][]string{"error": validationErrors})
			return
		}

		time, err := converter.StringToTime(bookR.PublishedDate)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{
				"error": err.Error(),
			})
			return
		}

		createdBook, err := h.Api.CreateBook(r.Context(), &models.Book{
			Title:         bookR.Title,
			PublishedDate: time,
			ImageURL:      bookR.ImageURL,
			Description:   bookR.Description,
		})

		if err != nil {
			h.Logger.Error().Err(err)
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"error": "StatusInternalServerError"})
			return
		}

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, createdBook)
	}
}

// GetBook godoc
// @Summary Get a Book
// @Description Get a book by its id.
// @Accept json
// @Produce json
// @Param id path int true "book ID"
// @Success 200 {object} models.Book
// @Router /book/{bookID} [get]
func (h *Handlers) GetBook() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bookID := chi.URLParam(r, "bookID")

		id, _ := strconv.ParseInt(bookID, 10, 64)

		book, err := h.Api.GetBook(r.Context(), id)
		if err != nil {

			if errors.As(err, &sql.ErrNoRows) {
				//h.Logger.With().Caller().Logger()
				h.Logger.Error().Msg(err.Error())
				render.Status(r, http.StatusBadRequest)
				render.JSON(w, r, map[string]string{"error": err.Error()})
			} else {
				render.Status(r, http.StatusInternalServerError)
				render.JSON(w, r, err.Error())
			}
			return
		}

		render.JSON(w, r, book)
		return
	}
}

// GetBook godoc
// @Summary Delete a Book
// @Description Delete a book by its id.
// @Accept json
// @Produce json
// @Param id path int true "book ID"
// @Success 200 "Ok"
// @Failure 400 "Bad request"
// @Failure 500 "Internal Server error"
// @Router /book/{bookID} [delete]
func (h *Handlers) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bookID := chi.URLParam(r, "bookID")
		id, _ := strconv.ParseInt(bookID, 10, 64)

		err := h.Api.Delete(r.Context(), id)

		if errors.As(err, &sql.ErrNoRows) {
			h.Logger.Error().Err(err)
			render.JSON(w, r, map[string]string{"error": err.Error()})
			render.Status(r, http.StatusBadRequest)
		} else {
			render.JSON(w, r, map[string]string{"error": "internal server error"})
			render.Status(r, http.StatusInternalServerError)
		}

		render.Status(r, http.StatusOK)
		return
	}
}
