package book

import (
	"net/url"

	"github.com/gmhafiz/go8/internal/middleware"
)

func Filter(queries url.Values) (*Request, bool) {
	var isSearch bool

	for key := range queries {
		if key == "search" {
			isSearch = true
			break
		}
	}

	if isSearch {
		f, err := middleware.Parse(queries)
		if err != nil {
			return nil, false
		}

		req := &Request{
			Title:       queries.Get("title"),
			Description: queries.Get("description"),
			Pagination:  f,
		}

		return req, true
	}
	return nil, false
}
