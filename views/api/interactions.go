package api

import (
	"github.com/boreq/blogs/database"
	bhttp "github.com/boreq/blogs/http"
	"github.com/boreq/blogs/http/api"
	"github.com/boreq/blogs/http/context"
	verrors "github.com/boreq/blogs/views/errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"time"
)

func callAndHandleErrors(w http.ResponseWriter, r *http.Request, p httprouter.Params, f api.Handle) {
	_, apiErr := f(r, p)
	if apiErr != nil {
		switch apiErr.GetCode() {
		case http.StatusBadRequest:
			verrors.BadRequest(w, r)
			return
		case http.StatusUnauthorized:
			bhttp.RedirectOrNext(w, r, "/")
			return
		default:
			verrors.InternalServerErrorWithStack(w, r, apiErr)
			return
		}
	}
	bhttp.RedirectOrNext(w, r, "/")
}

func internalSubscribe(r *http.Request, _ httprouter.Params) (interface{}, api.Error) {
	blog_id, err := strconv.ParseUint(r.FormValue("blog_id"), 10, 32)
	if err != nil {
		return nil, api.NewError(http.StatusBadRequest, "Invalid blog id.")
	}

	ctx := context.Get(r)
	if !ctx.User.IsAuthenticated() {
		return nil, api.NewError(http.StatusUnauthorized, "Sign in to subscribe to blogs.")
	}
	userId := ctx.User.GetUser().ID

	if _, err := database.DB.Exec(`
		INSERT INTO subscription(blog_id, user_id, date)
		SELECT $1, $2, $3
		WHERE NOT EXISTS(
			SELECT 1
			FROM subscription
			WHERE blog_id=$1 AND user_id=$2)`,
		blog_id, userId, time.Now().UTC()); err != nil {
		return nil, api.NewError(http.StatusInternalServerError, "Internal server error.")
	}

	return struct{}{}, nil
}

func subscribe(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	callAndHandleErrors(w, r, p, internalSubscribe)
}

func subscribeAjax(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	api.Call(w, r, p, internalSubscribe)
}

func internalUnsubscribe(r *http.Request, _ httprouter.Params) (interface{}, api.Error) {
	blogId, err := strconv.ParseUint(r.FormValue("blog_id"), 10, 32)
	if err != nil {
		return nil, api.NewError(http.StatusBadRequest, "Invalid blog id.")
	}

	ctx := context.Get(r)
	if !ctx.User.IsAuthenticated() {
		return nil, api.NewError(http.StatusUnauthorized, "Sign in to unsubscribe from blogs.")
	}
	userId := ctx.User.GetUser().ID

	if _, err := database.DB.Exec(`
		DELETE FROM subscription
		WHERE blog_id=$1 AND user_id=$2`,
		blogId, userId); err != nil {
		return nil, api.NewError(http.StatusInternalServerError, "Internal server error.")
	}

	return struct{}{}, nil
}

func unsubscribe(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	callAndHandleErrors(w, r, p, internalUnsubscribe)
}

func unsubscribeAjax(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	api.Call(w, r, p, internalUnsubscribe)
}

func internalStar(r *http.Request, _ httprouter.Params) (interface{}, api.Error) {
	postId, err := strconv.ParseUint(r.FormValue("post_id"), 10, 32)
	if err != nil {
		return nil, api.NewError(http.StatusBadRequest, "Invalid post id.")
	}
	ctx := context.Get(r)
	if !ctx.User.IsAuthenticated() {
		return nil, api.NewError(http.StatusUnauthorized, "Sign in to star posts.")
	}
	userId := ctx.User.GetUser().ID
	if _, err := database.DB.Exec(`
		INSERT INTO star(post_id, user_id, date)
		SELECT $1, $2, $3
		WHERE NOT EXISTS(
			SELECT 1
			FROM star
			WHERE post_id=$1 AND user_id=$2)`,
		postId, userId, time.Now().UTC()); err != nil {
		return nil, api.NewError(http.StatusInternalServerError, "Internal server error.")
	}
	return struct{}{}, nil
}

func star(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	callAndHandleErrors(w, r, p, internalStar)
}

func starAjax(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	api.Call(w, r, p, internalStar)
}

func internalUnstar(r *http.Request, _ httprouter.Params) (interface{}, api.Error) {
	postId, err := strconv.ParseUint(r.FormValue("post_id"), 10, 32)
	if err != nil {
		return nil, api.NewError(http.StatusBadRequest, "Invalid post id.")
	}
	ctx := context.Get(r)
	if !ctx.User.IsAuthenticated() {
		return nil, api.NewError(http.StatusUnauthorized, "Sign in to unstar posts.")
	}
	userId := ctx.User.GetUser().ID
	if _, err := database.DB.Exec(`
		DELETE FROM star
		WHERE post_id=$1 AND user_id=$2`,
		postId, userId); err != nil {
		return nil, api.NewError(http.StatusInternalServerError, "Internal server error.")
	}
	return struct{}{}, nil
}

func unstar(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	callAndHandleErrors(w, r, p, internalUnstar)
}

func unstarAjax(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	api.Call(w, r, p, internalUnstar)
}
