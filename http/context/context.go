// Package context provides easy access to a request bound structure holding
// additional parameters. Those commonly used parameters include, for example,
// information about the logged in user.
package context

import (
	"github.com/boreq/blogs/auth"
	"net/http"
	"sync"
)

var mutex sync.Mutex
var contexts = make(map[*http.Request]*Context)

type Context struct {
	User auth.User
}

func createContext(r *http.Request) *Context {
	ctx := &Context{
		User: auth.GetUser(r),
	}
	return ctx
}

// Get returns the context associated with a given request.
func Get(r *http.Request) *Context {
	mutex.Lock()
	defer mutex.Unlock()
	ctx, ok := contexts[r]
	if !ok {
		ctx = createContext(r)
		contexts[r] = ctx
	}
	return ctx
}

// Clear removes the context associated with a given request. This can be done
// to, for example, make sure that the next caller of get will receive a fresh
// context.
func Clear(r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()
	delete(contexts, r)
}

// ClearHandler is a middleware that removes a context once a request is
// finished. This is a recommended way of automatically removing contexts
// associated with requests.
func ClearHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer Clear(r)
		h.ServeHTTP(w, r)
	})
}
