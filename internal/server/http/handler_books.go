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
)

func (h *Handlers) GetAllBooks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		books, err := h.Api.GetAllBooks(r.Context())

		if err != nil {
			render.Status(r, http.StatusBadRequest)
			_ = render.Render(w, r, nil)
			return
		}

		render.JSON(w, r, books)
	}
}

func (h *Handlers) CreateBook() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type bookRequest struct {
			Title         string      `json:"title"`
			PublishedDate string      `json:"published_date"`
			ImageURL      null.String `json:"image_url"`
			Description   null.String `json:"description"`
		}
		var bookR bookRequest

		err := json.NewDecoder(r.Body).Decode(&bookR)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			return
		}

		time, err := converter.StringToTime(h.TimeConverter, bookR.PublishedDate)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{
				"error": err.Error(),
			})
			return
		}
		book := &models.Book{
			Title:         bookR.Title,
			PublishedDate: time,
			ImageURL:      bookR.ImageURL,
			Description:   bookR.Description,
		}

		createdBook, err := h.Api.CreateBook(r.Context(), book)
		if err == nil {
			h.Logger.Error().Err(err)
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"error": "StatusInternalServerError"})
			return
		}

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, createdBook)
	}
}

func (h *Handlers) GetBook() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bookID := chi.URLParam(r, "bookID")

		id, _ := strconv.ParseInt(bookID, 10, 64)

		book, err := h.Api.GetBook(r.Context(), id)
		if err != nil {

			if errors.As(err, &sql.ErrNoRows) {
				h.Logger.Error().Err(err)
				render.JSON(w, r, "no book found")
				render.Status(r, http.StatusBadRequest)
			} else {
				render.JSON(w, r, err.Error())
				render.Status(r, http.StatusInternalServerError)
			}
			return
		}

		h.Logger.Info().Msgf("here")
		render.JSON(w, r, book)
	}
}

func (h *Handlers) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bookID := chi.URLParam(r, "bookID")
		id, _ := strconv.ParseInt(bookID, 10, 64)

		err := h.Api.Delete(r.Context(), id)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			return
		}

		render.Status(r, http.StatusOK)
	}
}
