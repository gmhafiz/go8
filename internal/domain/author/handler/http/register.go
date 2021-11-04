package author

import (
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"

	"github.com/gmhafiz/go8/internal/domain/author/usecase"
	"github.com/gmhafiz/go8/internal/middleware"
)

type Options func(opts *Srv) error

type Srv struct {
	router    *chi.Mux
	validator *validator.Validate
	uc        usecase.UseCase
}

func RegisterHTTPEndPoints(opts ...Options) {
	s := &Srv{}
	for _, opt := range opts {
		err := opt(s)
		if err != nil {
			log.Fatalln(err)
		}
	}

	h := NewHandler(s.uc, s.validator)

	s.router.Route("/api/v1/author", func(router chi.Router) {
		router.Use(middleware.Json)
		router.Post("/", h.Create)
		router.Get("/", h.List)
		router.Get("/{id}", h.Get)
		router.Put("/{id}", h.Update)
		router.Delete("/{id}", h.Delete)
	})
}

func WithRouter(router *chi.Mux) func(s *Srv) error {
	return func(s *Srv) error {
		s.router = router
		return nil
	}
}

func WithValidator(v *validator.Validate) func(s *Srv) error {
	return func(s *Srv) error {
		s.validator = v
		return nil
	}
}

func WithUseCase(uc usecase.UseCase) func(s *Srv) error {
	return func(s *Srv) error {
		s.uc = uc
		return nil
	}
}
