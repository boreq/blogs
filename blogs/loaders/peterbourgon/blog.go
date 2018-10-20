package peterbourgon

import (
	"errors"
	"github.com/boreq/blogs/blogs/common"
	"github.com/boreq/blogs/blogs/loaders"
	htmlutils "github.com/boreq/blogs/utils/html"
	"golang.org/x/net/html"
	"strings"
	"time"
)

const domain = "https://peter.bourgon.org"
const homeURL = "https://peter.bourgon.org/blog/"

func New() loaders.Blog {
	return common.NewLoader(domain,
		homeURL,
		loadTitle,
		isArticleNode,
		populatePost)
}

func loadTitle() (string, error) {
	doc, err := common.DownloadAndParse(homeURL)
	if err != nil {
		return "", err
	}
	title := ""
	htmlutils.WalkAllNodes(doc, func(node *html.Node) {
		if htmlutils.IsTextNode(node) &&
			htmlutils.IsHtmlNode(node.Parent, "strong") &&
			htmlutils.IsHtmlNode(node.Parent.Parent, "p") &&
			htmlutils.IsHtmlNode(node.Parent.Parent.Parent, "div") &&
			htmlutils.HasAttrVal(node.Parent.Parent.Parent, "id", "header") {
			title = node.Data
		}
	})
	if title == "" {
		return "", errors.New("Could not load the title")
	}
	return title, nil
}

func isArticleNode(n *html.Node) bool {
	return htmlutils.IsHtmlNode(n, "p") &&
		htmlutils.HasAttrVal(n, "class", "bloglink")

}

func populatePost(n *html.Node, post *loaders.Post) {
	// Id
	if htmlutils.IsHtmlNode(n, "a") {
		if val, err := htmlutils.GetAttrVal(n, "href"); err == nil {
			post.Id = common.CleanPostId(val, domain)
		}
	}

	// Date
	if htmlutils.IsTextNode(n) &&
		htmlutils.IsHtmlNode(n.Parent, "p") {
		data := strings.TrimSpace(strings.Trim(strings.TrimSpace(n.Data), "—"))
		if t, err := time.Parse("2006 01 02", data); err == nil {
			post.Date = t
		}
	}

	// Title
	if htmlutils.IsTextNode(n) &&
		htmlutils.IsHtmlNode(n.Parent, "a") {
		post.Title = strings.TrimSpace(n.Data)
	}
}
