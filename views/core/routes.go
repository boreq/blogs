package core

import (
	"github.com/boreq/blogs/http"
	"github.com/julienschmidt/httprouter"
)

func Register(router *httprouter.Router) {
	router.GET("/", http.WithErrorHandling(index))
	router.GET("/blog/:id/:name", http.WithErrorHandling(blog))
}
