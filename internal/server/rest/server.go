package rest

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/go-chi/httplog"
	"github.com/rs/zerolog"

	"go8ddd/configs"
	"go8ddd/internal/middleware"
)

type Server struct {
	server   *http.Server
	router   *chi.Mux
	database *sql.DB
}

func (s Server) Start(log zerolog.Logger, cfg *configs.Configs) error {
	log.Info().Msgf("starting at %s:%s", cfg.Api.Host, cfg.Api.Port)

	printAllRegisteredRoutes(s.router, log)

	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s Server) GetRouter() *chi.Mux {
	return s.router
}

func New(log zerolog.Logger, cfg *configs.Configs, db *sql.DB) *Server {
	router := chi.NewRouter()

	router.Use(httplog.RequestLogger(log))
	router.Use(middleware.Cors)
	router.Use(chiMiddleware.Recoverer)

	swaggerServer(router)

	httpServer := &http.Server{
		Addr:              fmt.Sprintf("%s:%s", cfg.Api.Host, cfg.Api.Port),
		Handler:           router,
		TLSConfig:         nil,
		ReadTimeout:       cfg.Api.ApiReadTimeout * time.Second,
		ReadHeaderTimeout: cfg.Api.ApiReadHeaderTimeout * time.Second,
		WriteTimeout:      cfg.Api.ApiWriteTimeout * time.Second,
		IdleTimeout:       cfg.Api.ApiReadTimeout * time.Second,
	}

	return &Server{
		server:   httpServer,
		router:   router,
		database: db,
	}
}

func swaggerServer(router *chi.Mux) {
	root := "./docs"
	fs := http.FileServer(http.Dir(root))

	router.Get("/swagger", func(w http.ResponseWriter, r *http.Request) {
		if _, err := os.Stat(root + r.RequestURI); os.IsNotExist(err) {
			http.StripPrefix(r.RequestURI, fs).ServeHTTP(w, r)
		} else {
			fs.ServeHTTP(w, r)
		}
	})
}

func printAllRegisteredRoutes(router *chi.Mux, logger zerolog.Logger) {
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
