package middleware

import (
	"encoding/json"
	"log"
	"net/http"

	chiMiddleware "github.com/go-chi/chi/middleware"
)

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("PANIC")
		// notify to sentry/email/slack
		log.Printf("request: %s %s\n", r.Method, r.URL.RequestURI())
		body := r.GetBody
		b, _ := json.Marshal(body)
		log.Printf("host: %s\n", r.Host)
		log.Printf("request body: %s\n", b)
		chiMiddleware.Recoverer(next)

		next.ServeHTTP(w, r)
	})
}
