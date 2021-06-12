package author

import (
	"net/url"

	"github.com/gmhafiz/go8/internal/utility/filter"
)

type Filter struct {
	Base          filter.Filter
	FirstName         string `json:"first_name"`
	LastName         string `json:"last_name"`
}

func Filters(queries url.Values) *Filter {
	f := filter.New(queries)
	return &Filter{
		Base:          *f,
		FirstName:         queries.Get("first_name"),
		LastName:   queries.Get("last_name"),
	}
}
