package sql

import (
	"github.com/boreq/blogs/dto"
	"testing"
)

func TestOrder(t *testing.T) {
	result := Order(false)
	if result != "ASC" {
		t.Fatalf("got %s", result)
	}
}

func TestOrderReverse(t *testing.T) {
	result := Order(true)
	if result != "DESC" {
		t.Fatalf("got %s", result)
	}
}

func TestLimitOffsetFirstPage(t *testing.T) {
	page := dto.Page{
		Page:    1,
		PerPage: 10,
	}
	limit, offset := LimitOffset(page)
	if limit != 10 {
		t.Errorf("limit was %d", limit)
	}
	if offset != 0 {
		t.Errorf("offset was %d", offset)
	}
}

func TestLimitOffset(t *testing.T) {
	page := dto.Page{
		Page:    3,
		PerPage: 10,
	}
	limit, offset := LimitOffset(page)
	if limit != 10 {
		t.Errorf("limit was %d", limit)
	}
	if offset != 20 {
		t.Errorf("offset was %d", offset)
	}
}
