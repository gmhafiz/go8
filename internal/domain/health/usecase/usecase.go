package usecase

import "github.com/gmhafiz/go8/internal/domain/health"

type Health struct {
	healthRepo health.Repository
}

func NewHealthUseCase(health health.Repository) *Health {
	return &Health{
		healthRepo: health,
	}
}

func (u *Health) Readiness() error {
	return u.healthRepo.Readiness()
}
