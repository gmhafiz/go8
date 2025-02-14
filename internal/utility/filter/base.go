package filter

import (
	"net/url"
	"strconv"
	"strings"
)

const (
	paginationDefaultPage = 1
	paginationDefaultSize = 30
	paginationMaxSize     = 500

	queryParamPage          = "page"
	queryParamLimit         = "limit"
	queryParamOffset        = "offset"
	queryParamDisablePaging = "disable_paging"
	queryParamSort          = "sort"
)

type Filter struct {
	Page          int
	Offset        int
	Limit         int
	DisablePaging bool

	Sort   map[string]string
	Search bool
}

func New(queries url.Values) *Filter {
	var page, limit, offset int
	page, err := strconv.Atoi(queries.Get(queryParamPage))
	if err != nil {
		page = paginationDefaultPage
	}
	limit, err = strconv.Atoi(queries.Get(queryParamLimit))
	if err != nil {
		limit = paginationDefaultSize
	}
	if limit > paginationMaxSize {
		limit = paginationMaxSize
	}

	offset, err = strconv.Atoi(queries.Get(queryParamOffset))
	if err != nil {
		offset = limit * (page - 1) // calculates offset
	}

	disablePaging, _ := strconv.ParseBool(queries.Get(queryParamDisablePaging))

	sortKey := make(map[string]string)
	if queries.Has(queryParamSort) {
		s := queries[queryParamSort]
		for _, val := range s {
			key, order, found := strings.Cut(val, ",")
			if found {
				sortKey[key] = strings.ToUpper(order)
			} else {
				sortKey[key] = "ASC"
			}
		}
	}

	return &Filter{
		Page:          page,
		Offset:        offset,
		Limit:         limit,
		DisablePaging: disablePaging,
		Sort:          sortKey,
	}
}
