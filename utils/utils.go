// Package utils provides various utility functions.
package utils

import (
	"net/url"
	"strings"
	"time"
)

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
	if strings.HasPrefix(rv, "www.") {
		rv = strings.Replace(rv, "www.", "", 1)
	}
	return rv, nil
}
