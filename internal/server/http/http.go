package http

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"

	"eight/internal/api"
)

// Handlers struct has all the dependencies required for HTTP handlers
type Handlers struct {
	Api        *api.API
	Logger     zerolog.Logger
	Validation *validator.Validate
}

// HTTP struct holds all the dependencies required for starting HTTP server
type HTTP struct {
	server *http.Server
	cfg    *Config
	router *chi.Mux
}

// Config holds all the configuration required to start the HTTP server
type Config struct {
	Host         string        `yaml:"HOST"`
	Port         string        `yaml:"PORT"`
	ReadTimeout  time.Duration `yaml:"READ_TIMEOUT"`
	WriteTimeout time.Duration `yaml:"WRITE_TIMEOUT"`
	DialTimeout  time.Duration `yaml:"DIAL_TIMEOUT"`
}

func (h *HTTP) Start(logger zerolog.Logger) {
	logger.Info().Msgf("starting at %s:%s", h.cfg.Host, h.cfg.Port)

	PrintAllRegisteredRoutes(h.router, logger)

	if err := h.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Err(err)
		os.Exit(-1)
	}
}

func (h *HTTP) GetServer() *chi.Mux {
	return h.router
}

func NewService(cfg *Config, a *api.API, log zerolog.Logger, validation *validator.Validate) (*HTTP, error) {
	h := &Handlers{
		Api:        a,
		Logger:     log,
		Validation: validation,
	}

	serverHandler := Router(h, log)

	httpServer := &http.Server{
		Addr:              fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Handler:           serverHandler,
		TLSConfig:         nil,
		ReadTimeout:       cfg.ReadTimeout,
		ReadHeaderTimeout: cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.ReadTimeout * 2,
	}

	return &HTTP{
		server: httpServer,
		cfg:    cfg,
		router: serverHandler,
	}, nil
}
