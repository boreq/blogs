package core

import (
	"github.com/boreq/blogs/templates"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) error {
	var data = templates.GetDefaultData(r)
	return templates.RenderTemplate(w, "core/index.tmpl", data)
}
