package http

import (
	"net/http"

	"github.com/go-chi/render"

	"github.com/gmhafiz/go8/internal/domain/health"
)

type Handler struct {
	useCase health.UseCase
}

func NewHandler(useCase health.UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

func (h *Handler) Liveness(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusOK)
}

func (h *Handler) Readiness(w http.ResponseWriter, r *http.Request) {
	err := h.useCase.Readiness()
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		return
	}
	render.Status(r, http.StatusOK)
}
