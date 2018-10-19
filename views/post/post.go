package post

import (
	"github.com/boreq/blogs/http/api"
	"github.com/boreq/blogs/http/context"
	"github.com/boreq/blogs/logging"
	"github.com/boreq/blogs/service/post"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

var log = logging.New("views/post")

func New(prefix string, postService *post.PostService) *Post {
	rv := &Post{
		Prefix:      prefix,
		PostService: postService,
	}
	return rv
}

type Post struct {
	Prefix      string
	PostService *post.PostService
}

func (p *Post) Register(router *httprouter.Router) {
	router.POST(p.Prefix+"/:id/star", api.Wrap(p.star))
	router.POST(p.Prefix+"/:id/unstar", api.Wrap(p.unstar))
}

func (b *Post) star(r *http.Request, ps httprouter.Params) (api.Response, api.Error) {
	blogId, err := strconv.ParseUint(ps.ByName("id"), 10, 32)
	if err != nil {
		return nil, api.NewError(http.StatusBadRequest, "Invalid post id.")
	}

	ctx := context.Get(r)
	if !ctx.User.IsAuthenticated() {
		return nil, api.UnauthorizedError
	}
	userId := ctx.User.GetUser().ID

	if err := b.PostService.Star(uint(blogId), userId); err != nil {
		log.Error("star error", "err", err)
		return nil, api.InternalServerError
	}
	return api.NewResponseOk(nil), nil
}

func (b *Post) unstar(r *http.Request, ps httprouter.Params) (api.Response, api.Error) {
	blogId, err := strconv.ParseUint(ps.ByName("id"), 10, 32)
	if err != nil {
		return nil, api.NewError(http.StatusBadRequest, "Invalid post id.")
	}

	ctx := context.Get(r)
	if !ctx.User.IsAuthenticated() {
		return nil, api.UnauthorizedError
	}
	userId := ctx.User.GetUser().ID

	if err := b.PostService.Unstar(uint(blogId), userId); err != nil {
		log.Error("unstar error", "err", err)
		return nil, api.InternalServerError
	}
	return api.NewResponseOk(nil), nil
}
