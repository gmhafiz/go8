package health

type UseCase interface {
	Readiness() error
}
