package http

import (
	"fmt"
	"github.com/boreq/blogs/http/api"
	"github.com/boreq/blogs/logging"
	"github.com/boreq/blogs/service/context"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"github.com/rs/cors"
	"net/http"
	"runtime/debug"
)

var log = logging.New("http")

type Registerer interface {
	Register(router *httprouter.Router)
}

// New returns an http handler used by the blogs server.
func New(registerers []Registerer, contextService *context.ContextService) http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		notFoundError(w, r)
	})
	for _, registerer := range registerers {
		registerer.Register(router)
	}
	c := cors.AllowAll() // TODO
	return c.Handler(RecoverHandler(contextService.ClearContext(router)))
}

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
	log.Error("error", "url", r.URL.String(), "err", err)
	fmt.Println(string(debug.Stack()))
	callApiWithError(w, r, api.InternalServerError)
}

func notFoundError(w http.ResponseWriter, r *http.Request) {
	log.Warn("not found", "url", r.URL.String())
	callApiWithError(w, r, api.NotFoundError)
}

func callApiWithError(w http.ResponseWriter, r *http.Request, error api.Error) {
	api.Call(w, r, nil, func(r *http.Request, p httprouter.Params) (api.Response, api.Error) {
		return nil, error
	})
}
