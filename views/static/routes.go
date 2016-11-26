package static

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func Register(router *httprouter.Router) {
	router.ServeFiles("/static/*filepath", http.Dir("_static/"))
}
