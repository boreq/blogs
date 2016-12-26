package simblob

import (
	"github.com/boreq/blogs/blogs/common"
	"github.com/boreq/blogs/blogs/loaders"
	htmlutils "github.com/boreq/blogs/utils/html"
	"golang.org/x/net/html"
	"regexp"
	"strings"
	"time"
)

const domain = "simblob.blogspot.com"
const homeURL = "https://simblob.blogspot.com/"

var maxResultsRegexp = regexp.MustCompile("max-results=[0-9]+")

func New() loaders.Blog {
	return common.NewPaginatedLoader(domain,
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

func isNextPageLink(n *html.Node) bool {
	return htmlutils.IsHtmlNode(n, "a") &&
		htmlutils.HasAttrVal(n, "class", "blog-pager-older-link")
}

func getNextPageURL(n *html.Node) (string, error) {
	url, err := htmlutils.GetAttrVal(n, "href")
	if err != nil {
		return "", err
	}
	return maxResultsRegexp.ReplaceAllString(url, "max-results=50"), nil
}

func isArticleNode(n *html.Node) bool {
	return htmlutils.IsHtmlNode(n, "div") &&
		htmlutils.IsHtmlNode(n.Parent, "div") &&
		htmlutils.HasAttrVal(n.Parent, "class", "blog-posts hfeed")
}

func populatePost(n *html.Node, post *loaders.Post) {
	// Id
	if htmlutils.IsHtmlNode(n, "a") &&
		htmlutils.IsHtmlNode(n.Parent, "h2") {
		if val, err := htmlutils.GetAttrVal(n, "href"); err == nil {
			post.Id = common.CleanPostId(val, domain)
		}
	}

	// Date
	if htmlutils.IsTextNode(n) &&
		htmlutils.IsHtmlNode(n.Parent, "address") {
		if t, err := time.Parse("Monday, January 2, 2006", n.Data); err == nil {
			post.Date = t
		}

	}

	// Title
	if htmlutils.IsTextNode(n) &&
		htmlutils.IsHtmlNode(n.Parent, "a") &&
		htmlutils.IsHtmlNode(n.Parent.Parent, "h2") {
		post.Title = strings.TrimSpace(n.Data)
	}

	// Tags
	if htmlutils.IsTextNode(n) &&
		htmlutils.IsHtmlNode(n.Parent, "a") &&
		htmlutils.HasAttrVal(n.Parent.Parent, "class", "blogger-labels") {
		post.Tags = append(post.Tags, strings.TrimPrefix(strings.TrimSpace(n.Data), "#"))
	}
}
