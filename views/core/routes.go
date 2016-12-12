package core

import (
	"github.com/julienschmidt/httprouter"
)

func Register(router *httprouter.Router) {
	router.GET("/", index)
	router.GET("/blogs", blogs)
	router.GET("/posts", posts)
	router.GET("/tags", tags)
	router.GET("/updates", updates)
	router.GET("/blog/:id/:name", blog)
	router.GET("/profile/:id", profile_stars)
	router.GET("/profile/:id/stars", profile_stars)
	router.GET("/profile/:id/subscriptions", profile_subscriptions)
}
