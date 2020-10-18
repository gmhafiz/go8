package health

import (
	"github.com/go-chi/chi"
)

func initRoutes(router *chi.Mux, h *Handler) {
	router.Route("/health", func(router chi.Router) {
		router.Get("/liveness", h.Liveness())
		router.Get("/readiness", h.Readiness())
	})
}