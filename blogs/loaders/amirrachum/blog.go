package amirrachum

import (
	"github.com/boreq/blogs/blogs/common"
	"github.com/boreq/blogs/blogs/loaders"
	htmlutils "github.com/boreq/blogs/utils/html"
	"golang.org/x/net/html"
	"strings"
	"time"
)

const domain = "https://amir.rachum.com"
const homeURL = "https://amir.rachum.com/"

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
	return htmlutils.IsHtmlNode(n, "li") &&
		htmlutils.IsHtmlNode(n.Parent, "ul") &&
		htmlutils.IsHtmlNode(n.Parent.Parent, "section") &&
		htmlutils.IsHtmlNode(n.Parent.Parent.Parent, "main")
}

func populatePost(n *html.Node, post *loaders.Post) {
	// Id
	if htmlutils.IsHtmlNode(n, "a") &&
		htmlutils.IsHtmlNode(n.Parent, "div") &&
		htmlutils.HasAttrVal(n.Parent, "class", "title") {
		if val, err := htmlutils.GetAttrVal(n, "href"); err == nil {
			post.Id = common.CleanPostId(val, domain)
		}
	}

	// Date
	if htmlutils.IsTextNode(n) &&
		htmlutils.IsHtmlNode(n.Parent, "span") &&
		htmlutils.IsHtmlNode(n.Parent.Parent, "div") &&
		htmlutils.HasAttrVal(n.Parent.Parent, "class", "post-date") {
		if t, err := time.Parse("Jan 2, 2006", n.Data); err == nil {
			post.Date = t
		}

	}

	// Title
	if htmlutils.IsTextNode(n) &&
		htmlutils.IsHtmlNode(n.Parent, "a") {
		post.Title = strings.TrimSpace(n.Data)
	}
}
