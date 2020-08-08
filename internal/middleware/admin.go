package middleware

import (
	"net/http"
	
	"github.com/go-chi/render"
)

func AdminOnlyHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		admin := isAdmin(r)

		if !admin {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, map[string]string{
				"error": "unauthorized",
			})
			return
		}
		next.ServeHTTP(w, r)
	})
}

func isAdmin(r *http.Request) bool {
	header := r.Header.Get("Authorization")
	if header != "Bearer token" {
		return false
	}

	return true
}
