package yegor256

import (
	"github.com/boreq/blogs/blogs/common"
	"github.com/boreq/blogs/blogs/loaders"
	htmlutils "github.com/boreq/blogs/utils/html"
	"golang.org/x/net/html"
	"strings"
	"time"
)

const domain = "yegor256.com"
const homeURL = "http://www.yegor256.com/"

func New() loaders.Blog {
	return loaders.NewPaginated(domain,
		homeURL,
		loadTitle,
		isArticleNode,
		populatePost,
		isNextPageLink,
		getNextPageURL)
}

func loadTitle() (string, error) {
	return common.LoadTitle(homeURL)
}

func isNextPageLink(node *html.Node) bool {
	return htmlutils.IsHtmlNode(node, "a") &&
		htmlutils.IsHtmlNode(node.Parent, "div") &&
		htmlutils.HasAttrVal(node.Parent.Parent, "class", "pagination")
}

func getNextPageURL(n *html.Node) (string, error) {
	href, err := htmlutils.GetAttrVal(n.Parent, "href")
	if err != nil {
		return "", err
	}
	href = homeURL + strings.TrimLeft(href, "/")
	return href, nil
}

func isArticleNode(n *html.Node) bool {
	return htmlutils.HasAttrVal(n, "itemprop", "blogPosts")
}

func populatePost(n *html.Node, post *loaders.Post) {
	// Id
	if htmlutils.IsHtmlNode(n, "a") &&
		htmlutils.IsHtmlNode(n.Parent, "h1") &&
		htmlutils.IsHtmlNode(n.Parent.Parent, "header") {
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
			if t, err := time.Parse(time.RFC3339, val); err == nil {
				post.Date = t
			}
		}
	}

	// Title
	if htmlutils.IsTextNode(n) &&
		htmlutils.IsHtmlNode(n.Parent, "span") &&
		htmlutils.IsHtmlNode(n.Parent.Parent, "a") &&
		htmlutils.IsHtmlNode(n.Parent.Parent.Parent, "h1") {
		post.Title = strings.TrimSpace(n.Data)
	}

	// Summary
	if htmlutils.IsTextNode(n) &&
		htmlutils.IsHtmlNode(n.Parent, "p") &&
		htmlutils.IsHtmlNode(n.Parent.Parent, "div") &&
		htmlutils.HasAttrVal(n.Parent.Parent, "itemprop", "description") {
		post.Summary += strings.TrimSpace(n.Data)
	}
}
