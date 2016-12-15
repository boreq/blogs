// Package context provides access to a request-bound structure holding
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

func createContext(r *http.Request) (*Context, error) {
	user, err := auth.GetUser(r)
	if err != nil {
		return nil, err
	}

	ctx := &Context{
		User: user,
	}
	return ctx, nil
}

// Get returns the context associated with a given request. This function will
// panic if the context can't be retrieved.
func Get(r *http.Request) *Context {
	mutex.Lock()
	defer mutex.Unlock()

	ctx, ok := contexts[r]
	if !ok {
		var err error
		ctx, err = createContext(r)
		if err != nil {
			panic(err)
		}
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

// ClearContext is a middleware that removes a context once a request is
// finished. This is a recommended way of automatically removing contexts
// associated with requests.
func ClearContext(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer Clear(r)
		h.ServeHTTP(w, r)
	})
}
