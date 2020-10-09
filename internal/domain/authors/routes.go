package authors

import (
	"github.com/go-chi/chi"

	"go8ddd/internal/middleware"
)

func initRoutes(router *chi.Mux, handler *AuthorHandler) {
	router.Route("/api/v1/authors", func(router chi.Router) {
		router.With(middleware.Paginate).Get("/", handler.All())
	})
}
