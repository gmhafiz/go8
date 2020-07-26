package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"eight/internal/models"
)

func (a API) GetAllAuthors() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		authors, err := a.authors.AllAuthors(ctx)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			_ = render.Render(w, r, nil)
			return
		}

		render.JSON(w, r, authors)
	}
}

func (a API) CreateAuthor() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var authorRequest models.Author
		ctx := context.Background()

		author := &models.Author{
			FirstName:  authorRequest.FirstName,
			MiddleName: authorRequest.MiddleName,
			LastName:   authorRequest.LastName,
		}

		err := json.NewDecoder(r.Body).Decode(author)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			return
		}
		createdAuthor, err := a.authors.CreateAuthor(ctx, author)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, createdAuthor)
	}
}

func (a API) GetAuthor() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authorID := chi.URLParam(r, "authorID")
		ctx := context.Background()

		id, _ := strconv.ParseInt(authorID, 10, 64)

		author, err := a.authors.GetAuthor(ctx, id)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			return
		}

		render.JSON(w, r, author)
	}
}
