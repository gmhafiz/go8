package app

import (
	"net/http"

	"github.com/go-chi/chi"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()

	r.Get("/", s.handleIndex())

	r.Route("/admin", func(r chi.Router) {
		r.Use(s.AdminOnlyHandler)

		r.Get("/", s.handleAdminIndex())
	})

	r.Route("/api", func(r chi.Router) {
		r.Use(s.ContentTypeJsonHandler)

		r.Route("/v1", func(r chi.Router) {
			r.Get("/", s.getAllContact())
			r.Get("/something", s.handleSomething())
		})
	})

	return r
}




