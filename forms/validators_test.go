package forms

import (
	"testing"
)

func TestMaxLength(t *testing.T) {
	v := MaxLength(2)

	if len(v("")) != 0 {
		t.FailNow()
	}

	if len(v("a")) != 0 {
		t.FailNow()
	}

	if len(v("abc")) == 0 {
		t.FailNow()
	}
}

func TestMinLength(t *testing.T) {
	v := MinLength(2)

	if len(v("")) == 0 {
		t.FailNow()
	}

	if len(v("a")) == 0 {
		t.FailNow()
	}

	if len(v("abc")) != 0 {
		t.FailNow()
	}
}

func TestRegexp(t *testing.T) {
	v := Regexp("^[A-Za-z0-9]+$")

	if len(v("")) == 0 {
		t.FailNow()
	}

	if len(v("abc")) != 0 {
		t.FailNow()
	}

	if len(v("ab#c")) == 0 {
		t.FailNow()
	}
}
