package server

import (
	bookHandler "github.com/gmhafiz/go8/internal/domain/book/handler/http"
	bookRepo "github.com/gmhafiz/go8/internal/domain/book/repository/postgres"
	bookUseCase "github.com/gmhafiz/go8/internal/domain/book/usecase"
	healthHandler "github.com/gmhafiz/go8/internal/domain/health/handler/http"
	healthRepo "github.com/gmhafiz/go8/internal/domain/health/repository/postgres"
	healthUseCase "github.com/gmhafiz/go8/internal/domain/health/usecase"
)

func (s *Server) initDomains() {
	s.initHealth()
	s.initBook()
}

func (s *Server) initHealth() {
	newHealthRepo := healthRepo.NewHealthRepository(s.GetDB())
	newHealthUseCase := healthUseCase.NewHealthUseCase(newHealthRepo)
	healthHandler.RegisterHTTPEndPoints(s.router, newHealthUseCase)
}

func (s *Server) initBook() {
	newBookRepo := bookRepo.New(s.GetDB())
	newBookUseCase := bookUseCase.New(newBookRepo)
	bookHandler.RegisterHTTPEndPoints(s.router, newBookUseCase)
}
