package middleware

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gmhafiz/go8/internal/utility/respond"
)

type ID struct {
	Id int64 `json:"id"`
}

type key string

var (
	PaginationKey key = "pagination"
)

type Pagination struct {
	Page      int    `json:"page"`
	Size      int    `json:"size"`
	Direction string `json:"sort"`
}

func Paginate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var pagination Pagination
		if from := r.URL.Query().Get("page"); from != "" {
			p, err := strconv.Atoi(from)
			if err != nil {
				respond.Error(w, http.StatusBadRequest, err)
			}
			pagination.Page = p
		}

		if limit := r.URL.Query().Get("size"); limit != "" {
			l, err := strconv.Atoi(limit)
			if err != nil {
				respond.Error(w, http.StatusBadRequest, err)
			}
			pagination.Size = l
		}

		if direction := r.URL.Query().Get("direction"); direction != "" {
			pagination.Direction = direction
		}

		ctx := context.WithValue(r.Context(), PaginationKey, pagination)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
