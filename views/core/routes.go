package core

import (
	"github.com/julienschmidt/httprouter"
)

func Register(router *httprouter.Router) {
	router.GET("/", index)
	router.GET("/hello/:name", hello)
}
