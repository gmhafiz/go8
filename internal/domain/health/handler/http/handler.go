package http

import (
	"net/http"

	"github.com/gmhafiz/go8/internal/domain/health"
	"github.com/gmhafiz/go8/internal/utility/respond"
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
	respond.Render(w, http.StatusOK, nil)
}

func (h *Handler) Readiness(w http.ResponseWriter, r *http.Request) {
	err := h.useCase.Readiness()
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err)
		return
	}
	respond.Render(w, http.StatusOK, nil)
}
