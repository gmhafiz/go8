package middleware

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type key string

type Pagination struct {
	Total int `json:"total"`
	Page  int `json:"page"`
	Size  int `json:"size"`
}

type ID struct {
	Id int64 `json:"id"`
}

const (
	PaginationKey key = "pagination"
	IDKey         key = "id"
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
			pagination.Page = 1
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

		ctx := context.WithValue(r.Context(), PaginationKey, pagination)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func IDParam(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bookID := chi.URLParam(r, "id")
		idInt64, _ := strconv.ParseInt(bookID, 10, 64)

		ctx := context.WithValue(r.Context(), IDKey, idInt64)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
