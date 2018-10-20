// Package utils provides various utility functions and helpers for sorting,
// filtering and pagination.
package utils

import (
	"net/url"
	"regexp"
	"strings"
	"time"
)

var re = regexp.MustCompile("[^a-z0-9]+")

// Slugify turns a string into a slug containing alphanumeric characters and
// hyphens.
func Slugify(s string) string {
	return strings.Trim(re.ReplaceAllString(strings.ToLower(s), "-"), "-")
}

// ISO8601 returns an ISO8601 formatted string.
func ISO8601(t time.Time) string {
	return t.Format(time.RFC3339)
}

// CleanupUrl strips the scheme and "www." from the url.
func CleanupUrl(s string) (string, error) {
	u, err := url.Parse(s)
	if err != nil {
		return "", err
	}
	rv := strings.Replace(s, u.Scheme+"://", "", 1)
	if (strings.HasPrefix(rv, "www.")) {
		rv = strings.Replace(rv, "www.", "", 1)
	}
	return rv, nil
}
