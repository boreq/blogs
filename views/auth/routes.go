package auth

import (
	"github.com/boreq/blogs/http"
	"github.com/julienschmidt/httprouter"
)

func Register(router *httprouter.Router) {
	router.GET("/signup", http.NotAuthenticatedOnly(http.WithErrorHandling(register)))
	router.POST("/signup", http.NotAuthenticatedOnly(http.WithErrorHandling(register)))
	router.GET("/signin", http.NotAuthenticatedOnly(http.WithErrorHandling(login)))
	router.POST("/signin", http.NotAuthenticatedOnly(http.WithErrorHandling(login)))
	router.GET("/signout", logout)
}
