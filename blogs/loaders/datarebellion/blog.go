package datarebellion

import (
	"errors"
	"github.com/boreq/blogs/blogs/common"
	"github.com/boreq/blogs/blogs/loaders"
	"github.com/boreq/blogs/logging"
	htmlutils "github.com/boreq/blogs/utils/html"
	"golang.org/x/net/html"
	"strings"
	"sync"
	"time"
)

var log = logging.GetLogger("datarebellion")

const domain = "datarebellion.com/blog"
const homeURL = "https://datarebellion.com/blog/"

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
	doc, err := common.DownloadAndParse(homeURL)
	if err != nil {
		return "", err
	}
	title := ""
	// Walk the HTML tree emitting posts
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
	wg := &sync.WaitGroup{}
	startPageWorker(homeURL, postChan, errorChan, wg)
	wg.Wait()
	return nil
}

func startPageWorker(url string, postChan chan<- loaders.Post, errorChan chan<- error, wg *sync.WaitGroup) {
	log.Debugf("Starting a page worker for %s", url)
	wg.Add(1)
	go pageWorker(url, postChan, errorChan, wg)
}

func pageWorker(url string, postChan chan<- loaders.Post, errorChan chan<- error, wg *sync.WaitGroup) {
	doc, err := common.DownloadAndParse(url)
	if err != nil {
		errorChan <- err
		return
	}

	postsWg := &sync.WaitGroup{}
	// Walk the HTML tree emitting posts
	htmlutils.WalkAllNodes(doc, func(node *html.Node) {
		if htmlutils.IsHtmlNode(node, "article") {
			postsWg.Add(1)
			go yieldPost(node, postChan, errorChan, postsWg)
		}

		if isPaginationNextLink(node) {
			if href, err := htmlutils.GetAttrVal(node.Parent, "href"); err == nil {
				startPageWorker(href, postChan, errorChan, wg)
			}
		}
	})
	postsWg.Wait()

	wg.Done()
}

func isPaginationNextLink(node *html.Node) bool {
	return htmlutils.IsTextNode(node) &&
		htmlutils.IsHtmlNode(node.Parent, "a") &&
		htmlutils.IsHtmlNode(node.Parent.Parent, "li") &&
		htmlutils.HasAttrVal(node.Parent.Parent, "class", "pagination-next")
}

func yieldPost(n *html.Node, postChan chan<- loaders.Post, errorChan chan<- error, wg *sync.WaitGroup) {
	post := loaders.Post{}
	htmlutils.WalkAllNodes(n, func(node *html.Node) {
		populatePost(node, &post)
	})
	postChan <- post
	wg.Done()
}

func populatePost(n *html.Node, post *loaders.Post) {
	// Id
	if htmlutils.IsHtmlNode(n, "a") &&
		htmlutils.IsHtmlNode(n.Parent, "h2") {
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
		htmlutils.IsHtmlNode(n.Parent, "a") &&
		htmlutils.IsHtmlNode(n.Parent.Parent, "h2") {
		post.Title = strings.TrimSpace(n.Data)
	}
}
