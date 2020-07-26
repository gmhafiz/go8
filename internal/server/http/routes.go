package http

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"

	"eight/internal/middleware"
)

func Router(h *Handlers) *chi.Mux {
	r := chi.NewRouter()

	//r.Get("/", h.HandleLive())
	//r.Get("/healthz/liveness", h.HandleLive())
	//r.Get("/healthz/readiness", h.HandleReady())

	r.Route("/admin", func(r chi.Router) {
		r.Use(middleware.AdminOnlyHandler)
		//r.Get("/", s.handleAdminIndex())
	})

	r.Route("/api", func(r chi.Router) {
		//r.Use(s.ContentTypeJsonHandler)
		r.Use(middleware.AdminOnlyHandler)

		r.Route("/v1", func(r chi.Router) {
			r.Get("/books", h.GetAllBooks())
			r.Post("/book", h.CreateBook())
			r.Get("/book/{bookID}", h.GetBook())
			r.Delete("/book/{bookID}", h.Delete())

			//r.Get("/authors", a.GetAllAuthors())
			//r.Post("/author", a.CreateAuthor())
			//r.Get("/author/{id}", a.GetAuthor())
		})
	})

	return r
}

func PrintAllRegisteredRoutes(router *chi.Mux) {
	log.Println("All Registered Routes:")
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("%s %s\n", method, route)
		return nil
	}
	if err := chi.Walk(router, walkFunc); err != nil {
		log.Printf("Logging err: %s\n", err.Error())
	}
}
