package forms

import (
	"testing"
)

func TestTrimSpace(t *testing.T) {
	c := TrimSpace()

	if len(c("  abc ")) != 3 {
		t.FailNow()
	}
}
