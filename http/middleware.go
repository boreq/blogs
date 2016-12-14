package http

import (
	"errors"
	verrors "github.com/boreq/blogs/views/errors"
	"net/http"
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
				verrors.InternalServerErrorWithStack(w, r, err)
			}
			return

		}()
		h.ServeHTTP(w, r)
	})
}
