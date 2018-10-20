// Package context provides access to a request-bound structure holding
// additional parameters. Those commonly used parameters include, for example,
// information about the logged in user.
package context

import (
	"github.com/boreq/blogs/service/auth"
	"net/http"
	"sync"
)

type Context struct {
	User auth.User
}

func New(authService *auth.AuthService) *ContextService {
	rv := &ContextService{
		authService: authService,
		contexts:    make(map[*http.Request]*Context),
	}
	return rv
}

type ContextService struct {
	authService *auth.AuthService
	mutex       sync.Mutex
	contexts    map[*http.Request]*Context
}

func (c *ContextService) createContext(r *http.Request) (*Context, error) {
	user, err := c.authService.GetUser(r)
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
func (c *ContextService) Get(r *http.Request) *Context {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	ctx, ok := c.contexts[r]
	if !ok {
		var err error
		ctx, err = c.createContext(r)
		if err != nil {
			panic(err)
		}
		c.contexts[r] = ctx
	}
	return ctx
}

// Clear removes the context associated with a given request. This can be done
// to, for example, make sure that the next caller of get will receive a fresh
// context.
func (c *ContextService) Clear(r *http.Request) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.contexts, r)
}

// ClearContext is a middleware that removes a context once a request is
// finished. This is a recommended way of automatically removing contexts
// associated with requests.
func (c *ContextService) ClearContext(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer c.Clear(r)
		h.ServeHTTP(w, r)
	})
}
