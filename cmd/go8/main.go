package main

import (
	"database/sql"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"

	"go8ddd/configs"
	"go8ddd/internal/domain/authors"
	"go8ddd/internal/domain/books"
	"go8ddd/internal/domain/health"
	"go8ddd/internal/server/rest"
	"go8ddd/third_party/cache"
	"go8ddd/third_party/database"
	"go8ddd/third_party/logger"
	"go8ddd/third_party/validation"
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
	db := database.New(log, cfg)
	c := cache.New(cfg)
	s := rest.New(log, cfg, db)

	initializeBookDomain(s.GetRouter(), log, val, db, c)
	initializeAuthorDomain(s.GetRouter(), log, val, db)
	initializeHealthDomain(s.GetRouter(), log, db)

	err := s.Start(log, cfg)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
}

func initializeHealthDomain(router *chi.Mux, log zerolog.Logger, db *sql.DB) {
	health.New(router, log, db)
}

func initializeBookDomain(router *chi.Mux, log zerolog.Logger, validate *validator.Validate,
	db *sql.DB, cache *redis.Client) {
	repository := books.NewRepository(log, db, cache)
	useCase := books.NewUseCase(repository)
	books.NewHandler(router, validate, useCase)
}

func initializeAuthorDomain(router *chi.Mux, log zerolog.Logger, validate *validator.Validate,
	db *sql.DB) {
	repository := authors.NewRepository(log, db)
	useCase := authors.NewUseCase(repository)
	authors.NewHandler(router, validate, useCase)
}
