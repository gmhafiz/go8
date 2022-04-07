package middleware

import (
	"context"
	"net/http"
	"time"
)

type key string

const (
	AuditID key = "auditID"
)

type Action string

type Event struct {
	ActorID    int       `db:"actor_id" json:"actor_id,omitempty"`
	TableRowID int       `db:"table_row_id" json:"table_row_id,omitempty"`
	Table      string    `db:"table_name" json:"table,omitempty"`
	Action     Action    `db:"action" json:"action,omitempty"`
	OldValues  string    `db:"old_values" json:"old_values,omitempty"`
	NewValues  string    `db:"new_values" json:"new_values,omitempty"`
	HTTPMethod string    `db:"http_method" json:"http_method,omitempty"`
	URL        string    `db:"url" json:"url,omitempty"`
	IPAddress  string    `db:"ip_address" json:"ip_address,omitempty"`
	UserAgent  string    `db:"user_agent" json:"user_agent,omitempty"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

func Audit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ev := Event{
			ActorID:    getUserID(r),
			HTTPMethod: r.Method,
			URL:        r.RequestURI,
			IPAddress:  readUserIP(r),
			UserAgent:  r.UserAgent(),
		}

		ctx := context.WithValue(r.Context(), AuditID, ev)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserID(r *http.Request) int {
	val, _ := r.Context().Value(UserID).(int)
	return val
}

func readUserIP(r *http.Request) string {
	ipAddress := r.Header.Get("X-Real-Ip")
	if ipAddress == "" {
		ipAddress = r.Header.Get("X-Forwarded-For")
	}
	if ipAddress == "" {
		ipAddress = r.RemoteAddr
	}
	return ipAddress
}
