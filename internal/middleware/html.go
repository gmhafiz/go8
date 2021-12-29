package middleware

import "net/http"

func Html(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		h.ServeHTTP(w, r)
	}
}
