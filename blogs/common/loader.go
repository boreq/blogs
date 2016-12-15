package common

import (
	"errors"
	"github.com/boreq/blogs/blogs/loaders"
	"github.com/boreq/blogs/logging"
	htmlutils "github.com/boreq/blogs/utils/html"
	"golang.org/x/net/html"
	"sync"
)

type NodeCheckFunc func(n *html.Node) bool
type LoadTitleFunc func() (string, error)
type NextPageFunc func(n *html.Node) (string, error)
type PopulatePostFunc func(n *html.Node, post *loaders.Post)

// NewPaginatedLoader creates a loader which walks the HTML tree of the website
// downloaded from the home url and advances to the next pages.
//
//     for all nodes:
//         if isArticleNode(node)
//             for all node.children
//                 populatePost(node, post)
//         if isNextPageLink(node)
//             getNextPageURL(node)
//             repat the algorithm for all nodes of the next page
//
// All functions must be thread safe.
func NewPaginatedLoader(domain string, homeURL string, loadTitle LoadTitleFunc, isArticleNode NodeCheckFunc, populatePost PopulatePostFunc, isNextPageLink NodeCheckFunc, getNextPageURL NextPageFunc) loaders.Blog {
	rv := loader{
		domain:         domain,
		homeURL:        homeURL,
		loadTitle:      loadTitle,
		isArticleNode:  isArticleNode,
		populatePost:   populatePost,
		isNextPageLink: isNextPageLink,
		getNextPageURL: getNextPageURL,
		log:            logging.GetLogger(domain),
	}
	return rv
}

// NewLoader creates a loader which walks the HTML tree of the website
// downloaded from the home url.
//
//     for all nodes:
//         if isArticleNode(node)
//             for all node.children
//                 populatePost(node, post)
//
// All functions must be thread safe.
func NewLoader(domain string, homeURL string, loadTitle LoadTitleFunc, isArticleNode NodeCheckFunc, populatePost PopulatePostFunc) loaders.Blog {
	isNextPageLink := func(n *html.Node) bool {
		return false
	}
	getNextPageURL := func(n *html.Node) (string, error) {
		return "", errors.New("Not implemented")
	}
	rv := loader{
		domain:         domain,
		homeURL:        homeURL,
		loadTitle:      loadTitle,
		isArticleNode:  isArticleNode,
		populatePost:   populatePost,
		isNextPageLink: isNextPageLink,
		getNextPageURL: getNextPageURL,
		log:            logging.GetLogger(domain),
	}
	return rv
}

type loader struct {
	// Domain of the blog. A domain should have no protocol and no tailing
	// slashes. Example: "example.com/blog". The domain is used to generate
	// a link to the blog (simply by prefixing it with a protocol) and
	// generate post URLs from internal IDs (by joining the internal id
	// and the domain with a slash).
	domain string

	// The first page which is scanned to load the posts. The loader walks
	// the html tree; if isArticleNode returns true foe a tree,
	// populatePost is called on all its children. If isNextPageLink
	// returns true for a ndoe, getNextPageURL is called on that node to
	// advance to the next page.
	homeURL string

	// This function should load the title of the blog.
	loadTitle LoadTitleFunc

	// If this function returns true, populatePost will be called on the
	// node's children.
	isArticleNode NodeCheckFunc

	// Is called repeatedly with the same post as an argument on all
	// children of an article node and it should populate that post with
	// data.
	populatePost PopulatePostFunc

	// Should return true if this is the next page link. GetNextPageURL will
	// be called on the nodes for which this function returns true.
	isNextPageLink NodeCheckFunc

	// Should extract the URL of the next page from a node for which
	// isNextPageLink returns true.
	getNextPageURL NextPageFunc

	log logging.Logger
}

func (l loader) GetUrl() string {
	return l.domain
}

func (l loader) GetPostUrl(internalID string) string {
	return l.domain + "/" + internalID
}

func (l loader) LoadTitle() (string, error) {
	return l.loadTitle()
}

func (l loader) LoadPosts() (<-chan loaders.Post, <-chan error) {
	postChan := make(chan loaders.Post)
	errorChan := make(chan error)
	go func() {
		defer close(postChan)
		defer close(errorChan)
		if err := l.yieldPosts(postChan, errorChan); err != nil {
			errorChan <- err
		}
	}()
	return postChan, errorChan
}

func (l loader) yieldPosts(postChan chan<- loaders.Post, errorChan chan<- error) error {
	wg := &sync.WaitGroup{}
	l.startPageWorker(l.homeURL, postChan, errorChan, wg)
	wg.Wait()
	return nil
}

func (l loader) startPageWorker(url string, postChan chan<- loaders.Post, errorChan chan<- error, wg *sync.WaitGroup) {
	l.log.Debugf("Starting a page worker for %s", url)
	wg.Add(1)
	go l.pageWorker(url, postChan, errorChan, wg)
}

func (l loader) pageWorker(url string, postChan chan<- loaders.Post, errorChan chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()

	doc, err := DownloadAndParse(url)
	if err != nil {
		errorChan <- err
		return
	}

	postsWg := &sync.WaitGroup{}
	// Walk the HTML tree emitting posts
	htmlutils.WalkAllNodes(doc, func(node *html.Node) {
		if l.isArticleNode(node) {
			postsWg.Add(1)
			go l.yieldPost(node, postChan, errorChan, postsWg)
		}

		if l.isNextPageLink(node) {
			url, err := l.getNextPageURL(node)
			if err == nil {
				l.startPageWorker(url, postChan, errorChan, wg)
			}
		}
	})
	postsWg.Wait()
}

func (l loader) yieldPost(n *html.Node, postChan chan<- loaders.Post, errorChan chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()

	post := loaders.Post{}
	htmlutils.WalkAllNodes(n, func(node *html.Node) {
		l.populatePost(node, &post)
	})
	postChan <- post
}
