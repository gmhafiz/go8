package middleware

import (
	"context"
	"net/http"
)

const (
	UserID key = "userID"
)

// AuthN is an authentication middleware that checks user's identity
// Pass in dependencies as a parameters to AuthN()
func AuthN() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Retrieve user ID from Cookie / JWT and set it to context
			ctx := context.WithValue(r.Context(), UserID, 1)

			next.ServeHTTP(w, r.WithContext(ctx))

		})
	}
}
