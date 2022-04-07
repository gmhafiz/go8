package middleware

import (
	"net/http"
	"path/filepath"
)

func ContentType(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ext := filepath.Ext(r.RequestURI)
		switch ext {
		case ".png":
			w.Header().Set("Content-Type", "image/png")
		case ".css":
			w.Header().Set("Content-Type", "text/css")
		case ".js":
			w.Header().Set("Content-Type", "application/javascript")
		case ".json":
			w.Header().Set("Content-Type", "application/json")
		case ".ico":
			w.Header().Set("Content-Type", "image/icon")
		default:
			w.Header().Set("Content-Type", "text/html")
		}

		h.ServeHTTP(w, r)
	}
}
