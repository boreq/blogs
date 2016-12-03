package core

import (
	"github.com/julienschmidt/httprouter"
)

func Register(router *httprouter.Router) {
	router.GET("/", index)
	router.GET("/blogs", blogs)
	router.GET("/posts", posts)
	router.POST("/blog/subscribe", subscribe)
	router.POST("/blog/unsubscribe", unsubscribe)
	router.GET("/blog/:id/:name", blog)
}
