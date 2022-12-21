package health

type UseCase interface {
	Readiness() error
}

type Health struct {
	healthRepo Repository
}

func New(health Repository) *Health {
	return &Health{
		healthRepo: health,
	}
}

func (u *Health) Readiness() error {
	return u.healthRepo.Readiness()
}
