package http

import (
	"errors"
	"fmt"
	"github.com/boreq/blogs/http/api"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"runtime/debug"
)

// RecoverHandler is a middleware that recovers from panics and displays a 500
// error page.
func RecoverHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				err, ok := rec.(error)
				if !ok {
					err = errors.New("Recovered from a panic")
				}
				internalServerError(w, r, err)
			}
			return

		}()
		h.ServeHTTP(w, r)
	})
}

func internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	fmt.Printf("%s\n", err)
	fmt.Println(string(debug.Stack()))
	api.Call(w, r, nil, func(r *http.Request, p httprouter.Params) (api.Response, api.Error) {
		return nil, api.InternalServerError
	})
}
