package views

import (
	"github.com/boreq/blogs/views/core"
	"github.com/julienschmidt/httprouter"
)

func Register(router *httprouter.Router) {
	core.Register(router)
}
