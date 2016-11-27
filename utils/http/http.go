package http

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type ErrorHandler func(http.ResponseWriter, *http.Request, httprouter.Params) error

// WithErrorHandling converts an ErrorHandler into a normal httprouter handle
// adding error handling to it. If the wrapped function returns an error or
// panics an error page is displayed.
func WithErrorHandling(eh ErrorHandler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		if err := eh(w, r, p); err != nil {
			http.Error(w, err.Error(), 500)
		}
	}
}
