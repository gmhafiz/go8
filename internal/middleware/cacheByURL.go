package middleware

import (
	"context"
	"encoding/hex"
	"log"
	"net/http"

	"github.com/cespare/xxhash"
)

type cacheKey string

var CacheURL cacheKey

func CacheByURL(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// For cache purpose, we use request URI as the key for our result.
		// We save it into context so that we can pick it pick in our cache layer.
		// We check if key exists and valid to prevent panic.
		//
		// To retrieve:
		// 	  	url, ok := ctx.Value(middleware.CacheURL).(string)
		//      if !ok {
		//         call database layer
		//      }
		str, err := sum(r.URL.String())
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"message":"internal error"}`))
			return
		}

		ctx := context.WithValue(r.Context(), CacheURL, str)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func sum(url string) (string, error) {
	h := xxhash.New()
	_, err := h.Write([]byte(url))
	if err != nil {

		return "", err
	}
	sum := h.Sum(nil)
	str := hex.EncodeToString(sum)

	return str, err
}
