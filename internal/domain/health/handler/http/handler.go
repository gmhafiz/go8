package http

import (
	"net/http"

	"github.com/gmhafiz/go8/internal/domain/health"
	"github.com/gmhafiz/go8/internal/utility/presentation"
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
	presentation.Render(w, http.StatusOK, nil)
}

func (h *Handler) Readiness(w http.ResponseWriter, r *http.Request) {
	err := h.useCase.Readiness()
	if err != nil {
		presentation.Render(w, http.StatusInternalServerError, nil)
		return
	}
	presentation.Render(w, http.StatusOK, nil)
}
