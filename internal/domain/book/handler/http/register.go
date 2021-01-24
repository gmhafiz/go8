package http

import (
	"github.com/go-chi/chi"

	"github.com/gmhafiz/go8/internal/domain/book"
	"github.com/gmhafiz/go8/internal/middleware"
)

func RegisterHTTPEndPoints(router *chi.Mux, uc book.UseCase) {
	h := NewHandler(uc)

	router.Route("/api/v1/books", func(router chi.Router) {
		router.Use(middleware.Json)
		router.With(middleware.Paginate).Get("/", h.All)
		router.Post("/", h.Create)
		router.Put("/{bookID}", h.Update)
		router.Delete("/{bookID}", h.Delete)
	})
}
