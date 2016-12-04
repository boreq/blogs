package utils

import (
	"math"
	"net/http"
	"net/url"
	"strconv"
)

type Pagination struct {
	// Page is the current page number.
	Page int

	// AllPages is the number of all pages/last page.
	AllPages int

	// HasNext is true if there is a next page.
	HasNext bool

	// HasPrevious is true if there is a previous page.
	HasPrevious bool

	// Offset can be used as a parameter in a SQL query.
	Offset int

	// Limit can be used as a parameter in a SQL query.
	Limit int

	// URLQuery carries information about parameters preserved during page
	// changes.
	URLQuery string
}

// NewPagination uses the "page" query parameter to get the page number and
// initialize the struct.
func NewPagination(r *http.Request, allItems uint, itemsPerPage uint, preserveParams map[string]string) Pagination {
	page := getPageNumber(r)
	allPages := int(math.Ceil(float64(allItems) / float64(itemsPerPage)))
	if page > allPages {
		page = allPages
	}
	rv := Pagination{
		Page:        page,
		AllPages:    allPages,
		HasNext:     page < allPages,
		HasPrevious: page > 1,
		Offset:      int(itemsPerPage) * (page - 1),
		Limit:       int(itemsPerPage),
		URLQuery:    buildUrlQuery(preserveParams),
	}
	return rv
}

func buildUrlQuery(params map[string]string) string {
	v := url.Values{}
	for key, value := range params {
		v.Add(key, value)
	}
	rv := v.Encode()
	if rv != "" {
		rv += "&"
	}
	return "?" + rv
}

func getPageNumber(r *http.Request) int {
	pageParam, ok := r.URL.Query()["page"]
	if ok {
		p, err := strconv.ParseInt(pageParam[0], 10, 32)
		if err == nil && p >= 1 {
			return int(p)
		}
	}
	return 1
}
