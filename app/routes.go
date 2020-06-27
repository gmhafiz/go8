package app

import (
	"net/http"

	"github.com/go-chi/chi"
)

func (s *Server) Router() http.Handler {
	r := chi.NewRouter()

	r.Get("/", s.handleIndex())
	r.Get("/healthz/liveness", s.handleLive())
	r.Get("/healthz/readiness", s.handleReady())

	r.Route("/admin", func(r chi.Router) {
		r.Use(s.AdminOnlyHandler)

		r.Get("/", s.handleAdminIndex())
	})

	r.Route("/api", func(r chi.Router) {
		r.Use(s.ContentTypeJsonHandler)
		r.Use(s.AdminOnlyHandler)

		r.Route("/v1", func(r chi.Router) {
			r.Get("/books", s.getAllBooks())
			r.Get("/book/:id", s.getBook())
		})
	})

	return r
}




