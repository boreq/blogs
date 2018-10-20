package utils

import (
	"testing"
	"time"
)

func TestISO8601(t *testing.T) {
	d := time.Date(2018, 10, 20, 13, 6, 01, 0, time.UTC)
	iso8601 := ISO8601(d)
	if iso8601 != "2018-10-20T13:06:01Z" {
		t.Fatalf("got %s", iso8601)
	}
}

func TestCleanupUrl(t *testing.T) {
	urls := []string{
		"http://www.example.com",
		"http://example.com",
		"www.example.com",
	}

	for _, url := range urls {
		result, err := CleanupUrl(url)
		if err != nil {
			t.Errorf("failed for %s error %s", url, err)
		}
		if result != "example.com" {
			t.Errorf("failed for %s result %s", url, result)
		}
	}
}
