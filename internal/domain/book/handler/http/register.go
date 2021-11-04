package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"

	"github.com/gmhafiz/go8/internal/domain/book"
	"github.com/gmhafiz/go8/internal/middleware"
)

func RegisterHTTPEndPoints(router *chi.Mux, validator *validator.Validate, uc book.UseCase) *Handler {
	h := NewHandler(uc, validator)

	router.Route("/api/v1/book", func(router chi.Router) {
		router.Use(middleware.Json)

		router.Get("/", h.List)
		router.Get("/{bookID}", h.Get)
		router.Post("/", h.Create)
		router.Put("/{bookID}", h.Update)
		router.Delete("/{bookID}", h.Delete)
	})
	return h
}
