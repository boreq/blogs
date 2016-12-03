package utils

import (
	"math"
	"net/http"
	"strconv"
)

type Pagination struct {
	// Page is the current page number.
	Page uint

	// AllPages is the number of all pages/last page.
	AllPages uint

	// HasNext is true if there is a next page.
	HasNext bool

	// HasPrevious is true if there is a previous page.
	HasPrevious bool

	// Offset can be used as a parameter in a SQL query.
	Offset uint

	// Limit can be used as a parameter in a SQL query.
	Limit uint
}

// NewPagination uses the "page" query parameter to get the page number and
// initialize the struct.
func NewPagination(r *http.Request, allItems uint, itemsPerPage uint) Pagination {
	page := getPageNumber(r)
	allPages := uint(math.Ceil(float64(allItems) / float64(itemsPerPage)))
	if page > allPages {
		page = allPages
	}
	rv := Pagination{
		Page:        page,
		AllPages:    allPages,
		HasNext:     page < allPages,
		HasPrevious: page > 1,
		Offset:      itemsPerPage * (page - 1),
		Limit:       itemsPerPage,
	}
	return rv
}

func getPageNumber(r *http.Request) uint {
	pageParam, ok := r.URL.Query()["page"]
	if ok {
		p, err := strconv.ParseUint(pageParam[0], 10, 32)
		if err == nil && p >= 1 {
			return uint(p)
		}
	}
	return 1
}
