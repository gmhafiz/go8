package middleware

import (
	"go.uber.org/zap"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/render"
)

type logEntry struct {
	ReceivedTime      time.Time
	RequestMethod     string
	RequestURL        string
	RequestHeaderSize int64
	RequestBodySize   int64
	UserAgent         string
	Referer           string
	Proto             string

	RemoteIP string
	ServerIP string

	Status             int
	ResponseHeaderSize int64
	ResponseBodySize   int64
	Latency            time.Duration
}

//func RequestLog(log log.Logger) func(http.Handler) http.Handler {
func RequestLog(log *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			le := &logEntry{
				ReceivedTime:      start,
				RequestMethod:     r.Method,
				RequestURL:        r.URL.String(),
				RequestHeaderSize: headerSize(r.Header),
				UserAgent:         r.UserAgent(),
				Referer:           r.Referer(),
				Proto:             r.Proto,
				RemoteIP:          ipFromHostPort(r.RemoteAddr),
			}

			if addr, ok := r.Context().Value(http.LocalAddrContextKey).(net.Addr); ok {
				le.ServerIP = ipFromHostPort(addr.String())
			}

			zap.Fields()
			log.Info("",
				zap.String("request method", le.RequestMethod),
				zap.String("request url", le.RequestURL),
				zap.Int64("request header size", le.RequestHeaderSize),
				zap.String("user agent", le.UserAgent),
				zap.String("referer", le.Referer),
				zap.String("proto", le.Proto),
				zap.String("remote ip", le.RemoteIP),
				)

			// Call the next handler
			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}

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

type writeCounter int64

func (wc *writeCounter) Write(p []byte) (n int, err error) {
	*wc += writeCounter(len(p))
	return len(p), nil
}

func headerSize(h http.Header) int64 {
	var wc writeCounter
	_ = h.Write(&wc)
	return int64(wc) + 2 // for CRLF
}

func ipFromHostPort(hp string) string {
	h, _, err := net.SplitHostPort(hp)
	if err != nil {
		return ""
	}
	if len(h) > 0 && h[0] == '[' {
		return h[1 : len(h)-1]
	}
	return h
}