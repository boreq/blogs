package core

import (
	"github.com/boreq/blogs/templates"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	templates.RenderTemplate(w, "core/index.tmpl", nil)
}

func hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var data map[string]interface{} = make(map[string]interface{})
	data["name"] = ps.ByName("name")
	templates.RenderTemplate(w, "core/hello.tmpl", data)
}
