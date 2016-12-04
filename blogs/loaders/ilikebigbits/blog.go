package ilikebigbits

import (
	"github.com/boreq/blogs/blogs/common"
	"github.com/boreq/blogs/blogs/loaders"
	htmlutils "github.com/boreq/blogs/utils/html"
	"golang.org/x/net/html"
	"strings"
	"time"
)

const domain = "ilikebigbits.com"
const homeURL = "http://www.ilikebigbits.com/"

func New() loaders.Blog {
	rv := &blog{}
	return rv
}

type blog struct{}

func (b *blog) GetUrl() string {
	return domain
}

func (b *blog) GetPostUrl(internalID string) string {
	return domain + "/" + internalID
}

func (b *blog) LoadTitle() (string, error) {
	return common.LoadTitle(homeURL)
}

func (b *blog) LoadPosts() (<-chan loaders.Post, <-chan error) {
	postChan := make(chan loaders.Post)
	errorChan := make(chan error)
	go func() {
		defer close(postChan)
		defer close(errorChan)
		if err := b.yieldPosts(postChan, errorChan); err != nil {
			errorChan <- err
		}
	}()
	return postChan, errorChan
}

func (b *blog) yieldPosts(postChan chan<- loaders.Post, errorChan chan<- error) error {
	doc, err := common.DownloadAndParse(homeURL)
	if err != nil {
		return err
	}

	// Walk the HTML tree emitting posts
	htmlutils.WalkAllNodes(doc, func(node *html.Node) {
		if htmlutils.IsHtmlNode(node, "article") {
			yieldPost(node, postChan, errorChan)
		}
	})

	return nil
}

func yieldPost(n *html.Node, postChan chan<- loaders.Post, errorChan chan<- error) {
	post := loaders.Post{}
	htmlutils.WalkAllNodes(n, func(node *html.Node) {
		populatePost(node, &post)
	})
	postChan <- post
}

func populatePost(n *html.Node, post *loaders.Post) {
	// Id
	if htmlutils.IsHtmlNode(n, "a") &&
		htmlutils.IsHtmlNode(n.Parent, "h1") &&
		htmlutils.HasAttrVal(n.Parent, "class", "entry-title") {
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
			if t, err := time.Parse("2006-01-02", val); err == nil {
				post.Date = t
			}
		}
	}

	// Title
	if htmlutils.IsTextNode(n) && htmlutils.IsHtmlNode(n.Parent, "a") &&
		htmlutils.IsHtmlNode(n.Parent.Parent, "h1") &&
		htmlutils.HasAttrVal(n.Parent.Parent, "class", "entry-title") {
		post.Title = strings.TrimSpace(n.Data)
	}
}
