package posts

import (
	"github.com/boreq/blogs/http/api"
	"github.com/boreq/blogs/logging"
	"github.com/boreq/blogs/service/context"
	"github.com/boreq/blogs/service/posts"
	"github.com/boreq/blogs/service/tag"
	"github.com/boreq/blogs/views"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

var log = logging.New("views/posts")

func New(prefix string, postsService *posts.PostsService, tagService *tag.TagService, contextService *context.ContextService) *Posts {
	rv := &Posts{
		Prefix:         prefix,
		postsService:   postsService,
		contextService: contextService,
	}
	return rv
}

type Posts struct {
	Prefix         string
	tagService     *tag.TagService
	postsService   *posts.PostsService
	contextService *context.ContextService
}

func (p *Posts) Register(router *httprouter.Router) {
	router.GET(p.Prefix+"/list", api.Wrap(p.list))
	router.GET(p.Prefix+"/list/subscriptions", api.Wrap(p.listFromSubscriptions))
	router.GET(p.Prefix+"/list/starred", api.Wrap(p.listStarred))
}

func (p *Posts) listFromSubscriptions(r *http.Request, ps httprouter.Params) (api.Response, api.Error) {
	page := views.GetPage(r)
	reverse := views.GetSortReverse(r)
	sort, ok := map[string]posts.ListSort{
		"date":  posts.SortDate,
		"stars": posts.SortStars,
		"title": posts.SortTitle,
	}[views.GetSort(r)]
	if !ok {
		sort = posts.SortDate
	}

	ctx := p.contextService.Get(r)
	if !ctx.User.IsAuthenticated() {
		return nil, api.UnauthorizedError
	}
	userId := ctx.User.GetUser().ID

	posts, err := p.postsService.ListFromSubscriptions(page, sort, reverse, userId)
	if err != nil {
		log.Error("listFromSubscriptions error", "err", err)
		return nil, api.InternalServerError
	}
	return api.NewResponseOk(posts), nil
}

func (p *Posts) listStarred(r *http.Request, ps httprouter.Params) (api.Response, api.Error) {
	page := views.GetPage(r)
	reverse := views.GetSortReverse(r)
	sort, ok := map[string]posts.ListSort{
		"date":  posts.SortDate,
		"stars": posts.SortStars,
		"title": posts.SortTitle,
	}[views.GetSort(r)]
	if !ok {
		sort = posts.SortDate
	}

	ctx := p.contextService.Get(r)
	if !ctx.User.IsAuthenticated() {
		return nil, api.UnauthorizedError
	}
	userId := ctx.User.GetUser().ID

	posts, err := p.postsService.ListStarred(page, sort, reverse, userId)
	if err != nil {
		log.Error("listStarred error", "err", err)
		return nil, api.InternalServerError
	}
	return api.NewResponseOk(posts), nil
}

func (p *Posts) list(r *http.Request, ps httprouter.Params) (api.Response, api.Error) {
	page := views.GetPage(r)
	reverse := views.GetSortReverse(r)
	sort, ok := map[string]posts.ListSort{
		"date":  posts.SortDate,
		"stars": posts.SortStars,
		"title": posts.SortTitle,
	}[views.GetSort(r)]
	if !ok {
		sort = posts.SortDate
	}
	var userId *uint = nil
	ctx := p.contextService.Get(r)
	if ctx.User.IsAuthenticated() {
		userId = &ctx.User.GetUser().ID
	}
	posts, err := p.postsService.List(page, sort, reverse, userId)
	if err != nil {
		log.Error("list error", "err", err)
		return nil, api.InternalServerError
	}
	return api.NewResponseOk(posts), nil
}
