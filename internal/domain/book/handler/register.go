package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"

	"github.com/gmhafiz/go8/internal/domain/book/usecase"
)

func RegisterHTTPEndPoints(router *chi.Mux, validator *validator.Validate, uc usecase.Book) *Handler {
	h := NewHandler(uc, validator)

	router.Route("/api/v1/book", func(router chi.Router) {
		router.Get("/", h.List)
		router.Get("/{bookID}", h.Get)
		router.Post("/", h.Create)
		router.Put("/{bookID}", h.Update)
		router.Delete("/{bookID}", h.Delete)
	})
	return h
}
