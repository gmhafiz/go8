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

// Alive API
// @Summary Checks if API is up
// @Description Hits this API to see if API is running in the server
// @Success 200
// @Failure 500
// @Router /health/liveness [get]
func (h *Handler) Liveness(w http.ResponseWriter, r *http.Request) {
	respond.Render(w, http.StatusOK, nil)
}

// Alive Database
// @Summary Checks if both API and Database are up
// @Description Hits this API to see if both API and Database are running in the server
// @Success 200
// @Failure 500
// @Router /health/readiness [get]
func (h *Handler) Readiness(w http.ResponseWriter, r *http.Request) {
	err := h.useCase.Readiness()
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err)
		return
	}
	respond.Render(w, http.StatusOK, nil)
}
