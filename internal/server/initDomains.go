package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	authorHandler "github.com/gmhafiz/go8/internal/domain/author/handler"
	authorRepo "github.com/gmhafiz/go8/internal/domain/author/repository"
	authorUseCase "github.com/gmhafiz/go8/internal/domain/author/usecase"
	bookHandler "github.com/gmhafiz/go8/internal/domain/book/handler"
	bookRepo "github.com/gmhafiz/go8/internal/domain/book/repository"
	bookUseCase "github.com/gmhafiz/go8/internal/domain/book/usecase"
	"github.com/gmhafiz/go8/internal/domain/health"
	"github.com/gmhafiz/go8/internal/middleware"
	"github.com/gmhafiz/go8/internal/utility/respond"
)

type Domain struct {
	Health *health.Handler
	Book   *bookHandler.Handler
	Author *authorHandler.Handler
}

func (s *Server) InitDomains() {
	s.initVersion()
	s.initSwagger()
	s.initAuthor()
	s.initHealth()
	s.initBook()
}

func (s *Server) initVersion() {
	s.router.Route("/version", func(router chi.Router) {
		router.Use(middleware.Json)

		router.Get("/", func(w http.ResponseWriter, r *http.Request) {
			respond.Json(w, http.StatusOK, map[string]string{"version": s.Version})
		})
	})
}

func (s *Server) initHealth() {
	newHealthRepo := health.NewRepo(s.DB)
	newHealthUseCase := health.New(newHealthRepo)
	s.Domain.Health = health.RegisterHTTPEndPoints(s.router, newHealthUseCase)
}

func (s *Server) initSwagger() {
	if s.Config().Api.RunSwagger {
		fileServer := http.FileServer(http.Dir(swaggerDocsAssetPath))

		s.router.HandleFunc("/swagger", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/swagger/", http.StatusMovedPermanently)
		})
		s.router.Handle("/swagger/", http.StripPrefix("/swagger", middleware.ContentType(fileServer)))
		s.router.Handle("/swagger/*", http.StripPrefix("/swagger", middleware.ContentType(fileServer)))
	}
}

func (s *Server) initBook() {
	newBookRepo := bookRepo.New(s.DB)
	newBookUseCase := bookUseCase.New(newBookRepo)
	s.Domain.Book = bookHandler.RegisterHTTPEndPoints(s.router, s.validator, newBookUseCase)
}

func (s *Server) initAuthor() {
	newAuthorRepo := authorRepo.New(s.ent)
	newLRUCache := authorRepo.NewLRUCache(newAuthorRepo)
	newRedisCache := authorRepo.NewRedisCache(newAuthorRepo, s.cache)
	newAuthorSearchRepo := authorRepo.NewSearch(s.ent)

	newAuthorUseCase := authorUseCase.New(
		s.cfg.Cache,
		newAuthorRepo,
		newAuthorSearchRepo,
		newLRUCache,
		newRedisCache,
	)
	s.Domain.Author = authorHandler.RegisterHTTPEndPoints(s.router, s.validator, newAuthorUseCase)
}
