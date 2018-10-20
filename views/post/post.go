package post

import (
	"github.com/boreq/blogs/http/api"
	"github.com/boreq/blogs/http/context"
	"github.com/boreq/blogs/logging"
	"github.com/boreq/blogs/service/posts"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

var log = logging.New("views/post")
var invalidPostIdError = api.NewError(http.StatusBadRequest, "Invalid post id.")

func New(prefix string, postService *posts.PostsService) *Post {
	rv := &Post{
		prefix:       prefix,
		postsService: postService,
	}
	return rv
}

type Post struct {
	prefix       string
	postsService *posts.PostsService
}

func (p *Post) Register(router *httprouter.Router) {
	router.POST(p.prefix+"/:id/star", api.Wrap(p.star))
	router.POST(p.prefix+"/:id/unstar", api.Wrap(p.unstar))
}

func (b *Post) star(r *http.Request, ps httprouter.Params) (api.Response, api.Error) {
	postId, err := b.getPostId(ps)
	if err != nil {
		return nil, invalidPostIdError
	}

	ctx := context.Get(r)
	if !ctx.User.IsAuthenticated() {
		return nil, api.UnauthorizedError
	}
	userId := ctx.User.GetUser().ID

	if err := b.postsService.Star(postId, userId); err != nil {
		log.Error("star error", "err", err)
		return nil, api.InternalServerError
	}
	return api.NewResponseOk(nil), nil
}

func (b *Post) unstar(r *http.Request, ps httprouter.Params) (api.Response, api.Error) {
	postId, err := b.getPostId(ps)
	if err != nil {
		return nil, invalidPostIdError
	}

	ctx := context.Get(r)
	if !ctx.User.IsAuthenticated() {
		return nil, api.UnauthorizedError
	}
	userId := ctx.User.GetUser().ID

	if err := b.postsService.Unstar(postId, userId); err != nil {
		log.Error("unstar error", "err", err)
		return nil, api.InternalServerError
	}
	return api.NewResponseOk(nil), nil
}

func (b *Post) getPostId(ps httprouter.Params) (uint, error) {
	postId, err := strconv.ParseUint(ps.ByName("id"), 10, 32)
	return uint(postId), err
}
