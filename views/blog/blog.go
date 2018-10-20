package blog

import (
	"github.com/boreq/blogs/http/api"
	"github.com/boreq/blogs/logging"
	"github.com/boreq/blogs/service/blog"
	"github.com/boreq/blogs/service/context"
	"github.com/boreq/blogs/service/posts"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

var log = logging.New("views/blog")
var invalidBlogIdError = api.NewError(http.StatusBadRequest, "Invalid blog id.")

func New(prefix string, blogService *blog.BlogService, postsService *posts.PostsService, contextService *context.ContextService) *Blog {
	rv := &Blog{
		prefix:         prefix,
		blogService:    blogService,
		postsService:   postsService,
		contextService: contextService,
	}
	return rv
}

type Blog struct {
	prefix         string
	blogService    *blog.BlogService
	postsService   *posts.PostsService
	contextService *context.ContextService
}

func (b *Blog) Register(router *httprouter.Router) {
	router.GET(b.prefix+"/:id", api.Wrap(b.get))
	router.GET(b.prefix+"/:id/categories", api.Wrap(b.categories))
	router.GET(b.prefix+"/:id/tags", api.Wrap(b.tags))
	router.GET(b.prefix+"/:id/posts", api.Wrap(b.posts))
	router.POST(b.prefix+"/:id/subscribe", api.Wrap(b.subscribe))
	router.POST(b.prefix+"/:id/unsubscribe", api.Wrap(b.unsubscribe))
}

func (b *Blog) subscribe(r *http.Request, p httprouter.Params) (api.Response, api.Error) {
	blogId, err := getBlogId(p)
	if err != nil {
		return nil, invalidBlogIdError
	}

	ctx := b.contextService.Get(r)
	if !ctx.User.IsAuthenticated() {
		return nil, api.UnauthorizedError
	}
	userId := ctx.User.GetUser().ID

	if err := b.blogService.Subscribe(uint(blogId), userId); err != nil {
		log.Error("subscribe error", "err", err)
		return nil, api.InternalServerError
	}
	return api.NewResponseOk(nil), nil
}

func (b *Blog) unsubscribe(r *http.Request, p httprouter.Params) (api.Response, api.Error) {
	blogId, err := getBlogId(p)
	if err != nil {
		return nil, invalidBlogIdError
	}

	ctx := b.contextService.Get(r)
	if !ctx.User.IsAuthenticated() {
		return nil, api.UnauthorizedError
	}
	userId := ctx.User.GetUser().ID

	if err := b.blogService.Unsubscribe(uint(blogId), userId); err != nil {
		log.Error("unsubscribe error", "err", err)
		return nil, api.InternalServerError
	}
	return api.NewResponseOk(nil), nil
}

func (b *Blog) get(r *http.Request, p httprouter.Params) (api.Response, api.Error) {
	blogId, err := getBlogId(p)
	if err != nil {
		return nil, invalidBlogIdError
	}

	var userId *uint = nil
	ctx := b.contextService.Get(r)
	if ctx.User.IsAuthenticated() {
		userId = &ctx.User.GetUser().ID
	}

	blog, err := b.blogService.Get(uint(blogId), userId)
	if err != nil {
		log.Error("get error", "err", err)
		return nil, api.InternalServerError
	}
	return api.NewResponseOk(blog), nil
}

func (b *Blog) categories(r *http.Request, p httprouter.Params) (api.Response, api.Error) {
	blogId, err := getBlogId(p)
	if err != nil {
		return nil, invalidBlogIdError
	}

	categories, err := b.blogService.GetCategories(uint(blogId))
	if err != nil {
		log.Error("categories error", "err", err)
		return nil, api.InternalServerError
	}
	return api.NewResponseOk(categories), nil
}

func (b *Blog) tags(r *http.Request, p httprouter.Params) (api.Response, api.Error) {
	blogId, err := getBlogId(p)
	if err != nil {
		return nil, invalidBlogIdError
	}

	tags, err := b.blogService.GetTags(uint(blogId))
	if err != nil {
		log.Error("tags error", "err", err)
		return nil, api.InternalServerError
	}
	return api.NewResponseOk(tags), nil
}

func (b *Blog) posts(r *http.Request, p httprouter.Params) (api.Response, api.Error) {
	blogId, err := getBlogId(p)
	if err != nil {
		return nil, invalidBlogIdError
	}

	var userId *uint = nil
	ctx := b.contextService.Get(r)
	if ctx.User.IsAuthenticated() {
		userId = &ctx.User.GetUser().ID
	}

	posts, err := b.postsService.ListForBlog(uint(blogId), userId)
	if err != nil {
		log.Error("posts error", "err", err)
		return nil, api.InternalServerError
	}
	return api.NewResponseOk(posts), nil
}

func getBlogId(p httprouter.Params) (uint, error) {
	blogId, err := strconv.ParseUint(p.ByName("id"), 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(blogId), nil
}
