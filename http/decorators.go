package http

import (
	"github.com/boreq/blogs/http/context"
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

// NotAuthenticatedOnly redirects authenticated users to the home page.
func NotAuthenticatedOnly(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		ctx := context.Get(r)
		if ctx.User.IsAuthenticated() {
			http.Redirect(w, r, "/", 301)
			return
		}
		h(w, r, p)
	}
}
