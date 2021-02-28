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

func GetFilters(queries url.Values) *Filter {
	f := filter.New(queries)
	return &Filter{
		Base:          *f,
		Title:         queries.Get("title"),
		Description:   queries.Get("description"),
		PublishedDate: queries.Get("published_date"),
	}
}
