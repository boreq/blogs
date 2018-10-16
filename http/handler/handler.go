package handler

import (
	bhttp "github.com/boreq/blogs/http"
	"github.com/boreq/blogs/http/api"
	"github.com/boreq/blogs/http/context"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"net/http"
)

type Registerer interface {
	Register(router *httprouter.Router)
}

// New returns the http handler used by the blogs server.
func New(registerers []Registerer) http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api.Call(w, r, nil, func(r *http.Request, p httprouter.Params) (api.Response, api.Error) {
			return nil, api.NotFoundError
		})
	})
	for _, registerer := range registerers {
		registerer.Register(router)
	}
	c := cors.AllowAll()
	return c.Handler(bhttp.RecoverHandler(context.ClearContext(router)))
}
