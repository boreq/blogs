package golang

import (
	"github.com/boreq/blogs/blogs/common"
	"github.com/boreq/blogs/blogs/loaders"
	htmlutils "github.com/boreq/blogs/utils/html"
	"golang.org/x/net/html"
	"strings"
	"time"
)

const domain = "https://blog.golang.org"
const homeURL = "https://blog.golang.org/"
const archiveURL = "https://blog.golang.org/index"

func New() loaders.Blog {
	return common.NewLoader(domain,
		archiveURL,
		loadTitle,
		isArticleNode,
		populatePost)
}

func loadTitle() (string, error) {
	return common.LoadTitle(homeURL)
}

func isArticleNode(n *html.Node) bool {
	return htmlutils.IsHtmlNode(n, "p") &&
		htmlutils.HasAttrVal(n, "class", "blogtitle")
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
		htmlutils.IsHtmlNode(n.Parent, "span") &&
		htmlutils.HasAttrVal(n.Parent, "class", "date") {
		if t, err := time.Parse("2 January 2006", n.Data); err == nil {
			post.Date = t
		}
	}

	// Title
	if htmlutils.IsTextNode(n) &&
		htmlutils.IsHtmlNode(n.Parent, "a") {
		post.Title = strings.TrimSpace(n.Data)
	}

	// Tags
	if htmlutils.IsTextNode(n) &&
		htmlutils.IsHtmlNode(n.Parent, "span") &&
		htmlutils.HasAttrVal(n.Parent, "class", "tags") {
		tags := strings.Split(strings.TrimSpace(n.Data), " ")
		for _, tag := range tags {
			tag = strings.TrimSpace(tag)
			if tag != "" {
				post.Tags = append(post.Tags, tag)
			}
		}
	}
}
