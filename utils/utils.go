package utils

import (
	"regexp"
	"strings"
)

var re = regexp.MustCompile("[^a-z0-9]+")

// Slugify turns a string into a slug containing alphanumeric characters and
// hyphens.
func Slugify(s string) string {
	return strings.Trim(re.ReplaceAllString(strings.ToLower(s), "-"), "-")
}
