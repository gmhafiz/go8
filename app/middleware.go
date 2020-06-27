package app

import "net/http"

func (s *Server) AdminOnlyHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		admin := isAdmin(r)

		if !admin {

			http.NotFound(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// ContentTypeJsonHandler is a middleware to be used with chi Router
func (s *Server) ContentTypeJsonHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json;charset=utf8")
		next.ServeHTTP(w, r)
	})
}

func (s *Server) ContentTypeJson(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json;charset=utf8")

		next(w, r)
	}
}

func (s *Server) AdminOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		admin := isAdmin(r)

		if !admin {
			http.NotFound(w, r)
			return
		}
		next(w, r)
	}
}

func isAdmin(r *http.Request) bool {
	header := r.Header.Get("Authorization")
	if header != "Bearer token" {
		return false
	}

	return true
}