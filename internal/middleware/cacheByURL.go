package middleware

import (
	"context"
	"net/http"
)

type cacheKey string

var CacheURL cacheKey

func CacheByURL(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// For cache purpose, we use request URI as the key for our result.
		// We save it into context so that we can pick it pick in our cache layer.
		//
		// To retrieve:
		// 	  	url := ctx.Value(middleware.CacheURL).(string)
		//
		ctx := context.WithValue(r.Context(), CacheURL, r.URL.String())

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
