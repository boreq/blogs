package blog

import (
	"github.com/boreq/blogs/http/api"
	"github.com/boreq/blogs/http/context"
	"github.com/boreq/blogs/logging"
	"github.com/boreq/blogs/service/blog"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

var log = logging.New("views/blog")

func New(prefix string, blogService *blog.BlogService) *Blog {
	rv := &Blog{
		Prefix:      prefix,
		BlogService: blogService,
	}
	return rv
}

type Blog struct {
	Prefix      string
	BlogService *blog.BlogService
}

func (b *Blog) Register(router *httprouter.Router) {
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
