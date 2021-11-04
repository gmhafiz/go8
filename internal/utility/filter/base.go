package filter

import (
	"net/url"
	"strconv"
	"strings"
)

const (
	paginationDefaultPage = 1
	paginationDefaultSize = 30

	queryParamDisablePaging = "disable_paging"
	queryParamPage          = "page"
	queryParamSize          = "size"
	queryParamSort          = "sort"
)

type Filter struct {
	Offset        int               `json:"page"`
	Limit         uint              `json:"size"`
	DisablePaging bool              `json:"disable_paging"`
	Sort          map[string]string `json:"sort"`
	Search        bool
}

func New(queries url.Values) *Filter {
	var page, size int
	page, err := strconv.Atoi(queries.Get(queryParamPage))
	if err != nil {
		page = paginationDefaultPage
	}
	size, err = strconv.Atoi(queries.Get(queryParamSize))
	if err != nil {
		size = paginationDefaultSize
	}
	offset := size * (page - 1) // calculates offset

	disablePaging, _ := strconv.ParseBool(queries.Get(queryParamDisablePaging))

	sortKey := make(map[string]string)
	if queries.Has(queryParamSort) {
		s := queries[queryParamSort]
		for _, val := range s {
			split := strings.Split(val, ",")
			if len(split) == 2 {
				sortKey[split[0]] = split[1]
			} else {
				sortKey[split[0]] = "asc"
			}
		}
	}

	return &Filter{
		Offset:        offset,
		Limit:         uint(size),
		DisablePaging: disablePaging,
		Sort:          sortKey,
	}
}
