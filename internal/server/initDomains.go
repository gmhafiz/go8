package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	authorHandlerHTTP "github.com/gmhafiz/go8/internal/domain/author/handler/http"
	authorCache "github.com/gmhafiz/go8/internal/domain/author/repository/cache"
	authorRepo "github.com/gmhafiz/go8/internal/domain/author/repository/database"
	authorSearchRepo "github.com/gmhafiz/go8/internal/domain/author/repository/search"
	authorUseCase "github.com/gmhafiz/go8/internal/domain/author/usecase"
	bookHandlerHTTP "github.com/gmhafiz/go8/internal/domain/book/handler/http"
	bookRepo "github.com/gmhafiz/go8/internal/domain/book/repository/postgres"
	bookUseCase "github.com/gmhafiz/go8/internal/domain/book/usecase"
	healthHandlerHTTP "github.com/gmhafiz/go8/internal/domain/health/handler/http"
	healthRepo "github.com/gmhafiz/go8/internal/domain/health/repository/postgres"
	healthUseCase "github.com/gmhafiz/go8/internal/domain/health/usecase"
	"github.com/gmhafiz/go8/internal/middleware"
	"github.com/gmhafiz/go8/internal/utility/respond"
)

func (s *Server) InitDomains() {
	s.initVersion()
	s.initAuthor()
	s.initHealth()
	s.initBook()
}

func (s *Server) initVersion() {
	s.router.Route("/version", func(router chi.Router) {
		router.Use(middleware.Json)

		router.Get("/", func(w http.ResponseWriter, r *http.Request) {
			respond.Json(w, http.StatusOK, map[string]string{"version": s.version})
		})
	})
}

func (s *Server) initHealth() {
	newHealthRepo := healthRepo.New(s.DB())
	newHealthUseCase := healthUseCase.New(newHealthRepo)
	healthHandlerHTTP.RegisterHTTPEndPoints(s.router, newHealthUseCase)
}

func (s *Server) initBook() {
	newBookRepo := bookRepo.New(s.DB())
	newBookUseCase := bookUseCase.New(newBookRepo)
	_ = bookHandlerHTTP.RegisterHTTPEndPoints(s.router, s.validator, newBookUseCase)
}

func (s *Server) initAuthor() {
	newAuthorRepo := authorRepo.New(s.DB())
	newLRUCache := authorCache.NewLRUCache(newAuthorRepo)
	newRedisCache := authorCache.NewRedisCache(newAuthorRepo, s.Cache())
	newAuthorSearchRepo := authorSearchRepo.New(s.DB())

	newBookRepo := bookRepo.New(s.DB())

	newAuthorUseCase := authorUseCase.New(
		newAuthorRepo,
		newAuthorSearchRepo,
		newLRUCache,
		newRedisCache,
		newBookRepo,
	)
	authorHandlerHTTP.RegisterHTTPEndPoints(
		authorHandlerHTTP.WithRouter(s.router),
		authorHandlerHTTP.WithValidator(s.validator),
		authorHandlerHTTP.WithUseCase(newAuthorUseCase),
	)
}
