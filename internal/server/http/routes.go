package http

import (
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/httplog"
	"github.com/rs/zerolog"

	_ "eight/docs"
	"eight/internal/middleware"
)

func Router(h *Handlers, logger zerolog.Logger) *chi.Mux {
	r := chi.NewRouter()

	r.Use(httplog.RequestLogger(logger))
	r.Use(middleware.Cors)

	r.Get("/health/liveness", h.HandleLive())
	r.Get("/health/readiness", h.HandleReady())

	SwaggerServer(r)

	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.With(middleware.Paginate).Get("/books", h.GetAllBooks())
			r.Post("/book", h.CreateBook())
			r.Get("/book/{bookID}", h.GetBook())
			r.Delete("/book/{bookID}", h.Delete())

			r.With(middleware.Paginate).Get("/authors", h.GetAllAuthors())
			r.Post("/author", h.CreateAuthor())
			r.Get("/author/{authorID}", h.GetAuthor())
		})
	})

	return r
}

// PrintAllRegisteredRoutes prints all possible routes available
func PrintAllRegisteredRoutes(router *chi.Mux, logger zerolog.Logger) {
	walkFunc := func(method string, path string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		logger.
			Info().
			Dict("routes", zerolog.Dict().
				Str("method", method).
				Str("path", path)).
			Msg("")
		return nil
	}
	if err := chi.Walk(router, walkFunc); err != nil {
		logger.Err(err)
	}
}

// SwaggerServer is serving swagger.
func SwaggerServer(router *chi.Mux) {
	root := "./docs"
	fs := http.FileServer(http.Dir(root))

	router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		if _, err := os.Stat(root + r.RequestURI); os.IsNotExist(err) {
			http.StripPrefix(r.RequestURI, fs).ServeHTTP(w, r)
		} else {
			fs.ServeHTTP(w, r)
		}
	})
}
