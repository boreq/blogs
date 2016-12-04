package utils

import (
	"net/http"
	"net/url"
	"testing"
)

func makeRequest(key string, t *testing.T) *http.Request {
	urlString := "http://example.com/mypage"
	if key != "" {
		urlString += "?sort=" + key
	}
	url, err := url.Parse(urlString)
	if err != nil {
		t.Fatal(err)
	}
	return &http.Request{
		URL: url,
	}
}

func TestSort(t *testing.T) {
	const key = "key2"
	r := makeRequest(key, t)

	sort := NewSort(r, []SortParam{
		{"key1", "Label 1", "key1.id", false},
		{"key2", "Label 2", "key2.id", false},
		{"key3", "Label 3", "key3.id", false},
	})
	if sort.Query != "key2.id" {
		t.Fatal("Wrong query:", sort.Query)
	}
	if len(sort.Criteria) != 3 {
		t.Fatal("Wrong length:", len(sort.Criteria))
	}
	c := sort.Criteria[1]
	if c.GetKey() != "key2_desc" {
	}
	if c.Selected == false {
		t.Fatal("Not selected")
	}
}

func TestSortDesc(t *testing.T) {
	const key = "key2_desc"
	r := makeRequest(key, t)

	sort := NewSort(r, []SortParam{
		{"key1", "Label 1", "key1.id", false},
		{"key2", "Label 2", "key2.id", false},
		{"key3", "Label 3", "key3.id", false},
	})
	if sort.Query != "key2.id DESC" {
		t.Fatal("Wrong query:", sort.Query)
	}
	if len(sort.Criteria) != 3 {
		t.Fatal("Wrong length:", len(sort.Criteria))
	}
	c := sort.Criteria[1]
	if c.GetKey() != "key2" {
	}
	if c.Selected == false {
		t.Fatal("Not selected")
	}
}

func TestSortNoKey(t *testing.T) {
	r := makeRequest("", t)

	sort := NewSort(r, []SortParam{
		{"key1", "Label 1", "key1.id", false},
		{"key2", "Label 2", "key2.id", false},
		{"key3", "Label 3", "key3.id", false},
	})
	if sort.Query != "key1.id" {
		t.Fatal("Wrong query:", sort.Query)
	}
	if len(sort.Criteria) != 3 {
		t.Fatal("Wrong length:", len(sort.Criteria))
	}
}

func TestSortDirection(t *testing.T) {
	const key = "key2"
	r := makeRequest(key, t)

	sort := NewSort(r, []SortParam{
		{"key1", "Label 1", "key1.id", false},
		{"key2", "Label 2", "key2.id", true},
		{"key3", "Label 3", "key3.id", false},
	})
	if sort.Query != "key2.id DESC" {
		t.Fatal("Wrong query:", sort.Query)
	}
	c := sort.Criteria[1]
	if !c.Selected {
		t.Fatal("Not selected")
	}
	if c.GetKey() != "key2_desc" {
		t.Fatal("Wrong key")
	}
	if c.IsAsc() != false {
		t.Fatal("Wrong IsAsc")
	}
}
