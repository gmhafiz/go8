package filter

import (
	"net/url"
	"strconv"
)

type Filter struct {
	Page          int  `json:"page"`
	Size          int  `json:"size"`
	DisablePaging bool `json:"disable_paging"`
	Search        bool `json:"search"`
}

func New(queries url.Values) *Filter {
	page, _ := strconv.Atoi(queries.Get("page"))
	size, _ := strconv.Atoi(queries.Get("size"))
	disablePaging, _ := strconv.ParseBool(queries.Get("disable_paging"))
	isSearch := has(queries, "search")

	if !has(queries, "size") {
		size = 10
	}

	if !has(queries, "page") {
		page = 1
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
