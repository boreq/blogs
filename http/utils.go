package http

import (
	"errors"
	"net/http"
)

const nextQueryParamName = "next"

// getNext returns a value of the query string parameter named "next" or an
// error if that parameter is not present.
func getNext(r *http.Request) (string, error) {
	nextList, ok := r.URL.Query()[nextQueryParamName]
	if ok && len(nextList) > 0 {
		return nextList[0], nil
	}
	return "", errors.New("Query parameter not found")
}

// Redirect redirects to the url specified in the "next" query parameter or if
// it is not present to the specified url.
func RedirectOrNext(w http.ResponseWriter, r *http.Request, url string) {
	next, err := getNext(r)
	if err != nil {
		Redirect(w, r, url)
	} else {
		Redirect(w, r, next)
	}
}

// Redirect redirects to the specified url.
func Redirect(w http.ResponseWriter, r *http.Request, url string) {
	http.Redirect(w, r, url, 302)
}
