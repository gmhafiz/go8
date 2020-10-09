package main

import (
	"database/sql"
	"github.com/go-redis/redis/v8"
	"go8ddd/internal/library/cache"
	"go8ddd/internal/library/elasticsearch"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"

	"go8ddd/configs"
	"go8ddd/internal/domain/authors"
	"go8ddd/internal/domain/books"
	"go8ddd/internal/library/logger"
	"go8ddd/internal/library/validation"
	"go8ddd/internal/server/rest"
)

const Version = "v0.2.0"

// @title Go8
// @version 0.2.0
// @description Go + Postgres + Chi Router + SqlBoiler starter kit for API development.

// @contact.name Hafiz Shafruddin
// @contact.url http://www.gmhafiz.com/contact
// @contact.email gmhafiz@gmail.com

// @host localhost:3080
// @BasePath /api/v1
func main() {
	log := logger.New(Version)
	cfg := configs.New(log)
	val := validation.New()
	db := configs.NewDatabase(log, cfg)
	c, _ := cache.New(cfg)
	es := elasticsearch.New(cfg.Elasticsearch)
	s := rest.New(log, cfg, db)

	initializeBookDomain(s.GetRouter(), log, val, db, c, es)
	initializeAuthorDomain(s.GetRouter(), log, val, db)

	err := s.Start(log, cfg)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
}

func initializeBookDomain(router *chi.Mux, log zerolog.Logger, validate *validator.Validate, db *sql.DB, cache *redis.Client, es *elasticsearch.Es) {
	repository := books.NewRepository(log, db, cache, es)
	useCase := books.NewUseCase(repository)
	books.NewHandler(router, validate, useCase)
}

func initializeAuthorDomain(router *chi.Mux, log zerolog.Logger, validate *validator.Validate, db *sql.DB) {
	repository := authors.NewRepository(log, db)
	useCase := authors.NewUseCase(repository)
	authors.NewHandler(router, validate, useCase)
}
