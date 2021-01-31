package book

import "net/url"

type Filters struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Page        string `json:"page"`
	Size        string `json:"size"`
}

func Filter(queries url.Values) (*Filters, bool) {
	var isSearch bool

	for key := range queries {
		if key == "search" {
			isSearch = true
			break
		}
	}

	if isSearch {
		filter := &Filters{
			Title:       queries.Get("title"),
			Description: queries.Get("description"),
		}
		return filter, true
	}
	return nil, false
}
