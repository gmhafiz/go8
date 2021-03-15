package filter

import (
	"net/url"
	"strconv"
)

const (
	paginationDefaultPage = 1
	paginationDefaultSize = 30

	queryParamDisablePaging = "disable_paging"
	queryParamPage          = "page"
	queryParamSize          = "size"
	queryParamSearch        = "search"
)

type Filter struct {
	Page          int  `json:"page"`
	Size          int  `json:"size"`
	DisablePaging bool `json:"disable_paging"`
	Search        bool `json:"search"`
}

func New(queries url.Values) *Filter {
	page, _ := strconv.Atoi(queries.Get(queryParamPage))
	size, _ := strconv.Atoi(queries.Get(queryParamSize))
	disablePaging, _ := strconv.ParseBool(queries.Get(queryParamDisablePaging))
	isSearch := has(queries, queryParamSearch)

	if !has(queries, queryParamSize) {
		size = paginationDefaultSize
	}

	if !has(queries, queryParamPage) {
		page = paginationDefaultPage
	}
	page = size * (page - 1) // calculates offset

	return &Filter{
		Page:          page,
		Size:          size,
		DisablePaging: disablePaging,
		Search:        isSearch,
	}
}

func has(queries url.Values, param string) bool {
	return queries.Get(param) != ""
}
