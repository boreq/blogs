package lucumr

import (
	"errors"
	"github.com/boreq/blogs/blogs/common"
	"github.com/boreq/blogs/blogs/loaders"
	htmlutils "github.com/boreq/blogs/utils/html"
	"golang.org/x/net/html"
	"strings"
	"time"
)

const domain = "lucumr.pocoo.org"
const homeURL = "http://lucumr.pocoo.org/"

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
	titleChan := make(chan string)
	go func() {
		defer close(titleChan)
		htmlutils.WalkAllNodes(doc, func(node *html.Node) {
			if htmlutils.HasAttrVal(node, "class", "header") {
				title := ""
				htmlutils.WalkAllNodes(node, func(node *html.Node) {
					if htmlutils.IsTextNode(node) {
						title += node.Data
					}
				})
				titleChan <- strings.TrimSpace(title)
				return
			}
		})
	}()
	title, ok := <-titleChan
	if !ok {
		return "", errors.New("Title not found")
	}
	return title, nil
}

func isNextPageLink(node *html.Node) bool {
	return htmlutils.IsTextNode(node) &&
		strings.Contains(strings.ToLower(node.Data), "next") &&
		htmlutils.IsHtmlNode(node.Parent, "a") &&
		htmlutils.IsHtmlNode(node.Parent.Parent, "div") &&
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
	return htmlutils.HasAttrVal(n, "class", "entry-overview")
}

func populatePost(n *html.Node, post *loaders.Post) {
	// Id
	if htmlutils.IsHtmlNode(n, "a") &&
		htmlutils.IsHtmlNode(n.Parent, "h1") {
		if val, err := htmlutils.GetAttrVal(n, "href"); err == nil {
			post.Id = common.CleanPostId(val, domain)
		}
	}

	// Date
	if htmlutils.IsTextNode(n) &&
		htmlutils.HasAttrVal(n.Parent, "class", "date") {
		if t, err := time.Parse("Jan 2, 2006", n.Data); err == nil {
			post.Date = t
		}

	}

	// Title
	if htmlutils.IsTextNode(n) &&
		htmlutils.IsHtmlNode(n.Parent, "a") &&
		htmlutils.IsHtmlNode(n.Parent.Parent, "h1") {
		post.Title = strings.TrimSpace(n.Data)
	}

	// Summary
	if htmlutils.IsTextNode(n) &&
		htmlutils.IsHtmlNode(n.Parent, "p") &&
		htmlutils.IsHtmlNode(n.Parent.Parent, "div") &&
		htmlutils.HasAttrVal(n.Parent.Parent, "class", "summary") {
		post.Summary = strings.TrimSpace(n.Data)
	}

	// Follow the link and populate the tags
	if htmlutils.IsHtmlNode(n, "a") &&
		htmlutils.IsHtmlNode(n.Parent, "h1") {
		if href, err := htmlutils.GetAttrVal(n, "href"); err == nil {
			href = homeURL + "/" + strings.TrimLeft(href, "/")
			populatePostTags(href, post)
		}
	}
}

func populatePostTags(url string, post *loaders.Post) {
	if doc, err := common.DownloadAndParse(url); err == nil {
		htmlutils.WalkAllNodes(doc, func(n *html.Node) {
			if htmlutils.IsTextNode(n) &&
				htmlutils.IsHtmlNode(n.Parent, "a") &&
				htmlutils.IsHtmlNode(n.Parent.Parent, "p") &&
				htmlutils.HasAttrVal(n.Parent.Parent, "class", "tags") {
				post.Tags = append(post.Tags, strings.TrimSpace(n.Data))
			}
		})
	}
}
