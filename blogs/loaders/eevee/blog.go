package eevee

import (
	"github.com/boreq/blogs/blogs/common"
	"github.com/boreq/blogs/blogs/loaders"
	htmlutils "github.com/boreq/blogs/utils/html"
	"golang.org/x/net/html"
	"strings"
	"time"
)

const domain = "eev.ee"
const homeURL = "https://eev.ee/"
const archiveURL = "https://eev.ee/everything/archives/"

func New() loaders.Blog {
	return loaders.New(domain,
		archiveURL,
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
		htmlutils.IsHtmlNode(n.Parent, "h1") {
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

	// Category
	if htmlutils.IsHtmlNode(n, "h1") {
		if val, err := htmlutils.GetAttrVal(n, "class"); err == nil {
			if parts := strings.SplitN(val, "-", 2); len(parts) == 2 {
				post.Category = parts[1]
			}
		}
	}

	// Title
	if htmlutils.IsTextNode(n) && htmlutils.IsHtmlNode(n.Parent, "a") &&
		htmlutils.IsHtmlNode(n.Parent.Parent, "h1") {
		post.Title = strings.TrimSpace(n.Data)
	}

	// Tags
	if htmlutils.IsTextNode(n) &&
		htmlutils.IsHtmlNode(n.Parent, "a") &&
		htmlutils.HasAttrVal(n.Parent.Parent, "class", "tags") {
		post.Tags = append(post.Tags, strings.TrimPrefix(strings.TrimSpace(n.Data), "#"))
	}
}
