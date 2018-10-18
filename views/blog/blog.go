package blog

import (
	"github.com/boreq/blogs/http/api"
	"github.com/boreq/blogs/http/context"
	"github.com/boreq/blogs/logging"
	"github.com/boreq/blogs/service/blog"
	"github.com/boreq/blogs/service/posts"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

var log = logging.New("views/blog")

func New(prefix string, blogService *blog.BlogService, postsService *posts.PostsService) *Blog {
	rv := &Blog{
		Prefix:       prefix,
		BlogService:  blogService,
		PostsService: postsService,
	}
	return rv
}

type Blog struct {
	Prefix       string
	BlogService  *blog.BlogService
	PostsService *posts.PostsService
}

func (b *Blog) Register(router *httprouter.Router) {
	router.GET(b.Prefix+"/:id", api.Wrap(b.get))
	router.GET(b.Prefix+"/:id/categories", api.Wrap(b.categories))
	router.GET(b.Prefix+"/:id/tags", api.Wrap(b.tags))
	router.GET(b.Prefix+"/:id/posts", api.Wrap(b.posts))
	router.POST(b.Prefix+"/:id/subscribe", api.Wrap(b.subscribe))
	router.POST(b.Prefix+"/:id/unsubscribe", api.Wrap(b.unsubscribe))
}

func (b *Blog) subscribe(r *http.Request, p httprouter.Params) (api.Response, api.Error) {
	blogId, err := strconv.ParseUint(p.ByName("id"), 10, 32)
	if err != nil {
		return nil, api.NewError(http.StatusBadRequest, "Invalid blog id.")
	}

	ctx := context.Get(r)
	if !ctx.User.IsAuthenticated() {
		return nil, api.UnauthorizedError
	}
	userId := ctx.User.GetUser().ID

	if err := b.BlogService.Subscribe(uint(blogId), userId); err != nil {
		log.Error("BlogService subscribe error", "err", err)
		return nil, api.InternalServerError
	}
	return api.NewResponseOk(nil), nil
}

func (b *Blog) unsubscribe(r *http.Request, p httprouter.Params) (api.Response, api.Error) {
	blogId, err := strconv.ParseUint(p.ByName("id"), 10, 32)
	if err != nil {
		return nil, api.NewError(http.StatusBadRequest, "Invalid blog id.")
	}

	ctx := context.Get(r)
	if !ctx.User.IsAuthenticated() {
		return nil, api.UnauthorizedError
	}
	userId := ctx.User.GetUser().ID

	if err := b.BlogService.Unsubscribe(uint(blogId), userId); err != nil {
		log.Error("BlogService unsubscribe error", "err", err)
		return nil, api.InternalServerError
	}
	return api.NewResponseOk(nil), nil
}

func (b *Blog) get(r *http.Request, p httprouter.Params) (api.Response, api.Error) {
	blogId, err := strconv.ParseUint(p.ByName("id"), 10, 32)
	if err != nil {
		return nil, api.NewError(http.StatusBadRequest, "Invalid blog id.")
	}

	var userId *uint = nil
	ctx := context.Get(r)
	if ctx.User.IsAuthenticated() {
		userId = &ctx.User.GetUser().ID
	}

	blog, err := b.BlogService.Get(uint(blogId), userId)
	if err != nil {
		log.Error("BlogService get error", "err", err)
		return nil, api.InternalServerError
	}
	return api.NewResponseOk(blog), nil
}

func (b *Blog) categories(r *http.Request, p httprouter.Params) (api.Response, api.Error) {
	blogId, err := strconv.ParseUint(p.ByName("id"), 10, 32)
	if err != nil {
		return nil, api.NewError(http.StatusBadRequest, "Invalid blog id.")
	}

	categories, err := b.BlogService.GetCategories(uint(blogId))
	if err != nil {
		log.Error("BlogService get categories error", "err", err)
		return nil, api.InternalServerError
	}
	return api.NewResponseOk(categories), nil
}

func (b *Blog) tags(r *http.Request, p httprouter.Params) (api.Response, api.Error) {
	blogId, err := strconv.ParseUint(p.ByName("id"), 10, 32)
	if err != nil {
		return nil, api.NewError(http.StatusBadRequest, "Invalid blog id.")
	}

	tags, err := b.BlogService.GetTags(uint(blogId))
	if err != nil {
		log.Error("BlogService get tags error", "err", err)
		return nil, api.InternalServerError
	}
	return api.NewResponseOk(tags), nil
}

func (b *Blog) posts(r *http.Request, p httprouter.Params) (api.Response, api.Error) {
	blogId, err := strconv.ParseUint(p.ByName("id"), 10, 32)
	if err != nil {
		return nil, api.NewError(http.StatusBadRequest, "Invalid blog id.")
	}

	var userId *uint = nil
	ctx := context.Get(r)
	if ctx.User.IsAuthenticated() {
		userId = &ctx.User.GetUser().ID
	}

	posts, err := b.PostsService.ListForBlog(uint(blogId), userId)
	if err != nil {
		log.Error("BlogService get posts error", "err", err)
		return nil, api.InternalServerError
	}
	return api.NewResponseOk(posts), nil
}
