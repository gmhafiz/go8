package respond

import "net/http"

func Status(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
}
