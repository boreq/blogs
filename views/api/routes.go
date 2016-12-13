package api

import (
	"github.com/boreq/blogs/http/api"
	"github.com/julienschmidt/httprouter"
)

func Register(router *httprouter.Router) {
	router.GET("/api/updates/chart.json", api.Wrap(updatesChart))

	// Star, unstar
	router.POST("/nojs/post/star", star)
	router.POST("/nojs/post/unstar", unstar)
	router.POST("/api/post/star", starAjax)
	router.POST("/api/post/unstar", unstarAjax)

	// Subscribe, unsubscribe
	router.POST("/nojs/blog/subscribe", subscribe)
	router.POST("/nojs/blog/unsubscribe", unsubscribe)
	router.POST("/api/blog/subscribe", subscribeAjax)
	router.POST("/api/blog/unsubscribe", unsubscribeAjax)

	// Settings
	router.POST("/nojs/session/remove", removeSession)
}
