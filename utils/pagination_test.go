package utils

import (
	"net/http"
	"net/url"
	"testing"
)

func makePaginationRequest(key string, t *testing.T) *http.Request {
	urlString := "http://example.com/mypage"
	if key != "" {
		urlString += "?page=" + key
	}
	url, err := url.Parse(urlString)
	if err != nil {
		t.Fatal(err)
	}
	return &http.Request{
		URL: url,
	}
}

func TestBuildUrlQuery(t *testing.T) {
	if buildUrlQuery(nil) != "?" {
		t.Fatal("With no parameters given a query should end with ?")
	}

	params := map[string]string{"key": "value"}
	if buildUrlQuery(params) != "?key=value&" {
		t.Fatal("With parameters present a query should end with &")
	}
}

func TestNegativeOffset(t *testing.T) {
	r := makePaginationRequest("", t)
	pagination := NewPagination(r, 0, 5, nil)
	if pagination.Offset < 0 {
		t.Fatalf("Offset negative: %d", pagination.Offset)
	}
}
