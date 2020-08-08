package http

import (
	"net/http"

	"eight/internal/middleware"
	"github.com/go-chi/chi"
	"github.com/go-chi/httplog"
	"github.com/rs/zerolog"
)

//func Router(h *Handlers, logger log.Logger) *chi.Mux {
//func Router(h *Handlers, logger *zap.Logger) *chi.Mux {
func Router(h *Handlers, logger zerolog.Logger) *chi.Mux {
	r := chi.NewRouter()

	r.Use(httplog.RequestLogger(logger))

	r.Use(middleware.Cors)

	//r.Use(middleware.RequestLog(logger))

	//r.Get("/", h.HandleLive())
	r.Get("/health/liveness", h.HandleLive())
	r.Get("/health/readiness", h.HandleReady())

	r.Route("/admin", func(r chi.Router) {
		r.Use(middleware.AdminOnlyHandler)
		//r.Get("/", s.handleAdminIndex())
	})

	r.Route("/api", func(r chi.Router) {
		//r.Use(s.ContentTypeJsonHandler)
		//r.Use(middleware.AdminOnlyHandler)

		r.Route("/v1", func(r chi.Router) {
			r.With(middleware.Paginate).Get("/books", h.GetAllBooks())
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

func PrintAllRegisteredRoutes(router *chi.Mux, logger zerolog.Logger) {
	walkFunc := func(method string, path string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		logger.Info().
			Dict("routes", zerolog.Dict().Str("method", method).Str("path",
			path)).Msg("")

		return nil
	}
	if err := chi.Walk(router, walkFunc); err != nil {
		logger.Err(err)
		//logger.Error().Msgf("Logging err %s", err)
	}
}
