package author

import (
	"github.com/gmhafiz/go8/internal/utility/respond"
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/gmhafiz/go8/internal/domain/author"
)

type HTTP interface {
	Create(w http.ResponseWriter, r *http.Request)
	List(w http.ResponseWriter, r *http.Request)
	Read(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type Handler struct {
	useCase  author.UseCase
	validate *validator.Validate
}

func NewHandler(useCase author.UseCase) *Handler {
	return &Handler{
		useCase:  useCase,
		validate: validator.New(),
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	filters := author.Filters(r.URL.Query())

	authors, err := h.useCase.List(r.Context(), filters)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err)
		return
	}

	list, err := author.Resources(authors)

	respond.Render(w, http.StatusOK, list)
	return
}

func (h *Handler) Read(w http.ResponseWriter, r *http.Request) {
	authorID := respond.GetURLParamInt64(w, r, "id")

	b, err := h.useCase.ReadWithBooks(r.Context(), uint64(authorID))
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, nil)
		return
	}

	respond.Render(w, http.StatusOK, b)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}
