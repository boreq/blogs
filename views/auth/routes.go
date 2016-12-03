package auth

import (
	"github.com/boreq/blogs/http"
	"github.com/julienschmidt/httprouter"
)

func Register(router *httprouter.Router) {
	router.GET("/signup", http.NotAuthenticatedOnly(register))
	router.POST("/signup", http.NotAuthenticatedOnly(register))
	router.GET("/signin", http.NotAuthenticatedOnly(login))
	router.POST("/signin", http.NotAuthenticatedOnly(login))
	router.GET("/signout", logout)
}
