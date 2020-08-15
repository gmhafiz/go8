package main

import (
	"flag"

	"github.com/go-chi/httplog"

	"eight/internal/api"
	"eight/internal/configs"
	"eight/internal/datastore"
	"eight/internal/domain/authors"
	"eight/internal/domain/books"
	"eight/internal/server/http"
	"eight/pkg/redis"
	"eight/pkg/validation"
)

const Version = "v0.1.0"

var flagConfig = flag.String("config", "./config/dev.yml", "path to the config file")

// @title Go8
// @version 0.1.0
// @description Go + Postgres + Chi Router + SqlBoiler starter kit for API development.

// @contact.name Hafiz Shafruddin
// @contact.url http://www.gmhafiz.com/contact
// @contact.email gmhafiz@gmail.com

// @host localhost:3080
// @BasePath /api/v1
func main() {
	logger := httplog.NewLogger("go8", httplog.Options{
		JSON:    false, // switch to false for a human readable log format
		Concise: true,
		Tags:    map[string]string{"version": Version},
	})
	logger = logger.With().Caller().Logger()

	//cfg, err := configs.NewService("dev")
	cfg, err := configs.NewService(*flagConfig)
	if err != nil {
		logger.Error().Err(err)
		return
	}

	dataStoreCfg, err := cfg.DataStore()
	if err != nil {
		logger.Error().Err(err)
		return
	}

	db, err := datastore.NewService(dataStoreCfg)
	if err != nil {
		logger.Error().Err(err)
		return
	}

	cacheCfg, err := cfg.CacheStore()
	if err != nil {
		logger.Error().Err(err)
		return
	}

	rdb, err := redis.NewClient(cacheCfg)
	if err != nil {
		logger.Error().Err(err)
		return
	}

	bookService, err := books.NewService(db, logger, rdb)
	if err != nil {
		logger.Error().Err(err)
		return
	}

	authorService, err := authors.NewService(db, logger, rdb)
	if err != nil {
		logger.Error().Err(err)
		return
	}

	// additional microservice added here
	a, err := api.NewService(bookService, authorService)
	if err != nil {
		logger.Error().Err(err)
		return
	}

	httpCfg, err := cfg.HTTP()
	if err != nil {
		logger.Error().Err(err)
		return
	}

	val := validation.New()

	h, err := http.NewService(httpCfg, a, logger, val)
	if err != nil {
		logger.Error().Err(err)
		return
	}

	h.Start(logger)
}
