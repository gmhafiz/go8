package book

import (
	"net/url"

	"github.com/gmhafiz/go8/internal/utility/filter"
)

type Filter struct {
	Base          filter.Filter
	Title         string `json:"title"`
	Description   string `json:"description"`
	PublishedDate string `json:"published_date"`
}

func Filters(queries url.Values) *Filter {
	f := filter.New(queries)
	switch {
	case queries.Has("title"):
		fallthrough
	case queries.Has("description"):
		f.Search = true
	}

	return &Filter{
		Base:          *f,
		Title:         queries.Get("title"),
		Description:   queries.Get("description"),
		PublishedDate: queries.Get("published_date"),
	}
}
