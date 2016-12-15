package handler

import (
	bhttp "github.com/boreq/blogs/http"
	"github.com/boreq/blogs/http/context"
	"github.com/boreq/blogs/views"
	"github.com/boreq/blogs/views/errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// Get returns the http handler used by the blogs server.
func Get() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errors.NotFound(w, r)
	})
	views.Register(router)
	return bhttp.RecoverHandler(context.ClearContext(router))
}
