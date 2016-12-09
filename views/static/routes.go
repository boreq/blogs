package static

import (
	"github.com/boreq/blogs/config"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func Register(router *httprouter.Router) {
	router.ServeFiles("/static/*filepath", http.Dir(config.Config.StaticDirectory))
}
