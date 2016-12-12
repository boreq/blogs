package common

import (
	"errors"
	"fmt"
	htmlutils "github.com/boreq/blogs/utils/html"
	"golang.org/x/net/html"
	"math/rand"
	"net/http"
)

// IoS
var userAgents = []string{
	"Mozilla/5.0 (X11; Linux x86_64; rv:50.0) Gecko/20100101 Firefox/50.0",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/54.0.2840.100 Safari/537.36",
}

// DownloadAndParse downloads a web page using a GET request and parses the
// downloaded HTML.
func DownloadAndParse(url string) (*html.Node, error) {
	// Create the request (set a random user agent)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", userAgents[rand.Intn(len(userAgents))])

	// Download the page
	resp, err := http.DefaultClient.Do(request)
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
