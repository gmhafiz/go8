package author

import (
	"net/url"

	"github.com/gmhafiz/go8/internal/utility/filter"
)

type Filter struct {
	Base filter.Filter
	Name string `json:"name"`
}

func Filters(queries url.Values) *Filter {
	f := filter.New(queries)
	if queries.Has("name") {
		f.Search = true
	}
	return &Filter{
		Base: *f,
		Name: queries.Get("name"),
	}
}
