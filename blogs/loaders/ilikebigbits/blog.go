package ilikebigbits

import (
	"github.com/boreq/blogs/blogs/common"
	"github.com/boreq/blogs/blogs/loaders"
	htmlutils "github.com/boreq/blogs/utils/html"
	"golang.org/x/net/html"
	"strings"
	"time"
)

const domain = "ilikebigbits.com"
const homeURL = "http://www.ilikebigbits.com/"

func New() loaders.Blog {
	return common.NewLoader(domain,
		homeURL,
		loadTitle,
		isArticleNode,
		populatePost)
}

func loadTitle() (string, error) {
	return common.LoadTitle(homeURL)
}

func isArticleNode(n *html.Node) bool {
	return htmlutils.IsHtmlNode(n, "article")
}

func populatePost(n *html.Node, post *loaders.Post) {
	// Id
	if htmlutils.IsHtmlNode(n, "a") &&
		htmlutils.IsHtmlNode(n.Parent, "h1") &&
		htmlutils.HasAttrVal(n.Parent, "class", "entry-title") {
		if val, err := htmlutils.GetAttrVal(n, "href"); err == nil {
			val = strings.TrimPrefix(val, "https://"+domain+"/")
			val = strings.TrimPrefix(val, "http://"+domain+"/")
			val = strings.Trim(val, "/")
			post.Id = val
		}
	}

	// Date
	if htmlutils.IsHtmlNode(n, "time") {
		if val, err := htmlutils.GetAttrVal(n, "datetime"); err == nil {
			if t, err := time.Parse("2006-01-02", val); err == nil {
				post.Date = t
			}
		}
	}

	// Title
	if htmlutils.IsTextNode(n) && htmlutils.IsHtmlNode(n.Parent, "a") &&
		htmlutils.IsHtmlNode(n.Parent.Parent, "h1") &&
		htmlutils.HasAttrVal(n.Parent.Parent, "class", "entry-title") {
		post.Title = strings.TrimSpace(n.Data)
	}
}
