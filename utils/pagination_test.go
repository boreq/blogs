package utils

import (
	"testing"
)

func TestBuildUrlQuery(t *testing.T) {
	if buildUrlQuery(nil) != "?" {
		t.Fatal("With no parameters given a query should end with ?")
	}

	params := map[string]string{"key": "value"}
	if buildUrlQuery(params) != "?key=value&" {
		t.Fatal("With parameters present a query should end with &")
	}
}
