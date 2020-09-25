package middleware

import (
	"context"
	"net/http"
	"strconv"
)

type Pagination struct {
	Total int `json:"total"`
	Page  int `json:"page"`
	Size  int `json:"size"`
}

type key string
const (
	paginationKey key = ""
)

func Paginate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var pagination Pagination

		if from := r.URL.Query().Get("page"); from != "" {
			p, err := strconv.Atoi(from)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			}
			pagination.Page = p
		} else {
			pagination.Page = 0
		}

		if limit := r.URL.Query().Get("size"); limit != "" {
			l, err := strconv.Atoi(limit)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			}
			pagination.Size = l
		} else {
			pagination.Size = 10
		}

		ctx := context.WithValue(r.Context(), paginationKey, pagination)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
