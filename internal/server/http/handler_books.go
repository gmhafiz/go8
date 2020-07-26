package http

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/friendsofgo/errors"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	
	"eight/internal/models"
)

func (h *Handlers) GetAllBooks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		books, err := h.Api.GetAllBooks()

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
		var bookRequest models.Book
		ctx := context.Background()

		book := &models.Book{
			Title:       bookRequest.Title,
			Description: bookRequest.Description,
		}

		err := json.NewDecoder(r.Body).Decode(book)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			return
		}
		createdBook, err := h.Api.CreateBook(ctx, book)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, createdBook)
	}
}

func (h *Handlers) GetBook() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bookID := chi.URLParam(r, "bookID")
		ctx := context.Background()

		id, _ := strconv.ParseInt(bookID, 10, 64)

		book, err := h.Api.GetBook(ctx, id)
		if err != nil {

			if errors.As(err, &sql.ErrNoRows) {
				render.JSON(w, r, "no book found")
				render.Status(r, http.StatusBadRequest)
			} else {
				render.JSON(w, r, err.Error())
				render.Status(r, http.StatusInternalServerError)
			}
			return
		}

		render.JSON(w, r, book)
	}
}

func (h *Handlers) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bookID := chi.URLParam(r, "bookID")
		id, _ := strconv.ParseInt(bookID, 10, 64)

		ctx := context.Background()

		err := h.Api.Delete(ctx, id)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			return
		}

		render.Status(r, http.StatusOK)
	}
}
