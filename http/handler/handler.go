package handler

import (
	"github.com/boreq/blogs/http/context"
	"github.com/boreq/blogs/views"
	"github.com/boreq/blogs/views/errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func Get() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errors.NotFound(w, r)
	})
	views.Register(router)
	return context.ClearHandler(router)
}
