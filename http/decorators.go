package http

import (
	"github.com/boreq/blogs/http/context"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

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

// AuthenticatedOnly redirects not authenticated users to the home page.
func AuthenticatedOnly(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		ctx := context.Get(r)
		if !ctx.User.IsAuthenticated() {
			http.Redirect(w, r, "/", 301)
			return
		}
		h(w, r, p)
	}
}
