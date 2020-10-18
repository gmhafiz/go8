package books

import (
	"github.com/go-chi/chi"

	"go8ddd/internal/middleware"
)

func initRoutes(router *chi.Mux, h *Handler) {
	router.Route("/api/v1/books", func(router chi.Router) {
		router.With(middleware.Paginate).Get("/", h.All())
		router.Post("/", h.Create())
		router.With(middleware.IDParam).Get("/{id}", h.Get())
		router.With(middleware.IDParam).Put("/{id}", h.Update())
		router.With(middleware.IDParam).Delete("/{id}", h.Delete())
	})
}
