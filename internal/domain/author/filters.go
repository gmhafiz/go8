package author

import (
	"net/url"

	"github.com/gmhafiz/go8/internal/utility/filter"
)

type Filter struct {
	Base filter.Filter

	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name"`
	LastName   string `json:"last_name"`
}

func Filters(queries url.Values) *Filter {
	f := filter.New(queries)
	if queries.Has("first_name") || queries.Has("last_name") {
		f.Search = true
	}
	return &Filter{
		Base: *f,

		FirstName:  queries.Get("first_name"),
		MiddleName: queries.Get("middle_name"),
		LastName:   queries.Get("last_name"),
	}
}
