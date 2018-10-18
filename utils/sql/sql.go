package sql

import (
	"github.com/boreq/blogs/dto"
)

func Order(reverse bool) string {
	if reverse {
		return "DESC"
	} else {
		return "ASC"
	}
}

func LimitOffset(page dto.Page) (int, int) {
	offset := page.PerPage * (page.Page - 1)
	limit := page.PerPage
	return limit, offset
}
