package core

import (
	"github.com/julienschmidt/httprouter"
)

func Register(router *httprouter.Router) {
	router.GET("/", index)
	router.GET("/blogs", blogs)
	router.GET("/posts", posts)
	router.GET("/tags", tags)
	router.POST("/post/star", star)
	router.POST("/post/unstar", unstar)
	router.POST("/blog/subscribe", subscribe)
	router.POST("/blog/unsubscribe", unsubscribe)
	router.GET("/blog/:id/:name", blog)
	router.GET("/profile/:id", profile_stars)
	router.GET("/profile/:id/stars", profile_stars)
	router.GET("/profile/:id/subscriptions", profile_subscriptions)
}
