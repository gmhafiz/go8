package author

import (
	"github.com/go-chi/chi/v5"

	"github.com/gmhafiz/go8/internal/domain/author"
	"github.com/gmhafiz/go8/internal/middleware"
)

func RegisterHTTPEndPoints(router *chi.Mux, uc author.UseCase) {
	h := NewHandler(uc)

	router.Route("/api/v1/author", func(router chi.Router) {
		router.Use(middleware.Json)
		router.Post("/", h.Create)
		router.Get("/", h.List)
		router.Get("/{id}", h.Read)
		router.Put("/{id}", h.Update)
		router.Delete("/{id}", h.Delete)
	})
}
