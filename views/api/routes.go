package api

import (
	"github.com/boreq/blogs/http/api"
	"github.com/julienschmidt/httprouter"
)

func Register(router *httprouter.Router) {
	router.GET("/api/updates/chart.json", api.Wrap(updatesChart))
}
