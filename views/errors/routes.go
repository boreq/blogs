package errors

import (
	"github.com/julienschmidt/httprouter"
)

func Register(router *httprouter.Router) {
	router.HandlerFunc("GET", "/error/400", BadRequest)
	router.HandlerFunc("GET", "/error/404", NotFound)
	router.HandlerFunc("GET", "/error/500", InternalServerError)
}
