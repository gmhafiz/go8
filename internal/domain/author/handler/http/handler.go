package author

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/gmhafiz/go8/internal/domain/author"
	"github.com/gmhafiz/go8/internal/domain/author/usecase"
	"github.com/gmhafiz/go8/internal/utility/param"
	"github.com/gmhafiz/go8/internal/utility/respond"
	"github.com/gmhafiz/go8/internal/utility/validate"
)

type HTTP interface {
	Create(w http.ResponseWriter, r *http.Request)
	List(w http.ResponseWriter, r *http.Request)
	Read(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type Handler struct {
	useCase  usecase.UseCase
	validate *validator.Validate
}

func NewHandler(useCase usecase.UseCase, v *validator.Validate) *Handler {
	return &Handler{
		useCase:  useCase,
		validate: v,
	}
}

// Create creates a new author
// @Summary Create an Author
// @Description Create an author using JSON payload
// @Accept json
// @Produce json
// @Param Author body author.CreateRequest true "Create an author using the following format"
// @Success 201 {object} author.CreateResponse
// @Failure 400 {string} Bad Request
// @Failure 500 {string} Internal Server Error
// @router /api/v1/author [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req author.CreateRequest
	err := req.Bind(r.Body)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, err)
		return
	}

	errs := validate.Validate(h.validate, req)
	if errs != nil {
		respond.Errors(w, http.StatusBadRequest, errs)
		return
	}

	create, err := h.useCase.Create(r.Context(), req)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err)
		return
	}

	respond.Json(w, http.StatusCreated, create)
}

// List will fetch the authors based on given params
// @Summary Shows all authors
// @Description Lists all authors. By default, it gets first page with 30 items.
// @Accept json
// @Produce json
// @Param page query string false "page number"
// @Param size query string false "size of result"
// @Param name query string false "search by name"
// @Success 200 {object} respond.Standard
// @Failure 500 {string} Internal Server Error
// @router /api/v1/author [get]
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	filters := author.Filters(r.URL.Query())

	// For cache purpose, we use request URI as the key for our result.
	// We save it into context so that we can pick it pick in our cache layer.
	ctx := context.WithValue(r.Context(), author.CacheURL, r.URL.String())

	authors, total, err := h.useCase.List(ctx, filters)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err)
		return
	}

	list, err := author.Resources(authors)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err)
		return
	}

	respond.Json(w, http.StatusOK, respond.Standard{
		Data: list,
		Meta: respond.Meta{
			Size:  len(authors),
			Total: total,
		},
	})
}

// Get an author by its ID
// @Summary Get an Author
// @Description Get an author by its id.
// @Accept json
// @Produce json
// @Param id path int true "author ID"
// @Success 200 {object} author.ResWithBooks
// @Failure 400 {string} Bad Request
// @Failure 500 {string} Internal Server Error
// @router /api/v1/author/{id} [get]
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	authorID := param.Int64(w, r, "id")

	ctx := context.WithValue(r.Context(), author.CacheURL, r.URL.String())

	a, err := h.useCase.ReadWithBooks(ctx, uint64(authorID))
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err)
		return
	}

	res := author.ResourceWithBooks(a)

	respond.Render(w, http.StatusOK, res)
}

// Update an author
// @Summary Update n Author
// @Description Update an author by its model.
// @Accept json
// @Produce json
// @Param Author body author.Update true "Author Request"
// @Success 200 {object} []author.Res
// @Failure 400 {string} Bad Request
// @Failure 500 {string} Internal Server Error
// @router /api/v1/author/{id} [put]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := param.Int64(w, r, "id")
	if id == 0 {
		respond.Error(w, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	ctx := context.WithValue(r.Context(), author.CacheURL, r.URL.String())

	var req author.Update
	err := req.Bind(r.Body)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, err)
		return
	}
	req.ID = id

	updated, err := h.useCase.Update(ctx, &req)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err)
		return
	}

	res, err := author.ResourceUpdate(updated)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err)
	}

	respond.Render(w, http.StatusOK, res)
}

// Delete an author by its ID
// @Summary Delete an Author
// @Description Delete an author by its id.
// @Accept json
// @Produce json
// @Param id path int true "author ID"
// @Success 200 "Ok"
// @Failure 500 {string} Internal Server Error
// @router /api/v1/author/{id} [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := param.Int64(w, r, "id")
	if id == 0 {
		respond.Error(w, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	ctx := context.WithValue(r.Context(), author.CacheURL, r.URL.String())

	err := h.useCase.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, respond.ErrNoRecord) {
			respond.Error(w, http.StatusBadRequest, err)
			return
		}
		respond.Error(w, http.StatusInternalServerError, err)
		return
	}
}
