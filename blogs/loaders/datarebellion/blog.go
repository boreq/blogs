package datarebellion

import (
	"errors"
	"github.com/boreq/blogs/blogs/common"
	"github.com/boreq/blogs/blogs/loaders"
	htmlutils "github.com/boreq/blogs/utils/html"
	"golang.org/x/net/html"
	"strings"
	"time"
)

const domain = "https://datarebellion.com/blog"
const homeURL = "https://datarebellion.com/blog/"

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
	doc, err := common.DownloadAndParse(homeURL)
	if err != nil {
		return "", err
	}
	title := ""
	htmlutils.WalkAllNodes(doc, func(node *html.Node) {
		if htmlutils.IsTextNode(node) &&
			htmlutils.IsHtmlNode(node.Parent, "a") &&
			htmlutils.IsHtmlNode(node.Parent.Parent, "p") &&
			htmlutils.HasAttrVal(node.Parent.Parent, "class", "site-title") {
			title = node.Data
		}
	})
	if title == "" {
		return "", errors.New("Could not load the title")
	}
	return title, nil
}

func isNextPageLink(node *html.Node) bool {
	return htmlutils.IsTextNode(node) &&
		htmlutils.IsHtmlNode(node.Parent, "a") &&
		htmlutils.IsHtmlNode(node.Parent.Parent, "li") &&
		htmlutils.HasAttrVal(node.Parent.Parent, "class", "pagination-next")
}

func getNextPageURL(n *html.Node) (string, error) {
	href, err := htmlutils.GetAttrVal(n.Parent, "href")
	if err != nil {
		return "", err
	}
	return href, nil
}

func isArticleNode(n *html.Node) bool {
	return htmlutils.IsHtmlNode(n, "article")
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
	if htmlutils.IsHtmlNode(n, "time") {
		if val, err := htmlutils.GetAttrVal(n, "datetime"); err == nil {
			if t, err := time.Parse(time.RFC3339, val); err == nil {
				post.Date = t
			}
		}
	}

	// Title
	if htmlutils.IsTextNode(n) &&
		htmlutils.IsHtmlNode(n.Parent, "a") &&
		htmlutils.IsHtmlNode(n.Parent.Parent, "h2") {
		post.Title = strings.TrimSpace(n.Data)
	}
}
