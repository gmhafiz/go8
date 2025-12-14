package health

import (
	"net/http"

	"github.com/gmhafiz/go8/internal/utility/respond"
)

type Handler struct {
	useCase UseCase
}

func NewHandler(useCase UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

// Health checks if api is up
// @Summary Checks if API is up
// @Description Hits this API to see if API is running in the server
// @Success 200
// @Failure 500
// @router /api/health [get]
func (h *Handler) Health(w http.ResponseWriter, _ *http.Request) {
	respond.JSON(w, http.StatusOK, map[string]int{"status": 200})
}

// Readiness checks if database is alive
// @Summary Checks if both API and Database are up
// @Description Hits this API to see if both API and Database are running in the server
// @Success 200
// @Failure 500
// @router /api/health/readiness [get]
func (h *Handler) Readiness(w http.ResponseWriter, _ *http.Request) {
	err := h.useCase.Readiness()
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err)
		return
	}
	respond.JSON(w, http.StatusOK, map[string]int{"status": 200})
}
