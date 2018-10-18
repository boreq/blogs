package posts

import (
	"github.com/boreq/blogs/http/api"
	"github.com/boreq/blogs/http/context"
	"github.com/boreq/blogs/logging"
	"github.com/boreq/blogs/service/posts"
	"github.com/boreq/blogs/views"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

var log = logging.New("views/posts")

func New(prefix string, postsService *posts.PostsService) *Posts {
	rv := &Posts{
		Prefix:       prefix,
		PostsService: postsService,
	}
	return rv
}

type Posts struct {
	Prefix       string
	PostsService *posts.PostsService
}

func (p *Posts) Register(router *httprouter.Router) {
	router.GET(p.Prefix+"/list", api.Wrap(p.list))
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
	ctx := context.Get(r)
	if ctx.User.IsAuthenticated() {
		userId = &ctx.User.GetUser().ID
	}
	posts, err := p.PostsService.List(page, sort, reverse, userId)
	if err != nil {
		log.Error("list error", "err", err)
		return nil, api.InternalServerError
	}
	return api.NewResponseOk(posts), nil
}
