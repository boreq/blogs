package http

import (
	"github.com/boreq/blogs/http/api"
	"github.com/boreq/blogs/http/context"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// AuthenticatedOnly returns an UnauthorizedError if the user is unauthorized
// thereby from accessing the wrapped handle.
func AuthenticatedOnly(handle api.Handle) api.Handle {
	return func(r *http.Request, p httprouter.Params) (api.Response, api.Error) {
		ctx := context.Get(r)
		if !ctx.User.IsAuthenticated() {
			return nil, api.UnauthorizedError
		}
		return handle(r, p)
	}
}
