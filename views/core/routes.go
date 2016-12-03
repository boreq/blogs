package core

import (
	"github.com/julienschmidt/httprouter"
)

func Register(router *httprouter.Router) {
	router.GET("/", index)
	router.GET("/blogs", blogs)
	router.GET("/blog/:id/:name", blog)
}
