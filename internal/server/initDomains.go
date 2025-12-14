package server

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/gmhafiz/go8/internal/domain/authentication"
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

func (s *Server) InitDomains() {
	s.initVersion()
	s.initSwagger()
	s.initAuthentication()
	s.initAuthor()
	s.initHealth()
	s.initBook()
}

func (s *Server) initVersion() {
	s.router.Route("/version", func(router chi.Router) {
		router.Use(middleware.JSON)

		router.Get("/", func(w http.ResponseWriter, r *http.Request) {
			respond.JSON(w, http.StatusOK, map[string]string{"version": s.Version})
		})
	})
}

func (s *Server) initHealth() {
	newHealthRepo := health.NewRepo(s.sqlx)
	newHealthUseCase := health.New(newHealthRepo)
	health.RegisterHTTPEndPoints(s.router, newHealthUseCase)
}

//go:embed docs/*
var swaggerDocsAssetPath embed.FS

func (s *Server) initSwagger() {
	if s.Config().API.RunSwagger {
		docsPath, err := fs.Sub(swaggerDocsAssetPath, "docs")
		if err != nil {
			panic(err)
		}

		fileServer := http.FileServer(http.FS(docsPath))

		s.router.HandleFunc("/swagger", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/swagger/", http.StatusMovedPermanently)
		})
		s.router.Handle("/swagger/", http.StripPrefix("/swagger", middleware.ContentType(fileServer)))
		s.router.Handle("/swagger/*", http.StripPrefix("/swagger", middleware.ContentType(fileServer)))
	}
}

func (s *Server) initBook() {
	newBookRepo := bookRepo.New(s.sqlx)
	newBookUseCase := bookUseCase.New(newBookRepo)
	bookHandler.RegisterHTTPEndPoints(s.router, s.validator, newBookUseCase)
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
	authorHandler.RegisterHTTPEndPoints(s.router, s.validator, newAuthorUseCase)
}

func (s *Server) initAuthentication() {
	repo := authentication.NewRepo(s.ent, s.db, s.session)
	authentication.RegisterHTTPEndPoints(s.router, s.session, repo)
}
