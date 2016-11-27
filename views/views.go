package views

import (
	"github.com/boreq/blogs/views/auth"
	"github.com/boreq/blogs/views/core"
	"github.com/boreq/blogs/views/static"
	"github.com/julienschmidt/httprouter"
)

func Register(router *httprouter.Router) {
	auth.Register(router)
	core.Register(router)
	static.Register(router)
}
