package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"runtime/debug"

	"github.com/go-chi/chi/middleware"
)

// Recovery adapted from https://github.com/go-chi/chi/blob/master/middleware/recoverer.go
func Recovery(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil && rvr != http.ErrAbortHandler {
				defer r.Body.Close()

				logEntry := middleware.GetLogEntry(r)
				if logEntry != nil {
					logEntry.Panic(rvr, debug.Stack())
				} else {
					debug.PrintStack()
				}

				log.Printf("PANIC: %v", rvr)
				// send to centralised logging system
				log.Printf("request: %s %s\n", r.Method, r.URL.RequestURI())

				dump, err := httputil.DumpRequest(r, true)
				if err != nil {
					log.Println(err)
				}

				b, _ := json.Marshal(dump)
				log.Printf("host: %s\n", r.Host)
				log.Printf("request body: %s\n", b)

				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
