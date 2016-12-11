package utils

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
)

func makeRequestURL(key string, value string, t *testing.T) *http.Request {
	urlString := "http://example.com/mypage"
	if key != "" {
		urlString += fmt.Sprintf("?%s=%s", key, value)
	}
	url, err := url.Parse(urlString)
	if err != nil {
		t.Fatal(err)
	}
	return &http.Request{
		URL: url,
	}
}

func TestFilterKey(t *testing.T) {
	r := makeRequestURL("filter", "key2", t)

	f := NewFilter(r, []FilterParam{
		{"key1", "Label 1", "key1.query"},
		{"key2", "Label 2", "key2.query"},
		{"key3", "Label 3", "key3.query"},
	})
	if f.CurrentKey != "key2" {
		t.Fatal("Wrong key:", f.CurrentKey)
	}
}

func TestFilterKeyInvalid(t *testing.T) {
	r := makeRequestURL("filter", "invalidkey", t)

	f := NewFilter(r, []FilterParam{
		{"key1", "Label 1", "key1.query"},
		{"key2", "Label 2", "key2.query"},
		{"key3", "Label 3", "key3.query"},
	})
	if f.CurrentKey != "key1" {
		t.Fatal("Wrong key:", f.CurrentKey)
	}
}

func TestFilterMissing(t *testing.T) {
	r := makeRequestURL("invalid", "invalid", t)

	f := NewFilter(r, []FilterParam{
		{"key1", "Label 1", "key1.query"},
		{"key2", "Label 2", "key2.query"},
		{"key3", "Label 3", "key3.query"},
	})
	if f.CurrentKey != "key1" {
		t.Fatal("Wrong key:", f.CurrentKey)
	}
}
