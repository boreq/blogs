package views

import (
	"github.com/boreq/blogs/dto"
	"net/http"
	"strconv"
)

func GetPage(r *http.Request) dto.Page {
	return dto.Page{
		Page:    getPageNumber(r),
		PerPage: getPerPage(r, 10, 50),
	}
}

func GetSort(r *http.Request) string {
	sort, ok := r.URL.Query()["sort"]
	if ok {
		return sort[0]
	}
	return ""
}

func GetSortReverse(r *http.Request) bool {
	reverse, ok := r.URL.Query()["reverse"]
	if ok {
		return reverse[0] == "true"
	}
	return false
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

func getPerPage(r *http.Request, def int, max int) int {
	pageParam, ok := r.URL.Query()["perPage"]
	if ok {
		p, err := strconv.ParseInt(pageParam[0], 10, 32)
		if err == nil && p >= 1 && p <= int64(max) {
			return int(p)
		}
	}
	return def
}
