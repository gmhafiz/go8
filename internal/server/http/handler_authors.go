package http

import (
	"eight/internal/models"
	"eight/pkg/validation"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/volatiletech/null/v8"
	"net/http"
	"strconv"
)

type authorRequest struct {
	FirstName  string `json:"first_name" validate:"required"`
	MiddleName string `json:"middle_name" validate:""`
	LastName   string `json:"last_name" validate:"required"`
}

func (h *Handlers) GetAllAuthors() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authors, err := h.Api.GetAllAuthors(r.Context())

		if err != nil {
			render.Status(r, http.StatusBadRequest)
			return
		}

		render.JSON(w, r, authors)
	}
}

func (h *Handlers) CreateAuthor() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var authorRequest authorRequest

		err := json.NewDecoder(r.Body).Decode(&authorRequest)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{
				"error": err.Error(),
			})
			return
		}

		validationErrors := validation.Validate(h.Validation, authorRequest)
		if validationErrors != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string][]string{"error": validationErrors})
			return
		}
		
		createdAuthor, err := h.Api.CreateAuthor(r.Context(), &models.Author{
			FirstName: authorRequest.FirstName,
			MiddleName: null.String{
				String: authorRequest.MiddleName,
				Valid:  true,
			},
			LastName: authorRequest.LastName,
		})
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{
				"error": err.Error(),
			})
			return
		}

		render.JSON(w, r, createdAuthor)
	}
}

func (h *Handlers) GetAuthor() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authorID := chi.URLParam(r, "authorID")

		id, _ := strconv.ParseInt(authorID, 10, 64)

		author, err := h.Api.GetAuthor(r.Context(), id)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{
				"error": err.Error(),
			})
			return
		}

		render.JSON(w, r, author)
		return
	}
}
