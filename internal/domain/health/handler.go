package health

import "net/http"

type Health interface {
	Liveness(w http.ResponseWriter, r *http.Request)
	Readiness(w http.ResponseWriter, r *http.Request)
}
