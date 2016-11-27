package eevee

import (
	"errors"
	"github.com/boreq/blogs/blogs/loaders"
	htmlutils "github.com/boreq/blogs/utils/html"
	"golang.org/x/net/html"
	"net/http"
	"strings"
	"time"
)

const homeURL = "https://eev.ee/"
const archiveURL = "https://eev.ee/everything/archives/"

func New() loaders.Blog {
	rv := &blog{}
	return rv
}

type blog struct{}

func (b *blog) GetUrl() string {
	return "eev.ee"
}

func (b *blog) GetPostUrl(internalID string) string {
	return "eev.ee/" + internalID
}

func (b *blog) LoadTitle() (string, error) {
	// Get the page
	resp, err := http.Get(homeURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Parse the HTML response
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return "", err
	}

	// Walk the HTML tree looking for the title
	titleChan := make(chan string)
	go func() {
		defer close(titleChan)
		htmlutils.WalkAllNodes(doc, func(node *html.Node) {
			if htmlutils.IsTextNode(node) && htmlutils.IsHtmlNode(node.Parent, "title") {
				titleChan <- node.Data
			}
		})
	}()
	title, ok := <-titleChan
	if !ok {
		return "", errors.New("Title not found")
	}
	return title, nil
}

func (b *blog) LoadPosts() (<-chan loaders.Post, <-chan error) {
	postChan := make(chan loaders.Post)
	errorChan := make(chan error)
	go b.loadPosts(postChan, errorChan)
	return postChan, errorChan
}

func (b *blog) loadPosts(postChan chan<- loaders.Post, errorChan chan<- error) {
	defer close(postChan)
	defer close(errorChan)

	err := b.yieldPosts(postChan, errorChan)
	if err != nil {
		errorChan <- err
	}
}

func (b *blog) yieldPosts(postChan chan<- loaders.Post, errorChan chan<- error) error {
	// Get the page
	resp, err := http.Get(archiveURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Parse the HTML response
	doc, err := html.Parse(resp.Body)
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
	if htmlutils.IsHtmlNode(n, "a") && htmlutils.IsHtmlNode(n.Parent, "h1") {
		if val, err := htmlutils.GetAttrVal(n, "href"); err == nil {
			val = strings.TrimPrefix(val, "https://eev.ee/")
			val = strings.TrimPrefix(val, "http://eev.ee/")
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
	if htmlutils.IsTextNode(n) && htmlutils.IsHtmlNode(n.Parent, "a") && htmlutils.IsHtmlNode(n.Parent.Parent, "h1") {
		post.Title = strings.TrimSpace(n.Data)
	}

	// Tags
	if htmlutils.IsTextNode(n) && htmlutils.IsHtmlNode(n.Parent, "a") && htmlutils.HasAttrVal(n.Parent.Parent, "class", "tags") {
		post.Tags = append(post.Tags, strings.TrimPrefix(strings.TrimSpace(n.Data), "#"))
	}
}
