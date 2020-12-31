package health

type Repository interface {
	Readiness() error
}
