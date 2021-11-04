package health

import "net/http"

type Health interface {
	Health(w http.ResponseWriter, r *http.Request)
	Readiness(w http.ResponseWriter, r *http.Request)
}
