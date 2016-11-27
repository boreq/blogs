package forms

import (
	"strings"
)

// TrimSpace calls strings.TrimSpace on the value.
func TrimSpace() Cleaner {
	return func(v string) string {
		return strings.TrimSpace(v)
	}
}
