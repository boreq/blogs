package sql

import "testing"

func TestOrder(t *testing.T) {
	result := Order(false)
	if result != "DESC" {
		t.Fatalf("got %s", result)
	}
}

func TestOrderReverse(t *testing.T) {
	result := Order(true)
	if result != "ASC" {
		t.Fatalf("got %s", result)
	}
}
