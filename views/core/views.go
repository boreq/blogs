package core

import (
	"github.com/boreq/blogs/database"
	bhttp "github.com/boreq/blogs/http"
	"github.com/boreq/blogs/http/context"
	"github.com/boreq/blogs/views/errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func subscribe(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	blog_id, err := strconv.ParseUint(r.FormValue("blog_id"), 10, 32)
	if err != nil {
		errors.BadRequest(w, r)
		return
	}

	ctx := context.Get(r)
	if !ctx.User.IsAuthenticated() {
		bhttp.RedirectOrNext(w, r, "/")
		return
	}
	user_id := ctx.User.GetUser().ID

	if _, err := database.DB.Exec(`
		INSERT INTO subscription(blog_id, user_id) 
		SELECT $1, $2
		WHERE NOT EXISTS(
			SELECT 1
			FROM subscription
			WHERE blog_id=$1 AND user_id=$2)`,
		blog_id, user_id); err != nil {
		errors.InternalServerError(w, r)
		return
	}

	bhttp.RedirectOrNext(w, r, "/")
}

func unsubscribe(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	blog_id, err := strconv.ParseUint(r.FormValue("blog_id"), 10, 32)
	if err != nil {
		errors.BadRequest(w, r)
		return
	}

	ctx := context.Get(r)
	if !ctx.User.IsAuthenticated() {
		bhttp.RedirectOrNext(w, r, "/")
		return
	}
	user_id := ctx.User.GetUser().ID

	if _, err := database.DB.Exec(`
		DELETE FROM subscription
		WHERE blog_id=$1 AND user_id=$2`,
		blog_id, user_id); err != nil {
		errors.InternalServerError(w, r)
		return
	}

	bhttp.RedirectOrNext(w, r, "/")
}
