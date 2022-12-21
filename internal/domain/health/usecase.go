package usecase

import (
	"github.com/gmhafiz/go8/internal/domain/health/repository/postgres"
)

type UseCase interface {
	Readiness() error
}

type Health struct {
	healthRepo postgres.Repository
}

func New(health postgres.Repository) *Health {
	return &Health{
		healthRepo: health,
	}
}

func (u *Health) Readiness() error {
	return u.healthRepo.Readiness()
}
