package server

import (
	authorHandler "github.com/gmhafiz/go8/internal/domain/author/handler/http"
	authorRepo "github.com/gmhafiz/go8/internal/domain/author/repository/database"
	authorUseCase "github.com/gmhafiz/go8/internal/domain/author/usecase"
	"github.com/gmhafiz/go8/internal/domain/book/handler/http"
	bookRepo "github.com/gmhafiz/go8/internal/domain/book/repository/postgres"
	bookUseCase "github.com/gmhafiz/go8/internal/domain/book/usecase"
	healthHandler "github.com/gmhafiz/go8/internal/domain/health/handler/http"
	healthRepo "github.com/gmhafiz/go8/internal/domain/health/repository/postgres"
	healthUseCase "github.com/gmhafiz/go8/internal/domain/health/usecase"
)

func (s *Server) initDomains() {
    s.initAuthor()
	s.initHealth()
	s.initBook()
}

func (s *Server) initHealth() {
	newHealthRepo := healthRepo.New(s.DB())
	newHealthUseCase := healthUseCase.New(newHealthRepo)
	healthHandler.RegisterHTTPEndPoints(s.router, newHealthUseCase)
}

func (s *Server) initBook() {
	newBookRepo := bookRepo.New(s.DB())
	newBookUseCase := bookUseCase.New(newBookRepo)
	http.RegisterHTTPEndPoints(s.router, s.validator, newBookUseCase)
}


func (s *Server) initAuthor() {
	newAuthorRepo := authorRepo.New(s.DB())
	newBookRepo := bookRepo.New(s.DB())
	newAuthorUseCase := authorUseCase.New(newAuthorRepo, newBookRepo)
	authorHandler.RegisterHTTPEndPoints(s.router, newAuthorUseCase)
}