// Package html contains functions related to parsing HTML using the
// golang.org/x/net/html package.
package html

import (
	"errors"
	"golang.org/x/net/html"
)

// GetAttrVal returns a value of the specified attribute of this node. An error
// is returned if the attribute doesn't exist or if the node is nil.
func GetAttrVal(n *html.Node, key string) (string, error) {
	if n != nil {
		for _, attr := range n.Attr {
			if attr.Key == key {
				return attr.Val, nil
			}
		}
	}
	return "", errors.New("No such attribute")
}

// HasAttr returns true if this node has the specified attribute.  Returns false
// if the node is nil.
func HasAttr(n *html.Node, key string) bool {
	_, err := GetAttrVal(n, key)
	return err == nil
}

// HasAttrVal returns true if this node has the specified attribute of the
// specified value. Returns false if the node is nil.
func HasAttrVal(n *html.Node, key string, val string) bool {
	v, err := GetAttrVal(n, key)
	return err == nil && v == val
}

// IsHtmlNode returns true if the specified node is an html node of the
// specified name (a, table, article), returns false if the node is nil.
func IsHtmlNode(n *html.Node, name string) bool {
	if n == nil {
		return false
	}
	return n.Type == html.ElementNode && n.Data == name
}

// IsTextNode returns true if the specified node is a text node. Returns false
// if the node is nil.
func IsTextNode(n *html.Node) bool {
	if n == nil {
		return false
	}
	return n.Type == html.TextNode
}

type WalkFunc func(n *html.Node)

// WalkAllNodes traverses all nodes in depth-first order executing the provided
// function on them.
func WalkAllNodes(root *html.Node, walkF WalkFunc) {
	var f func(*html.Node)
	f = func(n *html.Node) {
		walkF(n)
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(root)
}
