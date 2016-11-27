package auth

import (
	"github.com/boreq/blogs/utils/http"
	"github.com/julienschmidt/httprouter"
)

func Register(router *httprouter.Router) {
	router.GET("/signup", http.WithErrorHandling(register))
	router.POST("/signup", http.WithErrorHandling(register))
	router.GET("/signin", http.WithErrorHandling(login))
	router.POST("/signin", http.WithErrorHandling(login))
}
