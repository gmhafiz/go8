package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/volatiletech/null/v8"

	"eight/internal/models"
	"eight/internal/util/converter"
	"eight/pkg/validation"
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

type AuthorResponse struct {
	AuthorID   int64       `json:"author_id"`
	FirstName  string      `json:"first_name"`
	MiddleName null.String `json:"middle_name"`
	LastName   string      `json:"last_name"`
	Rel        interface{} `json:"|rel|"`
}

type Book struct {
	BookID        int64       `json:"book_id"`
	Title         string      `json:"title"`
	PublishedDate time.Time   `json:"published_date"`
	ImageURL      null.String `json:"image_url"`
	Description   null.String `json:"description"`
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

		books, err := converter.DeepCopy(author.R)
		if err != nil {
			h.Logger.Error().Msg(err.Error())
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{
				"error": err.Error(),
			})
			return
		}

		resp := &AuthorResponse{
			AuthorID:   author.AuthorID,
			FirstName:  author.FirstName,
			MiddleName: author.MiddleName,
			LastName:   author.LastName,
			Rel:        books,
		}

		render.JSON(w, r, resp)
		return
	}
}
