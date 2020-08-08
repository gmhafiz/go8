package http

import (
	"fmt"
	"github.com/jinzhu/now"
	"github.com/rs/zerolog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"

	"eight/internal/api"
)

// Handlers struct has all the dependencies required for HTTP handlers
type Handlers struct {
	Api *api.API
	Logger zerolog.Logger
	TimeConverter now.Config
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
	//logger.Info("starting at port ", zap.String("host", h.cfg.Host) , zap.String("port", h.cfg.Port) )
	//logger.Info().Msgf("starting at port ", zap.String("host", h.cfg.Host) , zap.String("port",h.cfg.Port))

	logger.Info().Msgf("starting at %s:%s", h.cfg.Host, h.cfg.Port)

	//logger.Infof("starting at port %s:%s", h.cfg.Host, h.cfg.Port)
	//log.Printf("starting at port %s:%s", h.cfg.Host, h.cfg.Port)

	//logz.Info("starting at", zap.String("host", h.cfg.Host) , zap.String("port",h.cfg.Port))

	PrintAllRegisteredRoutes(h.router, logger)

	if err := h.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Err(err)
		//logger.Error().Err(err)
		os.Exit(-1)
	}
}

func (h *HTTP) GetServer() *chi.Mux {
	return h.router
}

func NewService(cfg *Config, a *api.API, log zerolog.Logger, timeConverter now.Config) (*HTTP, error) {
	h := &Handlers{
		Api: a,
		Logger: log,
		TimeConverter: timeConverter,
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
