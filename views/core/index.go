package core

import (
	"github.com/boreq/blogs/templates"
	"github.com/boreq/blogs/views/errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var data = templates.GetDefaultData(r)
	if err := templates.RenderTemplateSafe(w, "core/index.tmpl", data); err != nil {
		errors.InternalServerError(w, r)
		return
	}
}
