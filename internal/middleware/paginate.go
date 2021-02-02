package middleware

import (
	"context"
	"net/http"
	"net/url"
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
	Page      int    `json:"page" validate:"number"`
	Size      int    `json:"size" validate:"number"`
	Direction string `json:"sort" validate:"ascii"`
}

func NewPagination() Pagination {
	return Pagination{
		Page:      1,
		Size:      10,
		Direction: "asc",
	}
}

func Paginate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		queries := r.URL.Query()
		pagination, err := Parse(queries)
		if err != nil {
			respond.Error(w, http.StatusBadRequest, err)
			return
		}

		ctx := context.WithValue(r.Context(), PaginationKey, pagination)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func Parse(queries url.Values) (Pagination, error) {
	pagination := NewPagination()
	if from := queries.Get("page"); from != "" {
		p, err := strconv.Atoi(from)
		if err != nil {
			return pagination, err
		}
		pagination.Page = p
	}

	if limit := queries.Get("size"); limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return pagination, err
		}
		pagination.Size = l
	}

	if direction := queries.Get("direction"); direction != "" {
		pagination.Direction = direction
	}

	return pagination, nil
}
