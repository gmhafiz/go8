package http

import (
	"github.com/go-chi/chi/v5"

	"github.com/gmhafiz/go8/internal/domain/health"
	"github.com/gmhafiz/go8/internal/middleware"
)

func RegisterHTTPEndPoints(router *chi.Mux, uc health.UseCase) {
	h := NewHandler(uc)

	router.Route("/api/health", func(router chi.Router) {
		router.Use(middleware.Json)

		router.Get("/", h.Health)
		router.Get("/readiness", h.Readiness)
	})
}
