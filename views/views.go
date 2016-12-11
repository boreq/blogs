package views

import (
	"github.com/boreq/blogs/views/api"
	"github.com/boreq/blogs/views/auth"
	"github.com/boreq/blogs/views/core"
	"github.com/boreq/blogs/views/errors"
	"github.com/boreq/blogs/views/static"
	"github.com/julienschmidt/httprouter"
)

func Register(router *httprouter.Router) {
	api.Register(router)
	auth.Register(router)
	core.Register(router)
	errors.Register(router)
	static.Register(router)
}
