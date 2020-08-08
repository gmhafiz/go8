package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"eight/internal/models"
)

func (a API) GetAllAuthors() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authors, err := a.authors.AllAuthors(r.Context())
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
		createdAuthor, err := a.authors.CreateAuthor(r.Context(), author)
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

		id, _ := strconv.ParseInt(authorID, 10, 64)

		author, err := a.authors.GetAuthor(r.Context(), id)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			return
		}

		render.JSON(w, r, author)
	}
}
