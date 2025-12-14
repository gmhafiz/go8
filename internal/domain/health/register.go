package health

import (
	"github.com/go-chi/chi/v5"

	"github.com/gmhafiz/go8/internal/middleware"
)

func RegisterHTTPEndPoints(router *chi.Mux, uc UseCase) *Handler {
	h := NewHandler(uc)

	router.Route("/api/health", func(router chi.Router) {
		router.Use(middleware.JSON)

		router.Get("/", h.Health)
		router.Get("/readiness", h.Readiness)
	})

	return h
}
