package books

import (
	"github.com/go-chi/chi"

	"go8ddd/internal/middleware"
)

func initRoutes(router *chi.Mux, handler *BookHandler) {
	router.Route("/api/v1/books", func(router chi.Router) {
		router.With(middleware.Paginate).Get("/", handler.All())
		router.Post("/", handler.Create())
		router.With(middleware.IDParam).Get("/{id}", handler.Get())
		router.With(middleware.IDParam).Put("/{id}", handler.Update())
		router.With(middleware.IDParam).Delete("/{id}", handler.Delete())
	})
}
