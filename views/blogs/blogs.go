package blogs

import (
	"github.com/boreq/blogs/http/api"
	"github.com/boreq/blogs/logging"
	"github.com/boreq/blogs/service/blogs"
	"github.com/boreq/blogs/views"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

var log = logging.New("views/blogs")

func New(prefix string, blogsService *blogs.BlogsService) *Blogs {
	rv := &Blogs{
		Prefix:       prefix,
		BlogsService: blogsService,
	}
	return rv
}

type Blogs struct {
	Prefix       string
	BlogsService *blogs.BlogsService
}

func (b *Blogs) Register(router *httprouter.Router) {
	router.GET(b.Prefix+"/list", api.Wrap(b.list))
}

func (b *Blogs) list(r *http.Request, p httprouter.Params) (api.Response, api.Error) {
	page := views.GetPage(r)
	reverse := views.GetSortReverse(r)
	sort, ok := map[string]blogs.ListSort{
		"title":         blogs.SortTitle,
		"subscriptions": blogs.SortSubscribers,
		"lastPost":      blogs.SortLastPost,
	}[views.GetSort(r)]
	if !ok {
		sort = blogs.SortTitle
	}
	blogs, err := b.BlogsService.List(page, sort, reverse)
	if err != nil {
		log.Error("list error", "err", err)
		return nil, api.InternalServerError
	}
	return api.NewResponseOk(blogs), nil
}
