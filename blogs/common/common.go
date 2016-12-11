package common

import (
	"errors"
	"fmt"
	htmlutils "github.com/boreq/blogs/utils/html"
	"golang.org/x/net/html"
	"net/http"
)

// DownloadAndParse downloads a web page using a GET request and parses the
// downloaded HTML.
func DownloadAndParse(url string) (*html.Node, error) {
	// Download the page
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Status %d", resp.StatusCode)
	}

	// Parse the HTML response
	return html.Parse(resp.Body)
}

// LoadTitle downloads the specified page, parses it and returns its title.
func LoadTitle(url string) (string, error) {
	doc, err := DownloadAndParse(url)
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
