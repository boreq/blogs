package h2co3

import (
	"github.com/boreq/blogs/blogs/common"
	"github.com/boreq/blogs/blogs/loaders"
	htmlutils "github.com/boreq/blogs/utils/html"
	"golang.org/x/net/html"
	"strings"
	"time"
)

const domain = "h2co3.org/blog"
const homeURL = "http://h2co3.org/blog/"

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
	return htmlutils.IsHtmlNode(n, "header") &&
		htmlutils.IsHtmlNode(n.Parent, "div") &&
		htmlutils.IsHtmlNode(n.Parent.Parent, "div") &&
		htmlutils.HasAttrVal(n.Parent.Parent, "id", "main")
}

func populatePost(n *html.Node, post *loaders.Post) {
	// Id
	if htmlutils.IsHtmlNode(n, "a") &&
		htmlutils.IsHtmlNode(n.Parent, "h1") {
		if val, err := htmlutils.GetAttrVal(n, "href"); err == nil {
			val = strings.TrimPrefix(val, "https://"+domain+"/")
			val = strings.TrimPrefix(val, "http://"+domain+"/")
			val = strings.Trim(val, "/")
			post.Id = val
		}
	}

	// Date
	if htmlutils.IsTextNode(n) &&
		htmlutils.IsHtmlNode(n.Parent, "div") &&
		htmlutils.HasAttrVal(n.Parent, "class", "post-meta") {
		if strings.HasPrefix(n.Data, "On") {
			val := strings.TrimSpace(strings.TrimPrefix(n.Data, "On"))
			if t, err := time.Parse("2 Jan, 2006", val); err == nil {
				post.Date = t
			}
		}
	}

	// Title
	if htmlutils.IsTextNode(n) &&
		htmlutils.IsHtmlNode(n.Parent, "a") &&
		htmlutils.IsHtmlNode(n.Parent.Parent, "h1") {
		post.Title = strings.TrimSpace(n.Data)
	}
}
