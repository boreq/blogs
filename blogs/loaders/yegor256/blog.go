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
const archiveURL = "http://www.yegor256.com/contents.html"

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
	return htmlutils.IsHtmlNode(n, "div") &&
		htmlutils.HasAttrVal(n.Parent, "id", "all")
}

func populatePost(n *html.Node, post *loaders.Post) {
	// Id
	if htmlutils.IsHtmlNode(n, "a") &&
		htmlutils.IsHtmlNode(n.Parent, "div") {
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
		htmlutils.IsHtmlNode(n.Parent.Parent, "div") {
		post.Title = strings.TrimSpace(n.Data)
	}

	// Tags
	if htmlutils.IsTextNode(n) &&
		htmlutils.IsHtmlNode(n.Parent, "a") {
		if class, err := htmlutils.GetAttrVal(n.Parent, "class"); err == nil {
			if strings.Contains(class, "tag") {
				post.Tags = append(post.Tags, strings.TrimSpace(n.Data))
			}
		}
	}
}
