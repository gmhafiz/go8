package http

import (
	"github.com/gmhafiz/go8/internal/utility"
	"net/http"

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
	utility.Render(w, http.StatusOK, nil)
}

func (h *Handler) Readiness(w http.ResponseWriter, r *http.Request) {
	err := h.useCase.Readiness()
	if err != nil {
		utility.Render(w, http.StatusInternalServerError, nil)
		return
	}
	utility.Render(w, http.StatusOK, nil)
}
